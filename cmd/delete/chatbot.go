/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// chatbotCmd represents the chatbot command
var chatbotCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "Delete a ChatBot",
	Long: `Delete a ChatBot:

smarter delete chatbot -name --dry-run

The Smarter API will permanently delete the ChatBot with the specified name,
and all related chat history.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("chatbot")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			bodyJson, err := json.Marshal(body)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Response:", string(bodyJson))
			}
		}

	},
}

func init() {
	DeleteCmd.AddCommand(chatbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
