/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose shape (custom flags, non-"name" identifying
// parameter, or organic Short/Long text) pre-dates the v0.14 ResourceKind
// template and doesn't fit it, so they're kept as explicit ResourceSpecs.
//
// The ChatBot and chat history kinds were removed: the platform's cli api has
// no broker for them since v0.14 (LLMClient replaced ChatBot).
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:     "account",
		Short:   "Retrieve your Account manifest",
		Long:    "Retrieve your Account manifest:\n\nsmarter get account [flags]\n\nThe Smarter API will your Account manifest.",
		APIKind: "Account",
	},
	{
		Use:   "plugins",
		Short: "Retrieve a list of Plugins",
		Long: `Retrieves a list of Plugins:

smarter get plugins [flags]


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
		APIKind: "Plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "name", Shorthand: "n", Usage: "Name of the plugin"},
		Flags: []cmd.FlagSpec{
			{Name: "class", Shorthand: "c", Kind: cmd.FlagString, Usage: "Plugin class: static, sql, api", Deprecated: "the Smarter API ignores it; use the kind's own command, e.g. get sqlplugins or get apiplugins"},
		},
	},
	{
		Use:   "users",
		Short: "Retrieve a list of Users",
		Long: `Retrieves a list of Users:

smarter get users [flags]

The Smarter API will return a list of Users in the specified format,
or a manifest for a specific User.`,
		APIKind: "User",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "username", Param: "name", Shorthand: "u", Usage: "Smarter username"},
	},
}

func init() {
	cmd.RegisterResources(getCmd, legacySpecs, APIRequest, ConsoleOutput, cmd.ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.ResourceKinds))
	for i, k := range cmd.ResourceKinds {
		specs[i] = k.GetSpec()
	}
	cmd.RegisterResources(getCmd, specs, APIRequest, ConsoleOutput, cmd.ErrorOutput)
}
