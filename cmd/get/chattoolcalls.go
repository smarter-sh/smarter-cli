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
	Short: "Retrieve the chat tool calls for a session_key",
	Long: `Retrieve the chat tool calls for a session_key:

smarter get chat-tool-calls [session_key]

The Smarter API will return the chat tool calls for the session_key.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_key, _ := cmd.Flags().GetString("session_key")

		kwargs := map[string]string{
			"session_key": session_key,
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

	chatToolCalls.Flags().StringP("session_key", "s", "", "Chat session_key")
	if err := chatToolCalls.MarkFlagRequired("session_key"); err != nil {
		log.Fatalf("Error marking flag 'session_key' as required: %v", err)
	}
	if err := viper.BindPFlag("session_key", chatToolCalls.Flags().Lookup("session_key")); err != nil {
		log.Fatalf("Error binding flag 'session_key': %v", err)
	}

}
