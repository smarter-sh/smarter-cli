/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package logs

import (
	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/logs/<str:kind>/
	return cmd.APIRequest("logs/"+kind+"/", kwargs)

}

// ConsoleOutput adapts cmd.ConsoleOutput's error return to the OutputFunc
// signature RegisterResources expects.
func ConsoleOutput(bodyJson []byte) {
	if err := cmd.ConsoleOutput(bodyJson); err != nil {
		cmd.ErrorOutput(err)
	}
}

var logsCmd = &cobra.Command{
	Use:   "logs <kind> <name>",
	Short: "Returns the logs for a resource",
	Long: `Returns the logs for a resource:

smarter logs <kind> <name>

Returns the logs for the resource <kind> <name>.`,
}

func init() {
	cmd.RootCmd.AddCommand(logsCmd)
}
