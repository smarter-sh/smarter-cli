# smarter-cli Agent Skills

This folder holds the [Agent Skills](https://agentskills.io/specification) for
working on the Smarter command-line interface with an AI coding agent such as
Claude Code. This README is written mostly for people.

A skill is a folder of written know-how that an agent loads only when a task
needs it. Each one contains a `SKILL.md` file, and sometimes `references/` and
`scripts/` folders. These skills cover how this repository is built, tested,
released and committed. They are developer tooling, like the `Makefile` and
`.pre-commit-config.yaml`, and the cli never loads them.

They are adapted from the Smarter platform's skills, in
[smarter-sh/smarter/skills](https://github.com/smarter-sh/smarter/tree/main/skills),
and the Go SDK's, in
[smarter-sh/smarter-golang/skills](https://github.com/smarter-sh/smarter-golang/tree/main/skills),
so that the platform, its SDKs and its cli are developed the same way. The
decisions behind them are recorded in [`../doc/adr/`](../doc/adr/000-intro.md).
Load the platform's skills from `../smarter/skills/` when a cli task leads you
into the platform, e.g. to fix a route or its docs page.

## What's here

Start with `smarter-cli-development`. It sets out the rules that apply to every
change, the cli's architecture, and which narrower skill to load for each kind
of task.

| Skill                     | Use it when                                                                       |
| ------------------------- | --------------------------------------------------------------------------------- |
| `smarter-cli-development` | You're making any change. Start here.                                             |
| `smarter-cli-testing`     | You're writing or running the tests, or raising coverage.                         |
| `smarter-api-testing`     | You're running the cli or curl against the local platform, e.g. to confirm a fix. |
| `smarter-code-quality`    | You're about to finish a change and need to run the hooks.                        |
| `smarter-git-commits`     | You've been asked to write a commit message.                                      |
| `smarter-cli-release`     | You're touching versions, packaging or the build-release workflows.               |
| `monthly-maintenance`     | It's time to sync the cli with the platform's kinds and routes, and update deps.  |
| `smarter-contract-sync` | The platform's api contract changed, and this repository's copy must follow it now. |

The skills bundle no scripts. The repository's own `scripts/sync-contracts.sh`
(`make contracts`) refreshes the pinned copy of the platform's OpenAPI
document, which the contract test checks every command against (ADR-022).

## Claude Code

Claude Code looks for project skills in `.claude/skills/`. This repository
connects the two folders with a committed symlink, `.claude/skills -> ../skills`,
as the platform's and the SDKs' repositories do:

```bash
mkdir -p .claude
ln -s ../skills .claude/skills
```

Other agents that read `SKILL.md` folders can be pointed at this folder, or told
to read `skills/smarter-cli-development/SKILL.md` at the start of a session.

## When to add or change a skill

Add to a skill when you've explained the same thing to an agent twice, when a
mistake is costly (e.g. a deploy against a real platform, or a release
triggered by accident), when a pattern isn't obvious from the code, or when
checking a change takes a fragile series of steps that a script could run.
Don't add one for what the code, the linters or CI already enforce, for
one-off facts, or for personal preferences. Prefer editing an existing skill to
adding one, and keep each `SKILL.md` short.

When the platform or the SDKs change a convention that a skill here copies,
update the skill here too. When the cli fixes one of the rough edges a skill
or ADR describes, update both.

## Security considerations

A skill is a set of instructions that an agent follows, often without asking.
**Review a change to this folder as carefully as a change to CI.**

- **The cli changes real resources.** Its api keys can delete and deploy
  production resources, and the local platform's containers hold live cloud
  credentials. The skills confine agents to the mock server and the local
  platform, with a throwaway `$HOME` and key, and forbid `deploy`. Don't
  weaken that.
- **The release workflows publish to five channels** with maintainers'
  credentials. The skills forbid agents from triggering them. Keep it so.
- **Some instructions forbid things on purpose**, e.g. never running
  `git commit` or `make release`. Don't remove a rule like that unless you mean
  to give agents that power.
- **Skills can carry prompt injection.** Look for instructions unrelated to a
  skill's purpose, new scripts with network calls, and text that tells an agent
  to skip confirmations, hooks or safety rules.
- **Keep secrets out.** This repository is public: no api keys, tokens, account
  numbers, private hostnames or customer data in a skill, an example or a script.
