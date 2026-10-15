# ADR-018: One Package per Verb

## Status

Accepted

## Context

Each verb (ADR-001) has its own route, its own output (`get` prints a table, `delete` prints `deleted.`), and its own flags (`get` has `-i`, `--asc`, `--desc`). The shared code (configuration, transport, output, the kind registry) is used by all of them.

## Decision

- The package `cmd` holds the root command, the shared code and the commands without a kind: `root.go` and `config.go` (configuration), `api.go` (transport), `console.go` (output), `resource.go`, `resourcekind.go` and `resourcekinds.go` (the registry), and `apply`, `validate`, `configure`, `resources`, `status`, `version`, `whoami`.
- Each resource verb is a package under `cmd/`: `get`, `describe`, `delete`, `deploy`, `undeploy`, `logs`, `manifest`, `jsonschema` (the `json-schema` verb), and `chat`, which prompts an LLMClient. Each has the same two files:
  - `<verb>.go`: the verb's `*cobra.Command`, an `APIRequest(kind, kwargs)` that builds its route and calls `cmd.APIRequest`, and its `ConsoleOutput`.
  - `resources.go`: its `legacySpecs`, and an `init()` that registers them and the specs generated from `cmd.ResourceKinds`.
- `main.go` imports each verb package for its side effects (`_ "github.com/smarter-sh/smarter-cli/cmd/get"`), so their `init()` functions attach the verbs to `cmd.RootCmd`.

## Alternatives Considered

- One package for everything. The verbs' `APIRequest` and `ConsoleOutput` would need distinct names, and their tests would share one viper.
- A registry of verbs called from `main`. More explicit than blank imports, but a larger change for little gain.

## Consequences

- **Positive:**
  - A verb's differences live in one small package.
- **Negative:**
  - A new verb package does nothing until `main.go` imports it, and its tests need their own `TestMain` and helpers (ADR-009).
  - `chat` doesn't fit the pattern: its commands take an LLMClient's name by flag (`--llmclient`, with `--chatbot` a deprecated alias) and call the platform's `prompt/` routes.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-009: Testing](009-testing.md)
