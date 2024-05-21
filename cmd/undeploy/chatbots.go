/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package undeploy

import (
	"github.com/spf13/cobra"
)

var chatbotsCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "Deploy a ChatBot",
	Long: `Deploys a ChatBot:

smarter deploy chatbot <name> --json --yaml

The Smarter API will deploy the ChatBot.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
