/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the status command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Retrieve information about the api_key owner",
	Long: `Retrieve information about the api_key owner:

smarter whoami --json --yaml

Returns informtation about the Smarter user account that owns the
configured api_key.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPIResponse("whoami")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			bodyJson, err := json.Marshal(body)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Response:", string(bodyJson))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(whoamiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}