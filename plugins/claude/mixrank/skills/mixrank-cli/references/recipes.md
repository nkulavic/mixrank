# Recipes

Start with local diagnostics: `mixrank doctor`. Opt into a bounded API authentication check with `mixrank doctor --live`.

Resolve companies, preserving candidates:

```sh
mixrank api get-companies-match --url example.com
mixrank api get-person-match --first-name Alex --last-name Example --domain example.com --page-size 10
```

Company IDs, LinkedIn organization IDs and person IDs belong to different namespaces. Do not substitute one for another without a documented relationship.

Read a cached person, request contact add-ons, and refresh only when authorized:

```sh
mixrank api get-linkedin-profile --person.id 123 --strategy cached
mixrank api get-person-by-id --id 123 --enable b2b_emails,directdials
mixrank api get-linkedin-profile --person.id 123 --strategy fetch --refresh
```

These are illustrative IDs, not live-check fixtures. Contact append data does not establish current deliverability.

Validate a list through bulk jobs:

```sh
mixrank api post-email-validate-bulk-job --upload emails.txt --strategy cached
mixrank api get-email-validate-bulk-job
mixrank api get-email-validate-bulk-job-by-job-id --job-id JOB_ID
mixrank download --url-file download-url.txt --output validated.jsonl
```

Inspect `completed_at` and `download_url`; a blank URL means no file yet. `download-url.txt` contains the URL returned by the job, not an API key. Downloading does not delete the provider output. If the user requested private validation, use `--private t --strategy fetch` and explicitly delete the output once the saved file is verified:

```sh
mixrank api delete-email-validate-bulk-job-by-job-id --job-id JOB_ID
```

This deletion does not cancel a submitted job. The general toolkit does not submit arbitrary lists during setup or CI.

Search with DSL:

```sh
mixrank endpoints --mapping person2
mixrank api search-person2 --body people-query.json --output people.json
mixrank api get-companies-by-company-id-timeseries --company-id 123 --since 2024-01-01 --page-size 100 --max-pages 3
```

The JSON returned by Elasticsearch retains `_source`, hit scores, sort cursors, aggregation buckets, and total-count relation. `hits.total.relation=gte` is a lower bound; it is not an exact count.

## Duplicate company rows

Pass company arrays or Elasticsearch hits directly to `mixrank contacts --companies FILE` (MCP: `mixrank_company_contacts`). The shared workflow merges exact company-ID/domain matches before lookup, retaining alternate IDs, names and domains. Inspect `company_merge` for input and unique counts. No extra deduplication script is needed.

For a lead-ready response, add `--contactable-only` (MCP `contactable_only: true`) to keep only companies with at least one email or phone. If MixRank returns no person-level channel, opt into `--places-fallback` (MCP `places_fallback: true`) to add a matched business listing phone and website in the same `contacts` array. Configure the separate Places credential with `mixrank auth google-places login`; Places does not return individual emails, and ambiguous matches remain in `places_lookup`.
