---
name: smarter-git-commits
description: Use whenever the user asks you to "write a git commit" (or a commit message, or the git add) in the smarter-cli repo, or in ../smarter or an SDK while working on the cli. Covers the never-execute rule, authorship, the multiple -m command format, conventional-commit types and the cli's scopes, how commits drive semantic-release, and how thorough the body must be.
---

# Smarter Git Commits

## The rule: write, never execute

"Write a git commit", in any wording or tone, means **output the command as
text**. It never means run it.

- Never execute `git commit`, `git add`, `git reset`, `git restore --staged`,
  `git push`, `git stash`, `git tag`, or any other git command that writes, and
  never `make release`, which commits and pushes. The user stages and commits
  themselves.
- Read-only git commands (`git status`, `git diff`, `git log`) are fine for
  working out what changed.
- Say which release of the platform (https://github.com/smarter-sh/smarter/releases/)
  the change was made against, when it concerns the platform's routes or kinds
  (e.g. in the closing status paragraph).
- Don't volunteer commit messages. Write one only when asked.

## Authorship

You are the author. Lawrence McDaniel is the co-author.

- `--author="Claude <noreply@anthropic.com>"`
- The last `-m` is `Co-authored-by: Lawrence McDaniel <lpm0073@gmail.com>`.

## Command format

Exactly one `git commit` command in one code block:

```console
git commit --author="Claude <noreply@anthropic.com>" \
  -m "<type>(<scope>): <subject>" \
  -m "<body paragraph>" \
  -m "<body paragraph or bullet list>" \
  -m "Co-authored-by: Lawrence McDaniel <lpm0073@gmail.com>"
```

- One `-m` per paragraph; a bullet list goes in a single `-m` with real
  newlines between bullets. git joins `-m`s with blank lines.
- No heredoc, no `-F`, and no `git add` in the same block.
- Avoid backticks, `$` and `!` inside the double-quoted messages. The user's
  shell is zsh, where `!` in double quotes is history expansion and aborts the
  command. Reword instead.
- When the user asks for the `git add` too, write it as a separate code block
  that lists paths explicitly, and leaves out files that aren't part of the
  change (the user often has unrelated work in progress here).

## Subject line

`<type>(<scope>): <subject>`, checked by commitlint
(`@commitlint/config-conventional` + `@commitlint/config-angular`) and read by
semantic-release (`release.config.js`, ADR-015) to version the release.

- **type:** `feat`, `fix`, `test`, `docs`, `refactor`, `perf`, `style`,
  `build`, `ci`, or `chore`.
  - `feat` (a new command, kind, flag or output) makes a minor release, and
    `fix` a patch. Every user-visible change needs one of them, or the change
    is never released.
  - Adding the platform's new kinds to `ResourceKinds` is `feat`, as
    `b65fcb5 feat: add new v0.16 resources` was.
  - Skills, ADRs and docs are `docs`; packaging and workflows are `build` or `ci`.
  - Work the user asked for as a list of bugs is `fix`.
- **scope:** the area, in lower case: a verb (`get`, `apply`, `deploy`,
  `chat`, ...), `config`, `api`, `output`, `kinds`, `tests`, `debian`,
  `chocolatey`, `docker`, `ci`, `deps`, `contracts`, `skills`, `adr`. Omit it for changes
  that span the repository. In `../smarter`, the scope is the platform's
  Django app, as its own copy of this skill says.
- **subject:** imperative mood, lower case, no trailing period.
- **length:** the whole header must be **72 characters or fewer**. Measure it,
  and aim for 70:

  ```console
  printf %s "feat(kinds): add customdomain and imagesearchplugin resources" | wc -c
  ```

- `chore(release): ...` belongs to semantic-release. Never write it yourself.

## Body

Write the body for a reader with no context, such as a reviewer six months from
now. A terse four-bullet body has been rejected as "too brief".

- Start with the problem or motivation: what was wrong or missing, and how it
  showed up (e.g. "smarter manifest secret returned 404").
- Then go area by area, with a short heading paragraph per area (for example
  `Kind registry (cmd/resourcekinds.go):`), followed by one bullet per change
  saying what changed and why.
- Name files, functions, flags and commands precisely.
- Cover the tests: what's new, and what they prove.
- End with the current status: known skips, pre-existing failures, platform
  drift that the contract test still reports, follow-ups.

## Before writing

1. Run `git status` and `git diff --stat` (and `git diff` on anything unclear),
   so the message describes the actual working tree.
2. Ask, or check, which files belong to the commit, if that's ambiguous.
3. Measure the header with `wc -c`.
4. Start the response with a one-line `Written for:` audience note.

See [references/canonical-git-commit.md](references/canonical-git-commit.md)
for a complete, real example. It is a commit of the platform, but the format is
the same here.

A change that spans repositories (e.g. a cli command, the docs page in
`../smarter/docs/source/`, and the SDKs' config readers) gets one commit
command per repository.
