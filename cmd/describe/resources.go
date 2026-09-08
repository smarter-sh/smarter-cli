/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package describe

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"
)

// spec builds the common "describe <use> <name>" shape shared by most
// resources: a required positional name argument and a templated
// Short/Long.
func spec(use, kind, display string) cmd.ResourceSpec {
	a := cmd.Article(display)
	return cmd.ResourceSpec{
		Use:   use + " <name>",
		Short: fmt.Sprintf("Retrieve %s %s manifest by name", a, display),
		Long: fmt.Sprintf(`Retrieves a manifest for %s %s. For example:

	smarter describe %s <name> > my-plugin.yaml

This will generate a manifest for %s %s named <name> and write it to my-plugin.yaml in the current working directory.`, a, display, use, a, display),
		APIKind: kind,
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	}
}

var resourceSpecs = []cmd.ResourceSpec{
	{
		Use:   "account",
		Short: "Retrieve your Account manifest",
		Long: `Retrieves a manifest of your account. For example:

	smarter describe account > my-plugin.yaml

This will generate an example manifest for an account and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "account",
	},
	spec("apikey", "SmarterAuthToken", "SmarterAuthToken"),
	{
		Use:   "chatbot <name>",
		Short: "Retrieve a ChatBot manifest by name",
		Long: `Retrieves a manifest for a chatbot. For example:

	smarter describe chatbot <name> > my-plugin.yaml

This will generate a manifest for a chatbot named <name> and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "chatbot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "plugin <name>",
		Short: "Retrieve a manifest for a Plugin",
		Long: `Retrieves a manifest for a plugin. For example:

	smarter describe plugin <name> > my-plugin.yaml

This will retrieve the manifest for a plugin named <name> and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "user <username>",
		Short: "Retrieve a manifest for a User",
		Long: `Retrieves a manifest for a user. For example:

	smarter describe user <username> > my-plugin.yaml

This will retrieve a manifest for User <username> and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "user",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "username"},
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
	spec("secret", "secret", "Secret"),
	spec("vectorstore", "Vectorstore", "Vectorstore"),
}

func init() {
	for _, s := range resourceSpecs {
		cmd.RegisterResourceCmd(describeCmd, s, APIRequest, ConsoleOutput, ErrorOutput)
	}
}
