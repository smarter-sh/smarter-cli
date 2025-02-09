# Smarter Command-line Interface

The Smarter command-line interface for working with Smarter resources. Runs on Windows, macOS, Linux and DockerHub.
Download it at [https://smarter.sh/cli/](https://smarter.sh/cli/)

## Usage

```console
smarter --help
```

## Developers

### Build

#### Windows

```powershell
go get -v -t .
$VERSION = Get-Content -Path .\VERSION
$env:VERSION = $VERSION
go build -v -o smarter main.go -ldflags "-X main.Version=$env:VERSION" -o "./smarter-windows-${env:VERSION}.exe"
```

#### macOS / Linux

```bash
go get -v -t .
export VERSION=$(cat VERSION)
go build -v -o smarter main.go -ldflags "-X main.Version=$VERSION" -o "./smarter-linux-$VERSION"
```

### CI/CD

The GitHub Actions workflow [.github/workflows/build.yml](./.github/workflows/build.yml) publishes semantically-versioned releases to [https://github.com/smarter-sh/smarter-cli/releases](https://github.com/smarter-sh/smarter-cli/releases) which includes binaries for Windows, macOS, Linux and Docker.

Semantic version numbers are controlled by npm package [semantic-release](https://www.npmjs.com/package/semantic-release) which itself is governed by these git [commit comment guidelines](./doc/SEMANTIC_VERSIONING.md).

Package versions for Go lang, NPM and GitHub Actions are monitored by [Dependabot](https://docs.github.com/en/code-security/dependabot) and [Mergify](https://mergify.com/) and are automatically updated and merged to the [alpha branch](https://github.com/smarter-sh/smarter-cli/tree/alpha) of this repo.

### Cobra

This cli is built on the [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) frameworks for Go lang. See also:

- [Oscon 2017 - Building An Awesome CLI App In Go](https://spf13.com/presentation/building-an-awesome-cli-app-in-go-oscon/)
- [The Cobra Generator README](https://github.com/spf13/cobra-cli/blob/main/README.md)
- [The Cobra User Guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md).

## Contributing

Please see [./doc/CONTRIBUTING.md](./doc/CONTRIBUTING.md)
