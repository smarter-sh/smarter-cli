/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package logs

import (
	"github.com/spf13/cobra"
)

var chatsCmd = &cobra.Command{
	Use:   "chat <session_key>",
	Short: "Returns the logs for a Chat session_key",
	Long: `Returns the logs for a Chat:

smarter logs chat <session_key>

The Smarter API will return the logs for a Chat session_key.`,
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
