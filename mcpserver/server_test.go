package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank"
	"github.com/nkulavic/mixrank/catalog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPMCPDiscoveryCallAndAuth(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{"id":9007199254740993}`) }))
	defer upstream.Close()
	c, e := mixrank.New("synthetic", mixrank.Options{BaseURL: upstream.URL, LockDir: t.TempDir()})
	if e != nil {
		t.Fatal(e)
	}
	token := strings.Repeat("x", 32)
	h, e := Handler(func(context.Context) (*mixrank.Client, error) { return c, nil }, HTTPOptions{BearerToken: token, AllowedOrigins: []string{"https://agent.example"}})
	if e != nil {
		t.Fatal(e)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	for _, tc := range []struct {
		token, origin string
		status        int
	}{{"", "", 401}, {token, "https://evil.example", 403}} {
		r, _ := http.NewRequest("POST", server.URL+"/mcp", nil)
		if tc.token != "" {
			r.Header.Set("Authorization", "Bearer "+tc.token)
		}
		r.Header.Set("Origin", tc.origin)
		res, e := http.DefaultClient.Do(r)
		if e != nil {
			t.Fatal(e)
		}
		res.Body.Close()
		if res.StatusCode != tc.status {
			t.Fatal("guard", res.StatusCode)
		}
	}
	cl := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	session, e := cl.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: server.URL + "/mcp", HTTPClient: &http.Client{Transport: bearerTransport{token}}}, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer session.Close()
	tools, e := session.ListTools(context.Background(), nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(tools.Tools) != len(catalog.Read().Operations)+1+WorkflowToolCount {
		t.Fatal("parity", len(tools.Tools))
	}
	out, e := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "get_echo", Arguments: map[string]any{}})
	if e != nil || out.IsError {
		t.Fatal(e, out)
	}
	b, _ := json.Marshal(out)
	if !strings.Contains(string(b), "9007199254740993") {
		t.Fatal("precision loss", string(b))
	}
}

type bearerTransport struct{ token string }

func TestMCPContactWorkflow(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_search") {
			fmt.Fprintf(w, `{"hits":{"total":{"value":1,"relation":"eq"},"hits":[{"_id":"7","_source":{"person_id":7,"name":{"full":"Avery Morgan"},"updated_at":%q},"inner_hits":{"experience":{"hits":{"hits":[{"_source":{"company_id":123,"company_name":"Example","domain":"example.test","title":"Owner","is_current":true,"has_source_company_id":true}}]}}}}]}}`, time.Now().UTC().Format(time.RFC3339))
			return
		}
		if r.URL.Query().Get("enable") != "b2b_emails,directdials" {
			t.Error("contact add-ons missing")
		}
		fmt.Fprint(w, `{"id":7,"b2b_emails":[{"email":"avery@example.test"}],"directdials":["+12025550123"]}`)
	}))
	defer upstream.Close()
	c, _ := mixrank.New("synthetic", mixrank.Options{BaseURL: upstream.URL})
	s := New(func(context.Context) (*mixrank.Client, error) { return c, nil })
	a, b := mcp.NewInMemoryTransports()
	ss, e := s.Connect(context.Background(), a, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer ss.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "workflow-test", Version: "1"}, nil)
	cs, e := client.Connect(context.Background(), b, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer cs.Close()
	r, e := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "mixrank_company_contacts", Arguments: map[string]any{"companies": []any{map[string]any{"company_ids": []int{123}, "domain": "example.test"}}, "max_contacts_per_company": 1}})
	if e != nil || r.IsError {
		t.Fatal(e, r)
	}
	var report mixrank.ContactReport
	if e = json.Unmarshal([]byte(r.Content[0].(*mcp.TextContent).Text), &report); e != nil {
		t.Fatal(e)
	}
	if report.Companies[0].Contacts[0].BusinessEmails[0].Email != "avery@example.test" || report.RequestsMade != 2 {
		t.Fatal("workflow did not enrich contacts")
	}
}

func (b bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+b.token)
	return http.DefaultTransport.RoundTrip(r)
}
func TestMCPInMemoryErrorsAndMetadata(t *testing.T) {
	s := New(func(context.Context) (*mixrank.Client, error) { return nil, errors.New("credential unavailable") })
	a, b := mcp.NewInMemoryTransports()
	ctx := context.Background()
	ss, e := s.Connect(ctx, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer ss.Close()
	c := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, e := c.Connect(ctx, b, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer cs.Close()
	r, e := cs.CallTool(ctx, &mcp.CallToolParams{Name: "get_echo", Arguments: map[string]any{}})
	if e != nil || !r.IsError {
		t.Fatal("credential error not surfaced")
	}
	r, e = cs.CallTool(ctx, &mcp.CallToolParams{Name: "mixrank_catalog", Arguments: map[string]any{"index": "../secret"}})
	if e != nil || !r.IsError {
		t.Fatal("mapping path not validated")
	}
	list, _ := cs.ListTools(ctx, nil)
	for _, op := range list.Tools {
		if op.Name == "post_privacy_redaction_request" && (!*op.Annotations.DestructiveHint || op.Annotations.ReadOnlyHint) {
			t.Fatal("incorrect mutation hints")
		}
	}
}
