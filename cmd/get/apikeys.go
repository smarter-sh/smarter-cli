/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikeys",
	Short: "Retrieve a list of SmarterAuthTokens",
	Long: `Retrieves a list of SmarterAuthTokens:

smarter get apikey [flags]

The Smarter API will return a list of apikeys in the specified format,
or a manifest for a specific apikey.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/get/smarterauthtoken/
		bodyJson, err := APIRequest("SmarterAuthToken", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(apikeyCmd)

	apikeyCmd.Flags().StringP("name", "n", "", "SmarterAuthToken name")
	if err := viper.BindPFlag("name", apikeyCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

}
