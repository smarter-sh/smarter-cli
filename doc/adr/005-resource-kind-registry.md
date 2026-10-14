# ADR-005: The Resource Kind Registry

## Status

Accepted

## Context

The platform adds resource kinds every few releases (v0.14 added a dozen, v0.16 Budget, LLMHostCompute and WebsearchPlugin). Each kind needs a leaf command under `get`, `describe` and `manifest`, and deployable kinds under `deploy` and `undeploy`. Before v0.14, every verb package held its own copy of each kind's command name, api kind and help text, and the copies drifted.

## Decision

Kinds are data, declared once.

- `cmd.ResourceKinds` (`cmd/resourcekinds.go`) lists every kind that follows the standard shape: `Singular` and `Plural` command names, the `APIKind` sent to the platform, `Display` and `DisplayPlural` for help text, and `Deployable`.
- `ResourceKind` (`cmd/resourcekind.go`) generates each verb's `ResourceSpec` from an entry: `GetSpec`, `DescribeSpec`, `DeleteSpec`, `LogsSpec`, `ManifestSpec`, `JSONSchemaSpec`, `DeploySpec`, `UndeploySpec`.
- `Deployable` is true only for kinds whose platform broker implements deploy, the broker's `deployable` flag, which the platform's `resources/` command reports: SmarterAuthToken, CustomDomain, LLMClient, LLMHost, LLMHostCompute and Vectorstore, and User, which is a legacy spec here, as of October 2026. CustomDomain's and LLMHostCompute's brokers don't undeploy, so their `undeploy` commands return the platform's "not implemented" error. Until October 2026 the platform's table was hard-coded, and said only Prompt, whose broker doesn't deploy.
- `ResourceSpec` (`cmd/resource.go`) declares a leaf command: `Use`, `Short`, `Long`, `APIKind`, how it takes the resource's identifier (`NameArg`: none, positional, or a flag), and extra `Flags`. `RegisterResources` turns specs into cobra commands bound to viper.
- Kinds that don't fit the shape (a custom flag, a `username` identifier, or older help text worth keeping) are explicit `legacySpecs` in the verb's `resources.go`: `account`, `plugin` and `user`. The ChatBot and chat history commands were removed in October 2026: the platform has had no broker for them since v0.14.
- `ResourceKinds` and the legacy specs together cover the platform's `SAMKinds`, which its OpenAPI document lists as the `SAMKind` enum. The contract test fails when a `SAMKind` lacks a command, or a command sends a kind that isn't one (ADR-022).

Adding a kind is one line in `ResourceKinds`, plus its test expectations.

## Alternatives Considered

- Fetch the kinds from the platform's `resources/` route at startup. Help and completion would need the network, and a key.
- Generate the registry from the platform's schemas. Possible later; the list is short, and the display names need judgment.

## Consequences

- **Positive:**
  - A new kind gets consistent commands and help under every verb.
  - `TestResourceKindsAreWellFormed` catches empty fields and duplicate names.
- **Negative:**
  - The registry is updated by hand, after the contract test reports a new kind.
  - `Deployable` can drift from the platform's brokers unnoticed; the monthly maintenance checks it.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-018: One Package per Verb](018-package-per-verb.md)
- [ADR-022: A Pinned Copy of the Platform's OpenAPI Document](022-pinned-openapi.md)
