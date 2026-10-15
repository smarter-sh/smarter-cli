---
name: smarter-cli-development
description: Start here for any change to the smarter-cli repository, the Smarter platform's command-line interface (the `smarter` binary, Go, cobra and viper). Gives the non-negotiable rules, the cli's architecture (cmd.APIRequest and the /api/v1/cli/ routes, the ResourceKinds registry and ResourceSpec, one package per verb, the ~/.smarter/config.yaml contract shared with the SDKs, the output formats), how to add a resource kind, a verb or a flag, the pinned OpenAPI document and the contract test that keep the cli in step with the platform, and an index of the narrower skills to load for each kind of work.
---

# Smarter CLI Development

`smarter` is to the Smarter platform what `kubectl` is to Kubernetes: a
verb-noun command line (`smarter <verb> <kind> [name] [flags]`) for applying,
reading, deleting and deploying the resources that SAM manifests declare. The
platform is [smarter-sh/smarter](https://github.com/smarter-sh/smarter),
usually checked out beside this repository as `../smarter`. The SDKs
([Python](https://github.com/smarter-sh/smarter-python), [Node.js](https://github.com/smarter-sh/smarter-nodejs),
[Go](https://github.com/smarter-sh/smarter-golang)) cover the same api for
programs, and read this cli's configuration file.

The cli is a **thin client**: every command is one `POST` to the platform's
cli api, `/api/v1/cli/`, whose brokers do the work (`doc/adr/001-scope.md`,
`003-cli-api-only.md`). It holds no business rules about any kind.

## Non-negotiable rules

1. **Never execute git write commands**, here or in `../smarter` or the SDKs.
   "Write a git commit" means output the command text. See `smarter-git-commits`.
2. **Never trigger a release.** Don't run `make release` (it commits and
   pushes), `make choco-push`, `gh workflow run`, or `scripts/bump_version.sh`.
   See `smarter-cli-release`.
3. **Never point the cli at a real platform from a test or a probe.** Tests use
   the mock server. Manual checks use `--environment local` with a throwaway
   `$HOME` (`smarter-api-testing`). Never `deploy`, `undeploy` or `delete`
   against `alpha`, `beta`, `next` or `prod`: deploy creates billable
   infrastructure. The local platform's containers hold live cloud credentials
   too, so don't `deploy` there either unless the user asks.
4. **Never touch the user's `~/.smarter/`.** Run the binary with
   `HOME=$(mktemp -d)`, and write the config there.
5. **Run the tests yourself** after every change: `go test ./...`. Never ask the
   user to re-run them. Report failures honestly, including pre-existing ones.
6. **The pre-commit hooks are the lint standard.** See `smarter-code-quality`.
7. **Coverage:** at least 90% per package, 95% where it takes moderate effort.
8. **This repository is public.** No api keys, tokens, account numbers,
   private hostnames or customer data in code, tests, docs or skills. Tests use
   the placeholder key `0123456789abcdef` repeated four times (64 hex characters)
   and account `1234-5678-9012`.
9. **The configuration file is a contract.** `~/.smarter/config.yaml`'s keys,
   sections, environments and host rule are read by all three SDKs
   (`doc/adr/002-configuration.md`). Change them only with the user's
   agreement, and say which SDK files must follow
   (`../smarter-golang/config.go`, `../smarter-python`'s and `../smarter-nodejs`'s config modules).
10. **The platform's `openapi.json` is the authority on routes and kinds.** The
    cli pins a copy in `testdata/contracts/openapi.json` (`make contracts`), and
    `contract_test.go` fails when a command calls a route that isn't in it, a
    `/cli/` route has no command, or a kind lacks its commands (ADR-022).
    Never edit the copy, and never weaken the test to make it pass.
11. **New code returns errors**, and only a command's `Run` calls
    `cmd.ErrorOutput`. Don't add `log.Fatal`, `os.Exit` or `panic`
    (`doc/adr/006-error-handling.md`). stdout is for data only.

## Repository map

| Path                                      | What lives there                                                                  |
| ----------------------------------------- | --------------------------------------------------------------------------------- |
| `main.go`                                 | Loads `.env`, blank-imports every verb package, calls `cmd.Execute(Version)`      |
| `cmd/root.go`                             | `RootCmd`, global flags, `initConfig` (the configuration lookup)                  |
| `cmd/config.go`                           | `writeConfigFile`: only configuration keys, mode 0600                             |
| `cmd/api.go`                              | `APIRequest`: the only http code; `getAPIHost`, `SetAPIHostOverride` for tests    |
| `cmd/console.go`                          | `ConsoleOutput`, `DataOutput`, `ErrorOutput` (stderr, exit 1), `SetExitFunc`      |
| `cmd/resource.go`                         | `ResourceSpec`, `FlagSpec`, `NameArgSpec`, `RegisterResources`                    |
| `cmd/resourcekind.go`, `resourcekinds.go` | `ResourceKind` and its spec builders; the `ResourceKinds` registry                |
| `cmd/apply.go`, `validate.go`, ...        | Commands without a kind: `apply`, `validate`, `configure`, `resources`, ...       |
| `cmd/<verb>/<verb>.go`                    | A verb's command, its `APIRequest(kind, kwargs)` and `ConsoleOutput`              |
| `cmd/<verb>/resources.go`                 | The verb's `legacySpecs` and the `init()` that registers them and generated specs |
| `cmd/**/main_test.go`, `testutil_test.go` | Per-package fake `$HOME`, `withViperValue`, `newMockAPIServer`                    |
| `main_test.go`                            | Builds the binary; runs `version`, and `configure` on a fresh `$HOME`             |
| `contract_test.go`                        | Every command against the pinned OpenAPI document (ADR-022)                       |
| `testdata/contracts/openapi.json`         | The pinned copy of the platform's OpenAPI document. Never edit by hand            |
| `scripts/sync-contracts.sh`               | Copies it from `../smarter` (`make contracts [SMARTER_REF=<tag>]`)                |
| `scripts/bump_version.sh`                 | Release-time version stamping (semantic-release runs it)                          |
| `debian/`, `chocolatey/`, `Dockerfile`    | Packaging for the apt PPA, Chocolatey and Docker Hub                              |
| `.github/workflows/`                      | `test.yml`, semantic-release on `main`, the manual build-release workflow         |
| `doc/adr/`                                | The cli's architecture decision records                                           |
| `skills/`                                 | These skills                                                                      |

`venv/` (pre-commit) and `node_modules/` (semantic-release) are local tooling.
`go test ./...` walks into `venv/` and reports `[no test files]` there; ignore it.

The user guide is in the platform's repository:
`../smarter/docs/source/smarter-framework/smarter-cli.rst`.

## Architecture, and why

Read the ADR before you change one of these, and add an ADR when you make a
decision of that size.

- **One transport** (ADR-008): `cmd.APIRequest(slug, kwargs, body...)` POSTs
  to `<host>/api/v1/cli/<slug>/`, lower-cased with a trailing slash,
  `Authorization: Token <key>`, kwargs as the query string, with a
  `--timeout`, retried `--max_retries` times on transient failures as the
  SDKs do. It returns every failure as an error, with the platform's
  journal `error` description and context, for the command to report.
- **Kinds are data** (ADR-005): `cmd.ResourceKinds` holds one line per kind;
  `GetSpec`, `DescribeSpec`, `DeleteSpec`, `LogsSpec`, `ManifestSpec`,
  `JSONSchemaSpec`, `DeploySpec` and `UndeploySpec` generate the leaf
  commands. `account`, `plugin` and `user` are `legacySpecs` in a verb's
  `resources.go`. `Deployable` means the platform's broker implements deploy,
  as `smarter resources` reports (SmarterAuthToken, CustomDomain, LLMClient,
  LLMHost, LLMHostCompute, Vectorstore; User is a legacy spec).
- **One package per verb** (ADR-018), wired by blank imports in `main.go`.
- **cobra and viper** (ADR-004): every flag is bound to a viper key, and code
  reads viper, not the flag set. viper is a global singleton.
- **Configuration** (ADR-002): flags > the SDKs' `SMARTER_*` variables >
  `~/.smarter/config.yaml` > defaults. `--environment` > `SMARTER_ENVIRONMENT`
  > `config.environment` > `prod`; `--api_key` > `SMARTER_API_KEY` >
  > `<environment>.api_key`; the host is `localhost:9357` for
  > `local`, `<environment>.<root_domain>` for `alpha`/`beta`/`next`,
  > `<root_domain>` for `prod`. `initConfig` never exits; `writeConfigFile`
  > writes only configuration keys, owner-only.
- **Output** (ADR-013): responses are never modelled; they are reformatted as
  json (the whole envelope), yaml (`data` only, so `describe` output can be
  re-applied) or a table (`get`'s `titles` and `items`). `manifest` and
  `json-schema` print only `data` (`DataOutput`).

## Adding a resource kind (the common change)

When the platform adds a kind, `make contracts` brings it into the pinned
OpenAPI document's `SAMKind` enum, and the contract test fails with
`the SAMKind <Kind> has no "get" command`, and so on. Then:

1. Add one line to `cmd.ResourceKinds`, alphabetically: `Singular` and
   `Plural` are lower-case command names, `APIKind` is the kind exactly as in
   the `SAMKind` enum, `Display`/`DisplayPlural` are for help text.
2. Set `Deployable` exactly as the platform's `resources/` command reports the
   kind (`smarter resources -o json` against the local platform, or the
   broker's `deployable` flag in `../smarter`). Update
   `TestResourceKindsDeployableSubset`.
3. `go test ./...`, then check the help: `go run . get --help`,
   `go run . manifest <kind> --help`.
4. Mention the kind in `README.md` and the docs page if they list kinds.

## Adding a verb or a flag

- **A verb with kinds**: a new `cmd/<verb>/` package copying `cmd/jsonschema/`
  (the newest), with `<verb>.go`, `resources.go`, `main_test.go`,
  `testutil_test.go` and tests; a `<Verb>Spec()` on `ResourceKind` if every
  kind gets it; and a blank import in `main.go`, without which the verb
  silently doesn't exist. The contract test runs it automatically.
- **A command without kinds**: a file in `cmd/` copying `resources.go` (the
  command) or `validate.go` (with a manifest file). If it needs arguments
  beyond a name, add them to `commandArgs` in `contract_test.go`.
- **A flag**: a `FlagSpec` in the spec (bound to viper for you), or for a whole
  verb a persistent flag bound with `mustBindPFlag`. Flag names use
  underscores, like the configuration keys. Two commands that bind the same
  name share one viper key, and the later binding wins: `validate` reads its
  `-f` from its own flag set for that reason.
- Every new route must be an operation of the pinned `openapi.json`; if it
  isn't, the platform must publish it first (`smarter-api-contract` in
  `../smarter/skills/`).

## Keeping in step with the platform

Refresh the pinned document whenever the platform's api changes, and at the
monthly maintenance: `make contracts SMARTER_REF=<latest stable tag>` (which
records the version in `testdata/contracts/SOURCE`), then
`go test -run TestCommandsMatchTheOpenAPIDocument -v .`. Read the copy's diff:
it lists the api's changes. When you find a platform bug, report it rather
than working around it in silence.

A command must send only the query parameters that its route documents
(`name` for the kind routes); the contract test fails otherwise. The cli has
been bitten by parameters the platform silently ignored: `dry_run` (so
`apply --dry-run` applied), `username` (so the user commands found no one),
and `i`, `asc`, `desc`, `class`.

## Known rough edges

Mention these when relevant, and fix them when asked:

- `CONFIG_*` variables (viper's accidental prefix) still work beside the
  documented `SMARTER_*` ones, for old scripts (ADR-002).
- `undeploy customdomain` and `undeploy llmhostcompute` return the platform's
  "not implemented" error: their brokers deploy but don't undeploy.
- `SMARTER_ROOT_DOMAIN` is deliberately not read: the platform uses that name
  for its apex domain, and developers' shells export it (ADR-002).

## Which skill to load

| Task                                                        | Skill                  |
| ----------------------------------------------------------- | ---------------------- |
| Unit tests, the contract test, exits, coverage              | `smarter-cli-testing`  |
| Running the cli or curl against the local platform          | `smarter-api-testing`  |
| gofmt, go vet, goimports, golangci-lint, pre-commit         | `smarter-code-quality` |
| Writing a commit message                                    | `smarter-git-commits`  |
| Versions, semantic-release, the build-release workflow      | `smarter-cli-release`  |
| Monthly sync with the platform's kinds and routes, and deps | `monthly-maintenance`  |
| A platform contract change to follow now, between maintenances | `smarter-contract-sync` |

When a task leads into the platform (a route, a broker, the docs page), load
the platform's skills from `../smarter/skills/`, starting with
`smarter-development`; for the docs page, `smarter-sphinx-docs`.

## Definition of done

- The change follows an existing pattern here and the ADRs in `doc/adr/`.
- The contract test passes against the pinned `openapi.json`.
- Tests cover it, `go test ./...` passes, and `make coverage` passes (90% or
  more per package).
- The pre-commit hooks pass on the changed files.
- Help text (`Short`, `Long`) is written, and `README.md` and the docs page
  describe user-visible changes.
- Your summary says what changed, what you verified, what you skipped, and
  any platform drift or bug you found.
