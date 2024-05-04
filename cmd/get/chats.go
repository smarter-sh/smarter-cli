/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// chatsCmd represents the chats command
var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve a list of Chats or the history of a specific Chat by id",
	Long: `Retrieve a list of Chats or the history of a specific Chat by id:

smarter get chats --id --chatbot --json --yaml --csv --xml -n 10 --asc --desc --today --yesterday --this-week --last-week --this-month --last-month

The Smarter API will return a list of Chats in the specified format,
or a manifest for a specific Chat history.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("chats")
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
	GetCmd.AddCommand(chatsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
