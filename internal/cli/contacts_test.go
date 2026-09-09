package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nkulavic/mixrank"
)

func TestContactsCLIJSONAndNoOverwrite(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		fmt.Fprint(w, `{"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}`)
	}))
	defer s.Close()
	cl, _ := mixrank.New("synthetic", mixrank.Options{BaseURL: s.URL})
	factory := func(context.Context) (*mixrank.Client, error) { return cl, nil }
	dir := t.TempDir()
	input := filepath.Join(dir, "companies.json")
	out := filepath.Join(dir, "contacts.json")
	os.WriteFile(input, []byte(`{"companies":[{"company_id":123,"company_name":"Example","domain":"example.test"}]}`), 0600)
	cmd := newContactsCommand(factory)
	cmd.SetArgs([]string{"--companies", input, "--output", out})
	if e := cmd.Execute(); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(out)
	var result mixrank.ContactReport
	if json.Unmarshal(b, &result) != nil || result.Companies[0].Status != "no_matching_people" {
		t.Fatal(string(b))
	}
	cmd = newContactsCommand(factory)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--domain", "example.test", "--output", out})
	if e := cmd.Execute(); e == nil || calls != 1 {
		t.Fatal("existing file did not prevent calls")
	}
}

func TestContactsCLIValidationFlagsAndResume(t *testing.T) {
	submits, checks := 0, 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			submits++
			if r.URL.Query().Get("strategy") != "cached" {
				t.Error("validation strategy flag ignored")
			}
			fmt.Fprint(w, `{"id":"job","completed_at":null}`)
			return
		}
		checks++
		fmt.Fprint(w, `{"id":"job","completed_at":null}`)
	}))
	defer s.Close()
	cl, _ := mixrank.New("synthetic", mixrank.Options{BaseURL: s.URL})
	factory := func(context.Context) (*mixrank.Client, error) { return cl, nil }
	input := filepath.Join(t.TempDir(), "input.json")
	os.WriteFile(input, []byte(`{"complete":true,"companies":[{"contacts":[{"business_emails":[{"email":"avery@example.test"}]}]}]}`), 0600)
	cmd := newContactsCommand(factory)
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"validate", "--input", input, "--validation-wait", "0s", "--validation-strategy", "cached"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var report mixrank.ContactReport
	if json.Unmarshal(out.Bytes(), &report) != nil || report.Validation == nil || report.Validation.JobID != "job" || submits != 1 {
		t.Fatal("CLI validation did not submit a bulk job")
	}
	os.WriteFile(input, out.Bytes(), 0600)
	cmd = newContactsCommand(factory)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"validate", "--input", input, "--validation-wait", "0s"})
	if err := cmd.Execute(); err != nil || submits != 1 || checks != 1 {
		t.Fatal("CLI resume resubmitted", err)
	}
}
