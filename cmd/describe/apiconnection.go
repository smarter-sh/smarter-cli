/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var apiConnectionCmd = &cobra.Command{
	Use:   "apiconnection <name>",
	Short: "Retrieve a PluginDataApiConnection manifest by name",
	Long: `Retrieves a manifest for a PluginDataApiConnection. For example:

	smarter describe apiconnection <name> > my-plugin.yaml

This will generate a manifest for a PluginDataApiConnection named <name> and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/describe/plugindataapiconnection/
		bodyJson, err := APIRequest("PluginDataApiConnection", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(apiConnectionCmd)
}
