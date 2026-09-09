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
	"time"
)

// CompanyTarget identifies an account, without relying on an ambiguous name.
type CompanyTarget struct {
	Name          string   `json:"name,omitempty"`
	CompanyIDs    []string `json:"company_ids,omitempty"`
	Domain        string   `json:"domain,omitempty"`
	Qualification string   `json:"qualification_status,omitempty"`
}

type ContactOptions struct {
	Companies     []CompanyTarget `json:"companies"`
	Roles         []string        `json:"roles,omitempty"`
	MaxContacts   int             `json:"max_contacts_per_company,omitempty"`
	MaxCandidates int             `json:"max_candidates_per_company,omitempty"`
	MaxRequests   int             `json:"max_requests,omitempty"`
	EmailsOnly    bool            `json:"emails_only,omitempty"`
}

type BusinessEmail struct {
	Email      string `json:"email"`
	Domain     string `json:"domain"`
	Validation string `json:"validation"`
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
	Company                   CompanyTarget     `json:"company"`
	Status                    string            `json:"status"`
	PeopleMatched             int               `json:"people_matched"`
	CandidatesExamined        int               `json:"candidates_examined"`
	SearchLimited             bool              `json:"search_limited"`
	Contacts                  []BusinessContact `json:"contacts"`
	CandidatesWithoutContacts []BusinessContact `json:"candidates_without_contacts"`
	Issues                    []ContactIssue    `json:"issues"`
}

type ContactReport struct {
	RetrievedAt  string                 `json:"retrieved_at_utc"`
	Companies    []CompanyContactResult `json:"companies"`
	RequestsMade int                    `json:"requests_made"`
	Complete     bool                   `json:"complete"`
	Roles        []string               `json:"roles"`
	Notes        []string               `json:"notes"`
}

var defaultContactRoles = []string{"owner", "founder", "CEO", "president", "general manager", "operations manager", "office manager", "marketing director"}

// CompanyContacts discovers current relevant people and appends available B2B
// email addresses and direct dials. It never guesses emails, validates them,
// changes to a live-fetch strategy, sends messages, or requests consumer emails.
// Partial results are returned together with cancellation or permission errors.
func (c *Client) CompanyContacts(ctx context.Context, opts ContactOptions) (*ContactReport, error) {
	opts.Companies = append([]CompanyTarget(nil), opts.Companies...)
	for i := range opts.Companies {
		opts.Companies[i].CompanyIDs = append([]string(nil), opts.Companies[i].CompanyIDs...)
	}
	opts.Roles = append([]string(nil), opts.Roles...)
	if err := normalizeContactOptions(&opts); err != nil {
		return nil, err
	}
	r := &ContactReport{RetrievedAt: time.Now().UTC().Format(time.RFC3339), Complete: true, Roles: opts.Roles,
		Notes: []string{"Contacts are provider-reported, not independently verified. B2B emails are restricted to the matched employer domain; unrelated work emails are omitted.", "Email validation was not run. Direct dials are person-level numbers; their line type and association with a specific company are not verified.", "Current employment is the provider's cached observation. Missing, old, inferred or conflicting evidence is marked for review.", "Candidate and request limits bound coverage. A missing result is not evidence that a business has no contact information."}}
	for _, co := range opts.Companies {
		r.Companies = append(r.Companies, CompanyContactResult{Company: co, Status: "not_processed", Contacts: []BusinessContact{}, CandidatesWithoutContacts: []BusinessContact{}, Issues: []ContactIssue{}})
	}
	for i := range r.Companies {
		cr := &r.Companies[i]
		if err := ctx.Err(); err != nil {
			r.Complete = false
			return r, err
		}
		if r.RequestsMade >= opts.MaxRequests {
			cr.Status = "request_budget_exhausted"
			r.Complete = false
			continue
		}
		body, _ := json.Marshal(contactSearch(cr.Company, opts))
		r.RequestsMade++
		v, err := c.contactJSON(ctx, "search_person2", Request{Body: body})
		if err != nil {
			cr.Status = "search_error"
			cr.Issues = append(cr.Issues, ContactIssue{Reason: c.Redact(err.Error())})
			r.Complete = false
			if stopContactWorkflow(err) {
				return r, err
			}
			continue
		}
		hits := object(v["hits"])
		if hits == nil {
			cr.Status = "search_error"
			cr.Issues = append(cr.Issues, ContactIssue{Reason: "provider response has no hits object"})
			r.Complete = false
			continue
		}
		cr.PeopleMatched = integer(object(hits["total"])["value"])
		if _, ok := hits["total"].(json.Number); ok {
			cr.PeopleMatched = integer(hits["total"])
		}
		entries := array(hits["hits"])
		cr.SearchLimited = cr.PeopleMatched > len(entries) || str(object(hits["total"])["relation"]) == "gte"
		if timed, _ := v["timed_out"].(bool); timed || integer(object(v["_shards"])["failed"]) > 0 {
			cr.SearchLimited = true
			cr.Issues = append(cr.Issues, ContactIssue{Reason: "search timed out or some shards failed"})
			r.Complete = false
		}
		candidates := []BusinessContact{}
		seen := map[string]bool{}
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
		for _, candidate := range candidates {
			if len(cr.Contacts) >= opts.MaxContacts {
				break
			}
			if err = ctx.Err(); err != nil {
				r.Complete = false
				return r, err
			}
			if r.RequestsMade >= opts.MaxRequests {
				cr.Issues = append(cr.Issues, ContactIssue{Reason: "request budget exhausted"})
				r.Complete = false
				break
			}
			cr.CandidatesExamined++
			r.RequestsMade++
			enable := "b2b_emails"
			if !opts.EmailsOnly {
				enable += ",directdials"
			}
			profile, e := c.contactJSON(ctx, "get_person_by_id", Request{Parameters: map[string]string{"id": candidate.PersonID, "enable": enable}})
			if e != nil {
				cr.Issues = append(cr.Issues, ContactIssue{candidate.PersonID, c.Redact(e.Error())})
				r.Complete = false
				if stopContactWorkflow(e) {
					cr.Status = "enrichment_error"
					return r, e
				}
				continue
			}
			if id := str(profile["id"]); id != "" && id != candidate.PersonID {
				cr.Issues = append(cr.Issues, ContactIssue{candidate.PersonID, "person overview ID conflicts with search ID"})
				r.Complete = false
				continue
			}
			if redacted, _ := profile["privacy_redact"].(bool); redacted {
				cr.Issues = append(cr.Issues, ContactIssue{candidate.PersonID, "overview is privacy redacted"})
				continue
			}
			domain := cr.Company.Domain
			if domain == "" {
				domain = candidate.Employment.Domain
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
				if !ok || domain == "" || emailDomain != domain || seenEmail[strings.ToLower(email)] {
					continue
				}
				seenEmail[strings.ToLower(email)] = true
				candidate.BusinessEmails = append(candidate.BusinessEmails, BusinessEmail{email, emailDomain, "not_run"})
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
			if len(candidate.BusinessEmails) > 0 || len(candidate.DirectDials) > 0 {
				cr.Contacts = append(cr.Contacts, candidate)
			} else {
				cr.CandidatesWithoutContacts = append(cr.CandidatesWithoutContacts, candidate)
			}
		}
		cr.Status = "no_matching_people"
		if len(candidates) > 0 {
			cr.Status = "no_contact_data"
		}
		if len(cr.Contacts) > 0 {
			cr.Status = "review_required"
			for _, b := range cr.Contacts {
				if !b.NeedsReview {
					cr.Status = "contacts_found"
					break
				}
			}
		}
		if len(cr.Issues) > 0 && len(cr.Contacts) == 0 {
			cr.Status = "incomplete"
		}
	}
	return r, nil
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
	return c.Domain != "" && cleanDomain(str(e["domain"])) == c.Domain
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
	if len(o.Companies) == 0 || len(o.Companies) > 25 {
		return errors.New("contacts requires 1..25 companies per run")
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
	for i := range o.Companies {
		co := &o.Companies[i]
		if co.Domain != "" {
			co.Domain = cleanDomain(co.Domain)
			if !strings.Contains(co.Domain, ".") {
				return errors.New("company domain must be a valid domain name")
			}
		}
		if len(co.CompanyIDs) == 0 && co.Domain == "" {
			return errors.New("each company needs company_ids or domain; name alone is ambiguous")
		}
		if len(co.CompanyIDs) > 10 {
			return errors.New("maximum 10 IDs per company")
		}
		seen := map[string]bool{}
		ids := []string{}
		for _, id := range co.CompanyIDs {
			n, e := strconv.ParseInt(id, 10, 64)
			if e != nil || n <= 0 {
				return errors.New("company IDs must be positive integers")
			}
			id = strconv.FormatInt(n, 10)
			if !seen[id] {
				ids = append(ids, id)
				seen[id] = true
			}
		}
		co.CompanyIDs = ids
	}
	return nil
}

func contactSearch(co CompanyTarget, o ContactOptions) map[string]any {
	identifiers := []any{}
	if len(co.CompanyIDs) > 0 {
		identifiers = append(identifiers, map[string]any{"terms": map[string]any{"experience.company_id": co.CompanyIDs}})
	}
	if co.Domain != "" {
		identifiers = append(identifiers, map[string]any{"term": map[string]any{"experience.domain.keyword": co.Domain}})
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
