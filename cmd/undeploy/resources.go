/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package undeploy

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose Short/Long text pre-dates the v0.14
// ResourceKind template and is preserved verbatim rather than regenerated.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "chatbot <name>",
		Short: "Undo a ChatBot deployment.",
		Long: `Undo a ChatBot deployment. For example:

smarter undeploy chatbot <name>

This will reverse the effect of having deployed the ChatBot.`,
		APIKind: "ChatBot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
}

func init() {
	cmd.RegisterResources(undeployCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), ErrorOutput)

	// v0.14: only kinds flagged Deployable (e.g. Prompt, which replaces
	// ChatBot as the deployable unit) get an "undeploy" leaf command. Pulled
	// from the shared registry so its Use/APIKind/display text stays in
	// lockstep with its describe/get/manifest counterparts.
	var specs []cmd.ResourceSpec
	for _, k := range cmd.ResourceKinds {
		if k.Deployable {
			specs = append(specs, k.UndeploySpec())
		}
	}
	cmd.RegisterResources(undeployCmd, specs, APIRequest, cmd.AdaptOutput(ConsoleOutput), ErrorOutput)
}
