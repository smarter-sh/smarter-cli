/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Retrieve version information",
	Long: `Retrieve version information:

smarter version

Returns version information about this software.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := viper.GetBool("verbose")
		localVersion := []byte(`{"version":"` + Version + `"}`)
		if !verbose {
			fmt.Printf("Local version: %s\n", Version)
			return
		}
		kwargs := map[string]string{}

		// this request goes to /api/v1/cli/version/
		bodyJson, err := APIRequest("version", kwargs)
		if err != nil {
			ErrorOutput(err)
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
				ErrorOutput(err)
			}

			ConsoleOutput(combinedJson)
		}

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
