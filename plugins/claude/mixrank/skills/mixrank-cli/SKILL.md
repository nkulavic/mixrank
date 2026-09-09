---
name: mixrank-cli
description: Use the MixRank native CLI for API discovery, company and person data, bulk enrichment, email validation, exports, and troubleshooting.
---

Use `mixrank --help` and `mixrank endpoints` to discover the installed version. If `mixrank` is not on PATH, use the plugin's `scripts/run-mixrank` launcher (or `scripts/run-mixrank.ps1` on Windows) with the same arguments. The launcher installs a checksum-verified, version-pinned binary if the bundle does not already contain one.

Read [the generated operation reference](references/operations.md) for the relevant family. `mixrank endpoints OPERATION` shows its parameters, encoding, response fields, permissions and source. `mixrank api COMMAND --help` lists actual flags; command names replace operation underscores with hyphens. Read [recipes](references/recipes.md) for end-to-end examples.

- Authentication: let the host prompt for a plugin secret when supported. For CLI/Codex use `mixrank auth login` (masked OS-vault entry) or explicit stdin. Never put a key in arguments, a repository, a skill, or a client configuration file. Runtime `MIXRANK_API_KEY` overrides the MixRank vault entry. `mixrank auth status` reveals availability only.
- Pass JSON with `--body request.json` or `--body -`; query parameters with named flags or repeatable `--query key=value`. The SDK passes full Elasticsearch DSL unchanged and preserves integer precision.
- Bulk input is UTF-8 newline-separated text via `--upload`. Preserve job IDs; inspect status and download the returned result URL. Do not resubmit after timeout merely because the outcome is unknown.
- Use `--output filename` for binary/large responses; it refuses existing files. `--format jsonl` or `csv` exports object rows; CSV formulas are escaped. Offset pagination requires a finite `--max-pages`; it emits one JSONL envelope per page. Elasticsearch pagination stays explicit in DSL.
- Prefer bulk email validation for lists. Single validation is serialized per credential across local toolkit processes. Separate hosts need coordination. Do not evade the rule with shell parallelism, raw calls, or another client. Only run authorized checks; `doctor --live` uses a single echo call.
- LiveScan defaults to cached. A requested fetch/strict/besteffort strategy requires `--refresh`. Changing strategy is a user decision; no implicit fallback or retry. Posts also document refresh as of the catalog date.
- Keep partial outputs and job identifiers on failure. 403 can mean missing license; 512 is a temporary LiveScan fetch failure and 513 permanent. Report provider errors and source dates without exposing credentials. Never treat an empty object as equivalent to absent/null data.

Exit codes: 0 success; 2 local input/config error; 3 credentials or permission; 4 rate limit; 5 upstream failure; 6 uncertain side-effect outcome; 130 cancellation. A successful 202 means accepted, not completed. For asynchronous one-off validation, use `mixrank validation status` and acknowledge the matching signed callback with `validation ack --headers headers.json --body callback.json`; later one-offs remain blocked until then. Clear uncertain state only after independently confirming provider completion.
