/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

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
	Long: `Retrieve a list of ChatBots or a specific ChatBot by name:

smarter get chatbots --name --json --yaml

The Smarter API will return a list of ChatBots in the specified format,
or a manifest for a specific ChatBot.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		name := viper.GetString("name")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"name": name,
			"n":    strconv.Itoa(n),
		}

		bodyJson, err := GetAPIResponse("chatbots", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	GetCmd.AddCommand(chatbotsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatbotsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatbotsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pluginsCmd.Flags().String("name", "", "Name of the chatbot")
	pluginsCmd.Flags().Int("n", 10, "Number of chatbots to retrieve")

	if err := viper.BindPFlag("name", pluginsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	if err := viper.BindPFlag("n", pluginsCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
