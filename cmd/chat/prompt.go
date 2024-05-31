/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package chat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompt a ChatBot and echo its response to the console",
	Long: `Prompts a ChatBot and echos its response to the console:

	smarter chat prompt --chatbot <chatbot> [--new_session]

The Smarter API will send the prompt to the ChatBot and return its response.`,
	Run: func(cmd *cobra.Command, args []string) {

		prompt, _ := cmd.Flags().GetString("prompt")
		chatbot, _ := cmd.Flags().GetString("chatbot")
		new_session, _ := cmd.Flags().GetBool("new_session")
		uid := getUniqueID()

		if prompt == "" {
			log.Fatalf("The 'prompt' flag is required")
		}

		if chatbot == "" {
			log.Fatalf("The 'chatbot' flag is required")
		}

		kwargs := map[string]string{
			"uid":         uid,
			"new_session": fmt.Sprintf("%t", new_session),
		}

		dict := map[string]string{"prompt": prompt}
		fileContentsBytes, _ := json.Marshal(dict)
		fileContents := string(fileContentsBytes)

		// this request goes to /api/v1/cli/chat/config/<str:chatbot>/<str:uid>
		path := fmt.Sprintf("%s/", chatbot)
		bodyJson, err := APIRequest(path, kwargs, fileContents)
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
