/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
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

		// this request goes to /api/v1/cli/undeploy/chatbot/
		_, err := APIRequest("ChatBot", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput()
		}

	},
}

func init() {
	undeployCmd.AddCommand(chatbotsCmd)
}
