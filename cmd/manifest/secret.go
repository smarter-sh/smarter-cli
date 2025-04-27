/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secret [flags]",
	Short: "Generate an example manifest for a secret.",
	Long: `Generates an example manifest for a secret. For example:

	smarter manifest secret [flags] > my-secret.yaml

This will generate an example manifest for a secret and write it to my-secret.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/manifest/secret/
		bodyJson, err := APIRequest("secret", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(secretsCmd)
}
