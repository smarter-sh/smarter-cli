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

// chatsCmd represents the chats command
var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve a list of Chats or the history of a specific Chat by id",
	Long: `Generate an example manifest for a chat session. For example:

	smarter manifest chat > my-plugin.yaml

This will generate an example manifest a chat session and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		name := viper.GetString("name")
		session_id := viper.GetString("session")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"name":       name,
			"session_id": session_id,
			"n":          strconv.Itoa(n),
		}

		bodyJson, err := GetAPIResponse("chats", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	manifestCmd.AddCommand(chatsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chatsCmd.Flags().String("name", "", "Name of the chatbot")
	chatsCmd.Flags().String("session", "", "Chat session_id")
	chatsCmd.Flags().Int("n", 10, "Number of sessions to retrieve")

	if err := viper.BindPFlag("name", chatsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

	if err := viper.BindPFlag("session", chatsCmd.Flags().Lookup("session")); err != nil {
		log.Fatalf("Error binding flag 'session': %v", err)
	}

	if err := viper.BindPFlag("n", chatsCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
