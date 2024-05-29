/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"fmt"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("delete/"+kind+"/", kwargs, false)

}
func ConsoleOutput() {
	fmt.Println("deleted.")
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <kind> <name>",
	Short: "Permanently delete a Smarter resource",
	Long: `Permanently delete a Smarter resource:

smarter delete <kind> <name> --dry-run

The Smarter API will permanently delete the resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(deleteCmd)
}
