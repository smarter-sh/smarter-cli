/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

// chatbotCmd represents the chatbots command
var chatbotCmd = &cobra.Command{
	Use:   "chatbot --json --yaml",
	Short: "Generate an example manifest for a ChatBot",
	Long: `Generates an example manifest for a ChatBot resource. For example:

	smarter manifest chatbot --json --yaml > my-plugin.yaml

This will generate an example manifest for a chatbot and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		bodyJson, err := APIRequest("chatbot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(chatbotCmd)
}
