# Synthetic workflow evaluation

This is a manual application of the maintained research and knowledge skills to [synthetic fixtures](../tests/workflow-cases.json). It is separate from executable SDK tests and does not claim an automated model benchmark. No real people, prospect calls, messages or validation jobs were used.

| Scenario | Applied decision | Outcome |
|---|---|---|
| Ambiguous healthtech vertical | Distinguish care provision from software vendors before selecting the population. Use the hypothetical profile's stated exclusion. | Account 202 is the one resolved software account; its repeated row is deduplicated by ID/domain. Account 201 remains outside the qualified set pending evidence of a software business. Buying intent stays unknown. |
| Stale and crossed employment | Require current employer and relevant title on the same nested experience item. | Person 302 is a current engineering contact. Person 301 is historical at the target. Person 303's engineering title belongs to another company and does not qualify. An empty email add-on does not establish missing contacts or deliverability. |
| Conflicting identifiers and partial enrichment | Preserve tied match candidates and input row IDs. Do not collapse contradictory identity evidence. | Row r1 remains unresolved, r2 resolved. The accepted validation job is pending; no resubmission or completed-output claim is warranted. A missing historical headcount is left null. |

These cases exercise the intended workflow judgments: vertical disambiguation, entity deduplication, nested employment, disabled versus missing contact data, job uncertainty and sparse history. Automated tests cover the HTTP/client mechanisms that preserve the underlying candidates, nulls, numbers, pending statuses and partial results. Agent behavior should be evaluated again as skills evolve.
