/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve your Account manifest",
	Long: `Retrieve your Account manifest:

smarter get account [flags]

The Smarter API will your Account manifest.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/get/account/
		bodyJson, err := APIRequest("Account", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(accountCmd)
}
