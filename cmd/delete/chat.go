/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package delete

import (
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat <session_key>",
	Short: "Delete a chat history",
	Long: `Deletes a chat history:

smarter delete chat <session_key>

The Smarter API will permanently delete the chat history with the specified identifier.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{
			"session_key": args[0],
		}

		// this request goes to /api/v1/cli/delete/chat/
		_, err := APIRequest("chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	deleteCmd.AddCommand(chatCmd)
}
