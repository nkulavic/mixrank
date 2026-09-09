package mixrank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlacesSearchMasksFieldsAndMatchesCompany(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/places:searchText" {
			t.Fatalf("unexpected Places request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Goog-Api-Key") != "synthetic-places-key" {
			t.Fatal("Places key was not sent through the header")
		}
		if !strings.Contains(r.Header.Get("X-Goog-FieldMask"), "places.nationalPhoneNumber") || strings.Contains(r.Header.Get("X-Goog-FieldMask"), "*") {
			t.Fatal("field mask was not narrow")
		}
		var body map[string]any
		if json.NewDecoder(r.Body).Decode(&body) != nil || !strings.Contains(body["textQuery"].(string), "Target Roofing") {
			t.Fatal("company query was not sent")
		}
		fmt.Fprint(w, `{"places":[{"id":"place-1","displayName":{"text":"Target Roofing"},"formattedAddress":"1 Main St, Denver, CO","nationalPhoneNumber":"(303) 555-0100","websiteUri":"https://target.test","googleMapsUri":"https://maps.google.test/place-1","businessStatus":"OPERATIONAL","attributions":[{"provider":"Google Maps"}]}]}`)
	}))
	defer server.Close()
	p, err := NewPlacesClient("synthetic-places-key", PlacesOptions{BaseURL: server.URL + "/v1"})
	if err != nil {
		t.Fatal(err)
	}
	lookup, err := p.MatchCompany(context.Background(), CompanyTarget{Name: "Target Roofing", Domain: "target.test", Locality: "Denver", Region: "CO"})
	if err != nil || lookup.Status != "matched" || lookup.Matched == nil || lookup.Matched.PlaceID != "place-1" || lookup.Matched.NationalPhoneNumber == "" {
		t.Fatalf("unexpected Places match: %+v, %v", lookup, err)
	}
}

func TestCompanyContactsPlacesFallbackAndContactableOnly(t *testing.T) {
	mixrankServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			fmt.Fprint(w, `{"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}`)
			return
		}
		http.Error(w, "unexpected MixRank request", http.StatusNotFound)
	}))
	defer mixrankServer.Close()
	placesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if strings.Contains(body["textQuery"].(string), "Target Roofing") {
			fmt.Fprint(w, `{"places":[{"id":"place-1","displayName":{"text":"Target Roofing"},"formattedAddress":"1 Main St, Denver, CO","nationalPhoneNumber":"303-555-0100","websiteUri":"https://target.test","googleMapsUri":"https://maps.google.test/place-1"}]}`)
			return
		}
		fmt.Fprint(w, `{"places":[]}`)
	}))
	defer placesServer.Close()
	c := testClientWithBase(t, mixrankServer.URL)
	p, err := NewPlacesClient("synthetic-places-key", PlacesOptions{BaseURL: placesServer.URL + "/v1"})
	if err != nil {
		t.Fatal(err)
	}
	report, err := c.CompanyContacts(context.Background(), ContactOptions{
		Companies:       []CompanyTarget{{Name: "Target Roofing", Domain: "target.test", Locality: "Denver", Region: "CO"}, {Name: "No Contact", Domain: "none.test", Locality: "Denver", Region: "CO"}},
		PlacesFallback:  true,
		Places:          p,
		ContactableOnly: true,
		MaxContacts:     1,
		MaxRequests:     10,
		Concurrency:     2,
	})
	if err != nil || len(report.Companies) != 1 || report.CompaniesFiltered != 1 || report.ContactFilter != "any" {
		t.Fatalf("contactable-only result was not narrowed: %+v, %v", report, err)
	}
	b := report.Companies[0].Contacts[0]
	if b.ContactType != "business" || !b.HasPhone || len(b.BusinessPhones) != 1 || b.BusinessPhones[0] != "303-555-0100" || b.Source != "Google Places business listing" || b.PlaceID != "place-1" {
		t.Fatalf("business fallback row was not labeled correctly: %+v", b)
	}
}

func testClientWithBase(t *testing.T, base string) *Client {
	t.Helper()
	c, err := New("synthetic-test-credential", Options{BaseURL: base, LockDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	return c
}
