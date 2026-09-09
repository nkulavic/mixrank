package mixrank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func contactHit(id, company, title string, linked bool) map[string]any {
	return map[string]any{"_id": id, "_source": map[string]any{"person_id": json.Number(id), "name": map[string]any{"full": "Avery " + id}, "updated_at": time.Now().UTC().Format(time.RFC3339), "public_profile_url": "https://linkedin.com/in/synthetic-" + id}, "inner_hits": map[string]any{"experience": map[string]any{"hits": map[string]any{"hits": []any{map[string]any{"_source": map[string]any{"company_id": json.Number(company), "company_name": "Example Roofing", "domain": "example.test", "title": title, "is_current": true, "has_source_company_id": linked}}}}}}}
}

func TestCompanyContactsFindsUsableBusinessChannels(t *testing.T) {
	var appended []string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			b, _ := io.ReadAll(r.Body)
			var q map[string]any
			if json.Unmarshal(b, &q) != nil {
				t.Fatal("invalid DSL")
			}
			n := q["query"].(map[string]any)["bool"].(map[string]any)["must"].([]any)[0].(map[string]any)["nested"].(map[string]any)
			if n["path"] != "experience" || !strings.Contains(string(b), `"experience.is_current":true`) {
				t.Error("current employer and role were not constrained together")
			}
			wrong := contactHit("3", "999", "Owner", true)
			wrong["inner_hits"].(map[string]any)["experience"].(map[string]any)["hits"].(map[string]any)["hits"].([]any)[0].(map[string]any)["_source"].(map[string]any)["domain"] = "wrong.test"
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": map[string]any{"value": 3, "relation": "eq"}, "hits": []any{contactHit("1", "123", "Owner", true), contactHit("2", "123", "President", true), wrong}}})
			return
		}
		id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		appended = append(appended, id)
		if r.URL.Query().Get("enable") != "b2b_emails,directdials" {
			t.Error("wrong contact add-ons")
		}
		if id == "1" {
			fmt.Fprint(w, `{"id":1,"b2b_emails":null,"directdials":null}`)
			return
		}
		if id != "2" {
			t.Error("enriched contradictory employer", id)
		}
		fmt.Fprint(w, `{"id":2,"b2b_emails":[{"email":"avery@example.test"},{"email":"old@old-employer.test"},{"email":"avery@example.test"}],"directdials":["+12025550123","+12025550123"]}`)
	})
	r, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{{CompanyIDs: []string{"123"}, Domain: "https://www.example.test/"}}, MaxContacts: 1})
	if e != nil {
		t.Fatal(e)
	}
	co := r.Companies[0]
	if co.Status != "contacts_found" || len(co.Contacts) != 1 || len(co.CandidatesWithoutContacts) != 1 || r.RequestsMade != 3 {
		t.Fatalf("unexpected workflow result: %+v", co)
	}
	b := co.Contacts[0]
	if b.PersonID != "2" || b.NeedsReview || len(b.BusinessEmails) != 1 || b.BusinessEmails[0].Email != "avery@example.test" || b.BusinessEmails[0].Validation != "not_run" || len(b.DirectDials) != 1 {
		t.Fatalf("wrong contact result: %+v", b)
	}
	if strings.Join(appended, ",") != "1,2" {
		t.Fatal(appended)
	}
}

func TestCompanyContactsReviewAndBudget(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			h := contactHit("1", "123", "Owner", false)
			h["_source"].(map[string]any)["updated_at"] = "2020-01-01"
			h["_source"].(map[string]any)["headline"] = "Former owner at Example Roofing"
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": map[string]any{"value": 20, "relation": "gte"}, "hits": []any{h}}})
			return
		}
		if r.URL.Query().Get("enable") != "b2b_emails" {
			t.Error("requested phones in emails-only mode")
		}
		fmt.Fprint(w, `{"id":1,"b2b_emails":[{"email":"avery@example.test"}]}`)
	})
	opts := ContactOptions{Companies: []CompanyTarget{{CompanyIDs: []string{"123"}, Domain: "example.test"}}, EmailsOnly: true}
	r, e := c.CompanyContacts(context.Background(), opts)
	if e != nil {
		t.Fatal(e)
	}
	b := r.Companies[0].Contacts[0]
	if r.Companies[0].Status != "review_required" || !b.NeedsReview || len(b.ReviewReasons) != 3 || b.PhoneStatus != "not_requested" || !r.Companies[0].SearchLimited {
		t.Fatalf("lost review evidence: %+v", b)
	}
	opts.MaxRequests = 1
	r, e = c.CompanyContacts(context.Background(), opts)
	if e != nil || r.Complete || r.RequestsMade != 1 || len(r.Companies[0].Contacts) != 0 {
		t.Fatal("budget not enforced", e)
	}
}

func TestCompanyContactsPermissionPartialAndInvalidInput(t *testing.T) {
	calls := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) { calls++; http.Error(w, "forbidden", 403) })
	if _, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{{Name: "Ambiguous"}}}); e == nil || calls != 0 {
		t.Fatal("name-only request went upstream")
	}
	r, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{{Domain: "one.test"}, {Domain: "two.test"}}})
	if e == nil || r.Complete || calls != 1 || r.Companies[1].Status != "not_processed" {
		t.Fatal("permission failure not preserved")
	}
}

func TestParseCompanyTargetsFormatsAndPrecision(t *testing.T) {
	for _, input := range []string{`[{"company_ids":[9007199254740993],"name":"Example"}]`, `{"companies":[{"company_ids":"9007199254740993;2","company_name":"Example"}]}`, `{"hits":{"hits":[{"_source":{"company_id":9007199254740993,"company_name":"Example"}}]}}`} {
		cs, e := ParseCompanyTargets(strings.NewReader(input))
		if e != nil || cs[0].CompanyIDs[0] != "9007199254740993" || cs[0].Name != "Example" {
			t.Fatal(cs, e)
		}
	}
	if _, e := ParseCompanyTargets(strings.NewReader(`[] {}`)); e == nil {
		t.Fatal("trailing JSON accepted")
	}
}
