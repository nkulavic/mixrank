package mixrank

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func validationContactReport() *ContactReport {
	return &ContactReport{Complete: true, ContactFilter: "all", Companies: []CompanyContactResult{{Status: "contacts_found", Contacts: []BusinessContact{
		{PersonID: "1", HasEmail: true, EnrichmentStatus: "complete", BusinessEmails: []BusinessEmail{{Email: "avery@example.test", Validation: "not_run"}}},
		{PersonID: "2", HasEmail: true, EnrichmentStatus: "complete", BusinessEmails: []BusinessEmail{{Email: "avery@example.test", Validation: "not_run"}, {Email: "second@example.test", Validation: "not_run"}}},
	}}}}
}
func TestContactValidationSubmitResumeAndEvidence(t *testing.T) {
	submits, checks, downloads := 0, 0, 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			submits++
			if r.URL.Query().Get("strategy") != "besteffort" {
				t.Error("wrong strategy")
			}
			f, _, err := r.FormFile("file")
			if err != nil {
				t.Error(err)
				return
			}
			defer f.Close()
			b, _ := io.ReadAll(f)
			if string(b) != "avery@example.test\nsecond@example.test\n" {
				t.Error("emails not deduplicated", string(b))
			}
			fmt.Fprint(w, `{"id":"job-1","completed_at":null,"download_url":null}`)
			return
		}
		checks++
		fmt.Fprint(w, `{"id":"job-1","completed_at":"2026-09-09T00:00:00Z","download_url":"https://example.test/results"}`)
	})
	download := func(ctx context.Context, u string, w io.Writer, max int64) error {
		downloads++
		fmt.Fprintln(w, `{"email":"avery@example.test","validity":"valid","validated_at":"2026-09-09T00:00:00Z","is_domain_catchall":false}`)
		fmt.Fprintln(w, `{"email":"second@example.test","validity":"ambiguous","uses_greylist":true,"retry_after":"2026-09-10T00:00:00Z"}`)
		return nil
	}
	report := validationContactReport()
	report.ContactFilter = "valid-email"
	if err := c.validateContacts(context.Background(), report, ContactValidationOptions{}, download); err != nil {
		t.Fatal(err)
	}
	if submits != 1 || checks != 0 || report.Complete || report.Validation.Status != "pending" || report.Validation.JobID != "job-1" {
		t.Fatal("pending job not preserved", report)
	}
	filterContactReport(report, "valid-email")
	if report.ContactFilterApplied || len(report.Companies[0].Contacts) != 2 {
		t.Fatal("pending filter dropped resumable input")
	}
	if err := c.validateContacts(context.Background(), report, ContactValidationOptions{}, download); err != nil {
		t.Fatal(err)
	}
	if submits != 1 || checks != 1 || downloads != 1 || !report.Complete || report.RequestsMade != 3 || report.Validation.ResultsApplied != 2 || !report.ContactFilterApplied {
		t.Fatal("job not resumed", report)
	}
	emails := report.Companies[0].Contacts[1].BusinessEmails
	if emails[0].Validation != "valid" || emails[1].Validation != "ambiguous" || emails[1].ValidationDetails["uses_greylist"] != true {
		t.Fatal("provider evidence lost", emails)
	}
	if err := c.validateContacts(context.Background(), report, ContactValidationOptions{}, download); err != nil || submits != 1 || checks != 1 {
		t.Fatal("completed job was replayed", err)
	}
}
func TestContactValidationUncertainAndChangedInput(t *testing.T) {
	submits := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) { submits++; fmt.Fprint(w, `{"unexpected":"accepted"}`) })
	r := validationContactReport()
	if e := c.ValidateContacts(context.Background(), r, ContactValidationOptions{}); e == nil || r.Validation.Status != "submission_uncertain" {
		t.Fatal("uncertain submission lost")
	}
	if e := c.ValidateContacts(context.Background(), r, ContactValidationOptions{}); e == nil || submits != 1 {
		t.Fatal("uncertain submission replayed")
	}
	r.Validation = &ContactValidationReport{JobID: "known", Status: "pending", EmailFingerprint: "wrong"}
	if e := c.ValidateContacts(context.Background(), r, ContactValidationOptions{}); e == nil || submits != 1 {
		t.Fatal("changed email set accepted")
	}
}
func TestContactValidationIncompleteOutputAndBudget(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":"job","completed_at":"2026-09-09","download_url":"https://example.test/results"}`)
	})
	r := validationContactReport()
	download := func(ctx context.Context, u string, w io.Writer, max int64) error {
		fmt.Fprint(w, `[{"email":"avery@example.test","validity":"invalid"}]`)
		return nil
	}
	if e := c.validateContacts(context.Background(), r, ContactValidationOptions{MaxRequests: 1}, download); e == nil || r.RequestsMade != 1 || r.Validation.JobID != "job" {
		t.Fatal("download bypassed budget", e)
	}
	if e := c.validateContacts(context.Background(), r, ContactValidationOptions{MaxRequests: 2}, download); e == nil || r.Validation.Status != "partial" || r.Complete || r.Companies[0].Contacts[1].BusinessEmails[1].Validation != "not_returned" {
		t.Fatal("missing validation assumed valid", e)
	}
}
func TestContactValidationCancellationKeepsJob(t *testing.T) {
	polling := make(chan struct{})
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Fprint(w, `{"id":"job","completed_at":null}`)
			return
		}
		close(polling)
		<-r.Context().Done()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := validationContactReport()
	done := make(chan error, 1)
	go func() { done <- c.ValidateContacts(ctx, r, ContactValidationOptions{Wait: 10 * time.Second}) }()
	select {
	case <-polling:
	case <-ctx.Done():
		t.Fatal("job did not reach polling")
	}
	cancel()
	if e := <-done; e == nil || r.Validation.JobID != "job" {
		t.Fatal("cancellation lost submitted job")
	}
}
func TestContactValidationFileFormats(t *testing.T) {
	line := `{"email":"avery@example.test","validity":"valid","is_domain_catchall":false}`
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	w.Write([]byte(line + "\n"))
	w.Close()
	for _, input := range [][]byte{[]byte(line + "\n"), []byte("[" + line + "]"), gz.Bytes(), []byte("email,validity,is_domain_catchall\navery@example.test,valid,false\n")} {
		rows, e := parseContactValidation(input)
		if e != nil || len(rows) != 1 || rows[0]["is_domain_catchall"] != false {
			t.Fatal(rows, e)
		}
	}
	for _, input := range []string{"garbage", `[{"validity":"valid"}]`, line + " trailing", strings.Repeat("x", (8<<20)+1)} {
		if _, e := parseContactValidation([]byte(input)); e == nil {
			t.Fatal("bad output accepted")
		}
	}
}
func TestContactValidationJSONShape(t *testing.T) {
	r := validationContactReport()
	b, _ := json.Marshal(r)
	if !bytes.Contains(b, []byte(`"has_email":true`)) {
		t.Fatal("availability absent")
	}
}

func TestContactValidationResumeRequiresContacts(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("unexpected request") })
	r := &ContactReport{Validation: &ContactValidationReport{JobID: "job", Status: "pending"}}
	if err := c.ValidateContacts(context.Background(), r, ContactValidationOptions{}); err == nil || r.Validation.Status != "pending" || r.Validation.EmailFingerprint != "" {
		t.Fatal("missing input was mistaken for completed validation")
	}
}
