/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

// UserCmd represents the users command
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
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(UserCmd)
}
