package mixrank

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ContactValidationOptions struct {
	Strategy    string        `json:"strategy,omitempty"`
	MaxAge      string        `json:"maxage,omitempty"`
	Wait        time.Duration `json:"-"` // Zero submits/checks once; positive waits at most this duration.
	MaxRequests int           `json:"max_requests,omitempty"`
}
type ContactValidationReport struct {
	EmailFingerprint string   `json:"email_fingerprint,omitempty"`
	JobID            string   `json:"job_id,omitempty"`
	Status           string   `json:"status"`
	Strategy         string   `json:"strategy,omitempty"`
	EmailsSubmitted  int      `json:"emails_submitted"`
	ResultsApplied   int      `json:"results_applied"`
	CompletedAt      string   `json:"completed_at,omitempty"`
	Issues           []string `json:"issues"`
}

func normalizeContactValidationOptions(o *ContactValidationOptions) error {
	if o.Strategy == "" {
		o.Strategy = "besteffort"
	}
	switch o.Strategy {
	case "cached", "fetch", "strict", "besteffort":
	default:
		return errors.New("validation strategy must be cached, fetch, strict, or besteffort")
	}
	if o.MaxAge != "" {
		n, e := strconv.ParseFloat(o.MaxAge, 64)
		if e != nil || n < 0 {
			return errors.New("validation maxage must be nonnegative seconds")
		}
	}
	if o.Wait < 0 || o.Wait > 10*time.Minute {
		return errors.New("validation wait must be 0..10m")
	}
	if o.MaxRequests == 0 {
		o.MaxRequests = 30
	}
	if o.MaxRequests < 1 || o.MaxRequests > 250 {
		return errors.New("validation max requests must be 1..250")
	}
	return nil
}

// ValidateContacts attaches bulk deliverability evidence to existing contacts.
// It mutates report in place; callers must not use the report concurrently.
// A retained job ID resumes polling/download without submitting another job.
// Pending results are successful partial progress, not validation success.
func (c *Client) ValidateContacts(ctx context.Context, report *ContactReport, opts ContactValidationOptions) error {
	return c.validateContacts(ctx, report, opts, Download)
}
func (c *Client) validateContacts(ctx context.Context, r *ContactReport, o ContactValidationOptions, download func(context.Context, string, io.Writer, int64) error) error {
	if r == nil {
		return errors.New("contact report is required")
	}
	if err := normalizeContactValidationOptions(&o); err != nil {
		return err
	}
	refs := map[string][]*BusinessEmail{}
	emails := []string{}
	for i := range r.Companies {
		for j := range r.Companies[i].Contacts {
			for k := range r.Companies[i].Contacts[j].BusinessEmails {
				e := &r.Companies[i].Contacts[j].BusinessEmails[k]
				key := strings.ToLower(strings.TrimSpace(e.Email))
				if key == "" || strings.ContainsAny(key, "\r\n") {
					return errors.New("invalid email in contact report")
				}
				if _, ok := refs[key]; !ok {
					emails = append(emails, key)
				}
				refs[key] = append(refs[key], e)
			}
		}
	}
	if len(emails) > 2500 {
		return errors.New("contact validation supports at most 2500 unique emails per report")
	}
	if len(emails) == 0 && r.Validation != nil && r.Validation.JobID != "" {
		return errors.New("pending job report has no emails; restore its original contact data before resuming")
	}
	sort.Strings(emails)
	fingerprint := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(emails, "\n"))))
	if r.Validation != nil && r.Validation.Status != "completed" && r.Validation.EmailFingerprint != "" && r.Validation.EmailFingerprint != fingerprint {
		return errors.New("contact emails differ from the submitted job; restore the original report before resuming")
	}
	if r.Validation != nil && r.Validation.Status == "completed" {
		return nil
	}
	if r.Validation != nil && r.Validation.JobID == "" && r.Validation.Status == "submission_uncertain" {
		return errors.New("bulk submission outcome is uncertain; recover its job ID from the provider job list before resuming")
	}
	if r.Validation == nil || r.Validation.JobID == "" {
		r.Validation = &ContactValidationReport{Status: "not_started", Strategy: o.Strategy, EmailsSubmitted: len(emails), Issues: []string{}}
	}
	v := r.Validation
	v.EmailFingerprint = fingerprint
	if len(emails) == 0 {
		v.Status = "not_needed"
		return nil
	}
	used := 0
	reserve := func() error {
		if used >= o.MaxRequests {
			return errContactBudget
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		used++
		r.RequestsMade++
		return nil
	}
	fail := func(err error) error {
		r.Complete = false
		v.Issues = append(v.Issues, c.Redact(err.Error()))
		return err
	}
	var job map[string]any
	if v.JobID == "" {
		if err := reserve(); err != nil {
			return fail(err)
		}
		params := map[string]string{"strategy": o.Strategy}
		if o.MaxAge != "" {
			params["maxage"] = o.MaxAge
		}
		response, err := c.Call(ctx, "post_email_validate_bulk_job", Request{Parameters: params, Upload: strings.NewReader(strings.Join(emails, "\n") + "\n"), Filename: "contacts.txt"})
		if err != nil {
			var api *APIError
			if errors.As(err, &api) && api.Uncertain {
				v.Status = "submission_uncertain"
			} else {
				v.Status = "error"
			}
			return fail(err)
		}
		raw, err := response.JSON(1 << 20)
		response.Body.Close()
		if err != nil {
			v.Status = "submission_uncertain"
			return fail(errors.New("bulk submission accepted but response unreadable; recover job ID before resuming"))
		}
		job = object(raw)
		v.JobID = str(job["id"])
		if v.JobID == "" {
			v.Status = "submission_uncertain"
			return fail(errors.New("bulk submission returned no job ID; do not resubmit blindly"))
		}
		v.Status = "pending"
		r.Complete = false
		for _, list := range refs {
			for _, e := range list {
				e.Validation = "pending"
				e.ValidationDetails = nil
			}
		}
	} else {
		if err := reserve(); err != nil {
			return fail(err)
		}
		var err error
		job, err = c.contactJSON(ctx, "get_email_validate_bulk_job_by_job_id", Request{Parameters: map[string]string{"job_id": v.JobID}})
		if err != nil {
			return fail(err)
		}
	}
	deadline := time.Now().Add(o.Wait)
	pollCtx := ctx
	if o.Wait > 0 {
		var cancel context.CancelFunc
		pollCtx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}
	for str(job["completed_at"]) == "" || str(job["download_url"]) == "" {
		v.Status = "pending"
		r.Complete = false
		if o.Wait == 0 || time.Now().After(deadline) {
			return nil
		}
		delay := min(2*time.Second, time.Until(deadline))
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fail(ctx.Err())
		case <-timer.C:
		}
		if time.Now().After(deadline) {
			return nil
		}
		if err := reserve(); err != nil {
			return fail(err)
		}
		var err error
		job, err = c.contactJSON(pollCtx, "get_email_validate_bulk_job_by_job_id", Request{Parameters: map[string]string{"job_id": v.JobID}})
		if err != nil {
			if ctx.Err() == nil && pollCtx.Err() != nil {
				return nil
			}
			return fail(err)
		}
	}
	if err := reserve(); err != nil {
		return fail(err)
	}
	var buf bytes.Buffer
	if err := download(ctx, str(job["download_url"]), &buf, 8<<20); err != nil {
		v.Status = "download_error"
		return fail(err)
	}
	rows, err := parseContactValidation(buf.Bytes())
	if err != nil {
		v.Status = "download_error"
		return fail(err)
	}
	applied := map[string]bool{}
	for _, row := range rows {
		key := strings.ToLower(str(row["email"]))
		if len(refs[key]) == 0 {
			continue
		}
		validity := str(row["validity"])
		switch validity {
		case "valid", "ambiguous", "maybe_valid", "invalid", "uncached", "malformed", "temp_error", "perm_error":
		default:
			continue
		}
		if applied[key] {
			v.Status = "download_error"
			return fail(errors.New("duplicate validation result for an email; inspect job output"))
		}
		applied[key] = true
		for _, e := range refs[key] {
			e.Validation = validity
			e.ValidationDetails = row
		}
	}
	v.ResultsApplied = len(applied)
	v.CompletedAt = str(job["completed_at"])
	v.Status = "completed"
	for key, list := range refs {
		if !applied[key] {
			v.Status = "partial"
			for _, e := range list {
				e.Validation = "not_returned"
			}
		}
	}
	if v.Status == "partial" {
		return fail(errors.New("bulk output omitted usable results for some submitted emails"))
	}
	r.Complete = true
	for _, co := range r.Companies {
		if len(co.Issues) > 0 || co.Status == "not_processed" || co.Status == "request_budget_exhausted" || co.Status == "search_error" {
			r.Complete = false
		}
	}
	if r.ContactFilter != "" {
		filterContactReport(r, r.ContactFilter)
	}
	return nil
}

// Result files may be JSON arrays, JSONL, or CSV; gzip is detected by magic.
// Both compressed and decompressed input are bounded. Unknown shapes fail closed.
func parseContactValidation(data []byte) ([]map[string]any, error) {
	if len(data) > 8<<20 {
		return nil, errors.New("validation output exceeds 8 MiB")
	}
	if len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b {
		gz, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, errors.New("invalid compressed validation output")
		}
		defer gz.Close()
		data, err = io.ReadAll(io.LimitReader(gz, 8<<20+1))
		if err != nil || len(data) > 8<<20 {
			return nil, errors.New("decompressed validation output exceeds limit or is invalid")
		}
	}
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf}))
	if len(data) == 0 {
		return nil, errors.New("empty validation output")
	}
	rows := []map[string]any{}
	if data[0] == '[' {
		var items []map[string]any
		d := json.NewDecoder(bytes.NewReader(data))
		d.UseNumber()
		if err := d.Decode(&items); err != nil {
			return nil, errors.New("invalid validation JSON array")
		}
		if d.Decode(new(any)) != io.EOF {
			return nil, errors.New("trailing validation JSON")
		}
		rows = items
	} else if data[0] == '{' {
		scan := bufio.NewScanner(bytes.NewReader(data))
		scan.Buffer(make([]byte, 4096), 1<<20)
		for scan.Scan() {
			var row map[string]any
			d := json.NewDecoder(bytes.NewReader(scan.Bytes()))
			d.UseNumber()
			if d.Decode(&row) != nil || d.Decode(new(any)) != io.EOF {
				return nil, errors.New("invalid validation JSONL")
			}
			rows = append(rows, row)
		}
		if scan.Err() != nil {
			return nil, errors.New("validation row exceeds limit")
		}
	} else {
		reader := csv.NewReader(bytes.NewReader(data))
		header, err := reader.Read()
		if err != nil {
			return nil, errors.New("invalid validation CSV")
		}
		hasEmail, hasValidity := false, false
		for _, h := range header {
			hasEmail = hasEmail || h == "email"
			hasValidity = hasValidity || h == "validity"
		}
		if !hasEmail || !hasValidity {
			return nil, errors.New("validation CSV lacks email/validity columns")
		}
		for {
			values, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, errors.New("invalid validation CSV row")
			}
			row := map[string]any{}
			for i, h := range header {
				var val any = values[i]
				if values[i] == "true" {
					val = true
				} else if values[i] == "false" {
					val = false
				} else if values[i] == "" || values[i] == "null" {
					val = nil
				}
				row[h] = val
			}
			rows = append(rows, row)
		}
	}
	if len(rows) > 2500 {
		return nil, errors.New("validation output exceeds 2500 rows")
	}
	for _, row := range rows {
		if str(row["email"]) == "" {
			return nil, fmt.Errorf("validation output has no email identifier")
		}
	}
	return rows, nil
}
