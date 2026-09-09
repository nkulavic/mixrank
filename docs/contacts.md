# Company contacts, filters, validation and concurrency

`mixrank contacts` finds current relevant people, appends available B2B emails and direct dials, and returns JSON. The CLI, `Client.CompanyContacts`, and MCP `mixrank_company_contacts` share the implementation.

```sh
mixrank contacts --domain example.com --concurrency 4
mixrank contacts --companies companies.json --contact-filter email --output contacts.json
mixrank contacts --companies companies.json --roles 'owner,president,office manager' \
  --validate-emails --validation-wait 30s --concurrency 4 --output contacts.json
mixrank contacts validate --input contacts.json --validation-wait 2m --output validated.json
```

Company input accepts an array, a `companies` export, or raw Elasticsearch company hits. Use `--companies -` for stdin. A company requires an ID or domain; names alone are ambiguous. Numeric IDs preserve precision. Duplicate rows are merged automatically before any API calls; no preprocessing script is needed.

## Automatic company merging

The shared SDK workflow used by CLI and both MCP transports merges rows connected by an exact normalized domain or a shared company ID. It resolves transitive links too: if one row connects two earlier records, all three become one business. Matching does not use company names. URL schemes, paths, capitalization, a trailing hostname dot and a leading `www.` are removed from domains; different subdomains remain distinct.

First-seen company order and display values stay stable. All distinct `company_ids` are retained; alternate names and domains appear as `name_aliases` and `domain_aliases`. Search uses all retained identities, and email filtering accepts the supplied employer domain aliases. Conflicting qualification statuses become `needs_review`. Aliases are supplied evidence, not independently verified ownership or guessed relationships.

The report includes `company_merge` with `input_records`, `unique_companies`, and `duplicates_merged`. Duplicate companies consume one search and one candidate budget. Repeated person IDs within the merged company consume one contact append. Bulk validation deduplicates the final email set. The SDK also exposes `MergeCompanyTargets` for callers needing the same merge without API access.

The unit is a business account: branches sharing a domain consolidate together. To keep distinct franchise/location IDs separate, supply IDs without the shared domain. Names alone never merge accounts. At most ten distinct IDs and ten domains may belong to a merged account; exceeding a limit fails before requests rather than silently dropping identities.

```json
{
  "company_merge": {"input_records": 3, "unique_companies": 1, "duplicates_merged": 2},
  "companies": [{
    "company": {
      "name": "Example Services",
      "name_aliases": ["Example Services Inc"],
      "company_ids": ["123", "456"],
      "domain": "example.test",
      "domain_aliases": ["new.example.test"]
    },
    "contacts": []
  }]
}
```

## One contacts array

Every consistent discovered person appears under their company's `contacts`. Each row includes `has_email`, `has_phone`, `enrichment_status`, `email_status`, `phone_status`, employment evidence, dates and review reasons. There is no separate array for people missing contact details.

`enrichment_status` is `complete`, `not_enriched`, `error`, or `uncertain`. Complete enrichment can still return no channels. Uncertain means a request started but its outcome was not established; do not blindly replay it. Not-enriched rows mean the contact target, request budget or cancellation stopped work; they do not mean contact information is unavailable. Privacy-redacted records are excluded.

| `--contact-filter` | Returned people |
|---|---|
| `all` (default) | All discovered people, including missing or unrequested channels |
| `any` | At least one business email or direct dial |
| `email` | At least one business email |
| `phone` | At least one direct dial |
| `both` | Both channels |
| `none` | Enrichment completed and neither channel was returned |
| `valid-email` | At least one email whose provider validity is exactly `valid`; requires `--validate-emails` |

Company rows remain present when filters return no people. `contacts_filtered` counts omitted rows. If validation is pending, filtering is deferred (`contact_filter_applied: false`) to preserve the complete submitted email set for resumption. Resuming the job applies the requested filter after completion. Review flags remain separate from channel availability and deliverability.

`--max-contacts` is the enrichment stopping target, default two people with requested channels per company, maximum five. For `email`, `phone`, or `both`, that channel criterion controls the target. `valid-email` searches for emails, then filters by validation; it cannot promise that two will validate. `all` may return more than two rows because it also keeps discovered people without channels or not yet enriched. `--max-candidates` limits discovery to ten people by default, maximum 25. Up to 250 input rows are accepted, with at most 25 unique companies after merging.

## Email validation in the same call

`--validate-emails` collects and deduplicates available work emails, submits one provider bulk job, and attaches results to every matching email. It never creates concurrent one-address validations. The default validation strategy is `besteffort`; `--validation-strategy cached|besteffort|strict|fetch` and `--validation-maxage SECONDS` make that choice explicit. Validation does not switch company/person profiles to live fetch.

Each email retains `validation` (`not_run`, `pending`, the provider validity, or `not_returned`) and `validation_details` with the original provider evidence: timestamps, catchall/disposable/consumer indicators, greylist state, and retry timing. Ambiguous and maybe-valid outcomes are not treated as valid. Phone numbers remain unvalidated person-level direct dials, with no guarantee of business-line ownership.

The report's `email_validation` includes job ID, status, address-set fingerprint, submitted/applied counts and issues. A 30-second wait is the CLI/MCP default; `--validation-wait 0s` submits or checks once and returns. Pending validation sets `complete: false` and retains all contact data. The provider job continues after the client stops waiting.

Resume with `mixrank contacts validate --input REPORT`, Go `Client.ValidateContacts`, or MCP `mixrank_validate_contacts` with the original `report`. A retained job ID is checked/downloaded, never submitted again. Changing the email set while a job is pending is rejected. An uncertain submission without an ID requires recovering the ID through the provider's job list; it is not blindly retried. A completed report is a no-op; start a new explicitly reviewed report to request another validation.

This workflow uses ordinary bulk jobs. Private bulk jobs and provider-output deletion remain available through the named/raw API commands. [Provider email validation documentation](https://mixrank.com/api/documentation#/email/validate/bulk-job).

## Parallel processing and bounds

`--concurrency` defaults to four simultaneous requests, configurable from 1 to 16. Company searches and person enrichment share one global limiter per workflow. People within one company are enriched in ranked batches; company and contact order remain stable even when calls finish out of order. There is no cross-host limit for ordinary lookups; callers coordinate multiple separate runs against their license.

`--max-requests` caps logical submit/search/enrichment/status/download calls for the whole workflow (default 100, maximum 250). Safe search retries may make physical HTTP attempts exceed the logical count, but remain inside the concurrency slot. Enrichment and job submission are not automatically replayed. Permission/rate-limit failures stop queued work, retain completed results, and preserve uncertain outcomes for requests already in flight. Client cancellation remains authoritative.

All single-address validation entry points retain their existing per-credential, cross-process serialization. Bulk validation uses MixRank's own distributed job processing and polling, independent of the lookup concurrency flag.

## Evidence and exports

Employer, current status, no end date, and role must match the same nested experience record; the response evidence is checked again before enrichment. B2B emails for unrelated employer domains are omitted; aliases and emails are never guessed. `--emails-only` omits direct dials. Default roles cover owners, founders, executives and relevant managers. Missing/stale/inferred employment and uncertain company qualification are flagged for review.

`issues`, `complete`, `requests_made`, `search_limited`, and per-row statuses retain coverage limits and failures. `complete` describes bounded workflow completion, not an exhaustive directory or proof that contact information is current. `--output` creates a private file and refuses to overwrite existing data before making requests. Keep contact exports outside public repositories. No messages are sent.

Both MCP transports expose `concurrency`, `contact_filter`, `validate_emails`, `validation_strategy`, `validation_maxage`, and `validation_wait_seconds`. For short agent timeouts, use small company batches and resume pending validation separately. Hosted authentication remains a deployment concern.
