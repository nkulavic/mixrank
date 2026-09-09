# Validation and release status

The catalog was built from signed-in MixRank documentation on September 9, 2026 and the supplied provider Elasticsearch document. Duplicate sections were removed by method/path. It contains 97 API operations plus five searches. Extra account-visible families include LinkedIn profile/post bulk jobs, identity contribution, privacy redaction and Twitter queries. Posts now document cached/fetch strategies. Documentation access is not a blanket license grant.

Local tests use synthetic fixtures for request encoding, multipart upload, numeric precision, matching candidates, semantic retries, bounded responses, cancellation, pagination partial progress, webhook verification, one-off validation serialization across goroutines/processes, MCP transports/auth/discovery, exports, installer planning and profile review/version preservation. Generated SDK/reference/plugin parity is checked separately. A configured live check makes one `/echo` call; no load tests, arbitrary email lists or live-fetch campaigns are installation checks.

Client validation is reported separately from compilation. Manifest acceptance and a successful MCP handshake do not prove every GUI import flow, permission set or provider data condition. Desktop MCPB/Cowork UI imports and real remote OAuth logins need the corresponding user client environment. Native packages are cross-built for macOS/Linux/Windows amd64 and arm64. Windows/Linux OS-vault runtime tests depend on those operating systems and an available secret store.

Known boundaries for the initial release:

- Full documented request/response field access is generated; response data remains raw JSON with preserved numbers instead of hundreds of rigid structs that discard newly returned fields.
- Elasticsearch bodies and mappings are exposed unchanged. Server-version-specific PIT/search_after support is not invented; callers use supported DSL with explicit bounds.
- Source marketplaces use POSIX launchers; Windows native plugin ZIPs and standalone MCP have direct binary commands.
- Cowork's importable skill package is built; live connector authentication awaits the hosting phase.
- HTTP bearer development mode is tested. Production OAuth and provider deployment remain deferred as planned.
- Private personalization is an agent-assisted draft/review workflow. Merely registering source locations does not claim that the product has been analyzed.

## Catalog maintenance

Review newly signed-in documentation and provider revisions privately. Update factual names, types, parameters, operational notes and source dates in `catalog/endpoints.json` and mappings. Keep raw documentation snapshots outside the repository. Run `scripts/generate.py`, `scripts/package.py`, tests and their `--check` modes. A new undisclosed or account-inaccessible endpoint should be recorded as a limitation, not assigned a fabricated contract.
