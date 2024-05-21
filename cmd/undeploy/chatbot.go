/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package undeploy

import (
	"github.com/spf13/cobra"
)

var chatbotsCmd = &cobra.Command{
	Use:   "chatbot <name>",
	Short: "Undo a ChatBot deployment.",
	Long: `Undo a ChatBot deployment. For example:

smarter undeploy chatbot <name>

This will reverse the effect of having deployed the ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		bodyJson, err := APIRequest("chatbots", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	UndeployCmd.AddCommand(chatbotsCmd)
}