/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
		fmt.Println("chatbot called")
	},
}

func init() {
	deleteCmd.AddCommand(chatbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
