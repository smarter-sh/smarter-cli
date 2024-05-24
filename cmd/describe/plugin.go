/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugin <name>",
	Short: "Retrieve a manifest for a Plugin",
	Long: `Retrieves a manifest for a plugin. For example:

	smarter describe plugin <name> > my-plugin.yaml

This will retrieve the manifest for a plugin named <name> and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		bodyJson, err := APIRequest("plugin", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(pluginsCmd)
}
