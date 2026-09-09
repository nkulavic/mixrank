package mcpserver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank"
)

func addContactValidationTool(s *mcp.Server, factory ClientFactory) {
	schema := map[string]any{"type": "object", "required": []string{"report"}, "additionalProperties": false, "properties": map[string]any{
		"report":                  map[string]any{"type": "object", "description": "Existing contact report. Preserve email_validation.job_id to resume without resubmitting.", "required": []string{"companies"}},
		"validation_strategy":     map[string]any{"type": "string", "enum": []string{"cached", "besteffort", "strict", "fetch"}, "default": "besteffort"},
		"validation_wait_seconds": map[string]any{"type": "integer", "minimum": 0, "maximum": 600, "default": 30},
		"max_requests":            map[string]any{"type": "integer", "minimum": 1, "maximum": 250, "default": 30},
	}}
	s.AddTool(&mcp.Tool{Name: "mixrank_validate_contacts", Description: "Validate work emails in an existing contact report or resume its pending bulk job. Deduplicates addresses and attaches provider validity, dates, catchall and greylist evidence. Keeps pending/uncertain outcomes and never resubmits a retained job ID.", InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: ptr(false), IdempotentHint: false, OpenWorldHint: ptr(true)}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in struct {
			Report      *mixrank.ContactReport `json:"report"`
			Strategy    string                 `json:"validation_strategy"`
			Wait        *int                   `json:"validation_wait_seconds"`
			MaxRequests int                    `json:"max_requests"`
		}
		if json.Unmarshal(req.Params.Arguments, &in) != nil || in.Report == nil || in.Report.Companies == nil {
			return failure("expected a contact report"), nil
		}
		c, e := factory(ctx)
		if e != nil {
			return failure(e.Error()), nil
		}
		wait := 30 * time.Second
		if in.Wait != nil {
			wait = time.Duration(*in.Wait) * time.Second
		}
		err := c.ValidateContacts(ctx, in.Report, mixrank.ContactValidationOptions{Strategy: in.Strategy, Wait: wait, MaxRequests: in.MaxRequests})
		b, e := json.Marshal(in.Report)
		if e != nil || int64(len(b)) > MaxOutputBytes {
			return failure("report exceeds output bound; use CLI contacts validate --output"), nil
		}
		result := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: c.Redact(string(b))}}, IsError: err != nil}
		if err != nil {
			result.Content = append(result.Content, &mcp.TextContent{Text: c.Redact(err.Error())})
		}
		return result, nil
	})
}
