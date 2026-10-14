---
name: monthly-maintenance
description: Use when asked to do smarter-cli's monthly maintenance, or either half of it - (1) bring the cli back in sync with the Smarter platform by refreshing its pinned copy of the platform's OpenAPI document (make contracts), reading its diff and the platform's commits since the cli's last release, fixing what the contract test reports (new kinds, renamed or new /cli/ routes), checking which kinds' brokers deploy, and changing commands, tests and docs as needed, and (2) update the cli's Go modules, Go toolchain, npm release tooling and pre-commit hooks, retest, and, if the tests pass, write the git commits.
---

# smarter-cli Monthly Maintenance

The cli is a thin client of a platform that changes every week (ADR-001,
ADR-003). Its unit tests use a mock server, so nothing fails when the
platform adds a kind or renames a route: the cli just falls behind. Once a
month, bring it back in step, then update its dependencies. Two parts, in this
order, each with its own commit:

1. **Sync with the platform** (Part 1).
2. **Update dependencies and toolchain, and retest** (Part 2), so the retest
   covers the sync.

Load `smarter-cli-development` first. Its rules apply throughout: **never run
git write commands**, **never trigger a release**, **never point the cli at a
real platform except the local one**, and run the tests yourself.

Work from the smarter-cli repository root, with the platform at `../smarter`.

## Part 1: Sync with the platform

### 1. Find the baseline and the delta

The baseline is the date of the cli's last release tag:

```console
git -C ../smarter fetch --quiet
SINCE=$(git log -1 --format=%cI "$(git describe --tags --abbrev=0)")
git -C ../smarter log origin/main --since="$SINCE" --oneline -- \
  smarter/smarter/apps/api/v1/cli smarter/smarter/apps/api/v1/manifests \
  smarter/smarter/apps/api/v1/contract smarter/smarter/lib/manifest \
  docs/source/smarter-framework/smarter-cli.rst
```

Read the commit bodies of anything that touches the cli api, the kinds, or
the configuration file. Also skim the SDKs' commits that touch their config
readers (`../smarter-golang/config.go`, and the Python and Node.js
equivalents): a change there may mean the shared `config.yaml` contract moved
(ADR-002).

### 2. Refresh the pinned OpenAPI document

```console
make contracts SMARTER_REF=<the platform's latest stable release tag>   # records it in testdata/contracts/SOURCE
git diff --stat testdata/contracts/openapi.json
go test -run TestCommandsMatchTheOpenAPIDocument -v .
```

The diff of `testdata/contracts/openapi.json` is the list of the api's
changes; read the `/cli/` paths and the `SAMKind` enum. The contract test
(ADR-022) then says what to do:

| Failure                                            | Action                                                                                 |
| -------------------------------------------------- | -------------------------------------------------------------------------------------- |
| _the SAMKind X has no "…" command_                 | A new kind: add it to `ResourceKinds` (`smarter-cli-development`, "Adding a kind").    |
| _sent the kind "x", which is not a SAMKind_        | A removed or renamed kind: remove or rename its command.                               |
| _sent POST …, which is not an operation_           | **A broken command.** Fix the slug, with a unit test that asserts the documented path. |
| _publishes /cli/ operations that no command calls_ | A new route. List it for the user; build the command only when asked.                  |

### 3. Check what the test can't

- **Deployable kinds:** compare `ResourceKinds`' `Deployable` with the
  platform's `resources/` command (`smarter resources -o json` against the local
  platform), which reports each broker's `deployable` flag; the platform's
  `test_brokers.py` keeps that flag honest.
- **Help text and the docs page:** a renamed kind or a new verb needs
  `README.md` and `../smarter/docs/source/smarter-framework/smarter-cli.rst`
  updated.
- **Response shapes** that the output relies on (`get`'s `data.data.titles`
  and `items`, `validate`'s `data.valid` and `errors`, the journal's `error`
  object): try the commands you touched against the local platform.

### 4. Test, and commit

`make coverage`, the hooks (`smarter-code-quality`), and a smoke test against
the local platform of each command you changed (`smarter-api-testing`). Then
write a `feat(kinds): ...` commit for new kinds, or `fix(<verb>): ...` for
broken routes, or `test(contracts): ...` for a refresh that needed no change
(`smarter-git-commits`), naming the platform release you synced with.

## Part 2: Dependencies and toolchain

1. **Go:** if a new Go minor was released, update `go-version` in
   `.github/workflows/build-release.yml` and `default_language_version` in
   `.pre-commit-config.yaml`, the matrix in `.github/workflows/test.yml`, and
   the Dockerfile's `golang:` tag together (ADR-019). When the Go team drops a
   release, raise `go.mod`'s `go` line to the older supported one, and test
   with it (`GOTOOLCHAIN=go1.X.0 go test -short ./...`).
2. **Go modules:** `go get -u ./... && go mod tidy`. Read the changelogs of
   cobra and viper for behavior changes in flag binding, `IsSet`, or config
   writing, which the cli depends on.
3. **npm:** `npm outdated`, then `npm update` for the semantic-release
   packages. Dependabot may already have opened these PRs; check `git log
origin/alpha` before duplicating them.
4. **pre-commit:** `pre-commit autoupdate` (the weekly workflow may have done
   it already).
5. **Retest:** `make coverage`, `go vet ./...`, `golangci-lint run`, the hooks
   on every file (`pre-commit run --all-files`), and
   `go build -ldflags "-X main.Version=$(cat VERSION)"` then `version` with a
   throwaway `$HOME`.
6. If everything passes, write a `fix(deps): ...` commit (a `fix` so the
   updated binary is released). If a test fails, stop and report it with its
   output; don't commit around it.

## Report

Summarize: the platform release you synced with, what the contract test
reported before and after, what you changed, what you left for the user
(new routes, legacy kinds no longer served), the dependency versions that
moved, and the test results.
