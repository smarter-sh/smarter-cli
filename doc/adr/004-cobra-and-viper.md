# ADR-004: Cobra and Viper

## Status

Accepted

## Context

A verb-noun cli (ADR-001) needs nested commands, persistent and local flags, generated help, shell completion, and configuration merged from flags, a file and the environment. Writing that on `flag` would be most of the code.

## Decision

The cli is built on [cobra](https://github.com/spf13/cobra) for commands and flags, and [viper](https://github.com/spf13/viper) for configuration, as `kubectl`, `gh` and `hugo` are.

- `cmd.RootCmd` is the root. Each verb is a `*cobra.Command` added to it in its package's `init()` (ADR-018), and each kind a leaf command under the verb (ADR-005).
- Every flag is bound to a viper key (`mustBindPFlag`), and commands read values from viper, not from the flag set, so a flag, the configuration file and a default are read the same way.
- Global flags are persistent flags of the root: `--config`, `--environment`, `--api_key`, `--output_format`/`-o`, `--timeout`, `--verbose`/`-v`.
- Flag names use underscores (`--api_key`, `--output_format`, `--new_session`), as the configuration keys do; `--dry-run` is the exception.
- `cobra.OnInitialize(initConfig)` loads the configuration before every command; `PersistentPreRunE` validates `--environment` and `--output_format`.

## Alternatives Considered

- The standard library's `flag`. No subcommands, no help tree.
- `urfave/cli`. Less common for `kubectl`-style tools, and no viper integration.

## Consequences

- **Positive:**
  - Help, usage errors and `smarter completion` come for free.
- **Negative:**
  - viper is a process-wide singleton, and flags are bound to it at `init()` time. Tests must never call `viper.Reset()`; they override and restore single keys (ADR-009).
  - Two leaf commands that bind the same flag name share one viper key.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-002: Configuration](002-configuration.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-018: One Package per Verb](018-package-per-verb.md)
