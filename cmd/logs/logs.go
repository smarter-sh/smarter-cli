/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package logs

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("logs/"+kind+"/", kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}

// logsCmd represents the get command
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
