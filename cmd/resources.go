/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "List the resource kinds of the Smarter Platform",
	Long: `List the resource kinds of the Smarter Platform:

smarter resources [flags]

The Smarter API will return every kind of resource that a manifest can
declare, with its singular and plural names.`,
	Run: func(cmd *cobra.Command, args []string) {

		// this request goes to /api/v1/cli/resources/
		bodyJson, err := APIRequest("resources", map[string]string{})
		if err != nil {
			ErrorOutput(err)
			return
		}
		if err := ConsoleOutput(bodyJson); err != nil {
			ErrorOutput(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(resourcesCmd)
}
