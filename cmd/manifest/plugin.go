/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manifestCmd represents the manifest command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Generate an example Plugin manifest",
	Long: `Generate an example Plugin manifest. For example:

	smarter manifest plugin > my-plugin.yaml

This will generate an example manifest for a plugin resource and write it to my-plugin.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("manifest/plugin")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			if filepath, ok := body["filepath"].(string); ok {
				url := filepath
				fmt.Println("URL:", url)
				contents, err := GetAndPrintYAMLResponse(url, "plugin")
				if err != nil {
					fmt.Println("Error reading file:", err)
				} else {
					fmt.Println(contents)
				}
			} else {
				fmt.Println("Error: filepath not found or not a string")
			}
		}
	},
}

func init() {
	manifestCmd.AddCommand(pluginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
