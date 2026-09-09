// Package catalog contains the versioned, account-visible MixRank API contract.
package catalog

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed endpoints.json mappings/*.json
var Files embed.FS

type Parameter struct {
	Name        string   `json:"name"`
	In          string   `json:"in"`
	Required    bool     `json:"required"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     string   `json:"default,omitempty"`
	Min         string   `json:"min,omitempty"`
	Max         string   `json:"max,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}
type Field struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
type Operation struct {
	ID             string      `json:"id"`
	Method         string      `json:"method"`
	Path           string      `json:"path"`
	Summary        string      `json:"summary"`
	Encoding       string      `json:"encoding"`
	Parameters     []Parameter `json:"parameters"`
	ResponseFields []Field     `json:"response_fields"`
	Pagination     string      `json:"pagination"`
	SafeRetry      bool        `json:"safe_retry"`
	Permissions    string      `json:"permissions"`
	Sources        []string    `json:"sources"`
	Verified       string      `json:"verified"`
	Notes          string      `json:"notes,omitempty"`
	Webhook        string      `json:"webhook,omitempty"`
}
type Catalog struct {
	Version    string      `json:"version"`
	Operations []Operation `json:"operations"`
}

func Read() Catalog {
	b, _ := Files.ReadFile("endpoints.json")
	var c Catalog
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	return c
}
func Lookup(id string) (Operation, error) {
	for _, op := range Read().Operations {
		if op.ID == strings.ReplaceAll(id, "-", "_") {
			return op, nil
		}
	}
	return Operation{}, fmt.Errorf("unknown operation %q; run mixrank endpoints", id)
}
func Match(method, path string) (Operation, bool) {
	for _, op := range Read().Operations {
		if op.Method == method && op.Path == path {
			return op, true
		}
	}
	for _, op := range Read().Operations {
		if op.Method != method {
			continue
		}
		a, b := strings.Split(op.Path, "/"), strings.Split(path, "/")
		if len(a) != len(b) {
			continue
		}
		matched := true
		for i := range a {
			if a[i] != b[i] && !strings.HasPrefix(a[i], "{") {
				matched = false
				break
			}
		}
		if matched {
			return op, true
		}
	}
	return Operation{}, false
}
