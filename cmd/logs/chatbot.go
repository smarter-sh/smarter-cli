/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package logs

import (
	"github.com/spf13/cobra"
)

var chatbotCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Returns the logs for a ChatBot",
	Long: `Returns the logs for a ChatBot:

smarter logs chatbot <name>

The Smarter API will deploy the ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/logs/chatbot/
		bodyJson, err := APIRequest("ChatBot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	logsCmd.AddCommand(chatbotCmd)
}
