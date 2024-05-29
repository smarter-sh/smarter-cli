/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with a deployed ChatBot",
	Long: `Chat with a deployed ChatBot:

smarter chat <prompt> [flags]

The Smarter API will send the prompt to a deployed ChatBot and
then echo its response to the console.`,
	Run: func(cmd *cobra.Command, args []string) {

		session_key := viper.GetString("session_key")
		prompt := viper.GetString("prompt")

		kwargs := map[string]string{
			"session_key": session_key,
			"prompt":      prompt,
		}

		bodyJson, err := APIRequest("chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}
	},
}

func init() {
	RootCmd.AddCommand(chatCmd)

	chatCmd.PersistentFlags().StringP("session_key", "s", "", "Smarter Chat session_key to use")
	if err := viper.BindPFlag("session_key", chatCmd.PersistentFlags().Lookup("session_key")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	chatCmd.PersistentFlags().StringP("prompt", "p", "", "A prompt to send to the ChatBot")
	if err := viper.BindPFlag("prompt", chatCmd.PersistentFlags().Lookup("prompt")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
