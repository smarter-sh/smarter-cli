/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Retrieve a list of SmarterAuthTokens",
	Long: `Retrieves a list of SmarterAuthTokens:

smarter get apikey --name --json --yaml -n <10> --asc --desc

The Smarter API will return a list of apikeys in the specified format,
or a manifest for a specific apikey.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		kwargs := map[string]string{
			"name": name,
		}
		bodyJson, err := APIRequest("SmarterAuthToken", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(apikeyCmd)

	apikeyCmd.Flags().StringP("name", "n", "", "SmarterAuthToken name")
	if err := viper.BindPFlag("name", apikeyCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

}
