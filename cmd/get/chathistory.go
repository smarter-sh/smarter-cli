/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatHistoryCmd = &cobra.Command{
	Use:   "chat-history",
	Short: "Retrieve the chat history for a session_key",
	Long: `Retrieve the chat history for a session_key:

smarter get chat-history [session_key]

The Smarter API will return the chat history for the session_key.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_key, _ := cmd.Flags().GetString("session_key")

		kwargs := map[string]string{
			"session_key": session_key,
		}

		// this request goes to /api/v1/cli/get/chathistory/
		bodyJson, err := APIRequest("ChatHistory", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(chatHistoryCmd)

	chatHistoryCmd.Flags().StringP("session_key", "s", "", "Chat session_key")
	if err := chatHistoryCmd.MarkFlagRequired("session_key"); err != nil {
		log.Fatalf("Error marking flag 'session_key' as required: %v", err)
	}
	if err := viper.BindPFlag("session_key", chatHistoryCmd.Flags().Lookup("session_key")); err != nil {
		log.Fatalf("Error binding flag 'session_key': %v", err)
	}

	viper.AutomaticEnv() // This line is added
}
