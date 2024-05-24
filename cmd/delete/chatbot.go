/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

var chatbotCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Delete a ChatBot",
	Long: `Delete a ChatBot:

smarter delete chatbot <name> --dry-run

The Smarter API will permanently delete the ChatBot with the specified name,
and all related chat history.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"name": args[0],
		}
		bodyJson, err := APIRequest("chatbot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(chatbotCmd)
}
