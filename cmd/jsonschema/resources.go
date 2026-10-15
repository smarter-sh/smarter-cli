/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package jsonschema

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are the kinds that the other verbs also declare by hand, because
// their commands pre-date the v0.14 ResourceKind template.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:     "account",
		Short:   "Retrieve the JSON Schema of an Account manifest",
		Long:    "Retrieves the JSON Schema of an Account manifest:\n\n\tsmarter json-schema account > account.schema.json",
		APIKind: "Account",
	},
	{
		Use:     "plugin",
		Short:   "Retrieve the JSON Schema of a Plugin manifest",
		Long:    "Retrieves the JSON Schema of a Plugin manifest:\n\n\tsmarter json-schema plugin > plugin.schema.json",
		APIKind: "Plugin",
	},
	{
		Use:     "user",
		Short:   "Retrieve the JSON Schema of a User manifest",
		Long:    "Retrieves the JSON Schema of a User manifest:\n\n\tsmarter json-schema user > user.schema.json",
		APIKind: "User",
	},
}

func init() {
	cmd.RegisterResources(jsonSchemaCmd, legacySpecs, APIRequest, ConsoleOutput, cmd.ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.ResourceKinds))
	for i, k := range cmd.ResourceKinds {
		specs[i] = k.JSONSchemaSpec()
	}
	cmd.RegisterResources(jsonSchemaCmd, specs, APIRequest, ConsoleOutput, cmd.ErrorOutput)
}
