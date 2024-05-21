/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package deploy

import (
	"github.com/spf13/cobra"
)

var chatbotCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Deploy a ChatBot",
	Long: `Deploys a ChatBot:

smarter deploy chatbot <name> --json --yaml

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
	DeployCmd.AddCommand(chatbotCmd)
}
