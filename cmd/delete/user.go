/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user <username>",
	Short: "Delete a user from your account",
	Long: `Delete a user from your account:

smarter delete <username> --dry-run

The Smarter API will permanently delete the user with the specified name,
and dissassociate it from any Smarter resources. Your Smarter admin account
will replace the deleted user.`,
	Run: func(cmd *cobra.Command, args []string) {

		username := args[0]
		kwargs := map[string]string{
			"username": username,
		}
		bodyJson, err := APIRequest("user", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(userCmd)
}
