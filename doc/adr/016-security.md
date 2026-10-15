# ADR-016: Security

## Status

Accepted

The SDKs' ADR-016 adapt the platform's [ADR-016: Security](https://github.com/smarter-sh/smarter/blob/main/docs/source/adr/016-security.rst).

## Context

The cli holds api keys that can change production resources, and `deploy` creates billable infrastructure (DNS records, certificates, Kubernetes ingresses). The repository is public. Its release workflows hold credentials for GitHub, Docker Hub, Launchpad and Chocolatey.

## Decision

- **Keys** are stored in `~/.smarter/config.yaml`, per environment (ADR-002), readable only by its owner, and sent only in the `Authorization` header over https (the `local` environment uses http on localhost). `configure` masks them when it echoes them, and `--verbose` shows only their last four characters.
- **Tests never reach a real platform** (ADR-009), and use a placeholder key. Nothing in the repository, its tests or its docs holds a real key, account number, private hostname or customer data.
- **Deploy is explicit**: `deploy` and `delete` act on the named resource only, with no bulk form. Agents and tests never run them against a real platform unless the user asks.
- **Release secrets** live in GitHub Actions secrets, never in the repository; each workflow's header documents its secrets and their rotation. `.env` (git-ignored) is for local development only.
- The pre-commit hooks detect private keys and AWS credentials (ADR-014).

## Alternatives Considered

- The operating system's keychain for api keys. Better at rest; the SDKs read the plain file (ADR-002), so moving keys would break them.

## Consequences

- **Positive:**
  - A leaked repository or test log holds no credentials.
- **Negative:**
  - Keys are stored in plain text, protected only by the file's mode: `writeConfigFile` writes it `0600`, in a directory created `0700`, and restricts an older file the next time it writes it. Until October 2026 the cli wrote it `0644`, and `--verbose` printed the full key to stdout (ADR-007).

## Related ADRs

- [ADR-000: Introduction](000-intro.md)
- [ADR-002: Configuration](002-configuration.md)
- [ADR-007: Logging and Verbose Output](007-logging.md)
- [ADR-009: Testing](009-testing.md)
