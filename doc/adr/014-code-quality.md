# ADR-014: Code Quality

## Status

Accepted

Follows the platform's [ADR-028: Code Quality](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/028-code-quality.rst). The Go SDK's [ADR-014](https://github.com/smarter-sh/smarter-golang/blob/main/doc/adr/014-code-quality.md) adopts this repository's Go hooks.

## Context

Consistent code is easier to review, and a solo-maintained project needs its checks to run without anyone remembering them.

## Decision

`.pre-commit-config.yaml` is the standard:

- Go: `go-fmt`, `go-vet`, `go-imports`, `golangci-lint` (its default linters; there is no `.golangci.yml`), and `go-unit-tests`, from `dnephin/pre-commit-golang`.
- Generic: large files, case conflicts, JSON/TOML/XML/YAML syntax, merge conflicts, symlinks, AWS credentials, private keys, end of file, trailing whitespace.
- `commitlint` on commit messages, with `@commitlint/config-angular` (ADR-015).

Each Go file starts with the project's copyright comment. Errors are wrapped with `%w`. A linter is silenced only with a narrow `//nolint:<linter> // reason`.

## Alternatives Considered

- `go vet` alone. golangci-lint catches unchecked errors and dead code that vet doesn't.

## Consequences

- **Positive:**
  - The SDKs and the cli share a Go style.
- **Negative:**
  - pre-commit.ci skips the Go hooks and commitlint; the `test.yml` workflow runs gofmt, `go vet`, golangci-lint and the tests instead (ADR-010), but not commitlint.
  - Unlike the Go SDK, there is no codespell or Prettier hook, and no `.golangci.yml` enabling gosec, errorlint or revive.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-015: Semantic Versioning](015-versioning.md)
