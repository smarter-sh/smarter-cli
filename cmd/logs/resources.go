/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package logs

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "chat <session_key>",
		Short: "Returns the logs for a Chat session_key",
		Long: `Returns the logs for a Chat:

smarter logs chat <session_key>

The Smarter API will return the logs for a Chat session_key.`,
		APIKind: "Chat",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "chatbot <name>",
		Short: "Returns the logs for a ChatBot",
		Long: `Returns the logs for a ChatBot:

smarter logs chatbot <name>

The Smarter API will return the logs for the ChatBot.`,
		APIKind: "ChatBot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "chat-history <session_key>",
		Short: "Returns the logs for a ChatHistory session_key",
		Long: `Returns the logs for a ChatHistory:

smarter logs chat-history <session_key>

The Smarter API will return the logs for a ChatHistory session_key.`,
		APIKind: "ChatHistory",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
}

func init() {
	cmd.RegisterResources(logsCmd, legacySpecs, APIRequest, ConsoleOutput, cmd.ErrorOutput)
}
