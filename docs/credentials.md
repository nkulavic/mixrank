# Credential ownership

The toolkit never takes an API key as a command argument. `mixrank auth login` masks input; `--stdin` permits an explicit pipe. A nonempty `MIXRANK_API_KEY` is the runtime override. Status commands expose source and availability only.

| Client | Secret owner and injection |
|---|---|
| Codex CLI/Desktop, native CLI, generic local MCP | MixRank-owned OS-vault service `com.nkulavic.mixrank`, account `default`; environment override first |
| Claude Code plugin | Sensitive manifest `userConfig.api_key`, injected as `${user_config.api_key}` into `MIXRANK_API_KEY` |
| Claude Desktop MCPB | Sensitive `user_config.api_key`; Desktop manages encrypted storage and runtime injection |
| Cowork | Its own connector authentication, inside the Cowork runtime; no copied host credentials |
| Hosted MCP | Provider-managed upstream secret; client-owned remote access tokens are separate |

Claude Code documents macOS Keychain storage with a credentials-file fallback when unsupported/rejected. That fallback belongs to Claude; this toolkit never edits Claude's credential files. See [Claude user configuration](https://code.claude.com/docs/en/plugins-reference#user-configuration). Desktop handles its sensitive MCPB fields according to the [bundle manifest contract](https://github.com/modelcontextprotocol/mcpb/blob/main/MANIFEST.md).

For its own credential, MixRank uses macOS Keychain, Windows Credential Manager or Linux Secret Service. It fails if the store is missing or locked and never silently writes a plaintext substitute. Linux headless sessions need a working Secret Service session or explicit environment injection. macOS writes use the `security` interpreter's stdin to keep the key out of process arguments; its diagnostics are suppressed and the stored result is verified. The key is not embedded in plugins, configuration, binaries or docs.

Codex MCP OAuth storage controls remote login tokens; it is not an arbitrary upstream-key vault. See [Codex MCP](https://developers.openai.com/codex/mcp) and [plugin authentication](https://developers.openai.com/plugins/build/auth). Local plugins use MixRank's vault rather than inventing an undocumented Codex sensitive-configuration field.

MixRank places the API key in the URL path. The SDK therefore disables redirects for authenticated API calls, removes the key from its error diagnostics, and does not forward it to result download hosts. Avoid logging requests through custom HTTP transports outside the SDK's redaction boundary. Provider callback signatures use the API key as the raw HMAC secret; verify raw bytes, timestamp and event ID before processing, then deduplicate IDs.

`auth logout` deletes only the MixRank-owned entry. An injected environment variable still overrides it. Uninstall retains credentials and profiles unless explicitly selected for removal. Host-native plugin enable/import prompts remain host UI actions; installers do not bypass them or write host credential stores directly.
