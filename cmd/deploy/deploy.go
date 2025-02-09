/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package deploy

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/deploy/<str:kind>/
	return cmd.APIRequest("deploy/"+kind+"/", kwargs)

}
func ConsoleOutput() {
	fmt.Println("deployed.")
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

var deployCmd = &cobra.Command{
	Use:   "deploy <kind> <name>",
	Short: "Deploy a resource",
	Long: `Deploys a resource:

smarter deploy <kind> <name> [flags]

The Smarter API will deploy the resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(deployCmd)
}
