# MixRank toolkit

A Go SDK, native CLI, MCP server and portable agent skills for MixRank. Research companies and people, resolve identities, enrich records, validate email lists, inspect technology adoption and apply private product context to reusable workflows.

**102 catalog operations**: 97 distinct method/path pairs from the authorized account's API documentation and five Elasticsearch searches (`person2`, `companies`, `jobs`, `org_name`, `industries`), reviewed September 9, 2026. This is coverage of the documentation visible to that account, not a claim about undisclosed APIs or permission to call every operation.

## Install

Download `install.sh` (macOS/Linux) or `install.ps1` (Windows) from [releases](https://github.com/nkulavic/mixrank/releases), review it, and run it. The bootstrap verifies the pinned binary's SHA-256 checksum and starts the component wizard. Go is not required.

```sh
sh install.sh
# Or a reviewable, non-interactive installation:
sh install.sh --components cli,skills,codex-mcp --dry-run --yes
sh install.sh --components cli,skills,codex-mcp --yes
```

Add the reported binary directory to PATH. For a full Codex or Claude plugin, select `codex-plugin` or `claude-plugin` instead of that client's standalone MCP component. Codex CLI and Desktop share their user registration. The wizard builds a local plugin with a native binary for each OS. The source GitHub marketplaces use a POSIX launcher; Windows users can use native wizard setup or platform bundles. See [client installation](docs/installation.md).

```sh
mixrank auth login                  # masked input; MixRank-owned OS vault
mixrank auth status                 # no key displayed
mixrank doctor                     # local checks only
mixrank doctor --live               # one authenticated /echo request
mixrank endpoints get-person-match
mixrank api get-companies-match --url example.com
mixrank api search-companies --body query.json
```

[Operation reference](skills/mixrank-cli/references/operations.md) · [CLI recipes](skills/mixrank-cli/references/recipes.md) · [credentials](docs/credentials.md) · [validation status](docs/validation.md)

For follow-up, turn the company results into actual contact records:

```sh
mixrank contacts --companies companies.json --output contacts.json
mixrank contacts --domain example.com --max-contacts 2
```

This finds current relevant people and retrieves available business emails and direct dials, with source dates, missing-data statuses and review flags. The SDK and MCP expose the same workflow through `CompanyContacts` and `mixrank_company_contacts`. [Contact workflow](docs/contacts.md).

## SDK and MCP

```go
client, err := mixrank.New(os.Getenv("MIXRANK_API_KEY"), mixrank.Options{Retries: 2})
if err != nil { return err }
response, err := client.GetCompaniesMatch(ctx, mixrank.Request{
    Parameters: map[string]string{"url": "example.com"},
})
if err != nil { return err }
defer response.Body.Close()
value, err := response.JSON(mixrank.DefaultMaxResponseBytes)
```

Every catalog operation has a generated Go method, a named `mixrank api` command and an MCP tool. Raw HTTP/JSON access remains available. Close response bodies to release connections and validation locks. JSON uses `json.Number`; callers should avoid converting identifiers to floating point. Downloads stream to files rather than loading entire results into memory.

```sh
mixrank mcp                         # stdio
mixrank mcp --transport http --listen 127.0.0.1:8080
```

HTTP requires server-side `MIXRANK_API_KEY` and a separate `MIXRANK_MCP_TOKEN` of at least 32 characters, injected through the environment. It exposes `/mcp`, checks exact allowed origins and uses stateless Streamable HTTP. This bearer mode is for private development. The reusable handler accepts an authorization callback for a later OAuth deployment; it does not claim to implement production OAuth. [Hosting sequence](docs/hosting.md).

## Research and private personalization

Four maintained skills cover CLI use, dataset semantics, generic research and product personalization. They produce evidence-backed account lists, landscapes, contact maps, briefs, qualification evidence, outreach drafts and exports. They do not send messages or run campaigns.

```sh
mixrank personalize --name product --website https://example.com --repo owner/repo --path ./app
mixrank personalize --name another --website https://example.com --agent codex
mixrank profiles show product --draft
mixrank profiles diff product
mixrank profiles accept product      # after review
mixrank profiles use product
mixrank profiles refresh product
```

Without `--agent`, personalization writes a private task for your installed coding agent to analyze. It does not pretend that a source URL alone is an analyzed profile. The agent fills attributed evidence and targeting context; acceptance validates and activates a version and companion skill. No separate LLM subscription is needed. Profiles default to the OS user configuration directory, outside repositories. Product apps are not modified. [Personalization](docs/personalization.md).

## Operational behavior

- Lists use bulk validation. Single-address validation is serialized per credential across toolkit processes on one host, including raw requests. Asynchronous one-offs retain a pending gate until the matching signed completion webhook is acknowledged. Separate hosts must coordinate. HTTP instances use process-local coordination and need a shared coordinator before scaled single-address workloads.
- LiveScan and posts default to cached. `--strategy fetch --refresh` explicitly requests new data. The toolkit does not silently change strategy or replay uncertain enrichment/submission calls. Provider 512/513 fetch failures are preserved.
- Job IDs, status, blank download URLs, 202 acceptance and 204 pending responses remain distinct. Private email bulk output requires explicit deletion after download.
- Person/company match candidates and confidence are preserved. Contact append does not imply validated deliverability. Sparse history and absent/disabled fields remain unknown.
- CLI exit codes: 0 success, 2 input/config, 3 credential/permission, 4 rate limit, 5 upstream failure, 6 uncertain outcome, 130 cancelled safe request.

## Development

Go 1.26 or later and Python 3 are needed for development. Dependencies are pinned in `go.mod`/`go.sum`; releases use the official Go MCP SDK.

```sh
go test -race ./...
go vet ./...
python3 scripts/generate.py --check
python3 scripts/package.py --check
python3 scripts/package.py --release
```

CI uses synthetic data only. Live checks are separate and opt-in. Raw signed-in documentation, credentials, private source snapshots, profiles and customer outputs are excluded from public artifacts. API facts and original guidance are maintained in the catalog and skills; the vendor's raw documentation is not redistributed. MIT applies to original toolkit code, not MixRank data or service rights. This is an independent toolkit, not an official MixRank product.
