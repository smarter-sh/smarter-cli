/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Retrieve information about the api_key owner",
	Long: `Retrieve information about the api_key owner:

smarter whoami --json --yaml

Returns informtation about the Smarter user account that owns the
configured api_key.`,
	Run: func(cmd *cobra.Command, args []string) {
		kwargs := map[string]string{}
		bodyJson, err := APIRequest("whoami", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}
	},
}

func init() {
	RootCmd.AddCommand(whoamiCmd)
}
