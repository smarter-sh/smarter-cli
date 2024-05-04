/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve real-time status of the Smarter Platform",
	Long: `Retrieve real-time status of the Smarter Platform:

smarter get status --json --yaml

The Smarter API will return the current status of the Smarter Platform,
including the status of all services and resources by region.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPIResponse("status")
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
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
