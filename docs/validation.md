# Validation and release status

The catalog was built from signed-in MixRank documentation on September 9, 2026 and the supplied provider Elasticsearch document. Duplicate sections were removed by method/path. It contains 97 API operations plus five searches. Extra account-visible families include LinkedIn profile/post bulk jobs, identity contribution, privacy redaction and Twitter queries. Posts now document cached/fetch strategies. Documentation access is not a blanket license grant.

Local tests use synthetic fixtures for request encoding, multipart upload, numeric precision, matching candidates, semantic retries, bounded responses, cancellation, pagination partial progress, webhook verification, one-off validation serialization across goroutines/processes, MCP transports/auth/discovery, exports, installer planning and profile review/version preservation. Generated SDK/reference/plugin parity is checked separately. A configured live check makes one `/echo` call; no load tests, arbitrary email lists or live-fetch campaigns are installation checks.

Client validation is reported separately from compilation. Manifest acceptance and a successful MCP handshake do not prove every GUI import flow, permission set or provider data condition. Desktop MCPB/Cowork UI imports and real remote OAuth logins need the corresponding user client environment. Native packages are cross-built for macOS/Linux/Windows amd64 and arm64. Windows/Linux OS-vault runtime tests depend on those operating systems and an available secret store.

## Completed checks

- GitHub Actions on macOS, Windows and Linux: race-enabled Go tests, `go vet`, generated SDK/reference checks, synchronized plugin checks and catalog/skill invariants. See the [test runs](https://github.com/nkulavic/mixrank/actions/workflows/ci.yml) for commit-specific results.
- Six native targets compiled; 23 ZIP/MCPB archives checked for duplicate entries, unsafe paths, required license notices and matching release checksums.
- Codex personal plugin installed using its real cache directory; official SDK stdio handshake, discovery of 103 tools (102 operations plus catalog), and a catalog call passed from an unrelated working directory.
- Claude Code plugin and marketplace passed its strict manifest validator and the plugin installed. Its required sensitive API-key setting remains a host UI step: `/plugin configure mixrank@mixrank-toolkit`.
- Native macOS MCPB manifest passed the official MCPB validator and its bundled executable passed the same stdio test. The Desktop import UI was not exercised.
- MixRank-owned macOS Keychain save/read/replace/delete passed with an isolated synthetic entry. The authorized implementation credential was saved through the CLI workflow, status reported only its source, and one bounded live `/echo` check succeeded.
- Four portable skills passed structural validation. Research scenarios were reviewed against synthetic fixtures in the [workflow evaluation](workflow-evaluation.md); this was not an autonomous model benchmark.

No real validation lists, live-fetch campaigns, customer exports, product personalization or outbound messages were run during these checks. Windows Credential Manager and Linux Secret Service integration still need vault-enabled runtime checks; the portable test suite does not simulate a successful unlocked OS vault.

Known boundaries for the initial release:

- Full documented request/response field access is generated; response data remains raw JSON with preserved numbers instead of hundreds of rigid structs that discard newly returned fields.
- Elasticsearch bodies and mappings are exposed unchanged. Server-version-specific PIT/search_after support is not invented; callers use supported DSL with explicit bounds.
- Source marketplaces use POSIX launchers; Windows native plugin ZIPs and standalone MCP have direct binary commands.
- Cowork's importable skill package is built; live connector authentication awaits the hosting phase.
- HTTP bearer development mode is tested. Production OAuth and provider deployment remain deferred as planned.
- Private personalization is an agent-assisted draft/review workflow. Merely registering source locations does not claim that the product has been analyzed.

## Contact workflow, v0.2.0

The SDK, CLI and MCP share one company-to-contact implementation. Synthetic tests cover nested current-employer/role evidence, contradictory employers, inferred and stale associations, domain-filtered work emails, phone/email deduplication, continuing past people without contact data, request budgets, permission failures with partial reports, numeric precision in company imports, and output-file overwrite prevention. An MCP client invokes the complete workflow through its registered tool. Discovery now includes 104 tools: 102 provider operations, the catalog, and the contact workflow.

Contact enrichment is distinct from email validation. The contact workflow does not guess emails, request consumer email add-ons, send messages or silently trigger live refresh. The Illinois roofing live test uses private exports outside the repository; production contact values are not test fixtures.

## Catalog maintenance

Review newly signed-in documentation and provider revisions privately. Update factual names, types, parameters, operational notes and source dates in `catalog/endpoints.json` and mappings. Keep raw documentation snapshots outside the repository. Run `scripts/generate.py`, `scripts/package.py`, tests and their `--check` modes. A new undisclosed or account-inaccessible endpoint should be recorded as a limitation, not assigned a fabricated contract.

## Unified contacts and validation, v0.3.0

The output has one contacts array with explicit availability/enrichment fields, and optional filters. Tests use controlled overlapping HTTP handlers to prove person enrichment runs concurrently, output remains ranked, and a shared limit/request budget holds across multiple companies. Filters preserve company rows and distinguish not-enriched people from completed empty results.

Bulk validation tests assert deduplicated uploads, retained job IDs, no second submission on resume, provider evidence attachment, deferred filtering while pending, cancellation, uncertain submission protection, changed-address rejection, partial results, and bounded JSON/JSONL/CSV/gzip parsing. CLI and MCP integration tests cover the new flags/arguments and the resume tool. Current tool count is 105. Single-address serialization tests remain in the shared client suite.

Bounded live verification on September 9, 2026 ran the CLI against two supplied Illinois roofing companies with concurrency two. A separate bulk job validated two previously returned work emails; resume downloaded the provider's gzip JSONL and attached one `valid` and one `invalid` result with timestamps and other evidence. Live company/contact outputs remain in private user storage and are not fixtures or release assets.
