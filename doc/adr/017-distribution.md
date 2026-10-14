# ADR-017: Distribution Channels

## Status

Accepted

## Context

The cli's users are on Windows, macOS and Linux, and in containers. Each platform has a package manager its users expect; a download link alone isn't discoverable, and doesn't update.

## Decision

The cli is one statically linked Go binary per operating system, published to:

| Channel                                 | Built by                                                                                             | Users install with                    |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------- |
| GitHub Releases                         | `build-release-github.yml`: Windows, macOS and Linux binaries                                        | download from the release page        |
| Homebrew tap (`smarter-sh/tap`)         | `build-release-github.yml`: updates `Formula/smarter.rb`                                             | `brew install smarter`                |
| Ubuntu PPA (`lpm0073/smarter-cli`)      | `build-release-ubuntu-apt.yml`: a signed source package; `debian/rules` downloads the release binary | `apt-get install smarter-cli`         |
| Chocolatey                              | `build-release-windows-chocolatey.yml`, `chocolatey/smarter/`                                        | `choco install smarter`               |
| Docker Hub (`mcdaniel0073/smarter-cli`) | `build-release-dockerhub.yml`, `Dockerfile`                                                          | `docker run mcdaniel0073/smarter-cli` |

Assets are named `smarter-<os>-<version>` (`smarter-windows-<version>.exe`), and the package files that carry a version are rewritten by `scripts/bump_version.sh` (ADR-015). The Docker image runs as a non-root user.

## Alternatives Considered

- GoReleaser and its Homebrew, Scoop, deb and Docker publishers (ADR-010).
- `go install` only. Requires a Go toolchain.

## Consequences

- **Positive:**
  - Users install with the tool they already use.
- **Negative:**
  - The Homebrew formula and the Debian package download the binary from GitHub Releases, so the GitHub workflow must run first.
  - Only amd64 Linux is packaged for apt, and the macOS asset is built on whichever architecture `macos-latest` is.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-015: Semantic Versioning](015-versioning.md)
