/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

// apikeyCmd represents the chats command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Generate an example manifest for a SmarterAuthToken.",
	Long: `Generates an example manifest for a SmarterAuthToken. For example:

	smarter manifest apikey [flags] > my-plugin.yaml

This will generate an example manifest a SmarterAuthToken and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		bodyJson, err := APIRequest("SmarterAuthToken", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(apikeyCmd)
}
