/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package describe

import (
	"github.com/spf13/cobra"
)

// chatbotCmd represents the chatbots command
var chatbotCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Retrieve a ChatBot manifest by name",
	Long: `Retrieves a manifest for a chatbot. For example:

	smarter describe chatbot <name> > my-plugin.yaml

This will generate a manifest for a chatbot named <name> and write it to my-plugin.yaml in the current working directory.`,

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
	describeCmd.AddCommand(chatbotCmd)
}
