/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package undeploy

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are the kinds whose commands pre-date the ResourceKind template:
// User, which takes a username. Only account admins may undeploy users.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "user <username>",
		Short: "Deactivate a user of your account",
		Long: `Deactivates a user of your account, who can then no longer log in or use the api:

smarter undeploy user <username>`,
		APIKind: "User",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "username", Param: "name"},
	},
}

func init() {
	cmd.RegisterResources(undeployCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)

	// Only kinds flagged Deployable, whose platform broker implements deploy,
	// get a leaf command. Pulled from the shared registry so its
	// Use/APIKind/display text stays in lockstep with the other verbs. The
	// ChatBot command was removed: the platform has no ChatBot broker since
	// v0.14 (LLMClient replaced it).
	var specs []cmd.ResourceSpec
	for _, k := range cmd.ResourceKinds {
		if k.Deployable {
			specs = append(specs, k.UndeploySpec())
		}
	}
	cmd.RegisterResources(undeployCmd, specs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)
}
