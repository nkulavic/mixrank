# Client installation

The interactive bootstrap downloads a pinned native binary, verifies its checksum and invokes `mixrank setup`. Available components are `cli`, `skills`, `codex-plugin`, `codex-mcp`, `claude-plugin`, `claude-mcp`, `claude-desktop`, `cowork`, and `generic`. Select either the full plugin or standalone MCP for a client. No Go installation is needed by end users.

`--dry-run --yes --components ...` prints exact destinations and client commands without mutation. `--scope user|project` selects supported scope; Codex MCP/plugin uses shared user scope. `--destination` chooses a binary directory or standalone skills destination. `--upgrade` and `--repair` refresh selected owned components; use the newest release bootstrap for a newer binary. Backups of edited client configuration stay in private MixRank storage. Existing unrelated or modified files are preserved.

## Codex CLI and Desktop

```sh
codex plugin marketplace add https://github.com/nkulavic/mixrank.git
codex plugin add mixrank@mixrank-toolkit
```

The repository's `.agents/plugins/marketplace.json` resolves `plugins/mixrank`. It includes a portable root manifest/MCP file plus the required Codex compatibility manifest. Native wizard setup generates `mixrank-toolkit-local` in private user application storage, with the current binary bundled, and installs that local marketplace. The personal Codex marketplace uses the same plugin through its own entry when installed locally. Install only one marketplace copy. The default personal marketplace is discovered automatically; do not add it again with a second name.

For standalone MCP, run `mixrank setup --components cli,codex-mcp --yes`. It installs one absolute-path command in the user configuration shared by CLI/Desktop; GUI launches use the OS vault and do not depend on shell-exported keys. Restart/open a new agent session after plugin changes. [Codex plugin format](https://developers.openai.com/plugins/build/plugins).

## Claude Code

```sh
claude plugin marketplace add nkulavic/mixrank
claude plugin install mixrank@mixrank-toolkit
```

The GitHub marketplace references `plugins/claude/mixrank`, whose sensitive `userConfig` field lets Claude prompt for and inject the key. Use `/reload-plugins` or restart after installation. Do not put the key in `--config` arguments. Native plugin configuration differs from standalone `mixrank setup --components cli,claude-mcp --yes`, which uses MixRank's own vault.

## Windows plugin bundles

The GitHub source marketplaces use POSIX shell launchers for macOS/Linux. The wizard uses native binaries and local marketplaces on every OS. For manual native Windows installation, use the platform-specific Codex/Claude plugin ZIP, whose MCP configuration launches its bundled `.exe` directly, or use the standalone native MCP installer. Extract platform ZIPs into a trusted local plugin/marketplace location and use the client's local-plugin installation mechanism. This avoids requiring Git Bash or a Go toolchain. CLI and Desktop native MCPB packages also have Windows amd64/arm64 binaries.

## Claude Desktop

Download the `.mcpb` matching OS and CPU from releases. Open/import it in Desktop, enter the API key in the sensitive field, and enable it. Desktop encrypts and injects the value; do not manually copy it into `claude_desktop_config.json`. The bundle contains its native server. Install/import the skills/plugin package separately where the client exposes plugin support. Actual import UI and restart may require user interaction. [Desktop bundle manifest](https://github.com/modelcontextprotocol/mcpb/blob/main/MANIFEST.md).

## Cowork

Import `mixrank-cowork.zip` through Cowork's plugin UI. It contains the maintained skills and connector instructions. Cowork is a separate runtime: this package does not assume host binary paths, OS keychains or Desktop MCP registrations are available. Until hosted OAuth deployment, it supports supplied-file research and personalization; live API workflows need the later authenticated connector. No placeholder hosted endpoint is registered. [Cowork plugin behavior](https://code.claude.com/docs/en/plugins-reference#plugins-synced-from-claudeai).

## Other agents

Use `mixrank setup --components cli,generic,skills --yes` to create standard stdio MCP JSON with an absolute binary path and copy portable `SKILL.md` directories. Import that JSON using the agent's MCP configuration UI or file conventions. A custom skills path can be passed with a separate `setup --components skills --destination DIRECTORY` call. The HTTP endpoint is `/mcp`; remote authentication depends on the deployment's supported flow.

## Uninstall

`mixrank setup --uninstall --components cli,skills,codex-mcp --yes` removes selected installer-owned files/registrations. Modified files are retained. Credentials and private profiles survive by default. Add `--remove-credentials` or `--remove-profiles` only when you intend to delete that data. Plugin removal uses native client commands; the installer never edits host credential stores. A component installed outside the wizard should be removed through the same host mechanism that installed it.
