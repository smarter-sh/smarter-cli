/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

// apiConnectionCmd represents the chats command
var apiConnectionCmd = &cobra.Command{
	Use:   "apiconnection",
	Short: "Generate an example manifest for a PluginDataApiConnection.",
	Long: `Generates an example manifest for a PluginDataApiConnection. For example:

	smarter manifest apiconnection [flags] > my-plugin.yaml

This will generate an example manifest a PluginDataApiConnection and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		bodyJson, err := APIRequest("PluginDataApiConnection", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(apiConnectionCmd)
}
