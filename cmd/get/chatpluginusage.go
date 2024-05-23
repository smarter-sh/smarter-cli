/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatPluginUsage = &cobra.Command{
	Use:   "chat-plugin-usage",
	Short: "Retrieve the chat plugin usage for a session_id",
	Long: `Retrieve the chat plugin usage for a session_id:

smarter get chat-plugin-usage [session_id]

The Smarter API will return the chat plugin usage for the session_id.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_id, _ := cmd.Flags().GetString("session_id")

		kwargs := map[string]string{
			"session_id": session_id,
		}

		bodyJson, err := APIRequest("ChatPluginUsage", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(chatPluginUsage)

	chatPluginUsage.Flags().StringP("session_id", "s", "", "Chat session_id")
	if err := chatPluginUsage.MarkFlagRequired("session_id"); err != nil {
		log.Fatalf("Error marking flag 'session_id' as required: %v", err)
	}
	if err := viper.BindPFlag("session_id", chatPluginUsage.Flags().Lookup("session_id")); err != nil {
		log.Fatalf("Error binding flag 'session_id': %v", err)
	}

}
