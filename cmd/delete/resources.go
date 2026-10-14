/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package delete

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text or identifying parameter
// pre-dates the v0.14 ResourceKind template.
//
// The ChatBot and chat history kinds were removed: the platform's cli api has
// no broker for them since v0.14 (LLMClient replaced ChatBot).
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "plugin <name>",
		Short: "Delete a Plugin",
		Long: `Delete a Plugin:

smarter delete plugin <name>

The Smarter API will permanently delete the Plugin with the specified name,
and disassociate it from any ChatBots.`,
		APIKind: "plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "user <username>",
		Short: "Delete a user from your account",
		Long: `Delete a user from your account:

smarter delete user <username>

The Smarter API will permanently delete the user with the specified name,
and disassociate it from any Smarter resources. Your Smarter admin account
will replace the deleted user.`,
		APIKind: "user",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "username", Param: "name"},
	},
}

func init() {
	cmd.RegisterResources(deleteCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.ResourceKinds))
	for i, k := range cmd.ResourceKinds {
		specs[i] = k.DeleteSpec()
	}
	cmd.RegisterResources(deleteCmd, specs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)
}
