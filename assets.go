package mixrank

import "embed"

// SkillFiles is the portable maintained skill distribution embedded in binaries.
//
//go:embed skills
var SkillFiles embed.FS

// PluginManifests are the host-specific manifests used by native local setup.
//
//go:embed .codex-plugin/plugin.json .claude-plugin/plugin.json plugin.json
var PluginManifests embed.FS
