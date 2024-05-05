/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package delete

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func GetAPI(slug string) ([]byte, error) {

	return cmd.GetAPIResponse(slug)

}

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Permanently delete a Smarter resource",
	Long: `Permanently delete a Smarter resource:

smarter delete <kind> --dry-run

The Smarter API will permanently delete the resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(DeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
