/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "user --json --yaml",
	Short: "Generate an example manifest for a user.",
	Long: `Generate an example manifest for a user. For example:

	smarter manifest user --json --yaml > my-plugin.yaml

This will generate an example manifest for a user and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(usersCmd)
}
