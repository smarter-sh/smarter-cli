---
name: smarter-code-quality
description: Use before finishing any change to smarter-cli, to format, vet, lint and test it with the repository's pre-commit hooks (go-fmt, go-vet, go-imports, golangci-lint, go-unit-tests, the generic hooks, commitlint). Explains which tools are authoritative, where pre-commit lives, how to run the hooks on changed files under zsh, and the conventions the hooks don't check (copyright header, error wrapping, help text).
---

# Smarter CLI Code Quality

`.pre-commit-config.yaml` is the standard (`doc/adr/014-code-quality.md`).
The Go SDK copies its Go hooks from here. If a hook passes, the code meets the
standard; if it fails, fix the code. Never bypass a hook (`--no-verify`,
`SKIP=`, a bare `//nolint`).

## What the hooks enforce

| Hook          | Notes                                                                                        |
| ------------- | -------------------------------------------------------------------------------------------- |
| go-fmt        | `gofmt`. Let it reformat.                                                                    |
| go-imports    | Import grouping: standard library, then others.                                              |
| go-vet        | `go vet ./...`.                                                                              |
| golangci-lint | Its default linters (errcheck, govet, ineffassign, staticcheck, unused); no `.golangci.yml`. |
| go-unit-tests | `go test ./...`: a commit can't break the tests.                                             |

Also: large files, case conflicts, JSON/TOML/XML/YAML syntax, merge conflict
markers, symlinks, AWS credentials, private keys, end-of-file newline, trailing
whitespace, and commitlint (`@commitlint/config-angular`) on commit messages.

## How to run

pre-commit is installed in `venv/` at the repository root (`make python-init`
creates it; `make pre-commit-init` installs the hooks). The Go hooks need
`goimports` (`~/go/bin`) and `golangci-lint` on `PATH`. From the repository
root:

```console
gofmt -l cmd main.go && go vet ./... && golangci-lint run
source venv/bin/activate
git diff --name-only | xargs pre-commit run --files
git ls-files --others --exclude-standard | grep -v '^venv/\|^node_modules/' | xargs pre-commit run --files
```

- Use `xargs`. Under zsh, `pre-commit run --files $FILES` passes one argument,
  and pre-commit checks nothing.
- Hooks that modify files fail the first run, then pass. Rerun, and read the
  diff of what they changed.
- If `venv/` is missing or its Python is gone, use any pre-commit on `PATH`
  (e.g. the platform's `../smarter/venv/bin/pre-commit`); the hook
  environments are cached in `~/.cache/pre-commit`.
- `make pre-commit-run` runs every hook on every file. `make lint` runs
  golangci-lint only.

## Conventions the hooks don't check

- **Copyright header.** Every Go source file starts with the project's comment
  block (`Copyright © 2024 Lawrence McDaniel ...`), as the existing files do.
- **Errors:** wrap with `%w`; return errors from helpers and leave
  `ErrorOutput` to the command's `Run` (ADR-006). No new `log.Fatal`,
  `os.Exit` or `panic`.
- **stdout is data.** Diagnostics go through `log` (stderr) and only with
  `--verbose` (ADR-007). Never print the api key; mask it as `configure` does.
- **Help text:** `Short` is one line; `Long` shows the usage line and an
  example. Generated kinds get theirs from `ResourceKind`; don't hand-write a
  spec for a kind that fits the template (ADR-005, ADR-012).
- **Doc comments** on exported identifiers start with the name and say why, as
  in `cmd/resource.go`.
- **`interface{}`** appears only where the cli reformats untyped json
  (ADR-013). New code can write `any`.
- Match the surrounding code; don't reformat code you didn't otherwise change.
- Suppress a finding only with a narrow directive that names the linter and
  says why: `//nolint:errcheck // best effort close`.
