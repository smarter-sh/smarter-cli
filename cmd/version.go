/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the status command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Retrieve version information",
	Long: `Retrieve version information:

smarter version

Returns version information about this software.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		bodyJson := []byte(`{"version":"` + Version + `"}`)

		if jsonFlagValue {
			fmt.Println(string(bodyJson))
		} else if yamlFlagValue {
			bodyYaml, err := yaml.JSONToYAML(bodyJson)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(string(bodyYaml))
			}
		} else {
			fmt.Println(string(bodyJson))
		}

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
