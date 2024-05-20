/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve your Account manifest",
	Long: `Generate an example manifest for an account. For example:

	smarter manifest account > my-plugin.yaml

This will generate an example manifest for an account and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		kwargs := map[string]string{}

		bodyJson, err := GetAPIResponse("chats", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	manifestCmd.AddCommand(accountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
