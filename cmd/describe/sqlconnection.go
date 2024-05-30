/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var sqlConnectionCmd = &cobra.Command{
	Use:   "sqlconnection <name>",
	Short: "Retrieve a PluginDataSqlConnections manifest by name",
	Long: `Retrieves a manifest for a PluginDataSqlConnections. For example:

	smarter describe sqlconnection <name> > my-plugin.yaml

This will generate a manifest for a PluginDataSqlConnections named <name> and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/describe/plugindatasqlconnections/
		bodyJson, err := APIRequest("PluginDataSqlConnections", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(sqlConnectionCmd)
}
