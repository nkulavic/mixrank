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
