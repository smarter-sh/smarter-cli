/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

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

		// this request goes to /api/v1/cli/delete/plugin/
		_, err := APIRequest("plugin", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	deleteCmd.AddCommand(pluginCmd)
}
