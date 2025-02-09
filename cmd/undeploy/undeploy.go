/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package undeploy

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/undeploy/<str:kind>/
	return cmd.APIRequest("undeploy/"+kind+"/", kwargs)

}
func ConsoleOutput() {
	fmt.Println("undeployed.")
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

var undeployCmd = &cobra.Command{
	Use:   "undeploy <kind> <name>",
	Short: "Undo a Smarter resource deployment.",
	Long: `Undo a Smarter resource deployment. For example:

smarter undeploy <kind> <name>

The Smarter API will undo the deployment of the resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(undeployCmd)
}
