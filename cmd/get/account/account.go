/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package account

import (
	"github.com/QueriumCorp/smarter-cli/cmd/get"

	"github.com/spf13/cobra"
)

func GetAPIResponse(slug string) ([]byte, error) {

	kwargs := map[string]string{}

	return get.GetAPIResponse(slug, kwargs)

}
func ConsoleOutput(bodyJson []byte, jsonFlagValue bool, yamlFlagValue bool) {
	get.ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
}

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Retrieve your Account manifest",
	Long: `Retrieve your Account manifest:

smarter account --json --yaml

The Smarter API will return your Account manifest in the specified format.`,
}

func init() {
	get.GetCmd.AddCommand(accountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
