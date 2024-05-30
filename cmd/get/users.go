/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Retrieve a list of Users",
	Long: `Retrieves a list of Users:

smarter get users [flags]

The Smarter API will return a list of Users in the specified format,
or a manifest for a specific User.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("username")

		kwargs := map[string]string{
			"username": name,
		}

		// this request goes to /api/v1/cli/get/user/
		bodyJson, err := APIRequest("User", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(usersCmd)
	usersCmd.Flags().StringP("username", "u", "", "Smarter username")
	if err := viper.BindPFlag("username", usersCmd.Flags().Lookup("username")); err != nil {
		log.Fatalf("Error binding flag 'username': %v", err)
	}

}
