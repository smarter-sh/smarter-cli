/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// whoamiCmd represents the status command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Retrieve information about the api_key owner",
	Long: `Retrieve information about the api_key owner:

smarter whoami --json --yaml

Returns informtation about the Smarter user account that owns the
configured api_key.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		kwargs := map[string]string{}
		bodyJson, err := GetAPIResponseResponse("whoami", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}
	},
}

func init() {
	RootCmd.AddCommand(whoamiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
