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
	if len(tools.Tools) != len(catalog.Read().Operations)+1 {
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
