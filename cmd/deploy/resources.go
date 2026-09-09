/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package deploy

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text pre-dates the v0.14
// ResourceKind template and is preserved verbatim rather than regenerated.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "chatbot <name>",
		Short: "Deploy a ChatBot",
		Long: `Deploys a ChatBot:

smarter deploy chatbot <name> [flags]

The Smarter API will deploy the ChatBot.`,
		APIKind: "ChatBot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
}

func init() {
	cmd.RegisterResources(deployCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)

	// Only kinds flagged Deployable (e.g. Prompt, which replaces
	// ChatBot as the deployable unit) get a "deploy" leaf command. Pulled
	// from the shared registry so its Use/APIKind/display text stays in
	// lockstep with its describe/get/manifest counterparts.
	var specs []cmd.ResourceSpec
	for _, k := range cmd.ResourceKinds {
		if k.Deployable {
			specs = append(specs, k.DeploySpec())
		}
	}
	cmd.RegisterResources(deployCmd, specs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)
}
