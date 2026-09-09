package setup

import (
	"context"
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

func TestInstallPreserveModifiedAndUninstallSelection(t *testing.T) {
	d := t.TempDir()
	receipt := filepath.Join(d, "state")
	skills := filepath.Join(d, "skills")
	o := Options{Components: []string{"skills"}, Scope: "user", StateDir: receipt}
	steps := []Step{{Component: "skills", Destination: skills}}
	if _, e := Apply(context.Background(), o, steps); e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(skills, "mixrank-cli", "SKILL.md")
	original, e := os.ReadFile(file)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(file, append(original, []byte("\nPrivate customization\n")...), 0600); e != nil {
		t.Fatal(e)
	}
	o.Action = "repair"
	if _, e = Apply(context.Background(), o, steps); e == nil {
		t.Fatal("overwrote modified file")
	}
	o.Action = "uninstall"
	messages, e := Apply(context.Background(), o, steps)
	if e != nil {
		t.Fatal(e)
	}
	if len(messages) == 0 {
		t.Fatal("missing preservation report")
	}
	if _, e = os.Stat(file); e != nil {
		t.Fatal("deleted modified file")
	}
	if _, e = os.Stat(filepath.Join(skills, "mixrank-research", "SKILL.md")); !os.IsNotExist(e) {
		t.Fatal("owned file not removed")
	}
}
