/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve real-time status of the Smarter Platform",
	Long: `Retrieve real-time status of the Smarter Platform:

smarter get status --json --yaml

The Smarter API will return the current status of the Smarter Platform,
including the status of all services and resources by region.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		bodyJson, err := GetAPIResponse("status")
		if err != nil {
			panic(err)
		} else {
			switch {
			case jsonFlagValue:
				fmt.Println(string(bodyJson))
			case yamlFlagValue:
				bodyYaml, err := yaml.JSONToYAML(bodyJson)
				if err != nil {
					panic(err)
				} else {
					fmt.Println(string(bodyYaml))
				}
			default:
				fmt.Println(string(bodyJson))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
