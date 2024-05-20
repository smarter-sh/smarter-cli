/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatbotsCmd represents the chatbots command
var chatbotsCmd = &cobra.Command{
	Use:   "chatbots",
	Short: "Retrieve a list of ChatBots or a specific ChatBot by name",
	Long: `Generate an example manifest for a chatbot. For example:

	smarter manifest chatbot > my-plugin.yaml

This will generate an example manifest for a chatbot and write it to my-plugin.yaml in the current working directory.`,

	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"name": name,
			"n":    strconv.Itoa(n),
		}

		bodyJson, err := APIRequest("chatbots", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	manifestCmd.AddCommand(chatbotsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chatbotsCmd.Flags().String("name", "", "Name of the chatbot")
	chatbotsCmd.Flags().Int("n", 10, "Number of chatbots to retrieve")

	if err := viper.BindPFlag("name", chatbotsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

	if err := viper.BindPFlag("n", chatbotsCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
