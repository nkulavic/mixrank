package mixrank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestContactsParallelLimitAndStableOrder(t *testing.T) {
	var active, peak atomic.Int32
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": 2, "hits": []any{contactHit("1", "123", "Owner", true), contactHit("2", "123", "President", true)}}})
			return
		}
		n := active.Add(1)
		defer active.Add(-1)
		for {
			old := peak.Load()
			if n <= old || peak.CompareAndSwap(old, n) {
				break
			}
		}
		entered <- struct{}{}
		select {
		case <-release:
		case <-r.Context().Done():
			return
		}
		id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		fmt.Fprintf(w, `{"id":%s,"b2b_emails":[{"email":"person%s@example.test"}],"directdials":null}`, id, id)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan struct{})
	var report *ContactReport
	var err error
	go func() {
		defer close(done)
		report, err = c.CompanyContacts(ctx, ContactOptions{Companies: []CompanyTarget{{CompanyIDs: []string{"123"}, Domain: "example.test"}}, Concurrency: 2, MaxContacts: 2})
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-entered:
		case <-ctx.Done():
			t.Fatal("requests were not parallel")
		}
	}
	close(release)
	<-done
	if err != nil || peak.Load() != 2 || report.RequestsMade != 3 {
		t.Fatal("parallel execution failed", err, peak.Load())
	}
	cs := report.Companies[0].Contacts
	if len(cs) != 2 || cs[0].PersonID != "1" || cs[1].PersonID != "2" {
		t.Fatal("parallel completion changed ranking")
	}
}

func TestContactsGlobalBudgetAcrossCompanies(t *testing.T) {
	var active, peak, calls atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		n := active.Add(1)
		defer active.Add(-1)
		for {
			old := peak.Load()
			if n <= old || peak.CompareAndSwap(old, n) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		if strings.HasSuffix(r.URL.Path, "/_search") {
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": 3, "hits": []any{contactHit("1", "123", "Owner", true), contactHit("2", "123", "President", true), contactHit("3", "123", "CEO", true)}}})
			return
		}
		fmt.Fprint(w, `{"b2b_emails":null,"directdials":null}`)
	})
	companies := []CompanyTarget{}
	for i := 0; i < 8; i++ {
		companies = append(companies, CompanyTarget{CompanyIDs: []string{"123"}, Name: fmt.Sprint(i)})
	}
	r, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: companies, Concurrency: 3, MaxRequests: 10})
	if e != nil || r.Complete || calls.Load() != 10 || r.RequestsMade != 10 || peak.Load() > 3 || peak.Load() < 2 {
		t.Fatalf("budget/concurrency violated: calls=%d peak=%d report=%+v err=%v", calls.Load(), peak.Load(), r, e)
	}
	for i, co := range r.Companies {
		if co.Company.Name != fmt.Sprint(i) {
			t.Fatal("company order changed")
		}
	}
}

func TestContactsUnifiedRowsAndFilters(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": 3, "hits": []any{contactHit("1", "123", "Owner", true), contactHit("2", "123", "President", true), contactHit("3", "123", "CEO", true)}}})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/1") {
			fmt.Fprint(w, `{"b2b_emails":null,"directdials":null}`)
			return
		}
		fmt.Fprint(w, `{"b2b_emails":[{"email":"avery@example.test"}],"directdials":null}`)
	})
	for _, filter := range []string{"all", "any", "email", "phone", "none"} {
		r, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{{Domain: "example.test"}}, MaxContacts: 1, ContactFilter: filter})
		if e != nil {
			t.Fatal(e)
		}
		co := r.Companies[0]
		b, _ := json.Marshal(co)
		if strings.Contains(string(b), "candidates_without_contacts") {
			t.Fatal("legacy array exposed")
		}
		switch filter {
		case "all":
			if len(co.Contacts) != 3 || co.Contacts[2].EnrichmentStatus != "not_enriched" || co.Contacts[0].HasEmail {
				t.Fatal("unified evidence lost", co)
			}
		case "any", "email":
			if len(co.Contacts) != 1 || !co.Contacts[0].HasEmail || co.ContactsFiltered != 2 {
				t.Fatal("availability filter failed", co)
			}
		case "phone":
			if len(co.Contacts) != 0 || co.CandidatesExamined != 3 {
				t.Fatal("phone filter failed", co)
			}
		case "none":
			if len(co.Contacts) != 1 || co.Contacts[0].PersonID != "1" {
				t.Fatal("unknown mistaken for no data", co)
			}
		}
	}
}

func TestContactsFatalEnrichmentIsNotMarkedUnrequested(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			json.NewEncoder(w).Encode(map[string]any{"hits": map[string]any{"total": 1, "hits": []any{contactHit("1", "123", "Owner", true)}}})
			return
		}
		http.Error(w, "forbidden", 403)
	})
	r, e := c.CompanyContacts(context.Background(), ContactOptions{Companies: []CompanyTarget{{Domain: "example.test"}}, Concurrency: 1})
	if e == nil || r.Complete || r.Companies[0].Contacts[0].EnrichmentStatus != "error" || r.RequestsMade != 2 {
		t.Fatal("sent request mislabeled as unrequested", e, r)
	}
}
