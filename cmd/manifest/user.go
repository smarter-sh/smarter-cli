/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "user [flags]",
	Short: "Generate an example manifest for a user.",
	Long: `Generate an example manifest for a user. For example:

	smarter manifest user [flags] > my-plugin.yaml

This will generate an example manifest for a user and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/user/
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(usersCmd)
}
