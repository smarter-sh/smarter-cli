/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package chat

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompt a ChatBot and echo its response to the console",
	Long: `Prompts a ChatBot and echos its response to the console:

smarter chat prompt [flags]

The Smarter API will send the prompt to the ChatBot and return its response.`,
	Run: func(cmd *cobra.Command, args []string) {

		chatbot, _ := cmd.Flags().GetString("chatbot")
		session_key := fetchSessionKey()

		if chatbot == "" {
			log.Fatalf("The 'chatbot' flag is required")
		}

		kwargs := map[string]string{
			"session_key": session_key,
		}

		// this request goes to /api/v1/cli/chat/ which is immediately upstream
		// from here. Hence, we're not providing a slug
		bodyJson, err := APIRequest("", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	chatCmd.AddCommand(promptCmd)

	promptCmd.PersistentFlags().StringP("prompt", "p", "", "A prompt to send to the ChatBot")
	if err := viper.BindPFlag("prompt", promptCmd.PersistentFlags().Lookup("prompt")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
