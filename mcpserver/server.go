// Package mcpserver exposes the catalog through the official Go MCP SDK.
package mcpserver

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank"
	"github.com/nkulavic/mixrank/catalog"
	"net/http"
	"strings"
	"time"
)

type ClientFactory func(context.Context) (*mixrank.Client, error)

const MaxOutputBytes int64 = 1 << 20

func New(factory ClientFactory) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "mixrank", Version: "0.2.0"}, nil)
	for _, op := range catalog.Read().Operations {
		props := map[string]any{}
		required := []string{}
		for _, p := range op.Parameters {
			if p.In == "file" {
				continue
			}
			schema := map[string]any{"type": "string", "description": p.Description}
			if len(p.Enum) > 0 {
				schema["enum"] = p.Enum
			}
			props[p.Name] = schema
			if p.Required {
				required = append(required, p.Name)
			}
		}
		params := map[string]any{"type": "object", "properties": props, "additionalProperties": map[string]any{"type": "string"}}
		if len(required) > 0 {
			params["required"] = required
		}
		schema := map[string]any{"type": "object", "properties": map[string]any{"parameters": params, "body": map[string]any{"type": "object", "description": "Complete JSON body or form fields; Elasticsearch DSL is passed unchanged"}, "upload_text": map[string]any{"type": "string", "description": "UTF-8 lines for multipart bulk jobs"}, "refresh": map[string]any{"type": "boolean", "description": "Explicit authorization for a requested live-fetch strategy"}}, "additionalProperties": false}
		if len(required) > 0 {
			schema["required"] = []string{"parameters"}
		}
		if op.Encoding == "multipart" {
			schema["required"] = []string{"upload_text"}
		}
		s.AddTool(&mcp.Tool{Name: op.ID, Description: op.Summary + ". " + op.Notes, InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: op.SafeRetry, DestructiveHint: ptr(op.Method == "DELETE" || op.Path == "/privacy-redaction-request"), IdempotentHint: op.SafeRetry, OpenWorldHint: ptr(true)}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var in struct {
				Parameters map[string]string `json:"parameters"`
				Body       json.RawMessage   `json:"body"`
				UploadText *string           `json:"upload_text"`
				Refresh    bool              `json:"refresh"`
			}
			if err := json.Unmarshal(req.Params.Arguments, &in); err != nil {
				return failure("invalid arguments"), nil
			}
			if op.Encoding == "multipart" && in.UploadText == nil {
				return failure("upload_text is required"), nil
			}
			c, e := factory(ctx)
			if e != nil {
				return failure(e.Error()), nil
			}
			input := mixrank.Request{Parameters: in.Parameters, Body: in.Body, Refresh: in.Refresh}
			if in.UploadText != nil {
				input.Upload = strings.NewReader(*in.UploadText)
				input.Filename = "input.txt"
			}
			r, e := c.Call(ctx, op.ID, input)
			if e != nil {
				return failure(c.Redact(e.Error())), nil
			}
			defer r.Body.Close()
			v, e := r.JSON(MaxOutputBytes)
			if e != nil {
				return failure(e.Error() + "; use CLI --output for large or binary results"), nil
			}
			// RawMessage avoids converting large provider integers to float64.
			envelope := map[string]any{"status": r.StatusCode, "data": v}
			b, e := json.Marshal(envelope)
			if e != nil {
				return failure("invalid upstream response"), nil
			}
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
		})
	}
	s.AddTool(&mcp.Tool{Name: "mixrank_catalog", Description: "Inspect endpoint contracts or Elasticsearch field mappings without an API call.", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"operation": map[string]any{"type": "string"}, "index": map[string]any{"type": "string", "enum": []string{"person2", "companies", "jobs", "org_name", "industries"}}}}}, func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in struct{ Operation, Index string }
		if json.Unmarshal(req.Params.Arguments, &in) != nil {
			return failure("invalid arguments"), nil
		}
		var b []byte
		if in.Index != "" {
			if !validIndex(in.Index) {
				return failure("unknown index"), nil
			}
			b, _ = catalog.Files.ReadFile("mappings/" + in.Index + ".json")
		} else if in.Operation != "" {
			op, e := catalog.Lookup(in.Operation)
			if e != nil {
				return failure(e.Error()), nil
			}
			b, _ = json.Marshal(op)
		} else {
			b, _ = json.Marshal(catalog.Read().Operations)
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
	})
	addContactTool(s, factory)
	return s
}
func validIndex(s string) bool {
	for _, v := range []string{"person2", "companies", "jobs", "org_name", "industries"} {
		if s == v {
			return true
		}
	}
	return false
}
func ptr(v bool) *bool { return &v }
func failure(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: s}}}
}
func Stdio(ctx context.Context, factory ClientFactory) error {
	return New(factory).Run(ctx, &mcp.StdioTransport{})
}

// HTTPOptions supports a local bearer credential or an externally implemented
// OAuth verifier. A verifier must validate issuer, audience, expiry and scopes.
// A static bearer is for private development; it is not a claimed OAuth flow.
type HTTPOptions struct {
	BearerToken    string
	AllowedOrigins []string
	Authorize      func(*http.Request) error
}

func Handler(factory ClientFactory, opts HTTPOptions) (http.Handler, error) {
	if opts.BearerToken == "" && opts.Authorize == nil {
		return nil, errors.New("HTTP MCP requires authentication")
	}
	if opts.BearerToken != "" && len(opts.BearerToken) < 32 {
		return nil, errors.New("HTTP bearer token must contain at least 32 characters")
	}
	server := New(factory)
	base := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true, MaxRequestBodyBytes: 4 << 20, PropagateRequestCancellation: true})
	origins := map[string]bool{}
	for _, o := range opts.AllowedOrigins {
		if o == "*" || o == "null" {
			return nil, errors.New("explicit origins required")
		}
		origins[o] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.URL.Path != "/mcp" {
			http.NotFound(w, r)
			return
		}
		if o := r.Header.Get("Origin"); o != "" {
			if !origins[o] {
				http.Error(w, "origin denied", http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", o)
			w.Header().Add("Vary", "Origin")
		}
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, MCP-Protocol-Version, MCP-Session-Id")
			w.WriteHeader(204)
			return
		}
		authorized := false
		if opts.Authorize != nil {
			authorized = opts.Authorize(r) == nil
		} else {
			got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			authorized = strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") && subtle.ConstantTimeCompare([]byte(got), []byte(opts.BearerToken)) == 1
		}
		if !authorized {
			w.Header().Set("WWW-Authenticate", `Bearer realm="mixrank"`)
			http.Error(w, "unauthorized", 401)
			return
		}
		base.ServeHTTP(w, r)
	}), nil
}
func Serve(ctx context.Context, addr string, factory ClientFactory, opts HTTPOptions) error {
	h, e := Handler(factory, opts)
	if e != nil {
		return e
	}
	s := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32 << 10}
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = s.Shutdown(c)
		case <-done:
		}
	}()
	if e = s.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		return fmt.Errorf("HTTP listener failed: %w", e)
	}
	return nil
}
