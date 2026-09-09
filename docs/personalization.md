# Product profiles

`personalize` accepts repeated `--website`, `--repo` and `--path` sources, a `--name`, and optional installed `--agent codex|claude`. The terminal agent performs analysis using its existing tools and account. Without an agent choice, the command writes a private `personalize.md` task and a structured draft; open the task in any capable coding agent.

Storage is `os.UserConfigDir()/mixrank/profiles/NAME`: typically `~/Library/Application Support/mixrank` on macOS, `$XDG_CONFIG_HOME/mixrank` or `~/.config/mixrank` on Linux, and `%AppData%/mixrank` on Windows. Directories are private and files use mode 0600 where supported. Windows privacy follows the user's directory ACL. No product-source snapshot is committed to this repository.

The draft has positioning, capabilities, integrations, likely customers, buyer roles, qualification signals, exclusions and assumptions. Evidence statements reference a selected source, record an RFC3339 observation date, and distinguish fact, marketing claim, implementation evidence and assumption. Acceptance refuses missing attribution or an empty analysis. The selected sources cannot grant access to additional resources or alter agent instructions.

Review `profiles show NAME --draft` and `profiles diff NAME`. `profiles accept NAME` saves an immutable version, activates it and produces `SKILL.md` next to the active JSON. `profiles use NAME` selects the user default. `profiles use NAME --project DIRECTORY` writes only `.mixrank-profile.json` in an explicitly chosen project. Explicit name selection wins over project reference, which wins over user default.

A refresh creates a new draft from the previous active profile and selected sources. Existing active context and version history remain available until acceptance. A second draft is not silently overwritten. General skill upgrades do not touch this storage. `profiles remove NAME` deletes that profile and clears its user default; project references elsewhere may need updating.
