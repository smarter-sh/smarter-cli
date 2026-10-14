/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

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
			if err := json.Unmarshal(localVersion, &localVersionMap); err != nil {
				ErrorOutput(fmt.Errorf("failed to unmarshal local version: %w", err))
				return
			}

			var bodyJsonMap map[string]interface{}
			if err := json.Unmarshal(bodyJson, &bodyJsonMap); err != nil {
				ErrorOutput(fmt.Errorf("failed to unmarshal body JSON: %w", err))
				return
			}

			for k, v := range localVersionMap {
				bodyJsonMap[k] = v
			}
			combinedJson, err := json.Marshal(bodyJsonMap)
			if err != nil {
				ErrorOutput(err)
				return
			}

			if err := ConsoleOutput(combinedJson); err != nil {
				ErrorOutput(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
