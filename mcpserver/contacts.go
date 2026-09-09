package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank"
)

const WorkflowToolCount = 2

func addContactTool(s *mcp.Server, factory ClientFactory) {
	addContactValidationTool(s, factory)
	company := map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}, "company_ids": map[string]any{"type": "array", "maxItems": 10, "items": map[string]any{"type": []string{"string", "integer"}}}, "domain": map[string]any{"type": "string"}, "qualification_status": map[string]any{"type": "string"}}, "additionalProperties": false, "anyOf": []any{map[string]any{"required": []string{"company_ids"}}, map[string]any{"required": []string{"domain"}}}}
	bound := func(min, max, def int) map[string]any {
		return map[string]any{"type": "integer", "minimum": min, "maximum": max, "default": def}
	}
	validationProps := map[string]any{"validate_emails": map[string]any{"type": "boolean", "default": false}, "validation_strategy": map[string]any{"type": "string", "enum": []string{"cached", "besteffort", "strict", "fetch"}, "default": "besteffort"}, "validation_wait_seconds": bound(0, 600, 30), "validation_maxage": map[string]any{"type": "string"}}
	schema := map[string]any{"type": "object", "required": []string{"companies"}, "additionalProperties": false, "properties": map[string]any{"companies": map[string]any{"type": "array", "minItems": 1, "maxItems": 25, "items": company}, "roles": map[string]any{"type": "array", "maxItems": 20, "items": map[string]any{"type": "string", "minLength": 1, "maxLength": 100}}, "max_contacts_per_company": bound(1, 5, 2), "max_candidates_per_company": bound(1, 25, 10), "max_requests": bound(1, 250, 100), "emails_only": map[string]any{"type": "boolean", "default": false}}}
	props := schema["properties"].(map[string]any)
	for k, v := range validationProps {
		props[k] = v
	}
	props["concurrency"] = bound(1, 16, 4)
	props["contact_filter"] = map[string]any{"type": "string", "enum": []string{"all", "any", "email", "phone", "both", "none", "valid-email"}, "default": "all"}
	s.AddTool(&mcp.Tool{Name: "mixrank_company_contacts", Description: "Find current owners/executives/managers at supplied companies and retrieve actual available business emails and direct dials. Returns contacts, sources, employer evidence, missing data and review flags. Bounded parallel lookups with optional bulk email validation. One contacts array includes availability and enrichment status; contact_filter selects rows. No email guessing or message sending. Use this to turn a company list into a contact list.", InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: ptr(false), IdempotentHint: false, OpenWorldHint: ptr(true)}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in struct {
			Companies          json.RawMessage `json:"companies"`
			Roles              []string        `json:"roles"`
			MaxContacts        int             `json:"max_contacts_per_company"`
			MaxCandidates      int             `json:"max_candidates_per_company"`
			MaxRequests        int             `json:"max_requests"`
			Concurrency        int             `json:"concurrency"`
			ContactFilter      string          `json:"contact_filter"`
			ValidateEmails     bool            `json:"validate_emails"`
			ValidationStrategy string          `json:"validation_strategy"`
			ValidationWait     *int            `json:"validation_wait_seconds"`
			ValidationMaxAge   string          `json:"validation_maxage"`
			EmailsOnly         bool            `json:"emails_only"`
		}
		if json.Unmarshal(req.Params.Arguments, &in) != nil {
			return failure("invalid contact arguments"), nil
		}
		companies, e := mixrank.ParseCompanyTargets(bytes.NewReader(in.Companies))
		if e != nil {
			return failure(e.Error()), nil
		}
		c, e := factory(ctx)
		if e != nil {
			return failure(e.Error()), nil
		}
		wait := 30 * time.Second
		if in.ValidationWait != nil {
			wait = time.Duration(*in.ValidationWait) * time.Second
		}
		report, err := c.CompanyContacts(ctx, mixrank.ContactOptions{Companies: companies, Roles: in.Roles, MaxContacts: in.MaxContacts, MaxCandidates: in.MaxCandidates, MaxRequests: in.MaxRequests, EmailsOnly: in.EmailsOnly, Concurrency: in.Concurrency, ContactFilter: in.ContactFilter, ValidateEmails: in.ValidateEmails, ValidationOptions: mixrank.ContactValidationOptions{Strategy: in.ValidationStrategy, Wait: wait, MaxAge: in.ValidationMaxAge}})
		if report == nil {
			return failure(c.Redact(err.Error())), nil
		}
		b, e := json.Marshal(report)
		if e != nil {
			return failure("cannot encode contact report"), nil
		}
		if int64(len(b)) > MaxOutputBytes {
			return failure("contact report exceeds MCP output bound; use fewer companies or CLI --output"), nil
		}
		result := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: c.Redact(string(b))}}, IsError: err != nil}
		if err != nil {
			result.Content = append(result.Content, &mcp.TextContent{Text: c.Redact(err.Error())})
		}
		return result, nil
	})
}
