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

	return get.APIRequest(slug, kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	get.ConsoleOutput(bodyJson)
}
func ErrorOutput(err error) {
	get.ErrorOutput(err)
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
}
