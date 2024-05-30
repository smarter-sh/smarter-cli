/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugin [flags]",
	Short: "Generate an example manifest for a plugin.",
	Long: `Generates an example manifest for a plugin. For example:

	smarter manifest plugin [flags] > my-plugin.yaml

This will generate an example manifest for a plugin and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/plugin/
		bodyJson, err := APIRequest("plugin", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(pluginsCmd)
}
