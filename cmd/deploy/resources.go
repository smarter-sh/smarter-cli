/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package deploy

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are the kinds whose commands pre-date the ResourceKind template:
// User, which takes a username. Only account admins may deploy users.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "user <username>",
		Short: "Activate a user of your account",
		Long: `Activates a user of your account, who can then log in and use the api:

smarter deploy user <username>`,
		APIKind: "User",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "username", Param: "name"},
	},
}

func init() {
	cmd.RegisterResources(deployCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)

	// Only kinds flagged Deployable, whose platform broker implements deploy,
	// get a leaf command. Pulled from the shared registry so its
	// Use/APIKind/display text stays in lockstep with the other verbs. The
	// ChatBot command was removed: the platform has no ChatBot broker since
	// v0.14 (LLMClient replaced it).
	var specs []cmd.ResourceSpec
	for _, k := range cmd.ResourceKinds {
		if k.Deployable {
			specs = append(specs, k.DeploySpec())
		}
	}
	cmd.RegisterResources(deployCmd, specs, APIRequest, cmd.AdaptOutput(ConsoleOutput), cmd.ErrorOutput)
}
