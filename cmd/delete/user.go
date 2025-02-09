/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

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

		// this request goes to /api/v1/cli/delete/user/
		_, err := APIRequest("user", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	deleteCmd.AddCommand(userCmd)
}
