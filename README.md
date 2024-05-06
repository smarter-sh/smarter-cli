# Smarter Command-line Interface

The Smarter command-line interface for working with Smarter resources. Available on Windows, macOS, Linux and DockerHub.

## Usage

```console
smarter --help
```

## CI/CD

The GitHub Actions workflow [.github/workflows/build.yml](./.github/workflows/build.yml) publishes semantically-versioned releases to [https://github.com/QueriumCorp/smarter-cli/releases](https://github.com/QueriumCorp/smarter-cli/releases) which includes binaries for Windows, macOS, Linux and Docker.

Semantic version numbers are controlled by npm package [semantic-release](https://www.npmjs.com/package/semantic-release) which itself is governed by these git [commit comment guidelines](./doc/SEMANTIC_VERSIONING.md).

Package versions for Go lang, NPM and GitHub Actions are monitored by [Dependabot](https://docs.github.com/en/code-security/dependabot) and [Mergify](https://mergify.com/) and are automatically updated and merged to the [alpha branch](https://github.com/QueriumCorp/smarter-cli/tree/alpha) of this repo.

## Cobra

This cli is built on the [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) frameworks for Go lang. See also:

- [Oscon 2017 - Building An Awesome CLI App In Go](https://spf13.com/presentation/building-an-awesome-cli-app-in-go-oscon/)
- [The Cobra Generator README](https://github.com/spf13/cobra-cli/blob/main/README.md)
- [The Cobra User Guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md).

## Contributing

Please see [./doc/CONTRIBUTING.md](./doc/CONTRIBUTING.md)
