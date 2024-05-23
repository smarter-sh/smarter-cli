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
	Short: "Retrieve the chat history for a session_id",
	Long: `Retrieve the chat history for a session_id:

smarter get chat-history [session_id]

The Smarter API will return the chat history for the session_id.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_id, _ := cmd.Flags().GetString("session_id")

		kwargs := map[string]string{
			"session_id": session_id,
		}

		bodyJson, err := APIRequest("ChatHistory", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(chatHistoryCmd)

	chatHistoryCmd.Flags().StringP("session_id", "s", "", "Chat session_id")
	if err := chatHistoryCmd.MarkFlagRequired("session_id"); err != nil {
		log.Fatalf("Error marking flag 'session_id' as required: %v", err)
	}
	if err := viper.BindPFlag("session_id", chatHistoryCmd.Flags().Lookup("session_id")); err != nil {
		log.Fatalf("Error binding flag 'session_id': %v", err)
	}

	viper.AutomaticEnv() // This line is added
}
