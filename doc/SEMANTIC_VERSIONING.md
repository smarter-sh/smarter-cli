# Semantic Versioning Guide

See [Conventional Commits](https://www.conventionalcommits.org/)

The [pushMain.yml](../.github/workflows/pushMain.yml) workflow runs [semantic-release](https://semantic-release.gitbook.io/), configured in [release.config.js](../release.config.js), on every push to `main`. semantic-release reads the commit messages since the last tag and chooses the next version from them: `fix` and `perf` make a patch release, `feat` makes a minor release, and a `BREAKING CHANGE:` footer makes a major release. The other types make no release. It then updates `CHANGELOG.md` and `VERSION`, tags the release and publishes it on GitHub. The commitlint `commit-msg` hook in [.pre-commit-config.yaml](../.pre-commit-config.yaml) rejects a message that doesn't follow the format below. See [ADR-015](adr/015-versioning.md).

Work lands on `alpha` first, and reaches `main` through a pull request. Each commit message should be structured as follows:

> `<type>`: This represents the type of change made in the commit. Common types include feat (for a new feature), fix (for a bug fix), chore (for routine tasks like updating dependencies), docs (for documentation changes), style (for code style changes), refactor (for refactoring existing code), test (for adding or updating tests), and perf (for performance improvements).
>
> `[optional scope]`: This is an optional part that provides additional contextual information, like the part of the codebase the commit modifies.
>
> `<description>`: This is a brief description of the changes the commit makes.
>
> `[optional body]`: This is an optional part where you can provide a more detailed explanation of the changes.
>
> `[optional footer(s)]`: This is also optional and is often used to reference issue tracker IDs.

## Commit Message format

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries such as documentation generation

Example

```console
git commit -m "docs: add an example of a good commit message"
```
