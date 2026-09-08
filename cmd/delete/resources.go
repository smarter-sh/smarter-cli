/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package delete

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

var legacySpecs = []cmd.ResourceSpec{
	{
		Use:   "apiconnection <name>",
		Short: "Delete a PluginDataApiConnection",
		Long: `Deletes a PluginDataApiConnection:

smarter delete apiconnection <name>  [flags]

The Smarter API will permanently delete the PluginDataApiConnection with the specified name.`,
		APIKind: "PluginDataApiConnection",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "apikey <name>",
		Short: "Delete a SmarterAuthToken",
		Long: `Deletes a SmarterAuthToken:

smarter delete apikey <name> [flags]

The Smarter API will permanently delete the SmarterAuthToken with the specified name.`,
		APIKind: "SmarterAuthToken",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "chat <session_key>",
		Short: "Delete a chat history",
		Long: `Deletes a chat history:

smarter delete chat <session_key>

The Smarter API will permanently delete the chat history with the specified identifier.`,
		APIKind: "chat",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "session_key"},
	},
	{
		Use:   "chatbot <name>",
		Short: "Delete a ChatBot",
		Long: `Delete a ChatBot:

smarter delete chatbot <name> --dry-run

The Smarter API will permanently delete the ChatBot with the specified name,
and all related chat history.`,
		APIKind: "chatbot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "plugin <name>",
		Short: "Delete a Plugin",
		Long: `Delete a Plugin:

smarter delete plugin <name> --dry-run

The Smarter API will permanently delete the Plugin with the specified name,
and dissassociate it from any ChatBots.`,
		APIKind: "plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "sqlconnection <name>",
		Short: "Delete a PluginDataSqlConnection",
		Long: `Deletes a PluginDataSqlConnection:

smarter delete sqlconnection <name> [flags]

The Smarter API will permanently delete the PluginDataSqlConnection with the specified name.`,
		APIKind: "PluginDataSqlConnection",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "name"},
	},
	{
		Use:   "user <username>",
		Short: "Delete a user from your account",
		Long: `Delete a user from your account:

smarter delete <username> --dry-run

The Smarter API will permanently delete the user with the specified name,
and dissassociate it from any Smarter resources. Your Smarter admin account
will replace the deleted user.`,
		APIKind: "user",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgPositional, Kwarg: "username"},
	},
}

func init() {
	cmd.RegisterResources(deleteCmd, legacySpecs, APIRequest, cmd.AdaptOutput(ConsoleOutput), ErrorOutput)
}
