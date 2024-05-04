/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a Smarter manifest",
	Long: `Apply a Smarter manifest:

smarter apply -f <manifest.yaml> --json --yaml --dry-run

The Smarter API will apply the manifest to the Smarter account,
migrating the resource to the new state. The --json and --yaml
flags will output the manifest in the specified format. The
--dry-run flag will simulate the apply without making any changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPIResponse("apply")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			bodyStr, err := json.Marshal(body)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Response:", string(bodyStr))
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
