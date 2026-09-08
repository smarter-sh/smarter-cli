/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text pre-dates the v0.14
// ResourceKind template and is preserved verbatim rather than regenerated.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "account",
		Short: "Retrieve your Account manifest",
		Long: `Generate an example manifest for an Account. For example:

	smarter manifest account [flags] > my-plugin.yaml

This will generate an example manifest for an Account and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "Account",
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
		Use:   "user [flags]",
		Short: "Generate an example manifest for a user.",
		Long: `Generate an example manifest for a user. For example:

	smarter manifest user [flags] > my-plugin.yaml

This will generate an example manifest for a user and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "user",
	},
}

func init() {
	cmd.RegisterResources(manifestCmd, legacySpecs, APIRequest, ConsoleOutput, ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.ResourceKinds))
	for i, k := range cmd.ResourceKinds {
		specs[i] = k.ManifestSpec()
	}
	cmd.RegisterResources(manifestCmd, specs, APIRequest, ConsoleOutput, ErrorOutput)
}
