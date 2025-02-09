/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package chat

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIRequest(slug string, kwargs map[string]string, fileContents ...string) ([]byte, error) {
	var fileContent string
	if len(fileContents) > 0 {
		fileContent = fileContents[0]
	}

	// en route to either of:
	// 		/api/v1/cli/chat/<str:chatbot>/<str:uid>
	// 		/api/v1/cli/chat/config/<str:chatbot>/<str:uid>
	return cmd.APIRequest(fmt.Sprintf("chat/%s", slug), kwargs, fileContent)

}
func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "json")
	}
	cmd.ConsoleOutput(bodyJson)
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

func getUniqueID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var macAddr string
	for _, inter := range interfaces {
		if inter.HardwareAddr != nil {
			macAddr = inter.HardwareAddr.String()
			break
		}
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// URL encode the host and macAddr
	host = url.QueryEscape(host)
	macAddr = url.QueryEscape(macAddr)

	return fmt.Sprintf("%s-%s", host, macAddr)
}

var chatCmd = &cobra.Command{
	Use:   "chat [command] [flags]",
	Short: "Chat with a deployed ChatBot",
	Long: `Chat with a deployed ChatBot:

smarter chat <command> [flags]

The Smarter API will send the prompt to a deployed ChatBot and
then echo its response to the console.`,
}

func init() {
	cmd.RootCmd.AddCommand(chatCmd)

	chatCmd.PersistentFlags().StringP("chatbot", "c", "", "the name of a deployed ChatBot")
	if err := viper.BindPFlag("chatbot", chatCmd.PersistentFlags().Lookup("chatbot")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	chatCmd.PersistentFlags().BoolP("new_session", "n", false, "start a new session")
	if err := viper.BindPFlag("new_session", chatCmd.PersistentFlags().Lookup("new_session")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
