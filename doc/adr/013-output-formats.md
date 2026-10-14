# ADR-013: Output Formats

## Status

Accepted

## Context

People read lists as tables; scripts and `kubectl`-style workflows read manifests as yaml (`smarter describe ... > manifest.yaml`, edit, `smarter apply -f`), and parse json. The platform's cli api returns json, wrapped in `{"data": ..., "message": ...}`, and its `get` responses carry `titles` and `items` for tables.

## Decision

`cmd.ConsoleOutput(body)` (`cmd/console.go`) formats every successful response, by `--output_format`/`-o`:

- `json`: the whole response, indented.
- `yaml`: the response's `data` member only, converted with `ghodss/yaml`, so that `describe` output is a manifest that `smarter apply` accepts.
- `tabular`: the `titles` and `items` of a `get` response, aligned with `text/tabwriter`; `DateTimeField` columns are shown as `2006-Jan-02 15:04`.

`get` defaults to `tabular` and `chat` to `json` when `-o` isn't given. The root flag's default is `json`; `config.output_format` defaults to `yaml`.

`manifest` and `json-schema` print only the response's `data` (`cmd.DataOutput`), as yaml and json by default, so that their output is a file to save: an example manifest to edit and apply, or a schema for an editor. `validate` prints `manifest is valid.`, or each error with its path to stderr, and the response itself only when `-o` is given.

The cli **doesn't model responses**: it decodes them into `map[string]interface{}` only to reformat them, so a new field on the platform appears in the output without a cli change.

## Alternatives Considered

- Go structs per response, as the Go SDK has. Every platform change would need a cli release to be visible.
- `-o jsonpath`/templates, as `kubectl` has. Not asked for; `jq` and `yq` serve.

## Consequences

- **Positive:**
  - The cli shows whatever the platform returns.
- **Negative:**
  - `json` and `yaml` show different shapes (with and without the envelope).
  - `-o tabular` is rejected outside the `get` commands, whose responses are the only ones with tables. `configure` deliberately accepts only `json` and `yaml`, because its value becomes every command's default.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-006: Error Handling and Exit Codes](006-error-handling.md)
