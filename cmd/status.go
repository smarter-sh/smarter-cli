/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve real-time status of the Smarter Platform",
	Long: `Retrieve real-time status of the Smarter Platform:

smarter get status --json --yaml

The Smarter API will return the current status of the Smarter Platform,
including the status of all services and resources by region.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}
		bodyJson, err := APIRequest("status", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
