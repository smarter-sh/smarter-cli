/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

// chatbotsCmd represents the chatbots command
var chatbotsCmd = &cobra.Command{
	Use:   "chatbots",
	Short: "Retrieve a list of ChatBots or a specific ChatBot by name",
	Long: `Retrieve a list of ChatBots or a specific ChatBot by name:

smarter get chatbots --name --json --yaml

The Smarter API will return a list of ChatBots in the specified format,
or a manifest for a specific ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("chatbots")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response:", string(body))
		}

	},
}

func init() {
	GetCmd.AddCommand(chatbotsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
