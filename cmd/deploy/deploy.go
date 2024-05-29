/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package deploy

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("deploy/"+kind+"/", kwargs, false)

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
