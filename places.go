package mixrank

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultPlacesBaseURL = "https://places.googleapis.com/v1"

// PlacesOptions configures the optional Google Places (New) client. Places
// requests are billable and are deliberately not retried automatically.
type PlacesOptions struct {
	BaseURL    string
	HTTPClient *http.Client
	MaxResults int
}

// PlacesClient performs explicitly requested Google Places lookups. It does
// not cache provider content; callers may retain a place ID, which Google
// documents as the cache-safe identifier.
type PlacesClient struct {
	key        string
	base       *url.URL
	http       *http.Client
	maxResults int
}

// PlaceCandidate is the small, field-masked subset used by the contact
// fallback. It contains a business phone and website, never an individual
// person's contact details.
type PlaceCandidate struct {
	PlaceID                  string   `json:"place_id"`
	DisplayName              string   `json:"display_name,omitempty"`
	FormattedAddress         string   `json:"formatted_address,omitempty"`
	NationalPhoneNumber      string   `json:"national_phone_number,omitempty"`
	InternationalPhoneNumber string   `json:"international_phone_number,omitempty"`
	WebsiteURI               string   `json:"website_uri,omitempty"`
	GoogleMapsURI            string   `json:"google_maps_uri,omitempty"`
	BusinessStatus           string   `json:"business_status,omitempty"`
	Attributions             []string `json:"attributions,omitempty"`
}

// PlacesLookup records the outcome of one explicit company fallback lookup.
// The response is intentionally run-scoped; the toolkit does not persist
// Places content.
type PlacesLookup struct {
	Status      string           `json:"status"`
	Query       string           `json:"query,omitempty"`
	Candidates  []PlaceCandidate `json:"candidates,omitempty"`
	Matched     *PlaceCandidate  `json:"matched,omitempty"`
	RetrievedAt string           `json:"retrieved_at_utc,omitempty"`
	Error       string           `json:"error,omitempty"`
}

// NewPlacesClient creates a Google Places (New) client. The API key should be
// supplied by the OS credential store or environment, never as a CLI flag.
func NewPlacesClient(apiKey string, opts PlacesOptions) (*PlacesClient, error) {
	if strings.TrimSpace(apiKey) == "" || strings.ContainsAny(apiKey, "\r\n\x00") {
		return nil, errors.New("Google Places credential missing")
	}
	base := opts.BaseURL
	if base == "" {
		base = defaultPlacesBaseURL
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("invalid Google Places base URL")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "::1")) {
		return nil, errors.New("Google Places requires HTTPS (except loopback tests)")
	}
	h := opts.HTTPClient
	if h == nil {
		h = &http.Client{Timeout: 30 * time.Second}
	}
	copyClient := *h
	copyClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	max := opts.MaxResults
	if max == 0 {
		max = 3
	}
	if max < 1 || max > 5 {
		return nil, errors.New("Google Places max results must be 1..5")
	}
	return &PlacesClient{key: apiKey, base: u, http: &copyClient, maxResults: max}, nil
}

// SearchCompany performs one Text Search (New) request with a narrow field
// mask. It does not retry because every request can incur provider billing.
func (p *PlacesClient) SearchCompany(ctx context.Context, company CompanyTarget) ([]PlaceCandidate, string, error) {
	if p == nil {
		return nil, "", errors.New("Google Places client is not configured")
	}
	query := placesQuery(company)
	if query == "" {
		return nil, "", errors.New("company needs a name or domain for Google Places lookup")
	}
	body, err := json.Marshal(map[string]any{"textQuery": query, "pageSize": p.maxResults})
	if err != nil {
		return nil, query, errors.New("cannot encode Google Places query")
	}
	endpoint := *p.base
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/places:searchText"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, query, errors.New("cannot create Google Places request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mixrank-go/0.5.0")
	req.Header.Set("X-Goog-Api-Key", p.key)
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.formattedAddress,places.nationalPhoneNumber,places.internationalPhoneNumber,places.websiteUri,places.googleMapsUri,places.businessStatus,places.attributions")
	res, err := p.http.Do(req)
	if err != nil {
		return nil, query, errors.New("Google Places request failed")
	}
	defer res.Body.Close()
	b, readErr := io.ReadAll(io.LimitReader(res.Body, 4<<20+1))
	if readErr != nil || len(b) > 4<<20 {
		return nil, query, errors.New("Google Places response could not be read")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, query, fmt.Errorf("Google Places HTTP %d", res.StatusCode)
	}
	var raw struct {
		Places []struct {
			ID          string `json:"id"`
			DisplayName struct {
				Text string `json:"text"`
			} `json:"displayName"`
			FormattedAddress         string `json:"formattedAddress"`
			NationalPhoneNumber      string `json:"nationalPhoneNumber"`
			InternationalPhoneNumber string `json:"internationalPhoneNumber"`
			WebsiteURI               string `json:"websiteUri"`
			GoogleMapsURI            string `json:"googleMapsUri"`
			BusinessStatus           string `json:"businessStatus"`
			Attributions             []struct {
				Provider string `json:"provider"`
			} `json:"attributions"`
		} `json:"places"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, query, errors.New("Google Places response is not valid JSON")
	}
	out := make([]PlaceCandidate, 0, len(raw.Places))
	for _, v := range raw.Places {
		c := PlaceCandidate{PlaceID: v.ID, DisplayName: v.DisplayName.Text, FormattedAddress: v.FormattedAddress, NationalPhoneNumber: v.NationalPhoneNumber, InternationalPhoneNumber: v.InternationalPhoneNumber, WebsiteURI: v.WebsiteURI, GoogleMapsURI: v.GoogleMapsURI, BusinessStatus: v.BusinessStatus}
		for _, a := range v.Attributions {
			if strings.TrimSpace(a.Provider) != "" {
				c.Attributions = append(c.Attributions, a.Provider)
			}
		}
		out = append(out, c)
	}
	return out, query, nil
}

// MatchCompany runs one search and applies conservative local matching. An
// ambiguous result is preserved for review and never becomes a contact row.
func (p *PlacesClient) MatchCompany(ctx context.Context, company CompanyTarget) (*PlacesLookup, error) {
	candidates, query, err := p.SearchCompany(ctx, company)
	lookup := &PlacesLookup{Query: query, RetrievedAt: time.Now().UTC().Format(time.RFC3339), Candidates: candidates}
	if err != nil {
		lookup.Status = "error"
		lookup.Error = err.Error()
		return lookup, err
	}
	if len(candidates) == 0 {
		lookup.Status = "not_found"
		return lookup, nil
	}
	best, score, second := bestPlace(company, candidates)
	if best == nil || score < 45 || (second > 0 && score-second < 15) {
		lookup.Status = "ambiguous"
		return lookup, nil
	}
	lookup.Status = "matched"
	lookup.Matched = best
	return lookup, nil
}

func placesQuery(c CompanyTarget) string {
	parts := []string{}
	for _, v := range []string{c.Name, c.Domain, c.Locality, c.Region, c.CountryCode} {
		v = strings.TrimSpace(v)
		if v == "" || containsString(parts, v) {
			continue
		}
		parts = append(parts, v)
	}
	return strings.Join(parts, ", ")
}

func bestPlace(company CompanyTarget, candidates []PlaceCandidate) (*PlaceCandidate, int, int) {
	targetName := normalizedWords(company.Name)
	targetDomain := cleanDomain(company.Domain)
	if targetDomain == "" && len(company.DomainAliases) > 0 {
		for _, d := range company.DomainAliases {
			if targetDomain = cleanDomain(d); targetDomain != "" {
				break
			}
		}
	}
	bestIndex, bestScore, second := -1, 0, 0
	for i := range candidates {
		v := &candidates[i]
		score := 0
		name := normalizedWords(v.DisplayName)
		if targetName != "" && name != "" {
			switch {
			case name == targetName:
				score += 60
			case strings.Contains(name, targetName) || strings.Contains(targetName, name):
				score += 45
			case wordOverlap(targetName, name) >= 0.5:
				score += 30
			}
		}
		if targetDomain != "" && cleanDomain(v.WebsiteURI) == targetDomain {
			score += 45
		}
		address := strings.ToLower(v.FormattedAddress)
		if company.Locality != "" && strings.Contains(address, strings.ToLower(company.Locality)) {
			score += 12
		}
		if company.Region != "" && strings.Contains(address, strings.ToLower(company.Region)) {
			score += 8
		}
		if score > bestScore {
			second = bestScore
			bestScore = score
			bestIndex = i
		} else if score > second {
			second = score
		}
	}
	if bestIndex < 0 {
		return nil, 0, 0
	}
	return &candidates[bestIndex], bestScore, second
}

func normalizedWords(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastSpace := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastSpace = false
		} else if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}
	return strings.TrimSpace(b.String())
}

func wordOverlap(a, b string) float64 {
	aw, bw := strings.Fields(a), strings.Fields(b)
	if len(aw) == 0 || len(bw) == 0 {
		return 0
	}
	hit := 0
	for _, x := range aw {
		for _, y := range bw {
			if x == y {
				hit++
				break
			}
		}
	}
	return float64(hit) / float64(len(aw))
}
