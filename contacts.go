package mixrank

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// CompanyTarget identifies an account, without relying on an ambiguous name.
type CompanyTarget struct {
	Name          string   `json:"name,omitempty"`
	NameAliases   []string `json:"name_aliases,omitempty"`
	CompanyIDs    []string `json:"company_ids,omitempty"`
	Domain        string   `json:"domain,omitempty"`
	DomainAliases []string `json:"domain_aliases,omitempty"`
	Qualification string   `json:"qualification_status,omitempty"`
}

type ContactOptions struct {
	Companies         []CompanyTarget          `json:"companies"`
	Roles             []string                 `json:"roles,omitempty"`
	MaxContacts       int                      `json:"max_contacts_per_company,omitempty"`
	MaxCandidates     int                      `json:"max_candidates_per_company,omitempty"`
	MaxRequests       int                      `json:"max_requests,omitempty"`
	EmailsOnly        bool                     `json:"emails_only,omitempty"`
	Concurrency       int                      `json:"concurrency,omitempty"`
	ContactFilter     string                   `json:"contact_filter,omitempty"`
	ValidateEmails    bool                     `json:"validate_emails,omitempty"`
	ValidationOptions ContactValidationOptions `json:"validation_options,omitempty"`
}

type BusinessEmail struct {
	Email             string         `json:"email"`
	Domain            string         `json:"domain"`
	Validation        string         `json:"validation"`
	ValidationDetails map[string]any `json:"validation_details,omitempty"`
}

type ContactEmployment struct {
	CompanyID    string `json:"company_id"`
	CompanyName  string `json:"company_name"`
	Domain       string `json:"domain"`
	Title        string `json:"title"`
	IsCurrent    bool   `json:"is_current"`
	SourceLinked *bool  `json:"has_source_company_id"`
}

type BusinessContact struct {
	HasEmail         bool              `json:"has_email"`
	HasPhone         bool              `json:"has_phone"`
	EnrichmentStatus string            `json:"enrichment_status"`
	PersonID         string            `json:"person_id"`
	Name             string            `json:"name"`
	Title            string            `json:"title"`
	LinkedInURL      string            `json:"linkedin_url,omitempty"`
	Employment       ContactEmployment `json:"employment"`
	ProfileUpdatedAt string            `json:"profile_updated_at,omitempty"`
	BusinessEmails   []BusinessEmail   `json:"business_emails"`
	DirectDials      []string          `json:"direct_dials"`
	EmailStatus      string            `json:"email_status"`
	PhoneStatus      string            `json:"phone_status"`
	NeedsReview      bool              `json:"needs_review"`
	ReviewReasons    []string          `json:"review_reasons"`
	Source           string            `json:"source"`
}

type ContactIssue struct {
	PersonID string `json:"person_id,omitempty"`
	Reason   string `json:"reason"`
}

type CompanyContactResult struct {
	Company            CompanyTarget     `json:"company"`
	Status             string            `json:"status"`
	PeopleMatched      int               `json:"people_matched"`
	CandidatesExamined int               `json:"candidates_examined"`
	SearchLimited      bool              `json:"search_limited"`
	Contacts           []BusinessContact `json:"contacts"`
	ContactsFiltered   int               `json:"contacts_filtered"`
	Issues             []ContactIssue    `json:"issues"`
}

type ContactReport struct {
	CompanyMerge         *CompanyMergeSummary     `json:"company_merge,omitempty"`
	ContactFilterApplied bool                     `json:"contact_filter_applied"`
	Validation           *ContactValidationReport `json:"email_validation,omitempty"`
	Concurrency          int                      `json:"concurrency"`
	ContactFilter        string                   `json:"contact_filter"`
	RetrievedAt          string                   `json:"retrieved_at_utc"`
	Companies            []CompanyContactResult   `json:"companies"`
	RequestsMade         int                      `json:"requests_made"`
	Complete             bool                     `json:"complete"`
	Roles                []string                 `json:"roles"`
	Notes                []string                 `json:"notes"`
}

var defaultContactRoles = []string{"owner", "founder", "CEO", "president", "general manager", "operations manager", "office manager", "marketing director"}

// CompanyContacts discovers current relevant people and appends available B2B
// emails and direct dials with bounded parallel requests. All discovered people
// share one contacts array; availability and enrichment status stay explicit.
// Validation is opt-in and uses a deduplicated bulk job, never parallel one-offs.
func (c *Client) CompanyContacts(ctx context.Context, opts ContactOptions) (*ContactReport, error) {
	inputCompanies := len(opts.Companies)
	opts.Roles = append([]string(nil), opts.Roles...)
	if err := normalizeContactOptions(&opts); err != nil {
		return nil, err
	}
	r := &ContactReport{RetrievedAt: time.Now().UTC().Format(time.RFC3339), Complete: true, Roles: opts.Roles, Concurrency: opts.Concurrency, ContactFilter: opts.ContactFilter,
		CompanyMerge: &CompanyMergeSummary{InputRecords: inputCompanies, UniqueCompanies: len(opts.Companies), DuplicatesMerged: inputCompanies - len(opts.Companies)},
		Notes:        []string{"Contacts are provider-reported, not independently verified. B2B emails are restricted to the matched employer domain; unrelated work emails are omitted.", "Email validation is separate from contact availability. Direct dials are person-level numbers; line type and company association are not verified.", "Current employment is a cached observation. Missing, old, inferred or conflicting evidence is marked for review.", "All discovered people use contacts, including unavailable and not-enriched rows. Filters affect returned rows, not company-level coverage or errors."}}
	for _, co := range opts.Companies {
		r.Companies = append(r.Companies, CompanyContactResult{Company: co, Status: "not_processed", Contacts: []BusinessContact{}, Issues: []ContactIssue{}})
	}
	runCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	work := &contactWork{client: c, slots: make(chan struct{}, opts.Concurrency), max: int64(opts.MaxRequests), cancel: cancel}
	var wg sync.WaitGroup
	var next atomic.Int64
	for n := 0; n < min(opts.Concurrency, len(r.Companies)); n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(next.Add(1)) - 1
				if i >= len(r.Companies) {
					return
				}
				c.contactCompany(runCtx, &r.Companies[i], opts, work)
			}
		}()
	}
	wg.Wait()
	r.RequestsMade = int(work.calls.Load())
	for i := range r.Companies {
		cr := &r.Companies[i]
		if len(cr.Issues) > 0 || cr.Status == "not_processed" || cr.Status == "request_budget_exhausted" {
			r.Complete = false
		}
	}
	err := context.Cause(runCtx)
	if opts.ValidateEmails && err == nil {
		vo := opts.ValidationOptions
		remaining := opts.MaxRequests - r.RequestsMade
		if remaining < 1 {
			r.Complete = false
			r.Validation = &ContactValidationReport{Status: "not_started", Issues: []string{"request budget exhausted before email validation"}}
		} else {
			if vo.MaxRequests == 0 || vo.MaxRequests > remaining {
				vo.MaxRequests = remaining
			}
			err = c.ValidateContacts(ctx, r, vo)
		}
	}
	filterContactReport(r, opts.ContactFilter)
	return r, err
}

var errContactBudget = errors.New("request budget exhausted")
var errContactNotStarted = errors.New("queued request cancelled before it started")

type contactWork struct {
	client *Client
	slots  chan struct{}
	calls  atomic.Int64
	max    int64
	cancel context.CancelCauseFunc
}

func (w *contactWork) request(ctx context.Context, op string, in Request) (map[string]any, error) {
	select {
	case w.slots <- struct{}{}:
	case <-ctx.Done():
		return nil, errors.Join(errContactNotStarted, context.Cause(ctx))
	}
	defer func() { <-w.slots }()
	if err := ctx.Err(); err != nil {
		return nil, errors.Join(errContactNotStarted, context.Cause(ctx))
	}
	for {
		used := w.calls.Load()
		if used >= w.max {
			return nil, errContactBudget
		}
		if w.calls.CompareAndSwap(used, used+1) {
			break
		}
	}
	v, err := w.client.contactJSON(ctx, op, in)
	if stopContactWorkflow(err) {
		w.cancel(err)
	}
	return v, err
}
func (c *Client) contactCompany(ctx context.Context, cr *CompanyContactResult, opts ContactOptions, work *contactWork) {
	if ctx.Err() != nil {
		return
	}
	body, _ := json.Marshal(contactSearch(cr.Company, opts))
	v, err := work.request(ctx, "search_person2", Request{Body: body})
	if err != nil {
		cr.Status = "search_error"
		if errors.Is(err, errContactBudget) {
			cr.Status = "request_budget_exhausted"
		}
		cr.Issues = append(cr.Issues, ContactIssue{Reason: c.Redact(err.Error())})
		return
	}
	hits := object(v["hits"])
	if hits == nil {
		cr.Status = "search_error"
		cr.Issues = append(cr.Issues, ContactIssue{Reason: "provider response has no hits object"})
		return
	}
	cr.PeopleMatched = integer(object(hits["total"])["value"])
	if _, ok := hits["total"].(json.Number); ok {
		cr.PeopleMatched = integer(hits["total"])
	}
	cr.SearchLimited = cr.PeopleMatched > len(array(hits["hits"])) || str(object(hits["total"])["relation"]) == "gte"
	if timed, _ := v["timed_out"].(bool); timed || integer(object(v["_shards"])["failed"]) > 0 {
		cr.SearchLimited = true
		cr.Issues = append(cr.Issues, ContactIssue{Reason: "search timed out or some shards failed"})
	}
	cr.Contacts = contactCandidates(cr, opts, v)
	for i := range cr.Contacts {
		cr.Contacts[i].EnrichmentStatus = "not_enriched"
		cr.Contacts[i].EmailStatus = "not_requested"
		cr.Contacts[i].PhoneStatus = "not_requested"
	}
	available := 0
	for start := 0; start < len(cr.Contacts) && available < opts.MaxContacts && ctx.Err() == nil; {
		width := min(opts.Concurrency, opts.MaxContacts-available, len(cr.Contacts)-start)
		errs := make([]error, width)
		var wg sync.WaitGroup
		for j := 0; j < width; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				errs[j] = c.enrichContact(ctx, cr.Company, &cr.Contacts[start+j], opts, work)
			}(j)
		}
		wg.Wait()
		budget := false
		for j, e := range errs {
			b := &cr.Contacts[start+j]
			if b.EnrichmentStatus != "not_enriched" {
				cr.CandidatesExamined++
			}
			if e != nil {
				cr.Issues = append(cr.Issues, ContactIssue{b.PersonID, c.Redact(e.Error())})
				budget = budget || errors.Is(e, errContactBudget)
			}
			if contactMatches(*b, enrichmentFilter(opts.ContactFilter)) {
				available++
			}
		}
		start += width
		if budget {
			break
		}
	}
	kept := []BusinessContact{}
	for _, b := range cr.Contacts {
		if b.EnrichmentStatus == "redacted" {
			cr.Issues = append(cr.Issues, ContactIssue{b.PersonID, "privacy-redacted overview excluded"})
		} else {
			kept = append(kept, b)
		}
	}
	cr.Contacts = kept
	cr.Status = "no_matching_people"
	if len(cr.Contacts) > 0 {
		cr.Status = "no_contact_data"
	}
	for _, b := range cr.Contacts {
		if b.HasEmail || b.HasPhone {
			if cr.Status != "contacts_found" {
				cr.Status = "review_required"
			}
			if !b.NeedsReview {
				cr.Status = "contacts_found"
			}
		}
	}
	if len(cr.Issues) > 0 && cr.Status != "contacts_found" && cr.Status != "review_required" {
		cr.Status = "incomplete"
	}
}
func contactCandidates(cr *CompanyContactResult, opts ContactOptions, v map[string]any) []BusinessContact {
	candidates := []BusinessContact{}
	seen := map[string]bool{}
	entries := array(object(v["hits"])["hits"])
	if len(entries) > opts.MaxCandidates {
		entries = entries[:opts.MaxCandidates]
		cr.SearchLimited = true
	}
	for _, entry := range entries {
		h := object(entry)
		src := object(h["_source"])
		pid := str(src["person_id"])
		if pid == "" {
			pid = str(h["_id"])
		}
		if seen[pid] || pid == "" {
			continue
		}
		seen[pid] = true
		if redacted, _ := src["privacy_redact"].(bool); redacted {
			cr.Issues = append(cr.Issues, ContactIssue{pid, "privacy-redacted record excluded"})
			continue
		}
		matches := array(object(object(object(h["inner_hits"])["experience"])["hits"])["hits"])
		var chosen map[string]any
		for _, m := range matches {
			e := object(object(m)["_source"])
			current, _ := e["is_current"].(bool)
			if !current || str(e["end_date"]) != "" || !matchesContactCompany(cr.Company, e) || !matchesRole(str(e["title"]), opts.Roles) {
				continue
			}
			if chosen == nil || roleRank(str(e["title"]), opts.Roles) < roleRank(str(chosen["title"]), opts.Roles) {
				chosen = e
			}
		}
		if chosen == nil {
			cr.Issues = append(cr.Issues, ContactIssue{pid, "no consistent current employer and role evidence in the same experience record"})
			continue
		}
		b := BusinessContact{PersonID: pid, Name: str(object(src["name"])["full"]), Title: str(chosen["title"]), LinkedInURL: str(src["public_profile_url"]), ProfileUpdatedAt: str(src["updated_at"]), BusinessEmails: []BusinessEmail{}, DirectDials: []string{}, ReviewReasons: []string{}, Source: "MixRank person2 search + person overview (b2b_emails,directdials)", Employment: ContactEmployment{CompanyID: str(chosen["company_id"]), CompanyName: str(chosen["company_name"]), Domain: cleanDomain(str(chosen["domain"])), Title: str(chosen["title"]), IsCurrent: true}}
		if linked, ok := chosen["has_source_company_id"].(bool); ok {
			b.Employment.SourceLinked = &linked
		}
		if b.Employment.SourceLinked == nil || !*b.Employment.SourceLinked {
			b.ReviewReasons = append(b.ReviewReasons, "employer association is inferred or source linkage is unknown")
		}
		head := strings.ToLower(str(src["headline"]))
		coName := strings.ToLower(b.Employment.CompanyName)
		if (strings.Contains(head, "former") || strings.Contains(head, "retired")) && coName != "" && strings.Contains(head, coName) {
			b.ReviewReasons = append(b.ReviewReasons, "headline conflicts with current-employment flag")
		}
		if b.Name == "" || strings.EqualFold(b.Name, b.Employment.CompanyName) {
			b.ReviewReasons = append(b.ReviewReasons, "record does not establish a named individual")
		}
		if staleContactDate(b.ProfileUpdatedAt) {
			b.ReviewReasons = append(b.ReviewReasons, "profile update date is missing, unrecognized, or older than one year")
		}
		if cr.Company.Qualification == "needs_review" {
			b.ReviewReasons = append(b.ReviewReasons, "input company qualification needs review")
		}
		b.NeedsReview = len(b.ReviewReasons) > 0
		candidates = append(candidates, b)
	}
	sort.SliceStable(candidates, func(a, b int) bool {
		if candidates[a].NeedsReview != candidates[b].NeedsReview {
			return !candidates[a].NeedsReview
		}
		return roleRank(candidates[a].Title, opts.Roles) < roleRank(candidates[b].Title, opts.Roles)
	})

	return candidates

}
func (c *Client) enrichContact(ctx context.Context, company CompanyTarget, candidate *BusinessContact, opts ContactOptions, work *contactWork) error {
	enable := "b2b_emails"
	if !opts.EmailsOnly {
		enable += ",directdials"
	}
	profile, err := work.request(ctx, "get_person_by_id", Request{Parameters: map[string]string{"id": candidate.PersonID, "enable": enable}})
	if err != nil {
		if !errors.Is(err, errContactBudget) && !errors.Is(err, errContactNotStarted) {
			candidate.EnrichmentStatus = "error"
			var api *APIError
			if errors.As(err, &api) && api.Uncertain {
				candidate.EnrichmentStatus = "uncertain"
			}
			candidate.EmailStatus = "error"
			candidate.PhoneStatus = "error"
		}
		return err
	}
	if id := str(profile["id"]); id != "" && id != candidate.PersonID {
		candidate.EnrichmentStatus = "error"
		return errors.New("person overview ID conflicts with search ID")
	}
	if redacted, _ := profile["privacy_redact"].(bool); redacted {
		candidate.EnrichmentStatus = "redacted"
		candidate.EmailStatus = "redacted"
		candidate.PhoneStatus = "redacted"
		return nil
	}
	domains := companyDomains(company)
	if len(domains) == 0 && candidate.Employment.Domain != "" {
		domains = []string{candidate.Employment.Domain}
	}
	candidate.EmailStatus = contactFieldStatus(profile, "b2b_emails")
	seenEmail := map[string]bool{}
	for _, v := range array(profile["b2b_emails"]) {
		email := strings.TrimSpace(str(object(v)["email"]))
		parsed, e := mail.ParseAddress(email)
		if e != nil || parsed.Address != email {
			continue
		}
		_, emailDomain, ok := strings.Cut(email, "@")
		emailDomain = cleanDomain(emailDomain)
		if !ok || !containsString(domains, emailDomain) || seenEmail[strings.ToLower(email)] {
			continue
		}
		seenEmail[strings.ToLower(email)] = true
		candidate.BusinessEmails = append(candidate.BusinessEmails, BusinessEmail{Email: email, Domain: emailDomain, Validation: "not_run"})
	}
	if candidate.EmailStatus == "available" && len(candidate.BusinessEmails) == 0 {
		candidate.EmailStatus = "no_email_for_company_domain"
	}
	candidate.PhoneStatus = "not_requested"
	if !opts.EmailsOnly {
		candidate.PhoneStatus = contactFieldStatus(profile, "directdials")
		seenPhone := map[string]bool{}
		for _, v := range array(profile["directdials"]) {
			phone := strings.TrimSpace(str(v))
			if phone != "" && !seenPhone[phone] {
				candidate.DirectDials = append(candidate.DirectDials, phone)
				seenPhone[phone] = true
			}
		}
	}

	candidate.HasEmail = len(candidate.BusinessEmails) > 0
	candidate.HasPhone = len(candidate.DirectDials) > 0
	candidate.EnrichmentStatus = "complete"
	return nil
}
func enrichmentFilter(filter string) string {
	switch filter {
	case "email", "valid-email":
		return "email"
	case "phone", "both":
		return filter
	default:
		return "any"
	}
}
func contactMatches(b BusinessContact, filter string) bool {
	switch filter {
	case "all":
		return true
	case "any":
		return b.HasEmail || b.HasPhone
	case "email":
		return b.HasEmail
	case "phone":
		return b.HasPhone
	case "both":
		return b.HasEmail && b.HasPhone
	case "none":
		return b.EnrichmentStatus == "complete" && !b.HasEmail && !b.HasPhone
	case "valid-email":
		for _, e := range b.BusinessEmails {
			if e.Validation == "valid" {
				return true
			}
		}
	}
	return false
}
func filterContactReport(r *ContactReport, filter string) {
	r.ContactFilter = filter
	r.ContactFilterApplied = true
	if r.Validation != nil && r.Validation.Status != "completed" && r.Validation.Status != "not_needed" && filter != "all" {
		r.ContactFilterApplied = false
		return
	}
	for i := range r.Companies {
		cr := &r.Companies[i]
		kept := []BusinessContact{}
		for _, b := range cr.Contacts {
			if contactMatches(b, filter) {
				kept = append(kept, b)
			} else {
				cr.ContactsFiltered++
			}
		}
		cr.Contacts = kept
	}
}

func (c *Client) contactJSON(ctx context.Context, op string, in Request) (map[string]any, error) {
	r, e := c.Call(ctx, op, in)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("contact workflow received HTTP %d; result is not a completed profile", r.StatusCode)
	}
	v, e := r.JSON(4 << 20)
	if e != nil {
		return nil, e
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("contact workflow expected a JSON object")
	}
	return obj, nil
}

func stopContactWorkflow(e error) bool {
	var api *APIError
	return errors.Is(e, context.Canceled) || errors.Is(e, context.DeadlineExceeded) || (errors.As(e, &api) && (api.Status == 401 || api.Status == 403 || api.Status == 429))
}
func object(v any) map[string]any { m, _ := v.(map[string]any); return m }
func array(v any) []any           { a, _ := v.([]any); return a }
func str(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	}
	return ""
}
func integer(v any) int { n, _ := strconv.Atoi(str(v)); return n }
func cleanDomain(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if !strings.Contains(s, "://") {
		s = "https://" + s
	}
	u, e := url.Parse(s)
	if e != nil {
		return ""
	}
	return strings.TrimPrefix(strings.TrimSuffix(u.Hostname(), "."), "www.")
}
func matchesContactCompany(c CompanyTarget, e map[string]any) bool {
	id := str(e["company_id"])
	for _, v := range c.CompanyIDs {
		if id == v {
			return true
		}
	}
	return containsString(companyDomains(c), cleanDomain(str(e["domain"])))
}
func matchesRole(title string, roles []string) bool { return roleRank(title, roles) < len(roles) }
func roleRank(title string, roles []string) int {
	title = strings.ToLower(title)
	for i, r := range roles {
		if strings.Contains(title, strings.ToLower(r)) {
			return i
		}
	}
	return len(roles)
}
func staleContactDate(s string) bool {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999999", "2006-01-02"} {
		if t, e := time.Parse(layout, s); e == nil {
			return time.Since(t) > 365*24*time.Hour
		}
	}
	return true
}
func contactFieldStatus(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return "not_returned"
	}
	if v == nil {
		return "unavailable"
	}
	if a, ok := v.([]any); ok {
		if len(a) == 0 {
			return "unavailable"
		}
		return "available"
	}
	if _, ok := v.(map[string]any); ok {
		return "withheld_or_disabled"
	}
	return "unrecognized_response"
}

func normalizeContactOptions(o *ContactOptions) error {
	if o.Concurrency == 0 {
		o.Concurrency = 4
	}
	if o.Concurrency < 1 || o.Concurrency > 16 {
		return errors.New("concurrency must be 1..16")
	}
	if o.ContactFilter == "" {
		o.ContactFilter = "all"
	}
	switch o.ContactFilter {
	case "all", "any", "email", "phone", "both", "none", "valid-email":
	default:
		return errors.New("contact-filter must be all, any, email, phone, both, none, or valid-email")
	}
	if o.ContactFilter == "valid-email" && !o.ValidateEmails {
		return errors.New("valid-email filter requires email validation")
	}
	if o.EmailsOnly && (o.ContactFilter == "phone" || o.ContactFilter == "both") {
		return errors.New("phone filters cannot be combined with emails-only")
	}
	if o.ValidateEmails {
		if err := normalizeContactValidationOptions(&o.ValidationOptions); err != nil {
			return err
		}
	}
	if o.MaxContacts == 0 {
		o.MaxContacts = 2
	}
	if o.MaxCandidates == 0 {
		o.MaxCandidates = 10
	}
	if o.MaxRequests == 0 {
		o.MaxRequests = 100
	}
	if o.MaxContacts < 1 || o.MaxContacts > 5 || o.MaxCandidates < o.MaxContacts || o.MaxCandidates > 25 || o.MaxRequests < 1 || o.MaxRequests > 250 {
		return errors.New("contact limits: contacts 1..5, candidates contacts..25, requests 1..250")
	}
	if len(o.Roles) == 0 {
		o.Roles = append([]string{}, defaultContactRoles...)
	}
	if len(o.Roles) > 20 {
		return errors.New("maximum 20 role phrases")
	}
	for _, r := range o.Roles {
		if strings.TrimSpace(r) == "" || len(r) > 100 {
			return errors.New("role phrases must contain 1..100 characters")
		}
	}
	companies, err := MergeCompanyTargets(o.Companies)
	if err != nil {
		return err
	}
	if len(companies) > 25 {
		return errors.New("contacts supports at most 25 unique companies after merging")
	}
	o.Companies = companies
	return nil
}

func contactSearch(co CompanyTarget, o ContactOptions) map[string]any {
	identifiers := []any{}
	if len(co.CompanyIDs) > 0 {
		identifiers = append(identifiers, map[string]any{"terms": map[string]any{"experience.company_id": co.CompanyIDs}})
	}
	if domains := companyDomains(co); len(domains) > 0 {
		identifiers = append(identifiers, map[string]any{"terms": map[string]any{"experience.domain.keyword": domains}})
	}
	roles := []any{}
	for i, r := range o.Roles {
		roles = append(roles, map[string]any{"match_phrase": map[string]any{"experience.title": map[string]any{"query": r, "boost": len(o.Roles) - i}}})
	}
	nested := map[string]any{"path": "experience", "query": map[string]any{"bool": map[string]any{"filter": []any{map[string]any{"term": map[string]any{"experience.is_current": true}}, map[string]any{"bool": map[string]any{"should": identifiers, "minimum_should_match": 1}}}, "must": []any{map[string]any{"bool": map[string]any{"should": roles, "minimum_should_match": 1}}}, "must_not": []any{map[string]any{"exists": map[string]any{"field": "experience.end_date"}}}}}, "inner_hits": map[string]any{"size": 10, "_source": []string{"experience.company_id", "experience.company_name", "experience.domain", "experience.title", "experience.is_current", "experience.has_source_company_id", "experience.end_date"}}}
	return map[string]any{"size": o.MaxCandidates, "track_total_hits": true, "_source": []string{"person_id", "name", "headline", "public_profile_url", "updated_at", "privacy_redact"}, "query": map[string]any{"bool": map[string]any{"must": []any{map[string]any{"nested": nested}}, "must_not": []any{map[string]any{"term": map[string]any{"privacy_redact": true}}}}}, "sort": []any{map[string]any{"_score": "desc"}, map[string]any{"person_id": "asc"}}}
}

// ParseCompanyTargets accepts a companies array, a {companies: [...]} export,
// or an Elasticsearch company-search response. JSON numbers retain precision.
func ParseCompanyTargets(reader io.Reader) ([]CompanyTarget, error) {
	b, e := io.ReadAll(io.LimitReader(reader, 1<<20+1))
	if e != nil || len(b) > 1<<20 {
		return nil, errors.New("company input must fit within 1 MiB")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var v any
	if e = d.Decode(&v); e != nil {
		return nil, errors.New("invalid company JSON")
	}
	if e = d.Decode(new(any)); e != io.EOF {
		return nil, errors.New("company input must contain one JSON value")
	}
	entries := array(v)
	if obj := object(v); obj != nil {
		if a, ok := obj["companies"].([]any); ok {
			entries = a
		} else {
			entries = array(object(obj["hits"])["hits"])
		}
	}
	if len(entries) == 0 {
		return nil, errors.New("expected a company array, companies export, or Elasticsearch hits")
	}
	result := []CompanyTarget{}
	for _, v := range entries {
		m := object(v)
		if src := object(m["_source"]); src != nil {
			m = src
		}
		if co := object(m["company"]); co != nil {
			m = co
		}
		co := CompanyTarget{Name: str(m["name"]), Domain: str(m["domain"]), Qualification: str(m["qualification_status"])}
		for _, v := range array(m["name_aliases"]) {
			co.NameAliases = append(co.NameAliases, str(v))
		}
		for _, v := range array(m["domain_aliases"]) {
			co.DomainAliases = append(co.DomainAliases, str(v))
		}
		if co.Name == "" {
			co.Name = str(m["company_name"])
		}
		switch ids := m["company_ids"].(type) {
		case string:
			co.CompanyIDs = strings.FieldsFunc(ids, func(r rune) bool { return r == ';' || r == ',' })
		case []any:
			for _, id := range ids {
				co.CompanyIDs = append(co.CompanyIDs, str(id))
			}
		}
		if len(co.CompanyIDs) == 0 {
			id := str(m["company_id"])
			if id == "" {
				id = str(m["id"])
			}
			if id != "" {
				co.CompanyIDs = []string{id}
			}
		}
		result = append(result, co)
	}
	return result, nil
}
