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

smarter logs chatbot <name> --json --yaml

The Smarter API will deploy the ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		kwargs := map[string]string{
			"name": name,
		}

		bodyJson, err := APIRequest("chatbot", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	logsCmd.AddCommand(chatbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
