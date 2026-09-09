package mcpserver

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nkulavic/mixrank/catalog"
	"os"
	"os/exec"
	"testing"
	"time"
)

// Opt-in integration check against an actual installed plugin launcher/binary.
// Tool discovery and catalog lookup are local; this never calls MixRank.
func TestInstalledPluginStdio(t *testing.T) {
	path := os.Getenv("MIXRANK_TEST_LAUNCHER")
	if path == "" {
		t.Skip("set MIXRANK_TEST_LAUNCHER to an installed plugin executable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.Command(path, "mcp")
	cmd.Dir = t.TempDir()
	client := mcp.NewClient(&mcp.Implementation{Name: "installed-package-test", Version: "1"}, nil)
	session, e := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer session.Close()
	tools, e := session.ListTools(ctx, nil)
	if e != nil || len(tools.Tools) != len(catalog.Read().Operations)+1+WorkflowToolCount {
		t.Fatal("discovery", e)
	}
	r, e := session.CallTool(ctx, &mcp.CallToolParams{Name: "mixrank_catalog", Arguments: map[string]any{"operation": "get_person_match"}})
	if e != nil || r.IsError {
		t.Fatal("catalog call", e)
	}
}
