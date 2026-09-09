package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank"
)

const WorkflowToolCount = 1

func addContactTool(s *mcp.Server, factory ClientFactory) {
	company := map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}, "company_ids": map[string]any{"type": "array", "maxItems": 10, "items": map[string]any{"type": []string{"string", "integer"}}}, "domain": map[string]any{"type": "string"}, "qualification_status": map[string]any{"type": "string"}}, "additionalProperties": false, "anyOf": []any{map[string]any{"required": []string{"company_ids"}}, map[string]any{"required": []string{"domain"}}}}
	bound := func(min, max, def int) map[string]any {
		return map[string]any{"type": "integer", "minimum": min, "maximum": max, "default": def}
	}
	schema := map[string]any{"type": "object", "required": []string{"companies"}, "additionalProperties": false, "properties": map[string]any{"companies": map[string]any{"type": "array", "minItems": 1, "maxItems": 25, "items": company}, "roles": map[string]any{"type": "array", "maxItems": 20, "items": map[string]any{"type": "string", "minLength": 1, "maxLength": 100}}, "max_contacts_per_company": bound(1, 5, 2), "max_candidates_per_company": bound(1, 25, 10), "max_requests": bound(1, 250, 100), "emails_only": map[string]any{"type": "boolean", "default": false}}}
	s.AddTool(&mcp.Tool{Name: "mixrank_company_contacts", Description: "Find current owners/executives/managers at supplied companies and retrieve actual available business emails and direct dials. Returns contacts, sources, employer evidence, missing data and review flags. Bounded, cached workflow; no email guessing, validation or message sending. Use this to turn a company list into a contact list.", InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: ptr(false), IdempotentHint: false, OpenWorldHint: ptr(true)}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in struct {
			Companies     json.RawMessage `json:"companies"`
			Roles         []string        `json:"roles"`
			MaxContacts   int             `json:"max_contacts_per_company"`
			MaxCandidates int             `json:"max_candidates_per_company"`
			MaxRequests   int             `json:"max_requests"`
			EmailsOnly    bool            `json:"emails_only"`
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
		report, err := c.CompanyContacts(ctx, mixrank.ContactOptions{Companies: companies, Roles: in.Roles, MaxContacts: in.MaxContacts, MaxCandidates: in.MaxCandidates, MaxRequests: in.MaxRequests, EmailsOnly: in.EmailsOnly})
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
