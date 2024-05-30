/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package chat

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the React configuration for a ChatBot",
	Long: `Retrieves the React configuration for a ChatBot:

smarter chat config [flags]


The Smarter API will return a dict of the configuration that is provided
to the React chat application in the Smarter web console. This is the same
dict that is returned by the /chatapp/<chatbot>/config/ endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {

		chatbot, _ := cmd.Flags().GetString("chatbot")
		session_key := fetchSessionKey()

		if chatbot == "" {
			log.Fatalf("The 'chatbot' flag is required")
		}

		kwargs := map[string]string{
			"session_key": session_key,
		}

		// this request goes to /api/v1/cli/chat/config/
		bodyJson, err := APIRequest("config", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	chatCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("chatbot", "c", "", "the name of a deployed ChatBot")
	if err := viper.BindPFlag("chatbot", configCmd.Flags().Lookup("chatbot")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
