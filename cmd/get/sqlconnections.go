/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sqlConnectionCmd = &cobra.Command{
	Use:   "sqlconnections",
	Short: "Retrieve a list of PluginDataSqlConnections",
	Long: `Retrieve a list of PluginDataSqlConnections:

smarter get sqlconnection [flags]

The Smarter API will return a list of PluginDataSqlConnection in the specified format,
or a manifest for a specific PluginDataSqlConnection.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/get/plugindatasqlconnection/
		bodyJson, err := APIRequest("PluginDataSqlConnection", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(sqlConnectionCmd)

	sqlConnectionCmd.Flags().StringP("name", "n", "", "Name of the PluginDataSqlConnection")

	if err := viper.BindPFlag("name", sqlConnectionCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
