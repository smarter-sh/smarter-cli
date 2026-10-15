# ADR-019: Go Version

## Status

Accepted

The SDKs' ADR-019 decide their supported language versions; the Go SDK's [ADR-019](https://github.com/smarter-sh/smarter-golang/blob/main/doc/adr/019-go-versions.md) notes that the cli builds with the newest Go.

## Context

A cli is an application, not a library: users run its binary and never compile it, so it may use the newest Go. The Go team supports the two newest releases, and fixes security bugs only in them; a binary built with an older Go carries its bugs.

## Decision

- Releases are built with the newest stable Go: `go-version: '1.27'` in the release workflows, and `default_language_version: golang: '1.27'` in `.pre-commit-config.yaml`. These move together, when a new Go is released.
- The `go` line of `go.mod` is the oldest Go that can build the cli from source: the older of the two releases the Go team supports. CI tests both.
- The Dockerfile builds with the same pinned Go (`golang:1.27`).

## Alternatives Considered

- Build with the oldest supported Go, as a library would. Users would get fewer security fixes, for no benefit.

## Consequences

- **Positive:**
  - Every release has the newest Go's security fixes.
- **Negative:**
  - The `go` line of `go.mod` is `go 1.26.0`, the oldest Go the Go team supports in October 2026, and CI tests 1.26 and 1.27 (ADR-010). It must be raised when 1.26 goes out of support.
  - The Dockerfile's `golang:1.27` image must move with the release workflows.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-011: Dependency Management](011-dependencies.md)
