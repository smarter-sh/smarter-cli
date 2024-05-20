/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatsCmd represents the chats command
var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve a list of Chats or the history of a specific Chat by id",
	Long: `Retrieve a list of Chats or the history of a specific Chat by id:

smarter get chats --id --chatbot --json --yaml --csv --xml -n 10 --asc --desc --today --yesterday --this-week --last-week --this-month --last-month

The Smarter API will return a list of Chats in the specified format,
or a manifest for a specific Chat history.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")
		session_id := viper.GetString("session")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"name":       name,
			"session_id": session_id,
			"n":          strconv.Itoa(n),
		}

		bodyJson, err := APIRequest("chats", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(chatsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pluginsCmd.Flags().String("name", "", "Name of the chatbot")
	pluginsCmd.Flags().String("session", "", "Chat session_id")
	pluginsCmd.Flags().Int("n", 10, "Number of sessions to retrieve")

	if err := viper.BindPFlag("name", pluginsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

	if err := viper.BindPFlag("session", pluginsCmd.Flags().Lookup("session")); err != nil {
		log.Fatalf("Error binding flag 'session': %v", err)
	}

	if err := viper.BindPFlag("n", pluginsCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
