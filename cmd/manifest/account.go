/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve your Account manifest",
	Long: `Generate an example manifest for an Account. For example:

	smarter manifest account [flags] > my-plugin.yaml

This will generate an example manifest for an Account and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/account/
		bodyJson, err := APIRequest("Account", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(accountCmd)
}
