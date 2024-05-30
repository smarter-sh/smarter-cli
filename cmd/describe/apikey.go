/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey <name>",
	Short: "Retrieve a SmarterAuthToken manifest by name",
	Long: `Retrieves a manifest for a SmarterAuthToken. For example:

	smarter describe apikey <name> > my-plugin.yaml

This will generate a manifest for a SmarterAuthToken named <name> and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/describe/smarterauthtoken/
		bodyJson, err := APIRequest("SmarterAuthToken", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	describeCmd.AddCommand(apikeyCmd)
}
