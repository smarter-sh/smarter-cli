# ADR-002: Configuration

## Status

Accepted

The SDKs' ADR-002 ([Python](https://github.com/smarter-sh/smarter-python/blob/main/doc/adr/002-configuration.md), [Node.js](https://github.com/smarter-sh/smarter-nodejs/blob/main/doc/adr/002-configuration.md), [Go](https://github.com/smarter-sh/smarter-golang/blob/main/doc/adr/002-configuration.md)) read this file as the cli writes it.

## Context

A user works with several Smarter platforms: a local one, the `alpha`, `beta` and `next` pre-release platforms, and production, each with its own api key. A self-hosted installation has its own root domain. Commands must work without retyping credentials, and scripts must be able to override everything.

## Decision

The configuration lives in `~/.smarter/config.yaml` (`--config` to override), written by `smarter configure` and read by viper in `cmd/root.go`:

```yaml
config:
  account_number: "1234-5678-9012"
  environment: prod
  output_format: yaml
  root_domain: platform.smarter.sh
  username: example_user
local:
  api_key: "<64 hex characters>"
prod:
  api_key: "<64 hex characters>"
```

- **The environment** is one of `local`, `alpha`, `beta`, `next`, `prod`: the `--environment` flag, else `SMARTER_ENVIRONMENT`, else `config.environment`, else `prod`.
- **The api key** is per environment: `--api_key` when it is passed, else `SMARTER_API_KEY`, else `<environment>.api_key`.
- **The host** is `SMARTER_BASE_URL` (unless `--environment` is passed), else derived from the environment and `config.root_domain` (`getAPIHost` in `cmd/api.go`): `local` is `http://localhost:9357`, `alpha`/`beta`/`next` are `https://<environment>.<root_domain>`, and `prod` is `https://<root_domain>`.
- **The SDKs' environment variables** apply below the flags and above the file (`cmd/environment.go`): `SMARTER_API_KEY`, `SMARTER_BASE_URL`, `SMARTER_ENVIRONMENT`, `SMARTER_CONFIG_FILE`, `SMARTER_TIMEOUT` (seconds), `SMARTER_MAX_RETRIES` and `SMARTER_LOG` (`debug` or `info` turns on `--verbose`). An invalid value fails the command with the SDKs' message.
- **`SMARTER_ROOT_DOMAIN` is not read**, though the SDKs read it as `config.root_domain`: the platform's own settings use that name for its apex domain (`smarter.sh`, its marketing site), and developers' shells export it, which would send the cli's requests to the marketing site.
- On first run, `initConfig` creates `~/.smarter/` and writes the defaults; defaults also fill keys missing from an existing file. It never exits: a missing api key is reported by the first command that needs one, so `smarter configure` and `smarter version` work on a fresh install.
- The file is written by `writeConfigFile` (`cmd/config.go`), with only the keys above (and a top-level `output_format` that `configure` set), owner-only (ADR-016).
- `smarter configure` validates each value (account number `1234-5678-9012`, a 64-character hex api key, a domain name) before saving it, and masks the api key when it echoes it.
- `main.go` loads a `.env` file from the working directory (`godotenv`) for local development.

This file format and lookup are a **contract with the SDKs**, which read it. A change to a key, a section, the environments or the host rule must be made in the Python, Node.js and Go SDKs too.

## Alternatives Considered

- Environment variables only. Awkward for several platforms at once.
- One api key for all environments. A production key would be sent to pre-release platforms.

## Consequences

- **Positive:**
  - One `smarter configure` serves the cli and every SDK.
- **Negative:**
  - The host rule can't express a platform whose pre-release hosts don't follow `<environment>.<root_domain>`.
  - Rewriting the file drops keys the cli doesn't know. Until October 2026 the cli wrote every flag bound to viper into it (`i`, `asc`, `chatbot`, ...), and exited on a fresh install before `configure` could run; those files are cleaned up the next time `configure` writes them.
  - viper's `SetEnvPrefix("config")` with `AutomaticEnv()` also maps keys to `CONFIG_*` variables (`CONFIG_API_KEY`, ...). They predate the `SMARTER_*` ones, are undocumented, and are kept so that old scripts keep working.
  - The SDKs read `SMARTER_ROOT_DOMAIN`, which collides with the platform's variable of the same name; the cli deliberately doesn't.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-016: Security](016-security.md)
