/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatToolCalls = &cobra.Command{
	Use:   "chat-tool-calls",
	Short: "Retrieve the chat tool calls for a session_id",
	Long: `Retrieve the chat tool calls for a session_id:

smarter get chat-tool-calls [session_id]

The Smarter API will return the chat tool calls for the session_id.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_id, _ := cmd.Flags().GetString("session_id")

		kwargs := map[string]string{
			"session_id": session_id,
		}

		// this request goes to /api/v1/cli/get/chattoolcall/
		bodyJson, err := APIRequest("ChatToolCall", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(chatToolCalls)

	chatToolCalls.Flags().StringP("session_id", "s", "", "Chat session_id")
	if err := chatToolCalls.MarkFlagRequired("session_id"); err != nil {
		log.Fatalf("Error marking flag 'session_id' as required: %v", err)
	}
	if err := viper.BindPFlag("session_id", chatToolCalls.Flags().Lookup("session_id")); err != nil {
		log.Fatalf("Error binding flag 'session_id': %v", err)
	}

}
