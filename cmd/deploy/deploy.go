/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package deploy

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(slug string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("deploy/"+slug, kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}

// DeployCmd represents the get command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a resource",
	Long: `Deploys a resource:

smarter deploy <kind> <name> --json --yaml

The Smarter API will deploy the resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(DeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DeployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
