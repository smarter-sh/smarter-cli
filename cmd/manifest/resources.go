/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"
)

// spec builds the common "manifest <use> [flags]" shape shared by most
// resources: no identifying parameter, just a templated Short/Long.
func spec(use, kind, display string) cmd.ResourceSpec {
	a := cmd.Article(display)
	return cmd.ResourceSpec{
		Use:   use + " [flags]",
		Short: fmt.Sprintf("Generate an example manifest for %s %s.", a, display),
		Long: fmt.Sprintf(`Generates an example manifest for %s %s. For example:

	smarter manifest %s [flags] > my-plugin.yaml

This will generate an example manifest for %s %s and write it to my-plugin.yaml in the current working directory.`, a, display, use, a, display),
		APIKind: kind,
	}
}

var resourceSpecs = []cmd.ResourceSpec{
	{
		Use:   "account",
		Short: "Retrieve your Account manifest",
		Long: `Generate an example manifest for an Account. For example:

	smarter manifest account [flags] > my-plugin.yaml

This will generate an example manifest for an Account and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "Account",
	},
	{
		Use:   "apikey",
		Short: "Generate an example manifest for a SmarterAuthToken.",
		Long: `Generates an example manifest for a SmarterAuthToken. For example:

	smarter manifest apikey [flags] > my-plugin.yaml

This will generate an example manifest a SmarterAuthToken and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "SmarterAuthToken",
	},
	{
		Use:   "chat [flags]",
		Short: "Generate an example manifest for a chat session.",
		Long: `Generates an example manifest for a chat session. For example:

	smarter manifest chat [flags] > my-plugin.yaml

This will generate an example manifest a chat session and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "chat",
	},
	{
		Use:   "chatbot [flags]",
		Short: "Generate an example manifest for a ChatBot",
		Long: `Generates an example manifest for a ChatBot resource. For example:

	smarter manifest chatbot [flags] > my-plugin.yaml

This will generate an example manifest for a chatbot and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "chatbot",
	},
	{
		Use:   "plugin [flags]",
		Short: "Generate an example manifest for a plugin.",
		Long: `Generates an example manifest for a plugin. For example:

	smarter manifest plugin [flags] > my-plugin.yaml

This will generate an example manifest for a plugin and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "plugin",
	},
	{
		Use:   "secret [flags]",
		Short: "Generate an example manifest for a secret.",
		Long: `Generates an example manifest for a secret. For example:

	smarter manifest secret [flags] > my-secret.yaml

This will generate an example manifest for a secret and write it to my-secret.yaml in the current working directory.`,
		APIKind: "secret",
	},
	{
		Use:   "user [flags]",
		Short: "Generate an example manifest for a user.",
		Long: `Generate an example manifest for a user. For example:

	smarter manifest user [flags] > my-plugin.yaml

This will generate an example manifest for a user and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "user",
	},

	// v0.14 resources. apiconnection/sqlconnection supersede the old
	// PluginDataApiConnection/PluginDataSqlConnection kinds; the API kind
	// strings below are best-guess placeholders pending backend confirmation.
	spec("apiconnection", "ApiConnection", "ApiConnection"),
	spec("sqlconnection", "SqlConnection", "SqlConnection"),
	spec("apiplugin", "ApiPlugin", "ApiPlugin"),
	spec("sqlplugin", "SqlPlugin", "SqlPlugin"),
	spec("llmclient", "LlmClient", "LlmClient"),
	spec("prompt", "Prompt", "Prompt"),
	spec("promptconfig", "PromptConfig", "PromptConfig"),
	spec("provider", "Provider", "Provider"),
	spec("proxy", "Proxy", "Proxy"),
	spec("vectorstore", "Vectorstore", "Vectorstore"),
}

func init() {
	for _, s := range resourceSpecs {
		cmd.RegisterResourceCmd(manifestCmd, s, APIRequest, ConsoleOutput, ErrorOutput)
	}
}
