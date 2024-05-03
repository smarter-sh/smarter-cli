/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manifestCmd represents the manifest command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Generate an example Chat manifest",
	Long: `Generate an example Chat manifest. For example:

	smarter manifest chat > my-chat.yaml

This will generate an example manifest for your Chat and write it to my-chat.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		contents, err := getYamlFileContents("chat")
		if err != nil {
			fmt.Println("Error reading file:", err)
		} else {
			fmt.Println(contents)
		}
	},
}

func init() {
	manifestCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
