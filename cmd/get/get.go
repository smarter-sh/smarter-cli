/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package get

import (
	"fmt"

	"github.com/QueriumCorp/smarter-cli/cmd"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Generate a list of Smarter resources or a manifest for a specific resource",
	Long: `Generate a list of Smarter resources or a manifest for a specific resource:

smarter get <kind> --name --json --yaml --csv --xml -n 10 --asc --desc

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
