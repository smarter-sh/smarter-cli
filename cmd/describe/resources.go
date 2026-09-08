/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package describe

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text pre-dates the v0.14
// ResourceKind template and is preserved verbatim rather than regenerated.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "account",
		Short: "Retrieve your Account manifest",
		Long: `Retrieves a manifest of your account. For example:

	smarter describe account > my-plugin.yaml

This will generate an example manifest for an account and write it to my-plugin.yaml in the current working directory.`,
		APIKind: "account",
	},
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
}

func init() {
	cmd.RegisterResources(describeCmd, legacySpecs, APIRequest, ConsoleOutput, ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.V14ResourceKinds))
	for i, k := range cmd.V14ResourceKinds {
		specs[i] = k.DescribeSpec()
	}
	cmd.RegisterResources(describeCmd, specs, APIRequest, ConsoleOutput, ErrorOutput)
}
