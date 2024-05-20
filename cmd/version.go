/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"
)

// versionCmd represents the status command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Retrieve version information",
	Long: `Retrieve version information:

smarter version

Returns version information about this software.`,
	Run: func(cmd *cobra.Command, args []string) {

		localVersion := []byte(`{"version":"` + Version + `"}`)
		kwargs := map[string]string{}
		bodyJson, err := APIRequest("whoami", kwargs)
		if err != nil {
			panic(err)
		} else {
			var localVersionMap map[string]interface{}
			err := json.Unmarshal(localVersion, &localVersionMap)
			if err != nil {
				log.Fatalf("Failed to unmarshal local version: %v", err)
			}

			var bodyJsonMap map[string]interface{}
			err = json.Unmarshal(bodyJson, &bodyJsonMap)
			if err != nil {
				log.Fatalf("Failed to unmarshal body JSON: %v", err)
			}

			for k, v := range localVersionMap {
				bodyJsonMap[k] = v
			}
			combinedJson, err := json.Marshal(bodyJsonMap)
			if err != nil {
				panic(err)
			}

			ConsoleOutput(combinedJson)
		}

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
