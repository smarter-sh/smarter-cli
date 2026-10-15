# ADR-012: Help Text and Documentation

## Status

Accepted

Follows the platform's [ADR-029: Docstrings](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/029-docstrings.rst) and the SDKs' ADR-012.

## Context

Most users read the cli's documentation through `smarter --help`. The platform's documentation site, docs.smarter.sh, explains manifests and resources, and has a page for the cli.

## Decision

- **Help text is the reference.** Every command has a `Short` (one line, sentence case, no trailing period in new text) and a `Long` that shows the usage line and, where helpful, an example such as `smarter describe <kind> <name> > manifest.yaml`. Generated kinds get theirs from `ResourceKind` (ADR-005), so their wording is consistent.
- **Go doc comments** on exported identifiers start with the name, and explain why as well as what, as in `cmd/resource.go` and `cmd/resourcekind.go`.
- **README.md** covers installation, configuration, the command table, and building from source.
- **docs.smarter.sh** has the user guide, in the platform repository at `docs/source/smarter-framework/smarter-cli.rst`. A user-visible change to a command, a flag or the configuration is documented there too, following the platform's `smarter-sphinx-docs` skill.
- The documentation is evergreen: new features are woven into the existing text, with no "what's new" sections.

## Alternatives Considered

- Generated markdown docs (`cobra/doc`). Possible later; the generated pages would duplicate the help text.

## Consequences

- **Positive:**
  - `--help` is always current with the binary.
- **Negative:**
  - The README, the docs page and the help text overlap, and must be updated together.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-014: Code Quality](014-code-quality.md)
