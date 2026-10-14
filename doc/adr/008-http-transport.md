# ADR-008: HTTP Transport

## Status

Accepted

## Context

Every command is one request to the cli api (ADR-003), authenticated with an api key. The platform uses Django REST Framework token authentication ([django-rest-knox](https://jazzband.github.io/django-rest-knox/auth/)).

## Decision

`cmd.APIRequest` in `cmd/api.go` is the only place that sends http:

- `POST <host>/api/v1/cli/<slug>/`, where the host comes from the environment (ADR-002) and the slug from the verb and kind. The url is lower-cased and given a trailing slash.
- Headers: `Authorization: Token <api key>`, `Content-Type: application/json`, `User-Agent: smarter-cli/<version>` and `X-Smarter-Client: smarter-cli`, which identify the cli in the platform's logs as the SDKs identify themselves.
- The `kwargs` map becomes the query string, the only place the platform reads arguments from.
- Every request has a timeout, `--timeout`, two minutes by default: deploying an LLMHost returns when the cluster accepts its resources, and a prompt waits for the LLM.
- Retries follow the SDKs' rules: a request that got no response, or a 408, 429, 502 or 504, and a 503 that isn't journaled (the platform's own 503 means a resource isn't ready, which a retry won't change), is sent again, `--max_retries` times (2 by default), after exponential backoff from half a second to eight, with jitter, or the response's `Retry-After` up to a minute.
- The body is a manifest (`apply`) or a json prompt (`chat prompt`), else empty.
- A 200 response's body is returned as bytes, for the output formatter (ADR-013); any other status is an error (ADR-006).
- Tests point it at an `httptest.Server` with `cmd.SetAPIHostOverride(url)` (ADR-009).

## Alternatives Considered

- A generated client from `openapi.json`. Every cli route has the same shape; one function serves them all.
- Adopt the Go SDK's transport. Possible, and the SDK was written for it (ADR-001); it would change the cli's errors and output, and is not done yet.

## Consequences

- **Positive:**
  - One function to test, and to change.
- **Negative:**
  - A retried `apply` or `deploy` whose first attempt reached the platform before its response was lost is sent twice. The brokers' operations are idempotent (apply migrates to a state, deploy of a deployed resource does nothing), as the SDKs also assume.
  - Until October 2026 there were no retries.
  - Until October 2026 there was no timeout, and arguments were also copied into request headers.

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-002: Configuration](002-configuration.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-006: Error Handling and Exit Codes](006-error-handling.md)
