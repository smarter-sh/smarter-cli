# Smarter Command-line Interface

The Smarter command-line interface for working with Smarter resources. For example, with it you can create Smarter plugins, skills and guardrails, add
these to a ChatBot, deploy the ChatBot to a custom URL, interact with it from the command line, view chat log data, and manage your Smarter account. Runs on Windows, macOS, Linux and Docker.

Download it at [https://smarter.sh/cli/](https://smarter.sh/cli/)

## Usage

```console
smarter --help
```

### Configuration

Before running most commands, configure the cli with your Smarter account credentials. This writes to
`$HOME/.smarter/config.yaml` (override the location with `--config`):

```console
smarter configure
```

You'll be prompted for your account number, API key, username and root domain. Values can also be supplied as
persistent flags (`--api_key`, `--environment`, ...) or environment variables, and `--environment` selects which
Smarter environment to target (`local`, `alpha`, `beta`, `next`, `prod`; default is `prod`).

### Commands

| Command     | Description                                   |
| ----------- | ---------------------------------------------- |
| `configure` | Set or update your Smarter account credentials |
| `whoami`    | Show the currently configured account          |
| `get`       | List/retrieve Smarter resources                |
| `describe`  | Show details about a resource                  |
| `apply`     | Create or update a resource from a manifest    |
| `delete`    | Delete a resource                              |
| `deploy`    | Deploy a ChatBot                               |
| `undeploy`  | Undeploy a ChatBot                             |
| `chat`      | Interact with a deployed ChatBot               |
| `logs`      | View chat log data                             |
| `manifest`  | Work with resource manifests                   |
| `status`    | Show service status                            |
| `console`   | Open the Smarter web console                   |
| `version`   | Print the cli version                          |

Run `smarter <command> --help` for details and flags on any command.

## Installation

### GitHub Releases

Download the binary for your platform from the [releases page](https://github.com/smarter-sh/smarter-cli/releases), then make it executable and put it on your `PATH`.

```bash
chmod +x ./smarter-*
sudo mv ./smarter-* /usr/local/bin/smarter
```

### Homebrew (macOS / Linux)

```bash
brew tap smarter-sh/tap https://github.com/smarter-sh/homebrew-tap
brew trust smarter-sh/tap
brew install smarter
```

### Ubuntu APT

```bash
sudo add-apt-repository ppa:lpm0073/smarter-cli
sudo apt-get update
sudo apt-get install smarter-cli
```

### Chocolatey (Windows)

```powershell
choco install smarter
```

### Docker Hub

```bash
docker pull mcdaniel0073/smarter-cli:latest
docker run --rm mcdaniel0073/smarter-cli:latest --help
```

## Developers

### Build

Requires Go 1.21.6 or later (see [go.mod](./go.mod)).

#### Windows

```powershell
go get -v -t ./...
$VERSION = Get-Content -Path .\VERSION
$env:VERSION = $VERSION
go build -v -ldflags "-X main.Version=$env:VERSION" -o "./smarter-windows-${env:VERSION}.exe"
```

#### macOS / Linux

```bash
go get -v -t ./...
export VERSION=$(cat VERSION)
go build -v -ldflags "-X main.Version=$VERSION" -o "./smarter-linux-$VERSION"
```

Alternatively, `make build` builds a local `./smarter` binary; run `make help` to see all available developer targets.

### CI/CD

The GitHub Actions workflows in [.github/workflows/](./.github/workflows/) publish semantically-versioned releases to [https://github.com/smarter-sh/smarter-cli/releases](https://github.com/smarter-sh/smarter-cli/releases), covering binaries for Windows, macOS, Linux and Docker:

- [build-release-github.yml](./.github/workflows/build-release-github.yml) — GitHub Releases, and also publishes to Homebrew and Ubuntu APT
- [build-release-dockerhub.yml](./.github/workflows/build-release-dockerhub.yml) — Docker Hub
- [build-release-ubuntu-apt.yml](./.github/workflows/build-release-ubuntu-apt.yml) — Ubuntu APT / Launchpad PPA
- [build-release-windows-chocolatey.yml](./.github/workflows/build-release-windows-chocolatey.yml) — Chocolatey

Semantic version numbers are controlled by npm package [semantic-release](https://www.npmjs.com/package/semantic-release) which itself is governed by these git [commit comment guidelines](./doc/SEMANTIC_VERSIONING.md).

Package versions for Go lang, NPM and GitHub Actions are monitored by [Dependabot](https://docs.github.com/en/code-security/dependabot) and [Mergify](https://mergify.com/) and are automatically updated and merged to the [alpha branch](https://github.com/smarter-sh/smarter-cli/tree/alpha) of this repo.

### Release

A new version is cut automatically by semantic-release when commits land on `main` (see CI/CD above). To publish a
release manually, run the corresponding GitHub Action from the Actions tab: `Build-Release GitHub`,
`Build-Release DockerHub`, `Build-Release Chocolatey`, or `Build-Release Ubuntu APT`.

### Cobra

This cli is built on the [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) frameworks for Go lang. See also:

- [Oscon 2017 - Building An Awesome CLI App In Go](https://spf13.com/presentation/building-an-awesome-cli-app-in-go-oscon/)
- [The Cobra Generator README](https://github.com/spf13/cobra-cli/blob/main/README.md)
- [The Cobra User Guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md).

### Pre-commit

```console
brew install pre-commit
brew install golangci-lint
go install golang.org/x/tools/cmd/goimports@latest

pre-commit install
pre-commit autoupdate
pre-commit run --all-files
```

### Tests

```console
go test ./... -v 2>&1 | tail -30
```

## Contributing

Please see [./doc/CONTRIBUTING.md](./doc/CONTRIBUTING.md). By participating, you're expected to uphold our
[Code of Conduct](./CODE_OF_CONDUCT.md).

## License

[GNU Affero General Public License v3.0](./LICENSE)
