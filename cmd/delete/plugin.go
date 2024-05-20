/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Delete a Plugin",
	Long: `Delete a Plugin:

smarter delete plugin -name --dry-run

The Smarter API will permanently delete the Plugin with the specified name,
and dissassociate it from any ChatBots.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		kwargs := map[string]string{}
		bodyJson, err := GetAPI("plugin", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(pluginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
