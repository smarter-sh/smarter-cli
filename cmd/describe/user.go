/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user <username>",
	Short: "Retrieve a manifest for a User",
	Long: `Retrieves a manifest for a user. For example:

	smarter describe user <username> > my-plugin.yaml

This will retrieve a manifest for User <username> and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]

		kwargs := map[string]string{
			"username": username,
		}

		// this request goes to /api/v1/cli/describe/user/
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(UserCmd)
}
