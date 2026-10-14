---
name: smarter-cli-release
description: Use when a smarter-cli task touches versions or releases - semantic-release and release.config.js, VERSION, scripts/bump_version.sh, the build-release workflows (GitHub Releases, Homebrew tap, Ubuntu PPA, Chocolatey, Docker Hub), the Dockerfile, debian/ or chocolatey/ packaging, workflow secrets, or "why wasn't my change released". Explains how a commit becomes a release on five channels, what an agent may change and verify locally, and what it must never run.
---

# Smarter CLI Releases

Read `doc/adr/010-ci-cd.md`, `015-versioning.md` and `017-distribution.md`
first.

## How a commit becomes a release

1. Work lands on `alpha`. Dependabot and Mergify merge dependency updates
   there too.
2. A merge to `main` runs `.github/workflows/pushMain.yml`: semantic-release
   reads the commit messages since the last tag. `feat` makes a minor release,
   `fix` a patch, `BREAKING CHANGE:` a major; `docs`, `chore`, `test`, `ci`,
   `refactor`, `style` make none.
3. semantic-release runs `scripts/bump_version.sh <version>`, which writes
   `VERSION`, `VERIFICATION.txt`, the `Makefile`'s `choco-push` line and the
   download url in `debian/rules`; commits them with `CHANGELOG.md` as
   `chore(release): <version> [skip ci]`; tags `v<version>`; creates the
   GitHub release; and pushes `main` to `v<major>`. `alpha` and `beta` make
   prereleases.
4. A maintainer then runs the build-release workflows **by hand**, in this
   order, because the later ones download the GitHub release's binaries:
   - `Build-Release GitHub` (`build-release-github.yml`): Windows, macOS and
     Linux binaries, `smarter-<os>-<version>`, uploaded to the release; then the
     Homebrew formula in `smarter-sh/homebrew-tap` is updated with the macOS
     binary's url and sha256.
   - `Build-Release Ubuntu APT`: a signed source package uploaded to the
     Launchpad PPA, whose `debian/rules` downloads the Linux binary.
   - `Build-Release Chocolatey` and `Build-Release DockerHub`.

Each workflow's header comment lists the secrets it needs (`PAT`,
`DOCKERHUB_TOKEN`, `GPG_PRIVATE_KEY`, `GPG_PASSPHRASE`, ...) and how to
rotate them. Never print, echo or move a secret's value.

## What an agent never runs

- `make release` (an empty `fix:` commit, pushed), `make choco-pack`,
  `make choco-push`, `gh workflow run`, `gh release`, `dput`, `docker push`,
  `npx semantic-release`.
- `scripts/bump_version.sh`: it rewrites four tracked files in place, and the
  version is semantic-release's to choose.
- Editing `VERSION`, `CHANGELOG.md` or `debian/changelog` by hand.

Give the user the command instead, when they ask how to release.

## What you can do and verify locally

- **Why wasn't it released?** Check the commit types since the last tag:
  `git log --format='%s' $(git describe --tags --abbrev=0)..origin/main`. Only
  `feat`, `fix` and breaking changes release.
- **A versioned build**, as the workflows make it:
  `go build -ldflags "-X main.Version=$(cat VERSION)" -o "$SCRATCHPAD/smarter" . && "$SCRATCHPAD/smarter" version`
  (with a throwaway `$HOME`; see `smarter-cli-testing`). Without `-ldflags` it
  reports `local.dev`.
- **Cross-compiling:** `GOOS=windows GOARCH=amd64 go build -o "$SCRATCHPAD/smarter.exe" .`
  and likewise `darwin/arm64`, `linux/amd64`.
- **The Docker image:** `docker build -t smarter-cli-test .` then
  `docker run --rm smarter-cli-test version`. Remove the image afterwards.
  It builds from `golang:latest`, so it may use a newer Go than the release.
- **Workflow edits:** validate the YAML with the hooks
  (`smarter-code-quality`); a workflow can only really be tested by running
  it, which is the user's call.
- **`bump_version.sh` edits:** test them on copies in your scratchpad, never on
  the tracked files. It edits with `perl -pi`, because semantic-release runs it
  on the Ubuntu runner and maintainers on macOS, and `sed -i` differs between
  them: the BSD form silently broke releases 0.3 and 0.4 on CI.

## When you change packaging

- A new asset name, or a new file that carries the version, must be added to
  `scripts/bump_version.sh` **and** to `assets` in `release.config.js`, or
  semantic-release won't commit it.
- The Homebrew job seds `url`, `sha256` and `bin.install` in
  `Formula/smarter.rb`; renaming the macOS asset breaks it.
- The Go version in the workflows and in `.pre-commit-config.yaml`'s
  `default_language_version` move together (ADR-019).
- Write the change as `build(...)` or `ci(...)`; it won't make a release by
  itself.
