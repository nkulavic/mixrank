# MixRank v0.4.0

Company-contact workflows now merge duplicate input records before making API requests. CLI, Go SDK, and both MCP transports share the same logic: exact normalized domains and overlapping company IDs form one business, preserving alternate names, domains, IDs, and stable order. The report exposes merge counts. Search and email filtering use the retained aliases; names alone never establish identity.

Accept up to 250 input records resolving to at most 25 businesses. Merged identity limits are checked before network access. Conflicting qualification evidence stays marked for review. Shared-domain branches are grouped as a business account; distinct location IDs can be supplied without the shared domain when location separation is needed.

Regression coverage checks transitive merges, duplicate request elimination, alias-domain email filtering, numeric precision, stable order, unchanged caller inputs, identity bounds, and matching CLI/MCP behavior. Existing contact filters, bounded concurrency, and deduplicated bulk email validation remain available.
