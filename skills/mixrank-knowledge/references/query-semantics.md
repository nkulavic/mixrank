# Query construction and interpretation

The five search endpoints are `POST /elasticsearch/{index}/_search` for `person2`, `companies`, `jobs`, `org_name`, and `industries`. Inspect the mapping before selecting text versus keyword fields, nested paths, sorting or aggregations. Do not invent an HTTP `_mapping` endpoint; the CLI exposes bundled mappings locally.

`person2.person_id` is the person identifier. `job.company_id` and `experience.company_id` refer to employers; LinkedIn industry IDs join to the industries index. The company search index and profile APIs expose related IDs but their property names differ: inspect the mapping and returned records before joining. `org_name` provides name resolution, not a guarantee that similarly named organizations are identical.

The provider's text analyzer tokenizes and ignores case. It does not supply title synonym expansion: CEO and Chief Executive Officer need explicit alternatives. A `match` query defaults to OR between tokens; use `operator: "and"` where all tokens matter. Use exact fields only when the mapping documents them. Prefix/autocomplete queries must use mapped fields and their analyzer behavior.

Example: all experience constraints belong inside one nested clause, preventing a role at one employer from matching another employer on the same person:

```json
{
  "size": 25,
  "_source": ["person_id", "name", "experience", "updated_at"],
  "query": {
    "nested": {
      "path": "experience",
      "query": {
        "bool": {
          "filter": [
            {"term": {"experience.is_current": true}},
            {"term": {"experience.company_id": 123}}
          ],
          "must": [{"match": {"experience.title": {"query": "engineering director", "operator": "and"}}}]
        }
      }
    }
  }
}
```

The employer ID is illustrative. Start from matched company IDs for actual work. Add title alternatives only when justified by buyer-role criteria. `has_source_company_id` distinguishes directly linked from inferred employer associations where available.

For pagination, set explicit size and a deterministic sort supported by the mapping. Pass returned sort values to `search_after` only when the provider's deployed Elasticsearch supports it. Do not assume an unlimited result window, PIT support, or server version. Use bounded pages and record truncation. Count aggregations and result lists are different deliverables; preserve aggregation scope, missing buckets and lower-bound count relations.

Aggregation results over nested employment count experience items unless reverse nesting or person-level deduplication changes the unit. A person can have several concurrent roles. State the counting unit and avoid presenting role counts as unique people.
