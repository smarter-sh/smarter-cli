/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package chat

import (
	"fmt"
	"log"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIRequest(slug string, kwargs map[string]string) ([]byte, error) {

	if slug == "" {
		return cmd.APIRequest("chat", kwargs)
	} else {
		return cmd.APIRequest(fmt.Sprintf("chat/%s", slug), kwargs)
	}

}
func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "json")
	}
	cmd.ConsoleOutput(bodyJson)
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

func fetchSessionKey() string {
	environment := viper.GetString("config.environment")
	sessionKey := viper.GetString(fmt.Sprintf("%s.session_key", environment))

	return sessionKey
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with a deployed ChatBot",
	Long: `Chat with a deployed ChatBot:

smarter chat <prompt> [flags]

The Smarter API will send the prompt to a deployed ChatBot and
then echo its response to the console.`,
	Run: func(cmd *cobra.Command, args []string) {

		prompt := viper.GetString("prompt")

		kwargs := map[string]string{
			"prompt": prompt,
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
	cmd.RootCmd.AddCommand(chatCmd)

	promptCmd.Flags().StringP("chatbot", "c", "", "the name of a deployed ChatBot")
	if err := viper.BindPFlag("chatbot", promptCmd.Flags().Lookup("chatbot")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
