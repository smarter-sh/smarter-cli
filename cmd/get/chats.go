/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve a list of Chats",
	Long: `Retrieves a list of Chats:

smarter get chats [flags]

The Smarter API will return a list of Chats.`,
	Run: func(cmd *cobra.Command, args []string) {

		chatbot := viper.GetString("chatbot")
		session_key := viper.GetString("session_key")
		today := viper.GetBool("today")
		yesterday := viper.GetBool("yesterday")
		this_week := viper.GetBool("this-week")
		last_week := viper.GetBool("last-week")
		this_month := viper.GetBool("this-month")
		last_month := viper.GetBool("last-month")

		kwargs := map[string]string{
			"chatbot":     chatbot,
			"session_key": session_key,
			"today":       strconv.FormatBool(today),
			"yesterday":   strconv.FormatBool(yesterday),
			"this-week":   strconv.FormatBool(this_week),
			"last-week":   strconv.FormatBool(last_week),
			"this-month":  strconv.FormatBool(this_month),
			"last-month":  strconv.FormatBool(last_month),
		}

		// this request goes to /api/v1/cli/get/chat/
		bodyJson, err := APIRequest("Chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	getCmd.AddCommand(chatsCmd)

	chatsCmd.Flags().StringP("chatbot", "c", "", "Name of the chatbot")
	if err := viper.BindPFlag("chatbot", chatsCmd.Flags().Lookup("chatbot")); err != nil {
		log.Fatalf("Error binding flag 'chatbot': %v", err)
	}

	chatsCmd.Flags().StringP("session_key", "s", "", "Chat session_key")
	if err := viper.BindPFlag("session_key", chatsCmd.Flags().Lookup("session_key")); err != nil {
		log.Fatalf("Error binding flag 'session_key': %v", err)
	}

	chatsCmd.Flags().Bool("today", false, "Filter for today")
	if err := viper.BindPFlag("today", chatsCmd.Flags().Lookup("today")); err != nil {
		log.Fatalf("Error binding flag 'today': %v", err)
	}

	chatsCmd.Flags().Bool("yesterday", false, "Filter for yesterday")
	if err := viper.BindPFlag("yesterday", chatsCmd.Flags().Lookup("yesterday")); err != nil {
		log.Fatalf("Error binding flag 'yesterday': %v", err)
	}

	chatsCmd.Flags().Bool("this-week", false, "Filter for this week")
	if err := viper.BindPFlag("this-week", chatsCmd.Flags().Lookup("this-week")); err != nil {
		log.Fatalf("Error binding flag 'this-week': %v", err)
	}

	chatsCmd.Flags().Bool("last-week", false, "Filter for last week")
	if err := viper.BindPFlag("last-week", chatsCmd.Flags().Lookup("last-week")); err != nil {
		log.Fatalf("Error binding flag 'last-week': %v", err)
	}

	chatsCmd.Flags().Bool("this-month", false, "Filter for this month")
	if err := viper.BindPFlag("this-month", chatsCmd.Flags().Lookup("this-month")); err != nil {
		log.Fatalf("Error binding flag 'this-month': %v", err)
	}

	chatsCmd.Flags().Bool("last-month", false, "Filter for last month")
	if err := viper.BindPFlag("last-month", chatsCmd.Flags().Lookup("last-month")); err != nil {
		log.Fatalf("Error binding flag 'last-month': %v", err)
	}
}
