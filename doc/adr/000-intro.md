# ADR Introduction

Architectural Design Records (ADRs) are short documents that capture a single architectural or design decision, its context, and its consequences. Less is more: keep them brief.

These are the ADRs of the smarter cli, the Smarter platform's command-line interface. They follow the Smarter platform's own [ADRs](https://github.com/smarter-sh/smarter/tree/main/docs/source/adr), and sit beside the ADRs of the Smarter SDKs ([Python](https://github.com/smarter-sh/smarter-python/tree/main/doc/adr), [Node.js](https://github.com/smarter-sh/smarter-nodejs/tree/main/doc/adr), [Go](https://github.com/smarter-sh/smarter-golang/tree/main/doc/adr)). Where the cli and the SDKs make the same decision (scope, configuration, error handling, logging, testing, CI/CD, code quality, versioning, security), the cli's ADR takes the SDKs' number, so that the same subject has the same number in every repository. The cli is older than the SDKs, and several of their decisions were taken from it: the configuration file (ADR-002), the Go lint hooks (ADR-014), and the Go toolchain (ADR-019).

The platform's decisions about its own runtime (Django, Docker, Kubernetes, the database, Celery, React, templates, feature flags, internationalization) don't apply to a cli and are not repeated here.

These ADRs describe the cli as it is. Where the code doesn't yet meet a decision, the ADR says so under **Consequences**, rather than describing an ideal.

See: [Michael Nygard's ADRs](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html)

## Guidelines for writing ADRs

- **One decision per ADR.**
- **Be concise.** Focus on clarity and brevity.
- **State the context:** the forces that led to the decision.
- **List the alternatives**, and why they weren't chosen.
- **Explain the decision**, and **describe its consequences**, positive and negative.
- **Use simple language**, for future contributors without today's context.
- **Use the same format** as the others, and **reference related ADRs**.
- **Keep the numbering aligned** with the SDKs: a subject that the SDKs also decide takes their number. A decision made later that the SDKs and the cli share takes the next number free in every repository, which is why 020 and 021 (language-specific ADRs of the Node.js and Go SDKs) are unused here.

## Index

- [ADR-001: Scope and Purpose of the CLI](001-scope.md)
- [ADR-002: Configuration](002-configuration.md)
- [ADR-003: Only the CLI API](003-cli-api-only.md)
- [ADR-004: Cobra and Viper](004-cobra-and-viper.md)
- [ADR-005: The Resource Kind Registry](005-resource-kind-registry.md)
- [ADR-006: Error Handling and Exit Codes](006-error-handling.md)
- [ADR-007: Logging and Verbose Output](007-logging.md)
- [ADR-008: HTTP Transport](008-http-transport.md)
- [ADR-009: Testing](009-testing.md)
- [ADR-010: CI/CD](010-ci-cd.md)
- [ADR-011: Dependency Management](011-dependencies.md)
- [ADR-012: Help Text and Documentation](012-documentation.md)
- [ADR-013: Output Formats](013-output-formats.md)
- [ADR-014: Code Quality](014-code-quality.md)
- [ADR-015: Semantic Versioning](015-versioning.md)
- [ADR-016: Security](016-security.md)
- [ADR-017: Distribution Channels](017-distribution.md)
- [ADR-018: One Package per Verb](018-package-per-verb.md)
- [ADR-019: Go Version](019-go-version.md)
- [ADR-022: A Pinned Copy of the Platform's OpenAPI Document](022-pinned-openapi.md)
