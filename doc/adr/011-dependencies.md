# ADR-011: Dependency Management

## Status

Accepted

The SDKs' ADR-011 adapt the platform's [ADR-015: Dependency Management](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/015-dependency-management.rst).

## Context

Every dependency is compiled into a binary that users run with their api keys. The cli's job is small.

## Decision

- **Go modules**, pinned in `go.mod` and `go.sum`. Direct dependencies are few: cobra and viper (ADR-004), `ghodss/yaml` for json-to-yaml output (ADR-013), and `godotenv` for local `.env` files. A new one needs a reason the standard library can't meet.
- **npm** (`package.json`) holds only release tooling: semantic-release and its plugins. No npm code ships.
- **Python** (`venv/`, local only) holds pre-commit.
- Dependabot proposes monthly updates of Go modules, npm packages and GitHub Actions; Mergify merges them into `alpha`. `precommitVersionBumps.yml` updates the hook versions weekly.
- After an update: `go mod tidy`, `go test ./...`, and the hooks.

## Alternatives Considered

- Vendoring (`go mod vendor`). The module proxy and `go.sum` already make builds reproducible.

## Consequences

- **Positive:**
  - A small, auditable dependency tree.
- **Negative:**
  - viper pulls in many indirect modules (toml, afero, fsnotify, ...) for one configuration file.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-019: Go Version](019-go-version.md)
