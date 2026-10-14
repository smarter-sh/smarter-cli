---
name: smarter-cli-testing
description: Use when writing, fixing or running smarter-cli's Go tests, or raising coverage. Covers the per-package TestMain that fakes $HOME and config.yaml, withViperValue (never viper.Reset), newMockAPIServer and SetAPIHostOverride, capturing stdout safely, testing commands that exit (cmd.SetExitFunc and the exits helper, or a SMARTER_CLI_TEST_CRASHER subprocess), the contract test against the pinned OpenAPI document, the binary tests, make coverage and its 90% gate, a manual smoke test of the built binary, and what never to do in a test.
---

# Smarter CLI Testing

Tests use the standard `testing` package, one `<file>_test.go` per source
file, in the package under test (`doc/adr/009-testing.md`).

```console
go test ./...                          # everything, including the binary smoke test
go test -short ./...                   # skips building the binary
go test ./cmd/get/ -run TestAPIRequest -v
make coverage                          # per-package coverage, fails below 90% (CI runs it)
go test -run TestCommandsMatchTheOpenAPIDocument -v .   # the contract test
go test -coverprofile="$SCRATCHPAD/cover.out" ./cmd/ && go tool cover -func="$SCRATCHPAD/cover.out" | sort -k3 -n | head -30
```

Write coverage profiles to your scratchpad, not the repository (`*.out` is
git-ignored, but don't leave files behind). The target is 90% per package, 95%
where it takes moderate effort; `make coverage` enforces 90% in every package
but `main`, whose `main()` runs only in the binary tests. The
`go-unit-tests` pre-commit hook and the `test.yml` workflow run the tests.

## The five traps, and their fixes

### 1. `initConfig` runs on every `Execute`

`cobra.OnInitialize(initConfig)` fires whenever any test executes a command,
in any package that imports `cmd`. Without a config it creates `~/.smarter/`
in the **real** home directory, or exits on a missing api key and kills the
test binary.

**Fix:** every test package has a `main_test.go` whose `TestMain` creates a
temporary `$HOME` with a complete `config.yaml` (placeholder key
`testAPIKey`, 64 hex characters, account `1234-5678-9012`, environment `alpha`) before
`m.Run()`. A new test package copies `cmd/get/main_test.go` as is.

### 2. viper is a global singleton

Flags are bound to viper in `init()`. `viper.Reset()` destroys those bindings
for every later test.

**Fix:** never call `viper.Reset()`. Override one key at a time with
`withViperValue(t, key, value)` (`testutil_test.go`), which restores it in
`t.Cleanup`. Don't call `t.Parallel()` in tests that touch viper.

### 3. Requests must never leave the machine

**Fix:** `newMockAPIServer(t, handler)` starts an `httptest.Server`, points
`APIRequest` at it with `cmd.SetAPIHostOverride`, and sets `api_key` and
`environment`. Assert what the platform would receive:

```go
var gotPath, gotQuery string
newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
	gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"data":{}}`))
})
if _, err := APIRequest("Secret", map[string]string{"name": "foo"}); err != nil {
	t.Fatalf("APIRequest() error = %v", err)
}
if gotPath != "/api/v1/cli/describe/secret/" { ... }
```

Paths arrive **lower-cased**. Write the expected path from the platform's
`openapi.json`, not from the code under test: a test that copies the code's
slug passes on a route the platform doesn't serve (that is how
`example_manifest/` survived). The contract test is the backstop.

### 4. Code that exits

`cmd.ErrorOutput` ends the process, and a command's `Run` calls it on any
error. Replace the exit with a panic that the test recovers: `exits(t, fn)`
(`cmd/config_test.go`, and copied in `cmd/chat/chat_test.go`) calls
`cmd.SetExitFunc`, runs `fn`, and reports whether it exited:

```go
func TestSomethingExits(t *testing.T) {
	if !exits(t, func() { someCmd.Run(someCmd, nil) }) {
		t.Error("someCmd did not exit on ...")
	}
}
```

Only `ErrorOutput`'s real `os.Exit` needs a subprocess
(`TestErrorOutputExitsProcess`): re-run the test binary with
`-test.run=^Name$` and `SMARTER_CLI_TEST_CRASHER=1`. Anchor the pattern, or
the subprocess runs every test whose name shares the prefix.

### 5. A panic inside `captureStdout`

The test binary prints its coverage line to `os.Stdout` when it ends. If a
panic (e.g. from the exit hook) escapes `captureStdout` before it restores
`os.Stdout`, the line goes into a dead pipe, and `make coverage` silently
skips the package. `captureStdout` in `cmd/console_test.go` restores stdout in
a `defer`; copy that version.

## The contract test

`contract_test.go` (package `main`) runs **every** command in-process against
a recording mock platform and checks the requests against the pinned
`testdata/contracts/openapi.json` (ADR-022). When it fails:

- _is not an operation of the platform's OpenAPI document_: the cli calls a
  route the platform doesn't serve. Fix the slug.
- _publishes /cli/ operations that no command calls_: the platform added a
  route. Add a command, or ask the user.
- _sent the query parameter "x", which ... doesn't document_: the platform
  ignores it. Send the documented parameter (e.g. `name`, with
  `NameArgSpec.Param`), or stop sending it.
- _the SAMKind X has no "get" command_: a new kind. Add it to `ResourceKinds`
  (`smarter-cli-development`).
- _X failed_ or _sent no request_: a command needs arguments the test doesn't
  give. Add them to `commandArgs`, or the command to `noRequest`. Rerun with
  `SMARTER_CLI_CONTRACT_DEBUG=1` to see its output.

Exceptions belong in `kindExceptions`, each with its reason. Options that
send a different request (`apply --dry-run`) are extra `invocations`.

`cmd`'s `TestMain` replaces `sleep`, so retries don't wait, and every
package's `TestMain` unsets `cmd.SmarterEnvironmentVariables`, so a
developer's `SMARTER_*` variables can't change a test's outcome. A test of
those variables sets them with `t.Setenv`. Never edit the
pinned document by hand; refresh it with `make contracts`.

## Other patterns

- **stdout:** `captureStdout(t, fn)` in `cmd/get/get_test.go` swaps
  `os.Stdout` for a pipe. Copy it into a package that needs it.
- **Cobra commands:** call a leaf's `Run` directly, or
  `cmd.RootCmd.SetArgs([]string{"get", "secrets"})` and `Execute()`, with
  viper values set first. Flags bound to viper keep their values between
  tests; reset the ones you set.
- **Kind registry:** `TestResourceKindsAreWellFormed` and
  `TestResourceKindsDeployableSubset` guard `ResourceKinds`. When the
  platform's brokers change which kinds deploy, update the latter.
- **Flags of a subcommand:** a leaf's `Flags()` doesn't include its parent's
  persistent flags until cobra merges them; call `c.InheritedFlags()` first
  when a test sets them directly (see `setFlags` in `cmd/chat/chat_test.go`).
- **The binary:** root `main_test.go` builds `smarter` and runs it with a
  throwaway `$HOME`: `version`, and `configure` on a fresh install, which
  must write an owner-only file without flags in it.
- **Helpers are copied per package** (`withViperValue`, `newMockAPIServer`,
  `testAPIKey`, `captureStdout`), because `_test.go` helpers can't be
  imported. Keep the copies identical.

## A manual smoke test of the binary

After a change to a command's wiring, run the real binary against the mock
behavior you expect, or the local platform (`smarter-api-testing`), never
against your own `~/.smarter`:

```console
go build -o "$SCRATCHPAD/smarter" .
export HOME="$SCRATCHPAD/home" && mkdir -p "$HOME/.smarter"
printf 'config:\n  environment: local\n  root_domain: platform.smarter.sh\nlocal:\n  api_key: "%s"\n' "$(cat "$SCRATCHPAD/key.txt")" > "$HOME/.smarter/config.yaml"
"$SCRATCHPAD/smarter" get secrets
"$SCRATCHPAD/smarter" describe secret <name> -o yaml
```

Run those `export`s in a subshell or a single command, so your shell's `HOME`
isn't changed for later commands.

## Never

- Never let a test reach `alpha`, `beta`, `next`, `prod` or the local platform.
- Never use a real api key, account number or hostname in a test.
- Never call `viper.Reset()` or write to the real `~/.smarter/`.
- Never weaken an assertion to make a test pass; if the code is wrong, say so.
