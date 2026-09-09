# From companies to actual contacts

`mixrank contacts` runs a shared SDK workflow: search current people at a company, rank relevant roles, request available business emails and direct dials, and return a single JSON report. The same workflow is the MCP tool `mixrank_company_contacts` and Go method `Client.CompanyContacts`.

```sh
mixrank contacts --domain example.com --name 'Example Company'
mixrank contacts --company-id 123 --domain example.com --max-contacts 2
mixrank contacts --companies companies.json --output contacts.json
mixrank contacts --companies - --roles 'owner,president,office manager' --emails-only
```

The company file may be a JSON array, an object containing `companies`, or an unmodified Elasticsearch company-search response. It understands `company_id`, arrays of `company_ids`, and semicolon-separated legacy IDs. IDs retain integer precision. A company needs an ID or domain; a name alone is ambiguous.

```json
{
  "companies": [
    {"name": "Example Company", "company_ids": ["123"], "domain": "example.com"}
  ]
}
```

The report has a row for every supplied company, including those with no result. `contacts` contains people with at least one returned email or direct dial. `candidates_without_contacts` records examined people whose requested channels were absent or withheld. `issues`, `complete`, `requests_made`, and `search_limited` preserve errors and coverage limits. A partial report is written even if a later call is cancelled or rejected for permissions.

Each contact includes the person's name, role, profile URL, matched current employment, profile update date, `business_emails`, `direct_dials`, channel status, and review reasons. The workflow requests `b2b_emails` and `directdials` by default; `--emails-only` omits phones. Consumer and educational email add-ons are never requested. B2B emails for unrelated employer domains are omitted. Email-domain aliases are not guessed.

Names, emails and numbers are provider-reported. Email `validation` is `not_run`; direct dials are person-level numbers whose line type and company association have not been independently verified. Inferred employer associations, missing or old profile dates, conflicting "former" headlines, and uncertain input companies are flagged for review. Availability is not a deliverability or identity guarantee. No emails are guessed or sent. If validation is requested, use a separate bulk validation job.

The default roles are owner, founder, CEO, president, general manager, operations manager, office manager, and marketing director. The search constrains employer, current status, no end date, and role on the same nested experience item. It checks that returned evidence again before enrichment. Source-linked records with fewer review concerns come first, then the requested role order; this is not a buying-intent score.

Each run accepts at most 25 companies. Defaults are two contacts per company, ten candidates per company, 100 upstream workflow calls, and a ten-minute CLI timeout. `--max-contacts`, `--max-candidates`, `--max-requests`, and `--timeout` set explicit bounds. Calls are sequential, use the existing semantic retry policy, and never silently switch to LiveScan fetch or resubmit uncertain enrichment. Search retries can make physical HTTP attempts exceed the logical workflow-call counter.

MCP accepts the same `companies` array and optional `roles`, `max_contacts_per_company`, `max_candidates_per_company`, `max_requests`, and `emails_only`. It uses the host's existing credential configuration. Both stdio and HTTP expose this tool; hosted authentication remains a separate deployment concern.

For agent calls, begin with one company at a time unless the client has a sufficiently long tool timeout. Larger CLI batches can take several minutes when upstream enrichment is slow. Client cancellation remains authoritative; do not assume a timed-out enrichment was never processed.

```go
report, err := client.CompanyContacts(ctx, mixrank.ContactOptions{
    Companies: []mixrank.CompanyTarget{{Domain: "example.com"}},
    MaxContacts: 2,
})
// Preserve report even when err is non-nil: earlier companies may have succeeded.
```

Keep contact exports outside public repositories. `--output` creates a private file and refuses to overwrite existing data before making requests.
