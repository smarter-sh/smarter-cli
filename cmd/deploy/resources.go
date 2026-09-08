/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package deploy

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

var resourceSpecs = []cmd.ResourceSpec{
	{
		Use:   "chatbot <name>",
		Short: "Deploy a ChatBot",
		Long: `Deploys a ChatBot:

smarter deploy chatbot <name> [flags]

The Smarter API will deploy the ChatBot.`,
		APIKind: "ChatBot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		// v0.14 resource. Prompt replaces ChatBot as the deployable unit;
		// the "Prompt" API kind string is a best-guess placeholder pending
		// backend confirmation.
		Use:   "prompt <name>",
		Short: "Deploy a Prompt",
		Long: `Deploys a Prompt:

smarter deploy prompt <name> [flags]

The Smarter API will deploy the Prompt.`,
		APIKind: "Prompt",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
}

// deployOutput adapts deploy's no-argument ConsoleOutput to the
// cmd.OutputFunc signature shared by RegisterResourceCmd.
func deployOutput(_ []byte) {
	ConsoleOutput()
}

func init() {
	for _, s := range resourceSpecs {
		cmd.RegisterResourceCmd(deployCmd, s, APIRequest, deployOutput, ErrorOutput)
	}
}
