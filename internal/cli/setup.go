package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank/internal/credentials"
	"github.com/nkulavic/mixrank/internal/profiles"
	"github.com/nkulavic/mixrank/internal/setup"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func setupCommand() *cobra.Command {
	var o setup.Options
	var components string
	var yes, upgrade, repair, uninstall, removeCredentials, removeProfiles bool
	c := &cobra.Command{Use: "setup", Short: "Interactive installation, dry runs, upgrade, repair and uninstall", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if (upgrade && repair) || (upgrade && uninstall) || (repair && uninstall) {
			return errors.New("select one of upgrade, repair or uninstall")
		}
		if (removeCredentials || removeProfiles) && !uninstall {
			return errors.New("private-data removal flags require --uninstall")
		}
		if upgrade {
			o.Action = "upgrade"
		}
		if repair {
			o.Action = "repair"
		}
		if uninstall {
			o.Action = "uninstall"
		}
		interactive := !yes && term.IsTerminal(int(os.Stdin.Fd()))
		reader := bufio.NewReader(c.InOrStdin())
		ask := func(label string) string {
			fmt.Fprint(c.ErrOrStderr(), label)
			s, _ := reader.ReadString('\n')
			return strings.TrimSpace(s)
		}
		if components == "" && interactive {
			fmt.Fprintf(c.ErrOrStderr(), "MixRank setup — %s/%s\n", runtime.GOOS, runtime.GOARCH)
			for _, name := range []string{"codex", "claude"} {
				_, e := exec.LookPath(name)
				fmt.Fprintf(c.ErrOrStderr(), "%s detected: %t\n", name, e == nil)
			}
			fmt.Fprintln(c.ErrOrStderr(), "Components: cli, codex-plugin, claude-plugin, codex-mcp, claude-mcp, claude-desktop, cowork, skills, generic")
			components = ask("Select comma-separated components [cli,skills]: ")
			if components == "" {
				components = "cli,skills"
			}
		}
		if components == "" {
			return errors.New("non-interactive setup requires --components")
		}
		for _, s := range strings.Split(components, ",") {
			o.Components = append(o.Components, strings.TrimSpace(s))
		}
		steps, e := setup.Plan(o)
		if e != nil {
			return e
		}
		if o.DryRun {
			return jsonOut(c, map[string]any{"action": o.Action, "steps": steps, "remove_credentials": removeCredentials, "remove_profiles": removeProfiles})
		}
		messages, e := setup.Apply(c.Context(), o, steps)
		if e != nil {
			return e
		}
		for _, m := range messages {
			fmt.Fprintln(c.ErrOrStderr(), m)
		}
		if uninstall {
			if removeCredentials {
				if e = credentials.Delete(); e != nil {
					return e
				}
			}
			if removeProfiles {
				store, err := profiles.Default()
				if err != nil {
					return err
				}
				names, err := store.List()
				if err != nil {
					return err
				}
				for _, name := range names {
					if err = store.Remove(name); err != nil {
						return err
					}
				}
			}
			return jsonOut(c, map[string]any{"uninstalled": true, "private_data_preserved": !removeCredentials})
		}
		if interactive {
			_, _, e := credentials.Resolve()
			ownAuth := false
			for _, comp := range o.Components {
				if comp == "cli" || comp == "skills" || comp == "codex-plugin" || comp == "codex-mcp" || comp == "claude-mcp" || comp == "generic" {
					ownAuth = true
				}
			}
			if ownAuth && e != nil && strings.EqualFold(ask("Configure the MixRank OS-vault API key now? [y/N]: "), "y") {
				cmd := auth()
				cmd.SetArgs([]string{"login"})
				cmd.SetIn(c.InOrStdin())
				cmd.SetOut(c.OutOrStdout())
				cmd.SetErr(c.ErrOrStderr())
				if e = cmd.ExecuteContext(c.Context()); e != nil {
					return e
				}
			}
			if strings.EqualFold(ask("Prepare optional private product personalization? [y/N]: "), "y") {
				fmt.Fprintln(c.ErrOrStderr(), "Run mixrank personalize --website URL --repo OWNER/REPO --path PATH --agent codex (or claude). Sources can be combined.")
			}
		}
		d := doctor()
		d.SetOut(c.OutOrStdout())
		d.SetErr(c.ErrOrStderr())
		return d.ExecuteContext(c.Context())
	}}
	c.Flags().StringVar(&components, "components", "", "Comma-separated components")
	c.Flags().StringVar(&o.Scope, "scope", "user", "user or project")
	c.Flags().StringVar(&o.Destination, "destination", "", "CLI binary directory or portable skills destination")
	c.Flags().BoolVar(&o.DryRun, "dry-run", false, "Show exact actions without changes")
	c.Flags().BoolVarP(&yes, "yes", "y", false, "Non-interactive installation; no secret prompt")
	c.Flags().BoolVar(&upgrade, "upgrade", false, "Upgrade owned files and selected components")
	c.Flags().BoolVar(&repair, "repair", false, "Repair selected owned components")
	c.Flags().BoolVar(&uninstall, "uninstall", false, "Remove installer-owned files and registrations")
	c.Flags().BoolVar(&removeCredentials, "remove-credentials", false, "Delete MixRank OS credential during uninstall")
	c.Flags().BoolVar(&removeProfiles, "remove-profiles", false, "Delete all private MixRank profiles during uninstall")
	return c
}
