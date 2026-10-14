# ADR-003: Only the CLI API

## Status

Accepted

The SDKs' ADR-003 ([Go](https://github.com/smarter-sh/smarter-golang/blob/main/doc/adr/003-writes-through-manifests.md)) send their writes through the same api, "exactly as the cli sends them".

## Context

The platform serves two apis: the manifest api under `/api/v1/cli/`, built for this cli, where one broker per kind implements every verb; and REST routes per resource family under `/api/v1/`. The platform's [ADR-005: Manifests](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/005-manifests.rst) and [ADR-024: Broker Model](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/024-broker-model.rst) make the manifest the unit of change.

## Decision

The cli calls only `/api/v1/cli/`, through one function, `cmd.APIRequest(slug, kwargs, body)`:

| Command                          | Route                                                  |
| -------------------------------- | ------------------------------------------------------ |
| `apply -f <file>`                | `apply/` with the manifest as the body                 |
| `get <kinds>`                    | `get/<Kind>/`                                          |
| `describe <kind> <name>`         | `describe/<Kind>/`                                     |
| `delete`, `deploy`, `undeploy`   | `delete/<Kind>/`, `deploy/<Kind>/`, `undeploy/<Kind>/` |
| `logs <kind> <name>`             | `logs/<Kind>/`                                         |
| `manifest <kind>`                | `example-manifest/<Kind>/`                             |
| `json-schema <kind>`             | `json-schema/<Kind>/`                                  |
| `validate -f <file>`             | `validate/` with the manifest as the body              |
| `chat prompt`, `chat config`     | `prompt/<name>/`, `prompt/config/<name>/`              |
| `whoami`, `status`, `version -v` | `whoami/`, `status/`, `version/`                       |
| `resources`                      | `resources/`                                           |

- Every request is a `POST`. Arguments go in the query string, and only those that the route documents: `name` for the kind routes, `uid` and `new_session` for the prompt routes. A manifest or prompt goes in the body.
- `apply --dry-run` posts to `validate/`, which checks a manifest as `apply` does first, without saving it: the platform's `apply` has no dry run.
- Flags whose arguments the platform ignores (`get --i`, `--asc`, `--desc`, `get plugins --class`) stay accepted for old scripts, marked deprecated, and are never sent.
- The `<Kind>` segment is the platform's SAM kind; the whole url is lower-cased, and the platform matches kinds case-insensitively.
- The platform's `openapi.json` (`smarter/smarter/apps/api/v1/contract/` in the platform repository) is the authority on these routes.

## Alternatives Considered

- The REST routes. Each family has its own shapes; the cli would need code per kind.
- A mix. Two error shapes, two auth paths, and drift.

## Consequences

- **Positive:**
  - One transport, one error shape, and kinds that need no code (ADR-005).
- **Negative:**
  - Until October 2026 the cli sent arguments that no route reads: `dry_run`, so `apply --dry-run` applied the manifest while it reported that nothing changed; `username`, so `describe user`, `delete user` and `get users -u` never found the user; and `i`, `asc`, `desc` and `class`. The contract test now fails on an undocumented query parameter (ADR-022).
  - The cli breaks when the platform renames a route, until it follows. Its unit tests use mock servers, so the contract test checks every command against a pinned copy of the platform's `openapi.json` (ADR-022); that is how `manifest`'s `example_manifest/` and `chat`'s `chat/` were found, in October 2026.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-008: HTTP Transport](008-http-transport.md)
- [ADR-022: A Pinned Copy of the Platform's OpenAPI Document](022-pinned-openapi.md)
