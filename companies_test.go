package mixrank

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestMergeCompanyTargetsTransitiveStableAndNonMutating(t *testing.T) {
	input := []CompanyTarget{
		{Name: "Original", CompanyIDs: []string{"0001"}, Domain: "https://WWW.Example.test/path", Qualification: "qualified"},
		{Name: "Separate", CompanyIDs: []string{"9"}, Domain: "separate.test"},
		{Name: "Renamed", CompanyIDs: []string{"2"}, Domain: "new.test", Qualification: "needs_review"},
		{Name: "Bridge", CompanyIDs: []string{"1", "2"}, Domain: "example.test."},
		{Name: "Original", CompanyIDs: []string{"9007199254740993"}, Domain: "other.test"},
	}
	before, _ := json.Marshal(input)
	merged, err := MergeCompanyTargets(input)
	after, _ := json.Marshal(input)
	if err != nil || string(before) != string(after) || len(merged) != 3 {
		t.Fatal("merge changed input or lost identities", merged, err)
	}
	want := CompanyTarget{Name: "Original", NameAliases: []string{"Renamed", "Bridge"}, CompanyIDs: []string{"1", "2"}, Domain: "example.test", DomainAliases: []string{"new.test"}, Qualification: "needs_review"}
	if !reflect.DeepEqual(merged[0], want) || merged[1].Name != "Separate" || merged[2].CompanyIDs[0] != "9007199254740993" {
		t.Fatal("aliases, precision, order or conservative qualification lost", merged)
	}
	twice, err := MergeCompanyTargets(merged)
	if err != nil || !reflect.DeepEqual(twice, merged) {
		t.Fatal("merge is not idempotent", twice, err)
	}
}

func TestCompanyContactsMergesBeforeSearchAndAppend(t *testing.T) {
	searches, appends := 0, 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			searches++
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), `"experience.company_id":["123","456"]`) || !strings.Contains(string(body), `"experience.domain.keyword":["example.test","new.test"]`) {
				t.Error("query lost merged identities", string(body))
			}
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": 2, "hits": []any{contactHit("7", "123", "Owner", true), contactHit("7", "456", "Owner", true)}}})
			return
		}
		appends++
		fmt.Fprint(w, `{"id":7,"b2b_emails":[{"email":"avery@example.test"},{"email":"avery@new.test"},{"email":"old@unrelated.test"}],"directdials":[]}`)
	})
	r, err := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{
		{Name: "Example", CompanyIDs: []string{"123"}, Domain: "example.test"},
		{Name: "Example Inc", CompanyIDs: []string{"456"}, Domain: "www.example.test"},
		{Name: "New name", CompanyIDs: []string{"456"}, Domain: "new.test"},
	}})
	if err != nil || searches != 1 || appends != 1 || r.RequestsMade != 2 || len(r.Companies) != 1 {
		t.Fatal("duplicate upstream work", err, searches, appends, r)
	}
	if r.CompanyMerge.InputRecords != 3 || r.CompanyMerge.UniqueCompanies != 1 || r.CompanyMerge.DuplicatesMerged != 2 {
		t.Fatal("incorrect merge counts", r.CompanyMerge)
	}
	contacts := r.Companies[0].Contacts
	if len(contacts) != 1 || len(contacts[0].BusinessEmails) != 2 || contacts[0].BusinessEmails[1].Domain != "new.test" {
		t.Fatal("alias-domain contacts lost or unrelated email accepted", contacts)
	}
}

func TestCompanyMergeLimitsBeforeNetwork(t *testing.T) {
	calls := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		fmt.Fprint(w, `{"hits":{"total":0,"hits":[]}}`)
	})
	duplicates := make([]CompanyTarget, 250)
	for i := range duplicates {
		duplicates[i] = CompanyTarget{CompanyIDs: []string{"1"}}
	}
	r, err := c.CompanyContacts(context.Background(), ContactOptions{Companies: duplicates})
	if err != nil || calls != 1 || r.CompanyMerge.DuplicatesMerged != 249 {
		t.Fatal("limit applied before deduplication", r, err)
	}
	distinct := make([]CompanyTarget, 26)
	for i := range distinct {
		distinct[i] = CompanyTarget{CompanyIDs: []string{fmt.Sprint(i + 1)}}
	}
	tooManyAliases := make([]CompanyTarget, 11)
	for i := range tooManyAliases {
		tooManyAliases[i] = CompanyTarget{CompanyIDs: []string{fmt.Sprint(i + 1)}, Domain: "example.test"}
	}
	for _, input := range [][]CompanyTarget{nil, append(duplicates, duplicates[0]), distinct, tooManyAliases, {{Name: "name only"}}, {{CompanyIDs: []string{"0"}}}, {{Domain: "invalid"}}} {
		if _, err := c.CompanyContacts(context.Background(), ContactOptions{Companies: input}); err == nil || calls != 1 {
			t.Fatal("invalid input reached upstream", err)
		}
	}
}

func TestParseCompanyAliasesRoundTrip(t *testing.T) {
	input := `{"companies":[{"company":{"name":"Example","name_aliases":["Example Inc"],"company_ids":[123,456],"domain":"example.test","domain_aliases":["new.test"]}}]}`
	parsed, err := ParseCompanyTargets(strings.NewReader(input))
	if err != nil || len(parsed) != 1 || !reflect.DeepEqual(parsed[0].DomainAliases, []string{"new.test"}) || !reflect.DeepEqual(parsed[0].NameAliases, []string{"Example Inc"}) {
		t.Fatal("report round trip lost aliases", parsed, err)
	}
}
