package cli

import (
	"bytes"
	"encoding/json"
	"github.com/nkulavic/mixrank/catalog"
	"strings"
	"testing"
)

func TestAllCLICommandsAndFlags(t *testing.T) {
	root := New()
	for _, op := range catalog.Read().Operations {
		cmd, _, e := root.Find([]string{"api", strings.ReplaceAll(op.ID, "_", "-")})
		if e != nil || cmd.Name() != strings.ReplaceAll(op.ID, "_", "-") {
			t.Fatal(op.ID, e)
		}
	}
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"endpoints", "get-person-match"})
	if e := root.Execute(); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(buf.String(), "first_name") {
		t.Fatal("reference missing")
	}
}
func TestCSVPrecisionAndFormulaEscape(t *testing.T) {
	buf := &bytes.Buffer{}
	e := export(buf, []any{map[string]any{"id": json.Number("9007199254740993"), "name": "=1+1", "missing": nil}}, "csv")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(buf.String(), "9007199254740993") || !strings.Contains(buf.String(), "'=1+1") {
		t.Fatal(buf.String())
	}
}
func TestSetupDryRunAndNoSecretFlags(t *testing.T) {
	root := New()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"setup", "--components", "cli,codex-mcp", "--dry-run", "--yes"})
	if e := root.Execute(); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(buf.String(), "codex-mcp") {
		t.Fatal(buf.String())
	}
	root = New()
	root.SetArgs([]string{"auth", "login", "--api-key", "synthetic"})
	if e := root.Execute(); e == nil {
		t.Fatal("secret argument accepted")
	}
}
