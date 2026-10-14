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

You'll be prompted for your account number, API key, username and root domain. The file is readable only by you, and
the Smarter SDKs read it too. `--environment` selects which Smarter environment to target (`local`, `alpha`, `beta`,
`next`, `prod`; default is `prod`).

Flags override the file, and so do the environment variables that the SDKs share:

| Variable              | Meaning                                                                                    |
| --------------------- | ------------------------------------------------------------------------------------------ |
| `SMARTER_API_KEY`     | The API key (below `--api_key`)                                                            |
| `SMARTER_ENVIRONMENT` | The environment (below `--environment`)                                                    |
| `SMARTER_BASE_URL`    | The platform's url, e.g. `http://localhost:9357`, instead of the environment's             |
| `SMARTER_CONFIG_FILE` | The configuration file (below `--config`)                                                  |
| `SMARTER_TIMEOUT`     | Seconds to wait for a response (below `--timeout`; default 120)                            |
| `SMARTER_MAX_RETRIES` | Retries of a request that failed for a transient reason (below `--max_retries`; default 2) |
| `SMARTER_LOG`         | `debug` or `info` turns on `--verbose`                                                     |

### Commands

| Command       | Description                                                                               |
| ------------- | ----------------------------------------------------------------------------------------- |
| `configure`   | Set or update your Smarter account credentials                                            |
| `whoami`      | Show the currently configured account                                                     |
| `get`         | List/retrieve Smarter resources                                                           |
| `describe`    | Show a resource's manifest                                                                |
| `apply`       | Create or update a resource from a manifest                                               |
| `validate`    | Check a manifest without applying it (as `apply --dry-run` does); exits 1 if it's invalid |
| `delete`      | Delete a resource                                                                         |
| `deploy`      | Deploy an LLMClient, LLMHost, LLMHostCompute, VectorStore, CustomDomain, API key or user  |
| `undeploy`    | Undo a deployment                                                                         |
| `chat`        | Prompt an LLMClient (`smarter chat prompt --llmclient <name>`)                            |
| `logs`        | View a resource's logs                                                                    |
| `manifest`    | Generate an example manifest for a resource kind                                          |
| `json-schema` | Retrieve the JSON Schema of a resource kind's manifest                                    |
| `resources`   | List the resource kinds of the platform                                                   |
| `status`      | Show service status                                                                       |
| `version`     | Print the cli version                                                                     |

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

Requires Go 1.26 or later (see [go.mod](./go.mod)). Releases are built with Go 1.27.

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

- [build-release.yml](./.github/workflows/build-release.yml) — GitHub Releases, the Homebrew tap, Ubuntu APT / Launchpad PPA, Chocolatey and Docker Hub

Semantic version numbers are controlled by npm package [semantic-release](https://www.npmjs.com/package/semantic-release) which itself is governed by these git [commit comment guidelines](./doc/SEMANTIC_VERSIONING.md).

Package versions for Go lang, NPM and GitHub Actions are monitored by [Dependabot](https://docs.github.com/en/code-security/dependabot) and [Mergify](https://mergify.com/) and are automatically updated and merged to the [alpha branch](https://github.com/smarter-sh/smarter-cli/tree/alpha) of this repo.

### Release

A new version is cut automatically by semantic-release when commits land on `main` (see CI/CD above). To publish a
release manually, run the `Build-Release` GitHub Action from the Actions tab. Its inputs choose the channels to
publish to, so a channel that failed can be published again on its own.

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
go test ./...      # unit tests, the binary tests and the contract test
make coverage      # per-package coverage; fails below 90%, as CI does
make contracts     # refresh testdata/contracts/openapi.json from a checkout of the platform at ../smarter
```

The tests never reach a Smarter platform. The contract test runs every command against a mock platform, and checks its requests against a pinned copy of the platform's OpenAPI document.

### Design decisions and agent skills

The cli's architecture decision records are in [doc/adr/](./doc/adr/000-intro.md). Agent Skills for developing it with an AI coding agent are in [skills/](./skills/README.md); start with [CLAUDE.md](./CLAUDE.md).

## Contributing

Please see [./doc/CONTRIBUTING.md](./doc/CONTRIBUTING.md). By participating, you're expected to uphold our
[Code of Conduct](./CODE_OF_CONDUCT.md).

## License

[GNU Affero General Public License v3.0](./LICENSE)
