/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Smarter resource",
	Long: `Deploy a Smarter resource:

smarter deploy <kind>

The Smarter API will deploy the resource.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}
		bodyJson, err := APIRequest("deploy/chatbot/test", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
