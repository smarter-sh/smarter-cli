/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package logs

import (
	"github.com/spf13/cobra"
)

var chatsCmd = &cobra.Command{
	Use:   "chat <session_id>",
	Short: "Returns the logs for a Chat session_id",
	Long: `Returns the logs for a Chat:

smarter logs chat <session_id>

The Smarter API will return the logs for a Chat session_id.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/logs/chat/
		bodyJson, err := APIRequest("Chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	logsCmd.AddCommand(chatsCmd)
}
