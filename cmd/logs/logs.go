/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package logs

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(slug string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("logs/"+slug, kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}

// logsCmd represents the get command
var logsCmd = &cobra.Command{
	Use:   "logs <kind> <name>",
	Short: "Returns the logs for a resource",
	Long: `Returns the logs for a resource:

smarter logs <kind> <name> --json --yaml

Returns the logs for the resource <kind> <name>.`,
}

func init() {
	cmd.RootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
