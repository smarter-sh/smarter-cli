/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var sqlConnectionCmd = &cobra.Command{
	Use:   "sqlconnection",
	Short: "Generate an example manifest for a PluginDataSqlConnection.",
	Long: `Generates an example manifest for a PluginDataSqlConnection. For example:

	smarter manifest sqlconnection [flags] > my-plugin.yaml

This will generate an example manifest a PluginDataSqlConnection and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/plugindatasqlconnection/
		bodyJson, err := APIRequest("PluginDataSqlConnection", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(sqlConnectionCmd)
}
