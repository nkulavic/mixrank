---
name: mixrank-knowledge
description: Interpret MixRank identities, datasets, Elasticsearch mappings, nested employment, data freshness, matching and contact classifications when constructing or assessing research queries.
---

Read only the mapping relevant to the query: [person2](references/person2.json), [companies](references/companies.json), [jobs](references/jobs.json), [org_name](references/org_name.json), [industries](references/industries.json). These are factual provider-supplied mappings; request bodies support the entire Elasticsearch DSL. See [query semantics](references/query-semantics.md) for joins, nested filters and examples.

Select the operation by purpose:

| Need | API family |
|---|---|
| Discover a population | Elasticsearch search or documented directories/segments |
| Resolve supplied identifiers | Person or company match; retain candidates and confidence |
| Detailed profile or explicit refresh | LiveScan; cached by default |
| Contact addresses/numbers | Person record with requested add-ons |
| Deliverability evidence | Email validation; bulk for lists |
| Professional activity context | Person/company posts, with publication and fetch dates |
| Historical observations | Company headcount/follower timeseries |
| Technology evidence | App SDKs or website tags, including install/uninstall dates |

Names alone are weak identifiers. Keep MixRank company/person IDs, LinkedIn org/profile/user IDs, domains and canonical URLs distinct. A match can be empty, ambiguous or contradicted by another identifier. Return the candidate set and evidence; do not silently take the first result. Person Match defaults to one result at the provider, so increase page_size when judging ambiguity.

Person `enable` add-ons include `b2b_emails`, `b2c_emails`, `edu_emails`, `directdials`, and documented social content. No add-ons are enabled by default. A disabled type may return `{}` when data exists; `null` denotes missing data. Company `company_financials` is a paid add-on, not part of every response.

Keep original data timestamps separate from retrieval time. A cached result can be stale. Incomplete or privacy-redacted records are not negative evidence. Current employment requires a current flag on the same nested experience item as employer/role criteria. Historical employment is useful context but not proof of a current buyer role.

Validation outcomes include valid, ambiguous, maybe_valid, invalid, uncached, malformed, temp_error and perm_error. Preserve `validated_at`, catchall/disposable/consumer flags, greylist status and `retry_after`; ambiguous is not valid. Contact classification and validation are separate evidence. A result can be valid at one time without guaranteeing future delivery.

Headcount and followers are sampled observations; gaps and different history lengths are expected. Changes in follower/headcount/SDK presence can support qualification but do not establish buying intent. Attribute optional external research separately. License and response shape can vary with account entitlements; consult the installed operation catalog.
