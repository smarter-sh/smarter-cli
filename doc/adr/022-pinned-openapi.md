# ADR-022: A Pinned Copy of the Platform's OpenAPI Document

## Status

Accepted

Follows the Python SDK's [ADR-022](https://github.com/smarter-sh/smarter-python/blob/main/doc/adr/022-pinned-openapi.md). The Node.js and Go SDKs share this decision, under the same number.

## Context

The platform publishes its api as an OpenAPI document, generated from its Pydantic models and committed as `smarter/smarter/apps/api/v1/contract/openapi.json`. It covers every route of the cli api (`/api/v1/cli/`), and lists every manifest kind as its `SAMKind` enum.

The cli's unit tests answer requests with `httptest` servers (ADR-009), which answer whatever path they are asked. Until October 2026 the cli shipped `smarter manifest`, which called `example_manifest/`, and `smarter chat`, which called `chat/`, both renamed by the platform, plus 17 commands for kinds the platform no longer serves (ChatBot, the chat history kinds, and the old PluginData connection kinds), and `deploy prompt`, which the platform doesn't implement, with every unit test passing.

## Decision

The cli keeps a **copy** of the platform's `openapi.json`, pinned to a platform version, in `testdata/contracts/openapi.json`, and tests against it:

- `scripts/sync-contracts.sh` (`make contracts`, or `make contracts SMARTER_REF=<tag>`) copies it from a checkout of the platform, at its working tree or at a git ref. The copy is committed, and is never edited by hand.
- `TestCommandsMatchTheOpenAPIDocument` (`contract_test.go`) runs **every** command of the cli in-process against a mock platform that records each request, and fails when:
  - a request isn't an operation of the document (method and path template),
  - a `/cli/` operation of the document has no command,
  - a command sends a kind that isn't one of the document's `SAMKind`s,
  - a request sends a query parameter that its operation doesn't document,
  - a `SAMKind` has no `get`, `describe`, `delete`, `logs`, `manifest` or `json-schema` command, apart from a named list of exceptions, each with its reason.
- The copy is pinned to the platform's latest **stable release** (`make contracts SMARTER_REF=v0.23.0`), which the cli on `main` must work with; `testdata/contracts/SOURCE` records the version.
- The test runs in CI and in the `go-unit-tests` pre-commit hook. Refreshing the copy is part of the monthly maintenance; a failing contract test after a refresh is how the cli learns what changed.

The document can't say which kinds' brokers implement deploy, so `Deployable` (ADR-005) is checked against the platform's `resources/` command, by hand.

## Alternatives Considered

- Fetch the document from a running platform during the tests. Tests would need the network and a platform, and would change under the cli without a commit.
- Compare the cli with a checkout of the platform with a script. That is what the cli did first; nobody runs a script that nothing requires.
- Live tests against the local platform in CI. The cli's CI has no platform.

## Consequences

- **Positive:**
  - A renamed route, a removed kind or a new kind fails a unit test, in CI, rather than reaching users.
  - The diff of the copy, at each refresh, is a precise list of the api's changes.
  - The four Smarter clients (three SDKs and the cli) check themselves against the same document.
- **Negative:**
  - The copy is only as current as its last refresh.
  - The test must know the extra arguments some commands need (`apply -f`, `chat prompt --prompt`); a new command of that sort must be added to `commandArgs`.
  - The document doesn't say which kinds deploy; the platform's `resources/` command does, from its brokers' `deployable` flags (until October 2026 it was a hard-coded table, and wrong).

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-009: Testing](009-testing.md)
