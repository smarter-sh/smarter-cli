/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatsCmd represents the chats command
var chatsCmd = &cobra.Command{
	Use:   "chats --session_id <session_id> --chatbot <name> --json --yaml -n <10> --asc --desc --today --yesterday --this-week --last-week --this-month --last-month",
	Short: "Retrieve a list of Chats",
	Long: `Retrieves a list of Chats:

smarter get chats --session_id <session_id> --chatbot <name> --json --yaml -n <10> --asc --desc --today --yesterday --this-week --last-week --this-month --last-month

The Smarter API will return a list of Chats.`,
	Run: func(cmd *cobra.Command, args []string) {

		chatbot := viper.GetString("chatbot")
		session_id := viper.GetString("session_id")
		today := viper.GetBool("today")
		yesterday := viper.GetBool("yesterday")
		this_week := viper.GetBool("this-week")
		last_week := viper.GetBool("last-week")
		this_month := viper.GetBool("this-month")
		last_month := viper.GetBool("last-month")

		kwargs := map[string]string{
			"chatbot":    chatbot,
			"session_id": session_id,
			"today":      strconv.FormatBool(today),
			"yesterday":  strconv.FormatBool(yesterday),
			"this-week":  strconv.FormatBool(this_week),
			"last-week":  strconv.FormatBool(last_week),
			"this-month": strconv.FormatBool(this_month),
			"last-month": strconv.FormatBool(last_month),
		}

		bodyJson, err := APIRequest("chats", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(chatsCmd)

	pluginsCmd.Flags().String("chatbot", "c", "Name of the chatbot")
	if err := viper.BindPFlag("chatbot", pluginsCmd.Flags().Lookup("chatbot")); err != nil {
		log.Fatalf("Error binding flag 'chatbot': %v", err)
	}

	pluginsCmd.Flags().String("session_id", "s", "Chat session_id")
	if err := viper.BindPFlag("session_id", pluginsCmd.Flags().Lookup("session_id")); err != nil {
		log.Fatalf("Error binding flag 'session_id': %v", err)
	}

	pluginsCmd.Flags().Bool("today", false, "Filter for today")
	if err := viper.BindPFlag("today", pluginsCmd.Flags().Lookup("today")); err != nil {
		log.Fatalf("Error binding flag 'today': %v", err)
	}

	pluginsCmd.Flags().Bool("yesterday", false, "Filter for yesterday")
	if err := viper.BindPFlag("yesterday", pluginsCmd.Flags().Lookup("yesterday")); err != nil {
		log.Fatalf("Error binding flag 'yesterday': %v", err)
	}

	pluginsCmd.Flags().Bool("this-week", false, "Filter for this week")
	if err := viper.BindPFlag("this-week", pluginsCmd.Flags().Lookup("this-week")); err != nil {
		log.Fatalf("Error binding flag 'this-week': %v", err)
	}

	pluginsCmd.Flags().Bool("last-week", false, "Filter for last week")
	if err := viper.BindPFlag("last-week", pluginsCmd.Flags().Lookup("last-week")); err != nil {
		log.Fatalf("Error binding flag 'last-week': %v", err)
	}

	pluginsCmd.Flags().Bool("this-month", false, "Filter for this month")
	if err := viper.BindPFlag("this-month", pluginsCmd.Flags().Lookup("this-month")); err != nil {
		log.Fatalf("Error binding flag 'this-month': %v", err)
	}

	pluginsCmd.Flags().Bool("last-month", false, "Filter for last month")
	if err := viper.BindPFlag("last-month", pluginsCmd.Flags().Lookup("last-month")); err != nil {
		log.Fatalf("Error binding flag 'last-month': %v", err)
	}
}
