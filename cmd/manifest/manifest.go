/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package manifest

import (
	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/example-manifest/<str:kind>/
	return cmd.APIRequest("example-manifest/"+kind+"/", kwargs)

}

// ConsoleOutput prints the example manifest itself, as yaml unless
// --output_format says otherwise, so that it can be saved and applied.
func ConsoleOutput(bodyJson []byte) {
	if err := cmd.DataOutput(bodyJson, "yaml"); err != nil {
		cmd.ErrorOutput(err)
	}
}

var manifestCmd = &cobra.Command{
	Use:   "manifest <kind> [flags]",
	Short: "Generate an example manifest for the resource kind",
	Long: `Generate an example manifest for the resource kind. For example:

	smarter manifest <kind> [flags] > manifest.yaml

This will generate an example manifest for the specified kind of resource and write it to manifest.yaml in the current working directory.`,
}

func init() {
	cmd.RootCmd.AddCommand(manifestCmd)
}
