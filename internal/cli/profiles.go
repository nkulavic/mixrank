package cli

import (
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank/internal/profiles"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func personalizeCommand() *cobra.Command {
	var sites, repos, paths []string
	var name, agent string
	c := &cobra.Command{Use: "personalize", Short: "Prepare a private product profile for an installed coding agent to review", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		s, e := profiles.Default()
		if e != nil {
			return e
		}
		sources := []profiles.Source{}
		for _, v := range sites {
			if !strings.HasPrefix(v, "https://") && !strings.HasPrefix(v, "http://") {
				return errors.New("website must be an HTTP(S) URL")
			}
			sources = append(sources, profiles.Source{Kind: "website", Location: v})
		}
		for _, v := range repos {
			sources = append(sources, profiles.Source{Kind: "repo", Location: v})
		}
		for _, v := range paths {
			p, e := filepath.Abs(v)
			if e != nil {
				return e
			}
			if _, e = os.Stat(p); e != nil {
				return e
			}
			sources = append(sources, profiles.Source{Kind: "path", Location: p})
		}
		task, e := s.Draft(name, sources)
		if e != nil {
			return e
		}
		if agent == "" {
			return jsonOut(c, map[string]any{"task": task, "profile": name, "state": "draft", "next": "Ask your coding agent to read the task and prepare the profile, then review and run mixrank profiles accept " + name})
		}
		if agent != "codex" && agent != "claude" {
			return errors.New("agent must be codex or claude")
		}
		cmd := exec.CommandContext(c.Context(), agent, "Read "+task+" and perform its MixRank personalization task. Keep the profile private and present it for review before activation.")
		cmd.Stdin = c.InOrStdin()
		cmd.Stdout = c.OutOrStdout()
		cmd.Stderr = c.ErrOrStderr()
		return cmd.Run()
	}}
	c.Flags().StringArrayVar(&sites, "website", nil, "Website URL (repeatable)")
	c.Flags().StringArrayVar(&repos, "repo", nil, "Authorized owner/repo (repeatable)")
	c.Flags().StringArrayVar(&paths, "path", nil, "Local source path (repeatable)")
	c.Flags().StringVar(&name, "name", "default", "Private profile name")
	c.Flags().StringVar(&agent, "agent", "", "Open installed codex or claude for analysis; omit to write an agent task")
	return c
}
func profilesCommand() *cobra.Command {
	root := &cobra.Command{Use: "profiles", Short: "Review and manage private, versioned personalization"}
	root.AddCommand(&cobra.Command{Use: "list", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		s, e := profiles.Default()
		if e != nil {
			return e
		}
		v, e := s.List()
		if e != nil {
			return e
		}
		return jsonOut(c, v)
	}})
	var draft bool
	show := &cobra.Command{Use: "show [name]", Args: cobra.MaximumNArgs(1), RunE: func(c *cobra.Command, a []string) error {
		s, e := profiles.Default()
		if e != nil {
			return e
		}
		name := ""
		if len(a) > 0 {
			name = a[0]
		}
		cwd, _ := os.Getwd()
		name, e = s.Selected(name, cwd)
		if e != nil {
			return e
		}
		p, e := s.Show(name, draft)
		if e != nil {
			return e
		}
		return jsonOut(c, p)
	}}
	show.Flags().BoolVar(&draft, "draft", false, "Show pending draft instead of active version")
	root.AddCommand(show)
	var project string
	use := &cobra.Command{Use: "use name", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, a []string) error {
		s, e := profiles.Default()
		if e != nil {
			return e
		}
		if e = s.Select(a[0], project); e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"profile": a[0], "project": project})
	}}
	use.Flags().StringVar(&project, "project", "", "Write a reference only in this project directory")
	root.AddCommand(use)
	for _, action := range []string{"accept", "refresh", "remove"} {
		action := action
		root.AddCommand(&cobra.Command{Use: action + " name", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, a []string) error {
			s, e := profiles.Default()
			if e != nil {
				return e
			}
			switch action {
			case "accept":
				e = s.Accept(a[0])
			case "remove":
				e = s.Remove(a[0])
			case "refresh":
				var task string
				task, e = s.Draft(a[0], nil)
				if e == nil {
					return jsonOut(c, map[string]any{"task": task, "active_preserved": true})
				}
			}
			if e != nil {
				return e
			}
			fmt.Fprintln(c.OutOrStdout(), "Profile "+action+" complete")
			return nil
		}})
	}
	diff := &cobra.Command{Use: "diff name", Short: "Show active and draft profiles side by side as JSON", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, a []string) error {
		s, e := profiles.Default()
		if e != nil {
			return e
		}
		old, _ := s.Show(a[0], false)
		next, e := s.Show(a[0], true)
		if e != nil {
			return e
		}
		return jsonOut(c, map[string]any{"active": old, "draft": next})
	}}
	root.AddCommand(diff)
	return root
}
