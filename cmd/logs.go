/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve a list of Log files or a specific Log file by name or identifier",
	Long: `Retrieve a list of Log files, or a specific Log file by name or identifier:

smarter get logs --name --json --yaml --csv --xml -n 10 --asc --desc

The Smarter API will return a list of Log files in the specified format,
or a manifest for a specific Log file.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPIResponse("logs")
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
	RootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
