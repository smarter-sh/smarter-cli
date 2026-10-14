/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompt an LLMClient and echo its response to the console",
	Long: `Prompts an LLMClient and echos its response to the console:

	smarter chat prompt --llmclient <name> --prompt <text> [--new_session]

The Smarter API will send the prompt to the LLMClient and return its response.
Prompts from this machine share a session until --new_session starts another.`,
	Run: func(c *cobra.Command, args []string) {

		prompt, _ := c.Flags().GetString("prompt")
		if prompt == "" {
			cmd.ErrorOutput(errors.New("the 'prompt' flag is required"))
			return
		}
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

		// a map of strings always marshals
		fileContentsBytes, _ := json.Marshal(map[string]string{"prompt": prompt})

		// this request goes to /api/v1/cli/prompt/<str:name>/
		bodyJson, err := APIRequest(fmt.Sprintf("%s/", name), kwargs, string(fileContentsBytes))
		if err != nil {
			cmd.ErrorOutput(err)
			return
		}
		ConsoleOutput(bodyJson)

	},
}

func init() {
	chatCmd.AddCommand(promptCmd)

	promptCmd.PersistentFlags().StringP("prompt", "p", "", "A prompt to send to the LLMClient")
	if err := viper.BindPFlag("prompt", promptCmd.PersistentFlags().Lookup("prompt")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
