/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manifestCmd represents the manifest command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Generate an example Account manifest",
	Long: `Generate an example Account manifest. For example:

	smarter manifest account > my-account.yaml

This will generate an example manifest for your Account and write it to my-account.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		contents, err := getYamlFileContents("account")
		if err != nil {
			fmt.Println("Error reading file:", err)
		} else {
			fmt.Println(contents)
		}
	},
}

func init() {
	manifestCmd.AddCommand(accountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}