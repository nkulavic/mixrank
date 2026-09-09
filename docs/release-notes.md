# MixRank v0.3.0

Company contact workflows now use one `contacts` array with availability and enrichment status. `--contact-filter` selects people with email, phone, both, either, no returned channels, or valid emails. `all` preserves people who were not enriched because the target or budget was reached. The previous `candidates_without_contacts` output field is removed.

`--concurrency` runs company searches and person enrichment in parallel under one bounded request limiter (default 4, maximum 16). Results retain input company order and ranked contact order. Request budgets, cancellation, and semantic retry rules remain shared across CLI, SDK and both MCP transports.

`--validate-emails` deduplicates returned work emails into a MixRank bulk-validation job and attaches provider validity, dates and other evidence. Pending jobs retain their ID and original address set. `mixrank contacts validate --input REPORT` and `mixrank_validate_contacts` resume without submitting another job. Single-address validation remains serialized per credential across local processes.

The CLI, SDK and agent skills expose the same workflow. MCP now has 105 tools: 102 catalog operations, catalog discovery, company contacts and contact-report validation. Native binaries, Codex/Claude plugins, Desktop MCPB packages, portable skills, installers and checksums are included. Start a new Codex task or restart Claude to load updated installed plugins. Cowork's hosted connector and remote OAuth remain deployment work.

Mock tests cover unified output/filter semantics, real overlapping requests and global concurrency/request bounds, result ordering, bulk submission deduplication, pending/resume behavior, cancellation, uncertain outcomes, partial validation results and supported bounded download formats. Race detection, static analysis, generated-file checks and native build checks run in CI. No messages are sent by these workflows.
