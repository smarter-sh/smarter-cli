/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

// sqlConnectionCmd represents the chatbot command
var sqlConnectionCmd = &cobra.Command{
	Use:   "sqlconnection <name>",
	Short: "Delete a PluginDataSqlConnections",
	Long: `Deletes a PluginDataSqlConnections:

smarter delete sqlconnection <name> [flags]

The Smarter API will permanently delete the PluginDataSqlConnections with the specified name.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"name": args[0],
		}
		bodyJson, err := APIRequest("PluginDataSqlConnections", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(sqlConnectionCmd)
}
