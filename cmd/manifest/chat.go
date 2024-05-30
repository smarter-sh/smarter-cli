/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var chatsCmd = &cobra.Command{
	Use:   "chat [flags]",
	Short: "Generate an example manifest for a chat session.",
	Long: `Generates an example manifest for a chat session. For example:

	smarter manifest chat [flags] > my-plugin.yaml

This will generate an example manifest a chat session and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/chat/
		bodyJson, err := APIRequest("chat", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(chatsCmd)
}
