/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package chat

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// APIRequest sends a chat request for the LLMClient named in slug.
func APIRequest(slug string, kwargs map[string]string, fileContents ...string) ([]byte, error) {
	var fileContent string
	if len(fileContents) > 0 {
		fileContent = fileContents[0]
	}

	// en route to either of:
	// 		/api/v1/cli/prompt/<str:name>/
	// 		/api/v1/cli/prompt/config/<str:name>/
	return cmd.APIRequest(fmt.Sprintf("prompt/%s", slug), kwargs, fileContent)

}

func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "json")
	}
	if err := cmd.ConsoleOutput(bodyJson); err != nil {
		cmd.ErrorOutput(err)
	}
}

// llmClientName returns --llmclient, or the deprecated --chatbot.
func llmClientName(c *cobra.Command) (string, error) {
	name, _ := c.Flags().GetString("llmclient")
	if name == "" {
		name, _ = c.Flags().GetString("chatbot")
	}
	if name == "" {
		return "", errors.New("the 'llmclient' flag is required")
	}
	return name, nil
}

// getUniqueID identifies this machine to the platform, which caches the chat
// session under it: the hostname and the first hardware address. Either
// part may be empty when the machine doesn't report it.
func getUniqueID() string {
	var macAddr string
	if interfaces, err := net.Interfaces(); err == nil {
		for _, inter := range interfaces {
			if inter.HardwareAddr != nil {
				macAddr = inter.HardwareAddr.String()
				break
			}
		}
	}

	host, err := os.Hostname()
	if err != nil {
		host = "unknown-host"
	}

	// URL encode the host and macAddr
	host = url.QueryEscape(host)
	macAddr = url.QueryEscape(macAddr)

	return fmt.Sprintf("%s-%s", host, macAddr)
}

var chatCmd = &cobra.Command{
	Use:   "chat [command] [flags]",
	Short: "Chat with an LLMClient",
	Long: `Chat with an LLMClient:

smarter chat <command> --llmclient <name> [flags]

The Smarter API will send the prompt to the LLMClient and
then echo its response to the console.`,
}

func init() {
	cmd.RootCmd.AddCommand(chatCmd)

	chatCmd.PersistentFlags().StringP("llmclient", "l", "", "the name of an LLMClient")
	if err := viper.BindPFlag("llmclient", chatCmd.PersistentFlags().Lookup("llmclient")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	chatCmd.PersistentFlags().StringP("chatbot", "c", "", "the name of an LLMClient (deprecated: use --llmclient)")
	if err := chatCmd.PersistentFlags().MarkDeprecated("chatbot", "use --llmclient"); err != nil {
		log.Fatalf("Error deprecating flag: %v", err)
	}

	chatCmd.PersistentFlags().BoolP("new_session", "n", false, "start a new session")
	if err := viper.BindPFlag("new_session", chatCmd.PersistentFlags().Lookup("new_session")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

}
