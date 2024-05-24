/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve your Account manifest",
	Long: `Retrieves a manifest of your account. For example:

	smarter describe account > my-plugin.yaml

This will generate an example manifest for an account and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		bodyJson, err := APIRequest("account", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(accountCmd)
}
