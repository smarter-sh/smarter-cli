/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Retrieve a list of Users",
	Long: `Generate an example manifest for a user. For example:

	smarter manifest user > my-plugin.yaml

This will generate an example manifest for a user and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("username")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"username": name,
			"n":        strconv.Itoa(n),
		}
		bodyJson, err := GetAPIResponse("users", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	usersCmd.Flags().String("username", "", "Smarter username")
	usersCmd.Flags().Int("n", 10, "Number of users to retrieve")

	if err := viper.BindPFlag("username", usersCmd.Flags().Lookup("username")); err != nil {
		log.Fatalf("Error binding flag 'username': %v", err)
	}

	if err := viper.BindPFlag("n", usersCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
