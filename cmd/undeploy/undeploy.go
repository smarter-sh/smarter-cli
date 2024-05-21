/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package undeploy

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(slug string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("undeploy/"+slug, kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}

// UndeployCmd represents the get command
var UndeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Generate a list of Smarter resources or a manifest for a specific resource",
	Long: `Generate a list of Smarter resources or a manifest for a specific resource:

smarter get <kind> --name --json --yaml --csv --xml -n 10 --asc --desc

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(UndeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// UndeployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// UndeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
