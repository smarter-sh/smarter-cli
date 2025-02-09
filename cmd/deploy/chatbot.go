/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package deploy

import (
	"github.com/spf13/cobra"
)

var chatbotCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Deploy a ChatBot",
	Long: `Deploys a ChatBot:

smarter deploy chatbot <name> [flags]

The Smarter API will deploy the ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		// this request goes to /api/v1/cli/deploy/chatbot/
		_, err := APIRequest("ChatBot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	deployCmd.AddCommand(chatbotCmd)
}
