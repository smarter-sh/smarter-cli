/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"encoding/json"

	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func getFilePath(slug string) (string, error) {

	bodyBytes, err := cmd.GetAPIResponse(slug)
	if err != nil {
		return "", err
	} else {
		var body map[string]interface{}
		err = json.Unmarshal(bodyBytes, &body)
		if err != nil {
			return "", err
		}

		if filepath, ok := body["filepath"].(string); ok {
			return filepath, nil
		} else {
			panic("filepath not found or not a string")
		}
	}

}

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Generate an example manifest for the resource kind",
	Long: `Generate an example manifest for the resource kind. For example:

	smarter manifest plugin > my-plugin.yaml

This will generate an example manifest for a plugin resource and write it to my-plugin.yaml in the current working directory.`,
}

func init() {
	cmd.RootCmd.AddCommand(manifestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
