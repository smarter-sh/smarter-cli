/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(slug string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("manifest/"+slug, kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Generate an example manifest for the resource kind",
	Long: `Generate an example manifest for the resource kind. For example:

	smarter manifest <kind> > my-plugin.yaml

This will generate an example manifest for the specified kind of resource and write it to my-plugin.yaml in the current working directory.`,
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
