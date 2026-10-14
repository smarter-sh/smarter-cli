# ADR-010: CI/CD

## Status

Accepted

The SDKs' ADR-010 adapt the platform's [ADR-014: CI/CD](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/014-ci-cd.rst).

## Context

A release must be versioned from its commits, built for three operating systems, and published to five channels (ADR-017), by one maintainer.

## Decision

GitHub Actions in `.github/workflows/`:

- `test.yml`: on pull requests and pushes to `main`, `alpha` and `beta`, checks formatting, runs `go vet` and golangci-lint, and runs the tests with `make coverage`, on Go 1.26 and 1.27 (ADR-019).
- `pushMain.yml`: on a push to `main`, semantic-release (ADR-015) computes the version from the commit messages, updates `CHANGELOG.md`, `VERSION` and the files that carry the version, tags, and creates the GitHub release; it then pushes `main` to a `v<major>` branch.
- `build-release-github.yml`, `build-release-dockerhub.yml`, `build-release-ubuntu-apt.yml`, `build-release-windows-chocolatey.yml`: build and publish each channel. They run **by hand** (`workflow_dispatch`), after the version is tagged. Each workflow's header lists the secrets it needs (`PAT`, `DOCKERHUB_TOKEN`, `GPG_PRIVATE_KEY`, ...) and how to rotate them.
- `precommitVersionBumps.yml`: weekly `pre-commit autoupdate`, committed to `alpha`.
- Dependabot opens monthly updates of Go modules, npm packages and actions; Mergify merges them into `alpha` (ADR-011).
- Branches follow the platform: work lands on `alpha`, and `main` is released.

## Alternatives Considered

- GoReleaser. One config for every binary, checksum and Homebrew formula; worth adopting, but the channels already work.
- Publishing automatically on every release. The channels have manual review steps (Chocolatey moderation, Launchpad builds).

## Consequences

- **Positive:**
  - Versions and the changelog follow from commit messages, with no manual step.
- **Negative:**
  - Releases need a maintainer to run four workflows in order.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-011: Dependency Management](011-dependencies.md)
- [ADR-015: Semantic Versioning](015-versioning.md)
- [ADR-017: Distribution Channels](017-distribution.md)
