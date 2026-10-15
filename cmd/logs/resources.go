/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package logs

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text pre-dates the v0.14
// ResourceKind template.
//
// The ChatBot and chat history kinds were removed: the platform's cli api has
// no broker for them since v0.14 (LLMClient replaced ChatBot).
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "plugin <name>",
		Short: "Returns the logs for a Plugin",
		Long: `Returns the logs for a Plugin:

smarter logs plugin <name>

The Smarter API will return the logs for the Plugin with the specified name.`,
		APIKind: "Plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
}

func init() {
	cmd.RegisterResources(logsCmd, legacySpecs, APIRequest, ConsoleOutput, cmd.ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.ResourceKinds))
	for i, k := range cmd.ResourceKinds {
		specs[i] = k.LogsSpec()
	}
	cmd.RegisterResources(logsCmd, specs, APIRequest, ConsoleOutput, cmd.ErrorOutput)
}
