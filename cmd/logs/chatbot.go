/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
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

		bodyJson, err := APIRequest("chatbot", kwargs)
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
