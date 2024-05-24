/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

var sqlConnectionCmd = &cobra.Command{
	Use:   "sqlconnection <name>",
	Short: "Delete a PluginDataSqlConnection",
	Long: `Deletes a PluginDataSqlConnection:

smarter delete sqlconnection <name> [flags]

The Smarter API will permanently delete the PluginDataSqlConnection with the specified name.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"name": args[0],
		}
		_, err := APIRequest("PluginDataSqlConnection", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	deleteCmd.AddCommand(sqlConnectionCmd)
}
