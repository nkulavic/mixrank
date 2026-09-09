package mixrank

import (
	"errors"
	"strconv"
	"strings"
)

// CompanyMergeSummary describes consolidation before any upstream requests.
type CompanyMergeSummary struct {
	InputRecords     int `json:"input_records"`
	UniqueCompanies  int `json:"unique_companies"`
	DuplicatesMerged int `json:"duplicates_merged"`
}

// MergeCompanyTargets merges records connected by an exact normalized domain
// or company ID. Names alone never identify a company. The first record controls
// ordering and display values; all distinct IDs, names and domains are retained.
// Shared domains identify business accounts, not individual franchise locations.
// This function makes no API calls and does not mutate the input.
func MergeCompanyTargets(input []CompanyTarget) ([]CompanyTarget, error) {
	if len(input) == 0 || len(input) > 250 {
		return nil, errors.New("company input requires 1..250 records before merging")
	}
	rows := make([]CompanyTarget, len(input))
	parent := make([]int, len(input))
	var find func(int) int
	find = func(i int) int {
		if parent[i] != i {
			parent[i] = find(parent[i])
		}
		return parent[i]
	}
	identities := map[string]int{}
	for i, co := range input {
		parent[i] = i
		n := CompanyTarget{Qualification: strings.TrimSpace(co.Qualification), Website: strings.TrimSpace(co.Website), Locality: strings.TrimSpace(co.Locality), Region: strings.TrimSpace(co.Region), CountryCode: strings.TrimSpace(co.CountryCode)}
		for _, name := range append([]string{co.Name}, co.NameAliases...) {
			addCompanyName(&n, strings.TrimSpace(name))
		}
		domainValues := append([]string{co.Domain}, co.DomainAliases...)
		if co.Domain == "" && co.Website != "" {
			domainValues = append(domainValues, co.Website)
		}
		for _, domain := range domainValues {
			if strings.TrimSpace(domain) == "" {
				continue
			}
			domain = cleanDomain(domain)
			if !strings.Contains(domain, ".") {
				return nil, errors.New("company domain must be a valid domain name")
			}
			addCompanyDomain(&n, domain)
		}
		for _, id := range co.CompanyIDs {
			value, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
			if err != nil || value <= 0 {
				return nil, errors.New("company IDs must be positive integers")
			}
			n.CompanyIDs = appendUnique(n.CompanyIDs, strconv.FormatInt(value, 10))
		}
		if len(n.CompanyIDs) == 0 && n.Domain == "" {
			return nil, errors.New("each company needs company_ids or domain; name alone is ambiguous")
		}
		rows[i] = n
		keys := []string{}
		for _, id := range n.CompanyIDs {
			keys = append(keys, "id:"+id)
		}
		for _, domain := range companyDomains(n) {
			keys = append(keys, "domain:"+domain)
		}
		for _, key := range keys {
			if previous, exists := identities[key]; exists {
				a, b := find(i), find(previous)
				parent[max(a, b)] = min(a, b)
			} else {
				identities[key] = i
			}
		}
	}
	result := []CompanyTarget{}
	positions := map[int]int{}
	for i, co := range rows {
		root := find(i)
		pos, exists := positions[root]
		if !exists {
			pos = len(result)
			positions[root] = pos
			result = append(result, CompanyTarget{})
		}
		target := &result[pos]
		for _, name := range append([]string{co.Name}, co.NameAliases...) {
			addCompanyName(target, name)
		}
		for _, domain := range companyDomains(co) {
			addCompanyDomain(target, domain)
		}
		for _, id := range co.CompanyIDs {
			target.CompanyIDs = appendUnique(target.CompanyIDs, id)
		}
		if target.Qualification == "" {
			target.Qualification = co.Qualification
		} else if co.Qualification != "" && target.Qualification != co.Qualification {
			target.Qualification = "needs_review"
		}
		if target.Website == "" {
			target.Website = co.Website
		}
		if target.Locality == "" {
			target.Locality = co.Locality
		}
		if target.Region == "" {
			target.Region = co.Region
		}
		if target.CountryCode == "" {
			target.CountryCode = co.CountryCode
		}
		if len(target.CompanyIDs) > 10 || len(companyDomains(*target)) > 10 {
			return nil, errors.New("merged company exceeds 10 IDs or 10 domains; narrow the input")
		}
	}
	return result, nil
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func appendUnique(values []string, value string) []string {
	if !containsString(values, value) {
		return append(values, value)
	}
	return values
}

func companyDomains(co CompanyTarget) []string {
	if co.Domain == "" {
		return co.DomainAliases
	}
	return append([]string{co.Domain}, co.DomainAliases...)
}

func addCompanyName(co *CompanyTarget, name string) {
	if name == "" || name == co.Name {
		return
	}
	if co.Name == "" {
		co.Name = name
	} else {
		co.NameAliases = appendUnique(co.NameAliases, name)
	}
}

func addCompanyDomain(co *CompanyTarget, domain string) {
	if co.Domain == domain {
		return
	}
	if co.Domain == "" {
		co.Domain = domain
	} else {
		co.DomainAliases = appendUnique(co.DomainAliases, domain)
	}
}
