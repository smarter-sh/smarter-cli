/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package chat

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the ReactJS app configuration for a ChatBot",
	Long: `Retrieves the ReactJS app configuration for a ChatBot:

smarter chat config --chatbot <chatbot> [--new_session]


The Smarter API will return a dict of the configuration that is provided
to the React chat application in the Smarter web console. This is the same
dict that is returned by the /chatapp/<chatbot>/config/ endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {

		chatbot, _ := cmd.Flags().GetString("chatbot")
		new_session, _ := cmd.Flags().GetBool("new_session")
		uid := getUniqueID()

		if chatbot == "" {
			log.Fatalf("The 'chatbot' flag is required")
		}

		kwargs := map[string]string{
			"uid":         uid,
			"new_session": fmt.Sprintf("%t", new_session),
		}

		// this request goes to /api/v1/cli/chat/config/<str:chatbot>/<str:uid>
		path := fmt.Sprintf("config/%s/", chatbot)
		bodyJson, err := APIRequest(path, kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	chatCmd.AddCommand(configCmd)

}
