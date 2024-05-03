# Smarter Command-line Interface

A Go lang / Cobra command-line interface for working with Smarter resources.

## Usage

### Apply

Apply a Smarter API resource manifest. See this [example](./data/manifests/example-configuration.yaml)

```console
smarter apply <./file/path/manifest.yaml>
```

### Get

```console
smarter get plugins --name --yaml --json
smarter get chatbots --name --yaml --json
smarter get chats --id --yaml --json
smarter get account --yaml --json
smarter get account users --username --yaml --json
```

### Status

```console
smarter status --yaml --json
```

### Deploy

```console
smarter deploy chatbot --name --yaml --json
```

### Logs

```console
smarter logs <kind> --name --yaml --json
```

### Delete

```console
smarter delete plugin --name --yaml --json
smarter delete chatbot --name --yaml --json
smarter delete chat --id --yaml --json
smarter delete user --username --yaml --json
```

scaffolding example

```console
cobra-cli add get
cobra-cli add plugins -p 'GetCmd'
```

## Project Startup

```console
go mod init github.com/QueriumCorp/smarter-cli
go mod tidy
go get -u github.com/spf13/cobra@latest
go install github.com/spf13/cobra-cli@latest
```

Cobra initialization

```console
cobra-cli init --author "Lawrence McDaniel lawrence@querium.com" --viper
```

To run in development environment

```console
go run main.go
```

## Cobra

See:

- [Oscon 2017 - Building An Awesome CLI App In Go](https://spf13.com/presentation/building-an-awesome-cli-app-in-go-oscon/)
- [https://github.com/spf13/cobra](https://github.com/spf13/cobra)

For complete details on using the Cobra-CLI generator, please read [The Cobra Generator README](https://github.com/spf13/cobra-cli/blob/main/README.md)

For complete details on using the Cobra library, please read the [The Cobra User Guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md).


## Smarter API

- https://api.smarter.sh/v0/cli/describe/: print the manifest
- https://api.smarter.sh/v0/cli/apply/: Apply a manifest
- https://api.smarter.sh/v0/cli/status/: Smarter platform status
- https://api.smarter.sh/v0/cli/deploy/: Deploy a resource

- https://api.smarter.sh/v0/cli/logs/: Get logs for a resource
- https://api.smarter.sh/v0/cli/delete/: Delete a resource
