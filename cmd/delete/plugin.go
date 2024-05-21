/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin <name>",
	Short: "Delete a Plugin",
	Long: `Delete a Plugin:

smarter delete plugin <name> --dry-run

The Smarter API will permanently delete the Plugin with the specified name,
and dissassociate it from any ChatBots.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"name": args[0],
		}
		bodyJson, err := APIRequest("plugin", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(pluginCmd)
}
