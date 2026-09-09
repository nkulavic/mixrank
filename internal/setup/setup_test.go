package setup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanNoWritesAndDedup(t *testing.T) {
	d := filepath.Join(t.TempDir(), "destination")
	steps, e := Plan(Options{Components: []string{"generic", "skills"}, Scope: "user", Destination: d, DryRun: true})
	if e != nil {
		t.Fatal(e)
	}
	if len(steps) != 3 {
		t.Fatal(steps)
	}
	if _, e = os.Stat(d); !os.IsNotExist(e) {
		t.Fatal("plan made changes")
	}
	if _, e = Plan(Options{Components: []string{"codex-plugin", "codex-mcp"}}); e == nil {
		t.Fatal("duplicate MCP allowed")
	}
	if _, e = Plan(Options{Components: []string{"unknown"}}); e == nil {
		t.Fatal("unknown component")
	}
}
