/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

// pluginsCmd represents the plugins command
var pluginsCmd = &cobra.Command{
	Use:   "plugin --json --yaml",
	Short: "Generate an example manifest for a plugin.",
	Long: `Generates an example manifest for a plugin. For example:

	smarter manifest plugin --json --yaml > my-plugin.yaml

This will generate an example manifest for a plugin and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		bodyJson, err := APIRequest("plugin", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(pluginsCmd)
}
