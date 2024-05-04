/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

// manifestCmd represents the manifest command
var chatbotCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "Generate an example ChatBot manifest",
	Long: `Generate an example ChatBot manifest. For example:

	smarter manifest chatbot > my-chatbot.yaml

This will generate an example manifest for your ChatBot and write it to my-chatbot.yaml in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		contents, err := getYamlFileContents("chatbot")
		if err != nil {
			fmt.Println("Error reading file:", err)
		} else {
			fmt.Println(contents)
		}
	},
}

func init() {
	manifestCmd.AddCommand(chatbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}