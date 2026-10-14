# ADR-015: Semantic Versioning

## Status

Accepted

Follows the platform's [ADR-020: Semantic Versioning](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/020-semantic-versioning.rst).

## Context

Users install the cli from five channels (ADR-017), each of which needs the version in its own file. Versions must be computed, not chosen by hand.

## Decision

- [semantic-release](https://semantic-release.gitbook.io/) (`release.config.js`) sets the version from [Conventional Commits](https://www.conventionalcommits.org/): `fix` is a patch, `feat` a minor, `BREAKING CHANGE` a major. commitlint enforces the format (ADR-014).
- `main`, `next` and `next-major` release; `alpha` and `beta` release prereleases.
- `scripts/bump_version.sh <version>` (with `perl -pi`, which behaves the same on the Ubuntu runner and on macOS) writes the version into `VERSION`, `VERIFICATION.txt`, the `Makefile`'s `choco-push` target and `debian/rules`; semantic-release commits them with `CHANGELOG.md` as `chore(release): <version> [skip ci]`.
- The binary learns its version at build time: `go build -ldflags "-X main.Version=$(cat VERSION)"`. A plain `go build` reports `local.dev`.
- The cli's version is independent of the platform's. It speaks the platform's `/api/v1/` (ADR-003).

## Alternatives Considered

- Matching the platform's version. The two release on different schedules.

## Consequences

- **Positive:**
  - Every channel is stamped with the same version.
- **Negative:**
  - Until October 2026 the script used BSD `sed -i ''`, which fails on the Ubuntu runner, so releases 0.3 and 0.4 left `VERIFICATION.txt` and the `choco-push` line at 0.2.1. The next release corrects them.
  - `smarter version` doesn't say which platform versions it supports.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-014: Code Quality](014-code-quality.md)
- [ADR-017: Distribution Channels](017-distribution.md)
