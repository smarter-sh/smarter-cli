# ADR-007: Logging and Verbose Output

## Status

Accepted

The SDKs' ADR-007 adapt the platform's [ADR-011: Logging](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/011-logging.rst).

## Context

A cli's normal output is its result, and nothing else. When a command fails, the user needs to see what was sent and what came back, without a debugger.

## Decision

- The cli is **silent by default**: a command prints its result, or `deployed.`, `deleted.`, `manifest applied.`, and nothing more.
- `--verbose`/`-v` turns on diagnostics: the configuration file used, the environment, the request url and body, and dumps of the request and response (`httputil.DumpRequestOut`, `DumpResponse` without its body), and the platform's stack trace on an error.
- Diagnostics use the standard `log` package, which writes to stderr, so `-v` doesn't corrupt data on stdout.
- The api key is never logged in full: `maskSecret` keeps its last four characters, in `initConfig`'s log and in the request dump's `Authorization` header.
- `smarter version -v` also asks the platform for its version.

## Alternatives Considered

- Levelled, structured logging (`log/slog`). A cli has one user, at one terminal; on/off is enough.

## Consequences

- **Positive:**
  - Bug reports can include `-v` output.
- **Negative:**
  - The request dump still shows the request body, e.g. a manifest with secret values. Review `-v` output before sharing it.
  - Until October 2026, `APIRequest` printed its dumps to stdout with the full api key; `TestAPIRequestVerboseMasksTheKeyOnStderr` guards against that.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-006: Error Handling and Exit Codes](006-error-handling.md)
- [ADR-016: Security](016-security.md)
