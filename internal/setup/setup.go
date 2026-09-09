// Package setup plans installations and makes owned, backed-up local changes.
package setup

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank"
	"github.com/nkulavic/mixrank/internal/profiles"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Options struct {
	Components                 []string
	Scope, Destination, Action string
	DryRun                     bool
	StateDir                   string // Optional isolated receipt directory for embedding and integration tests.
}
type Step struct {
	Component   string   `json:"component"`
	Destination string   `json:"destination,omitempty"`
	Command     []string `json:"command,omitempty"`
	Manual      string   `json:"manual,omitempty"`
}
type Receipt struct {
	Components map[string]string   `json:"components"`
	Files      map[string]string   `json:"files"`
	Commands   map[string][]string `json:"commands"`
}

func statePath() (string, error) {
	s, e := profiles.Default()
	return filepath.Join(s.Root, "install.json"), e
}
func Plan(o Options) ([]Step, error) {
	home, e := os.UserHomeDir()
	if e != nil {
		return nil, e
	}
	if o.Scope == "" {
		o.Scope = "user"
	}
	if o.Scope != "user" && o.Scope != "project" {
		return nil, errors.New("scope must be user or project")
	}
	if len(o.Components) == 0 {
		return nil, errors.New("select at least one component")
	}
	seen := map[string]bool{}
	for _, v := range o.Components {
		seen[v] = true
	}
	if seen["codex-plugin"] && seen["codex-mcp"] {
		return nil, errors.New("select Codex plugin or standalone MCP; both register the same server")
	}
	if seen["claude-plugin"] && seen["claude-mcp"] {
		return nil, errors.New("select Claude plugin or standalone MCP")
	}
	dest := o.Destination
	if dest == "" {
		dest = filepath.Join(home, ".local", "bin")
	}
	dest, e = filepath.Abs(dest)
	if e != nil {
		return nil, e
	}
	binary := filepath.Join(dest, "mixrank")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	steps := []Step{}
	configStore, e := profiles.Default()
	if e != nil {
		return nil, e
	}
	localRoot := filepath.Join(configStore.Root, "agent-packages")
	needsBinary := seen["cli"] || seen["codex-mcp"] || seen["claude-mcp"] || seen["generic"] || seen["skills"]
	if needsBinary {
		steps = append(steps, Step{Component: "cli", Destination: binary})
	}
	for _, comp := range o.Components {
		switch comp {
		case "cli":
		case "skills":
			d := filepath.Join(home, ".agents", "skills")
			if o.Scope == "project" {
				wd, _ := os.Getwd()
				d = filepath.Join(wd, ".agents", "skills")
			}
			if o.Destination != "" {
				d = dest
			}
			steps = append(steps, Step{Component: comp, Destination: d})
		case "codex-mcp":
			if o.Scope != "user" {
				return nil, errors.New("Codex standalone MCP setup supports shared user scope")
			}
			steps = append(steps, Step{Component: comp, Command: []string{"codex", "mcp", "add", "mixrank", "--", binary, "mcp"}})
		case "claude-mcp":
			steps = append(steps, Step{Component: comp, Command: []string{"claude", "mcp", "add", "--scope", o.Scope, "mixrank", "--", binary, "mcp"}})
		case "codex-plugin":
			if o.Scope != "user" {
				return nil, errors.New("Codex plugin setup supports shared user scope")
			}
			d := filepath.Join(localRoot, "codex")
			steps = append(steps, Step{Component: comp, Destination: d}, Step{Component: comp, Command: []string{"codex", "plugin", "marketplace", "add", d, "--json"}}, Step{Component: comp, Command: []string{"codex", "plugin", "add", "mixrank@mixrank-toolkit-local", "--json"}})
		case "claude-plugin":
			d := filepath.Join(localRoot, "claude")
			steps = append(steps, Step{Component: comp, Destination: d}, Step{Component: comp, Command: []string{"claude", "plugin", "marketplace", "add", d}}, Step{Component: comp, Command: []string{"claude", "plugin", "install", "mixrank@mixrank-toolkit-local", "--scope", o.Scope}})
		case "claude-desktop":
			steps = append(steps, Step{Component: comp, Manual: "Download the matching platform .mcpb from https://github.com/nkulavic/mixrank/releases; open it in Claude Desktop, enter the sensitive API-key field, then enable. Desktop owns encrypted storage. Restart if prompted."})
		case "cowork":
			steps = append(steps, Step{Component: comp, Manual: "Import mixrank-cowork.zip from https://github.com/nkulavic/mixrank/releases in Cowork's plugin UI. Skills run inside Cowork. Connect remote MCP through Cowork's connector authentication after hosted OAuth deployment; host credentials and Desktop configuration are not copied."})
		case "generic":
			p, e := profiles.Default()
			if e != nil {
				return nil, e
			}
			steps = append(steps, Step{Component: comp, Destination: filepath.Join(p.Root, "mcp.json"), Manual: "Import the generated MCP configuration into your agent. It uses an absolute native binary and MixRank-owned OS credentials."})
		default:
			return nil, fmt.Errorf("unknown component %q", comp)
		}
	}
	return steps, nil
}
func Apply(ctx context.Context, o Options, steps []Step) ([]string, error) {
	path, e := statePath()
	if o.StateDir != "" {
		path = filepath.Join(o.StateDir, "install.json")
		e = nil
	}
	if e != nil {
		return nil, e
	}
	r := Receipt{Files: map[string]string{}, Commands: map[string][]string{}}
	if b, e := os.ReadFile(path); e == nil {
		if e = json.Unmarshal(b, &r); e != nil {
			return nil, e
		}
	}
	if r.Files == nil {
		r.Files = map[string]string{}
	}
	if r.Commands == nil {
		r.Commands = map[string][]string{}
	}
	if r.Components == nil {
		r.Components = map[string]string{}
	}
	selected := map[string]bool{}
	for _, comp := range o.Components {
		selected[comp] = true
	}
	currentComponent := ""
	messages := []string{}
	if o.DryRun {
		return messages, nil
	}
	save := func() error { return profiles.WritePrivate(path, r) }
	write := func(p string, b []byte, mode fs.FileMode) error {
		if existing, e := os.ReadFile(p); e == nil {
			if string(existing) == string(b) {
				r.Files[p] = fmt.Sprintf("%x", sha256.Sum256(b))
				r.Components[p] = currentComponent
				return nil
			}
			if old, ok := r.Files[p]; !ok || fmt.Sprintf("%x", sha256.Sum256(existing)) != old {
				return fmt.Errorf("refusing to replace unowned or modified file %s", p)
			}
			if e = os.WriteFile(p+".backup-"+time.Now().UTC().Format("20060102T150405"), existing, 0600); e != nil {
				return e
			}
		}
		if e = os.MkdirAll(filepath.Dir(p), 0700); e != nil {
			return e
		}
		temp, e := os.CreateTemp(filepath.Dir(p), ".mixrank-install-*")
		if e != nil {
			return e
		}
		defer os.Remove(temp.Name())
		if e = temp.Chmod(mode); e != nil {
			temp.Close()
			return e
		}
		if _, e = temp.Write(b); e != nil {
			temp.Close()
			return e
		}
		if e = temp.Close(); e != nil {
			return e
		}
		if e = os.Rename(temp.Name(), p); e != nil {
			return e
		}
		r.Files[p] = fmt.Sprintf("%x", sha256.Sum256(b))
		r.Components[p] = currentComponent
		return save()
	}
	if o.Action == "uninstall" {
		for p, hash := range r.Files {
			if !selected[r.Components[p]] {
				continue
			}
			if r.Components[p] == "cli" {
				for _, dep := range []string{"codex-mcp", "claude-mcp"} {
					if r.Commands[dep] != nil && !selected[dep] {
						return messages, errors.New("CLI is still required by an installed MCP component; include that component in uninstall")
					}
				}
			}
			b, e := os.ReadFile(p)
			if os.IsNotExist(e) {
				delete(r.Files, p)
				continue
			}
			if e != nil {
				return messages, e
			}
			if fmt.Sprintf("%x", sha256.Sum256(b)) != hash {
				messages = append(messages, "Preserved modified file: "+p)
				continue
			}
			if e = os.Remove(p); e != nil {
				return messages, e
			}
			delete(r.Files, p)
		}
		for component, cmd := range r.Commands {
			if !selected[component] {
				continue
			}
			remove := []string{}
			switch component {
			case "codex-mcp":
				remove = []string{"codex", "mcp", "remove", "mixrank"}
			case "claude-mcp":
				remove = []string{"claude", "mcp", "remove", "--scope", o.Scope, "mixrank"}
			case "codex-plugin":
				remove = []string{"codex", "plugin", "remove", "mixrank@mixrank-toolkit-local"}
			case "claude-plugin":
				remove = []string{"claude", "plugin", "uninstall", "mixrank@mixrank-toolkit-local", "--scope", o.Scope}
			}
			if i := index(cmd, "--scope"); i >= 0 && i+1 < len(cmd) {
				for j := range remove {
					if remove[j] == "--scope" && j+1 < len(remove) {
						remove[j+1] = cmd[i+1]
					}
				}
			}
			if len(remove) > 0 {
				if e = run(ctx, remove); e != nil {
					return messages, e
				}
				delete(r.Commands, component)
			}
		}
		return messages, save()
	}
	binary := ""
	skipComponents := map[string]bool{}
	for _, comp := range o.Components {
		if (comp == "codex-plugin" || comp == "claude-plugin") && r.Commands[comp] == nil {
			host := "codex"
			if comp == "claude-plugin" {
				host = "claude"
			}
			b, err := exec.CommandContext(ctx, host, "plugin", "list", "--json").Output()
			if err == nil && strings.Contains(string(b), "mixrank@") {
				skipComponents[comp] = true
				messages = append(messages, "MixRank is already installed in "+host+"; preserving its marketplace registration. Update it through that client instead of adding a duplicate.")
			}
			if exec.CommandContext(ctx, host, "mcp", "get", "mixrank").Run() == nil {
				skipComponents[comp] = true
				messages = append(messages, "Existing MixRank MCP registration in "+host+" was preserved; skipped duplicate plugin registration.")
			}
		}
	}
	for _, step := range steps {
		if skipComponents[step.Component] {
			continue
		}
		currentComponent = step.Component
		if step.Manual != "" {
			messages = append(messages, step.Manual)
		}

		if (step.Component == "codex-plugin" || step.Component == "claude-plugin") && step.Destination != "" {
			if e = writeNativePlugin(step.Component, step.Destination, write); e != nil {
				return messages, e
			}
			continue
		}
		switch step.Component {
		case "cli":
			binary = step.Destination
			source, e := os.Executable()
			if e != nil {
				return messages, e
			}
			source, _ = filepath.EvalSymlinks(source)
			if source == binary {
				messages = append(messages, "CLI already at "+binary)
				continue
			}
			b, e := os.ReadFile(source)
			if e != nil {
				return messages, e
			}
			if e = write(binary, b, 0755); e != nil {
				return messages, e
			}
			messages = append(messages, "CLI installed at "+binary+"; add its directory to PATH")
		case "skills":
			e = fs.WalkDir(mixrank.SkillFiles, "skills", func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				b, e := mixrank.SkillFiles.ReadFile(p)
				if e != nil {
					return e
				}
				return write(filepath.Join(step.Destination, strings.TrimPrefix(p, "skills/")), b, 0644)
			})
			if e != nil {
				return messages, e
			}
		case "generic":
			b, _ := json.MarshalIndent(map[string]any{"mcpServers": map[string]any{"mixrank": map[string]any{"command": binary, "args": []string{"mcp"}}}}, "", "  ")
			if e = write(step.Destination, b, 0600); e != nil {
				return messages, e
			}
		default:
			if len(step.Command) > 0 {
				if _, e = exec.LookPath(step.Command[0]); e != nil {
					return messages, fmt.Errorf("%s is not installed; choose its import package or install the client first", step.Command[0])
				}
				if strings.Contains(step.Component, "mcp") && r.Commands[step.Component] == nil {
					check := []string{step.Command[0], "mcp", "get", "mixrank"}
					if exec.CommandContext(ctx, check[0], check[1:]...).Run() == nil {
						return messages, errors.New("an unowned mixrank MCP registration exists; inspect it before repair")
					}
				}
				if e = backupClientConfig(step.Command[0]); e != nil {
					return messages, e
				}
				if strings.Contains(step.Component, "mcp") && r.Commands[step.Component] != nil {
					remove := []string{step.Command[0], "mcp", "remove", "mixrank"}
					if step.Command[0] == "claude" {
						remove = []string{"claude", "mcp", "remove", "--scope", o.Scope, "mixrank"}
					}
					if e = run(ctx, remove); e != nil {
						return messages, e
					}
				}
				if e = run(ctx, step.Command); e != nil {
					return messages, e
				}
				if !contains(step.Command, "marketplace") {
					r.Commands[step.Component] = step.Command
					if e = save(); e != nil {
						return messages, e
					}
				}
			}
		}
	}
	return messages, save()
}
func contains(a []string, v string) bool {
	for _, x := range a {
		if x == v {
			return true
		}
	}
	return false
}
func run(ctx context.Context, a []string) error {
	cmd := exec.CommandContext(ctx, a[0], a[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if e := cmd.Run(); e != nil {
		return fmt.Errorf("client setup command failed: %s", strings.Join(a, " "))
	}
	return nil
}

func index(a []string, v string) int {
	for i, x := range a {
		if x == v {
			return i
		}
	}
	return -1
}
func backupClientConfig(client string) error {
	home, e := os.UserHomeDir()
	if e != nil {
		return e
	}
	var paths []string
	if client == "codex" {
		base := os.Getenv("CODEX_HOME")
		if base == "" {
			base = filepath.Join(home, ".codex")
		}
		paths = []string{filepath.Join(base, "config.toml")}
	} else {
		paths = []string{filepath.Join(home, ".claude.json"), filepath.Join(home, ".claude", "settings.json")}
		wd, _ := os.Getwd()
		paths = append(paths, filepath.Join(wd, ".mcp.json"), filepath.Join(wd, ".claude", "settings.json"))
	}
	for _, p := range paths {
		b, e := os.ReadFile(p)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return errors.New("client configuration backup failed")
		}
		config, e := profiles.Default()
		if e != nil {
			return e
		}
		dir := filepath.Join(config.Root, "backups", time.Now().UTC().Format("20060102T150405.000000000"))
		if e = os.MkdirAll(dir, 0700); e != nil {
			return e
		}
		sum := fmt.Sprintf("%x", sha256.Sum256([]byte(p)))
		if e = os.WriteFile(filepath.Join(dir, sum[:12]+"-"+filepath.Base(p)), b, 0600); e != nil {
			return e
		}
	}
	return nil
}

func writeNativePlugin(component, root string, write func(string, []byte, fs.FileMode) error) error {
	codex := component == "codex-plugin"
	meta := ".claude-plugin"
	variable := "${CLAUDE_PLUGIN_ROOT}"
	if codex {
		meta = ".codex-plugin"
		variable = "."
	}
	pluginRoot := filepath.Join(root, "plugins", "mixrank")
	b, e := mixrank.PluginManifests.ReadFile(meta + "/plugin.json")
	if e != nil {
		return e
	}
	if e = write(filepath.Join(pluginRoot, meta, "plugin.json"), b, 0644); e != nil {
		return e
	}
	binary := "mixrank"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	source, e := os.Executable()
	if e != nil {
		return e
	}
	b, e = os.ReadFile(source)
	if e != nil {
		return e
	}
	if e = write(filepath.Join(pluginRoot, "bin", binary), b, 0755); e != nil {
		return e
	}
	config := map[string]any{"command": variable + "/bin/" + binary, "args": []string{"mcp"}}
	if codex {
		config["env_vars"] = []string{"MIXRANK_API_KEY"}
		config["cwd"] = "."
	} else {
		config["env"] = map[string]string{"MIXRANK_API_KEY": "${user_config.api_key}"}
	}
	b, _ = json.MarshalIndent(map[string]any{"mcpServers": map[string]any{"mixrank": config}}, "", "  ")
	if e = write(filepath.Join(pluginRoot, ".mcp.json"), b, 0644); e != nil {
		return e
	}
	if codex {
		portable, e := mixrank.PluginManifests.ReadFile("plugin.json")
		if e != nil {
			return e
		}
		if e = write(filepath.Join(pluginRoot, "plugin.json"), portable, 0644); e != nil {
			return e
		}
		pmcp, _ := json.MarshalIndent(map[string]any{"$schema": "https://agent-plugins.org/schemas/1.0.0/mcp.schema.json", "mcpServers": map[string]any{"mixrank": map[string]any{"type": "stdio", "command": "./bin/" + binary, "args": []string{"mcp"}}}}, "", "  ")
		if e = write(filepath.Join(pluginRoot, "mcp.json"), pmcp, 0644); e != nil {
			return e
		}
	}
	e = fs.WalkDir(mixrank.SkillFiles, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, e := mixrank.SkillFiles.ReadFile(path)
		if e != nil {
			return e
		}
		return write(filepath.Join(pluginRoot, path), b, 0644)
	})
	if e != nil {
		return e
	}
	var marketplace any
	var marketpath string
	if codex {
		marketpath = filepath.Join(root, ".agents", "plugins", "marketplace.json")
		marketplace = map[string]any{"name": "mixrank-toolkit-local", "interface": map[string]string{"displayName": "MixRank Local Toolkit"}, "plugins": []any{map[string]any{"name": "mixrank", "source": map[string]string{"source": "local", "path": "./plugins/mixrank"}, "policy": map[string]string{"installation": "AVAILABLE", "authentication": "ON_INSTALL"}, "category": "Productivity"}}}
	} else {
		marketpath = filepath.Join(root, ".claude-plugin", "marketplace.json")
		marketplace = map[string]any{"name": "mixrank-toolkit-local", "description": "Native MixRank tools and research skills", "owner": map[string]string{"name": "Nick Kulavic"}, "plugins": []any{map[string]string{"name": "mixrank", "source": "./plugins/mixrank"}}}
	}
	b, _ = json.MarshalIndent(marketplace, "", "  ")
	return write(marketpath, b, 0644)
}
