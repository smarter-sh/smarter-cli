# ADR-006: Error Handling and Exit Codes

## Status

Accepted

The SDKs' ADR-006 adapt the platform's [ADR-012: Error Handling](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/012-error-handling.rst) to a library, which returns errors; a cli reports them and exits.

## Context

Scripts and CI pipelines test a command's exit status, and read its stdout as data. An error message mixed into stdout corrupts a `smarter describe ... > manifest.yaml`. The platform returns errors as json with `description`, `context` and, in debug, `stacktrace`.

## Decision

- **stdout carries data only.** Errors go to stderr.
- **Exit status is 0 on success and 1 on any error**, an invalid manifest under `validate` included. `cmd.ErrorOutput(err)` prints `Error: <message>` to stderr and exits 1; `cmd.Execute` exits 1 when cobra reports a usage error. Tests replace the exit with `cmd.SetExitFunc` (ADR-009).
- **Functions below the command return errors**, wrapped with `%w` (`APIRequest`, `getAPIHost`, `readManifest`, `ConsoleOutput`, `validateChoices`, `writeConfigFile`, ...), and the command's `Run` passes them to `ErrorOutput`. Verb packages hand `cmd.ErrorOutput` to `RegisterResources` as their `ErrFunc`.
- **A non-200 response** is reported with its status, url, the platform's `description` and `context` (from the journal's `error` object, or Django REST Framework's `detail`), and with `--verbose` its `errorClass` and `stacktrace`.
- New code returns errors and leaves exiting to the command; it doesn't add `log.Fatal`, `os.Exit` or `panic`.

## Alternatives Considered

- Distinct exit codes per error class (usage, auth, not found). Useful, but no caller has asked, and it would be a breaking change to add them inconsistently.
- `RunE` everywhere, letting cobra print errors. Cobra also prints the usage text on every error, which buries the message.

## Consequences

- **Positive:**
  - `smarter ... > file` never captures an error message.
- **Negative:**
  - `log.Fatalf` remains only in `init()` functions, for programming errors such as binding a flag that doesn't exist, which no input can cause.
  - Until October 2026, `APIRequest`, `initConfig`, `getAPIHost` and the chat commands exited or panicked themselves; tests of those paths ran in subprocesses.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-007: Logging and Verbose Output](007-logging.md)
- [ADR-008: HTTP Transport](008-http-transport.md)
- [ADR-009: Testing](009-testing.md)
