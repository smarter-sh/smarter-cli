/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"fmt"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Retrieve a list of Users",
	Long: `Retrieve a list of Users, or a specific User by username:

smarter get users --name --json --yaml --csv --xml -n 10 --asc --desc

The Smarter API will return a list of Users in the specified format,
or a manifest for a specific User.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("users")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response:", string(body))
		}
	},
}

func init() {
	accountCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
