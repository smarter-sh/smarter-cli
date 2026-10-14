/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package jsonschema

import (
	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/json-schema/<str:kind>/
	return cmd.APIRequest("json-schema/"+kind+"/", kwargs)

}

// ConsoleOutput prints the JSON Schema itself, as json unless
// --output_format says otherwise, so that it can be saved and used by an editor.
func ConsoleOutput(bodyJson []byte) {
	if err := cmd.DataOutput(bodyJson, "json"); err != nil {
		cmd.ErrorOutput(err)
	}
}

var jsonSchemaCmd = &cobra.Command{
	Use:   "json-schema <kind> [flags]",
	Short: "Retrieve the JSON Schema of a resource kind's manifest",
	Long: `Retrieves the JSON Schema of a resource kind's manifest. For example:

	smarter json-schema <kind> > <kind>.schema.json

This is the schema that the Smarter API validates manifests of that kind against.`,
}

func init() {
	cmd.RootCmd.AddCommand(jsonSchemaCmd)
}
