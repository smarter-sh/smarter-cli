# ADR-009: Testing

## Status

Accepted

The SDKs' ADR-009 adapt the platform's [ADR-013: Testing](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/013-testing.rst).

## Context

The cli's tests must not touch the developer's real `~/.smarter/config.yaml`, must not reach a real platform, and must cover code that calls `os.Exit`. cobra runs `initConfig` on every `Execute`, and viper is a global singleton bound to flags at `init()` time.

## Decision

- Tests use the standard `testing` package, one `<file>_test.go` per source file, in the package under test, and run with `go test ./...` (`make test`, and the `go-unit-tests` pre-commit hook).
- **A fake home.** Each package's `TestMain` (`main_test.go`) creates a temporary `$HOME` with a complete `config.yaml` and a placeholder 64-hex api key before any test runs.
- **A mock platform.** `newMockAPIServer(t, handler)` (`testutil_test.go`) starts an `httptest.Server`, points `APIRequest` at it with `cmd.SetAPIHostOverride`, and sets the api key and environment; tests assert the method, path, query and body that the platform would receive.
- **Viper isolation.** `withViperValue(t, key, value)` overrides one key and restores it in `t.Cleanup`. Tests never call `viper.Reset()`.
- **Exits.** `cmd.SetExitFunc` replaces `ErrorOutput`'s `os.Exit` with a panic that the test recovers (`exits(t, fn)`). Code that must exit for real (`ErrorOutput` itself) is tested in a subprocess: the test re-executes its own binary with `-test.run=<name>` and `SMARTER_CLI_TEST_CRASHER=1`.
- **The contract.** `contract_test.go` runs every command against a recording mock platform, and checks every request against a pinned copy of the platform's `openapi.json` (ADR-022).
- **The binary.** `main_test.go` at the root builds `smarter` and runs it with a throwaway `$HOME`: `version`, and `configure` on a fresh install; `-short` skips it.
- **Coverage:** at least 90% per package, 95% where it takes moderate effort. `make coverage` fails below 90%, in CI too; the `main` package is exempt, because `main()` runs only in the binary tests.

## Alternatives Considered

- Commands that return errors (cobra's `RunE`). Cobra would print the usage text with every error (ADR-006); the exit hook keeps `Run` as it is.
- Live tests against a platform in CI. The cli's CI has no platform; the contract test checks the requests against the platform's published api instead, and a manual smoke test against the local platform is described in the `smarter-api-testing` skill.

## Consequences

- **Positive:**
  - Tests are hermetic and fast, and safe on a developer's machine.
- **Negative:**
  - `withViperValue`, `newMockAPIServer` and the `TestMain` fixture are copied into every package, because Go test helpers aren't importable across packages.
  - A mock server can't catch a route the platform renamed; the contract test does, against the pinned copy (ADR-022), but only as of its last refresh.
  - `captureStdout` must restore `os.Stdout` even when its function panics: the test binary prints its coverage there.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-006: Error Handling and Exit Codes](006-error-handling.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-022: A Pinned Copy of the Platform's OpenAPI Document](022-pinned-openapi.md)
