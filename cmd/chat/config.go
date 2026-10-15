/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package chat

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the chat configuration of an LLMClient",
	Long: `Retrieves the chat configuration of an LLMClient:

smarter chat config --llmclient <name> [--new_session]

The Smarter API will return the configuration that it provides to the
React chat application in the Smarter web console, including the
session's history.`,
	Run: func(c *cobra.Command, args []string) {

		name, err := llmClientName(c)
		if err != nil {
			cmd.ErrorOutput(err)
			return
		}
		newSession, _ := c.Flags().GetBool("new_session")

		kwargs := map[string]string{
			"uid":         getUniqueID(),
			"new_session": fmt.Sprintf("%t", newSession),
		}

		// this request goes to /api/v1/cli/prompt/config/<str:name>/
		bodyJson, err := APIRequest(fmt.Sprintf("config/%s/", name), kwargs)
		if err != nil {
			cmd.ErrorOutput(err)
			return
		}
		ConsoleOutput(bodyJson)

	},
}

func init() {
	chatCmd.AddCommand(configCmd)

}
