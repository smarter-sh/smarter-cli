/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatbotsCmd = &cobra.Command{
	Use:   "chatbots",
	Short: "Retrieve a list of ChatBots",
	Long: `Retrieve a list of ChatBots:

smarter get chatbots [flags]

The Smarter API will return a list of ChatBots in the specified format,
or a manifest for a specific ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/get/chatbot/
		bodyJson, err := APIRequest("Chatbot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(chatbotsCmd)

	chatbotsCmd.Flags().StringP("name", "n", "", "Name of the chatbot")

	if err := viper.BindPFlag("name", chatbotsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
