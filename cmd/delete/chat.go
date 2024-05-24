/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat <session_id>",
	Short: "Delete a chat history",
	Long: `Deletes a chat history:

smarter delete chat <session_id>

The Smarter API will permanently delete the chat history with the specified identifier.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"session_id": args[0],
		}
		bodyJson, err := APIRequest("chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	DeleteCmd.AddCommand(chatCmd)
}
