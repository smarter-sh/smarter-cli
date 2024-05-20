/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user from your account",
	Long: `Delete a user from your account:

smarter delete account user --username --dry-run

The Smarter API will permanently delete the user with the specified name,
and dissassociate it from any Smarter resources. Your Smarter admin account
will replace the deleted user.`,
	Run: func(cmd *cobra.Command, args []string) {

		username := viper.GetString("username")
		kwargs := map[string]string{
			"username": username,
		}
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
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
	accountCmd.Flags().String("username", "", "Name of the plugin")

	if err := viper.BindPFlag("username", accountCmd.Flags().Lookup("username")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

}
