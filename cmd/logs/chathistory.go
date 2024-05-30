/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package logs

import (
	"github.com/spf13/cobra"
)

var chatHistoryCmd = &cobra.Command{
	Use:   "chat-history <session_id>",
	Short: "Returns the logs for a ChatHistory session_id",
	Long: `Returns the logs for a ChatHistory:

smarter logs chat-history <session_id>

The Smarter API will return the logs for a ChatHistory session_id.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/logs/chathistory/
		bodyJson, err := APIRequest("ChatHistory", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	logsCmd.AddCommand(chatHistoryCmd)
}
