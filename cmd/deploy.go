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

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Smarter resource",
	Long: `Deploy a Smarter resource:

smarter deploy <kind>

The Smarter API will deploy the resource.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		bodyJson, err := GetAPIResponse("deploy/chatbot/test")
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
	RootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
