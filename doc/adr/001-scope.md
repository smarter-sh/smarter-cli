# ADR-001: Scope and Purpose of the CLI

## Status

Accepted

## Context

Smarter resources (LLMClients, Plugins, Guardrails, Secrets, ...) are declared in YAML manifests, the Smarter Application Manifest (SAM), as Kubernetes resources are. People who work at a terminal, and scripts and CI pipelines, need to apply, read, delete and deploy those resources without the web console, the way `kubectl` manages Kubernetes. The platform's [APIs and SDKs](https://docs.smarter.sh/en/latest/smarter-framework/apis-and-sdks.html) page lists the cli beside the SDKs.

## Decision

smarter-cli is the platform's official command-line interface: one static binary named `smarter`, for Windows, macOS and Linux, and a Docker image.

- It uses `kubectl`'s **verb-noun** grammar: `smarter <verb> <kind> [name] [flags]`, with the verbs `apply`, `validate`, `get`, `describe`, `delete`, `deploy`, `undeploy`, `logs`, `manifest`, `json-schema`, `chat` (prompt an LLMClient), and the commands `configure`, `whoami`, `resources`, `status`, `version`.
- It is a **thin client**. Every command is one request to the platform's cli api (ADR-003), whose brokers do the work. The cli holds no business rules about a kind, and adding a kind on the platform needs only a registry entry here (ADR-005).
- It serves **people and scripts**: output is json, yaml or a table (ADR-013), errors go to stderr with exit status 1 (ADR-006).
- It is **not a library**. Go programs that want the same api import the [Go SDK](https://github.com/smarter-sh/smarter-golang), which reads the cli's configuration file (ADR-002) and was written so that the cli could adopt it later.

## Alternatives Considered

- Only the web console. Not scriptable, and no place in CI.
- Only the SDKs. Every user would write a program for what one command does.
- Put the logic of each kind in the cli. It would drift from the platform with every release.

## Consequences

- **Positive:**
  - A new kind on the platform reaches the cli with one line (ADR-005).
  - The binary is small and has no runtime dependencies.
- **Negative:**
  - The cli depends on the platform's cli api; when a route changes, the cli breaks until it follows (see the `monthly-maintenance` skill).
  - Its http code is in the `cmd` package, tied to cobra and viper, so other programs can't reuse it; that is the Go SDK's job.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-002: Configuration](002-configuration.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
