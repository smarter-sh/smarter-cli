/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"log"

	"github.com/QueriumCorp/smarter-cli/cmd/get"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Retrieve a list of Users",
	Long: `Retrieve a list of Users, or a specific User by username:

smarter get account users --name --json --yaml -n <10> --asc --desc

The Smarter API will return a list of Users in the specified format,
or a manifest for a specific User.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("username")

		kwargs := map[string]string{
			"username": name,
		}
		bodyJson, err := get.APIRequest("users", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	accountCmd.AddCommand(usersCmd)
	accountCmd.Flags().StringP("username", "u", "", "Smarter username")
	if err := viper.BindPFlag("username", accountCmd.Flags().Lookup("username")); err != nil {
		log.Fatalf("Error binding flag 'username': %v", err)
	}

}
