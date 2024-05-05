/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user from your account",
	Long: `Delete a user from your account:

smarter delete user --name --dry-run

The Smarter API will permanently delete the user with the specified name,
and dissassociate it from any Smarter resources. Your Smarter admin account
will replace the deleted user.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		bodyJson, err := GetAPI("user")
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
	accountCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
