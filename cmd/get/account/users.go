/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"log"
	"strconv"

	"github.com/QueriumCorp/smarter-cli/cmd/get"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		name := viper.GetString("username")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"username": name,
			"n":        strconv.Itoa(n),
		}
		bodyJson, err := get.APIRequest("users", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
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
	accountCmd.Flags().String("username", "", "Smarter username")
	accountCmd.Flags().Int("n", 10, "Number of users to retrieve")

	if err := viper.BindPFlag("username", accountCmd.Flags().Lookup("username")); err != nil {
		log.Fatalf("Error binding flag 'username': %v", err)
	}

	if err := viper.BindPFlag("n", accountCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
