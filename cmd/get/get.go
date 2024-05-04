/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

func GetAPI(slug string) (map[string]interface{}, error) {

	return cmd.GetAPIResponse(slug)

}

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Generate a list of Smarter resources or a manifest for a specific resource",
	Long: `Generate a list of Smarter resources or a manifest for a specific resource:

smarter get <kind> --name --json --yaml --csv --xml -n 10 --asc --desc

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(GetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// GetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
