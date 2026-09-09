# HTTP and deployment sequence

The Go MCP SDK serves stateless Streamable HTTP at `/mcp`. The exported `mcpserver.Handler` can be embedded in a Go HTTP function or service. It requires an auth callback or a distinct bearer token, validates explicit origins, enforces request/output limits, propagates cancellation, and sets `Cache-Control: no-store`. Upstream credentials are resolved from server environment variables for CLI HTTP mode; no local credential vault or persistent disk is required. Process-local one-off validation coordination does not cover multiple instances.

If the contact tool enables `places_fallback`, configure `GOOGLE_PLACES_API_KEY` as a separate provider-managed secret alongside `MIXRANK_API_KEY`. The server performs one bounded, non-retried Places Text Search per fallback company and never persists the returned Places content.

A private bearer configuration is available for local integration testing. This is not OAuth and must not be advertised as a production Claude/Cowork connector. Production activation requires OAuth 2.1-compatible client flows, protected-resource metadata, issuer/audience/expiry/scope validation, appropriate registration/PKCE behavior, and tests in the actual target clients. The exported Authorize callback is the integration boundary, not an implemented authorization server.

Build and validation precede deployment. The user's paid Vercel account is the first candidate; inspect current Go runtime support, function duration, streaming behavior, size limits and account settings before selecting a deployment. Vercel configuration and provider-managed secrets must be created only in that later phase. Long operations should submit a job and later inspect/download it instead of holding a function open. Receiver hosting and webhook persistence belong to deployment or the consuming application.

A scaled service must coordinate all one-off email validations using the same upstream credential across instances or disable that operation in favor of bulk jobs. Pending callback state also needs a shared store for asynchronous one-offs. In-process stateless behavior is not a distributed lock. Hosted secrets are injected by the provider, never uploaded from local keychains. Remote clients own their login tokens independently of the upstream MixRank API key.

References: [Vercel Go runtime](https://vercel.com/docs/functions/runtimes/go), [MCP authorization](https://modelcontextprotocol.io/specification/latest/basic/authorization), [OpenAI plugin authentication](https://developers.openai.com/plugins/build/auth).

## Vercel readiness review, September 9, 2026

Vercel documents a beta Go runtime on all plans. Its Go framework preset detects a root `go.mod` and `main.go`, `cmd/api/main.go`, or `cmd/server/main.go`; a server entrypoint must listen on `PORT`. The existing exported MCP handler can be used by that later entrypoint. [Go runtime documentation](https://vercel.com/docs/functions/runtimes/go).

The general function limits list a 4.5 MB request/response payload limit and 250 MB uncompressed bundle limit. The toolkit's HTTP input bound is 4 MiB and individual tool output bound is 1 MiB; large uploads and downloads should stay outside the MCP function. The published 800-second Pro duration table specifically describes Node.js, Bun and Python, so it is not evidence of the Go runtime's exact duration or streaming behavior. Verify those on the selected project before activation. [Function limits](https://vercel.com/docs/functions/limitations).

Vercel CLI 59.13.1 reported logged out on the implementation machine. The user has stated that their account is paid, but account settings and deployment behavior were not verified. No Vercel project, deployment, server secret or connector was created. The next phase needs account authentication, a selected OAuth provider and client login tests before production activation.
