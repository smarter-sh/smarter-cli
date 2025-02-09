/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve real-time status of the Smarter Platform",
	Long: `Retrieve real-time status of the Smarter Platform:

smarter status [flags]

The Smarter API will return the current status of the Smarter Platform,
including the status of all services and resources by region.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/status/
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
