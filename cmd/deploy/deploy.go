/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package deploy

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("deploy/"+kind+"/", kwargs)

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
}
