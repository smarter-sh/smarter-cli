/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatPluginUsage = &cobra.Command{
	Use:   "chat-plugin-usage",
	Short: "Retrieve the chat plugin usage for a session_key",
	Long: `Retrieve the chat plugin usage for a session_key:

smarter get chat-plugin-usage [session_key]

The Smarter API will return the chat plugin usage for the session_key.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_key, _ := cmd.Flags().GetString("session_key")

		kwargs := map[string]string{
			"session_key": session_key,
		}

		// this request goes to /api/v1/cli/get/chatpluginusage/
		bodyJson, err := APIRequest("ChatPluginUsage", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(chatPluginUsage)

	chatPluginUsage.Flags().StringP("session_key", "s", "", "Chat session_key")
	if err := chatPluginUsage.MarkFlagRequired("session_key"); err != nil {
		log.Fatalf("Error marking flag 'session_key' as required: %v", err)
	}
	if err := viper.BindPFlag("session_key", chatPluginUsage.Flags().Lookup("session_key")); err != nil {
		log.Fatalf("Error binding flag 'session_key': %v", err)
	}

}
