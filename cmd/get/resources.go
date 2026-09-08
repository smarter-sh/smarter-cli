/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"github.com/smarter-sh/smarter-cli/cmd"
)

// legacySpecs are resources whose shape (custom flags, non-"name" identifying
// parameter, or organic Short/Long text) pre-dates the v0.14 ResourceKind
// template and doesn't fit it, so they're kept as explicit ResourceSpecs.
var legacySpecs = []cmd.ResourceSpec{
	{
		Use:     "account",
		Short:   "Retrieve your Account manifest",
		Long:    "Retrieve your Account manifest:\n\nsmarter get account [flags]\n\nThe Smarter API will your Account manifest.",
		APIKind: "Account",
	},
	{
		Use:   "chatbots",
		Short: "Retrieve a list of ChatBots",
		Long: `Retrieve a list of ChatBots:

smarter get chatbots [flags]

The Smarter API will return a list of ChatBots in the specified format,
or a manifest for a specific ChatBot.`,
		APIKind: "Chatbot",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "name", Shorthand: "n", Usage: "Name of the chatbot"},
	},
	{
		Use:   "chat-history",
		Short: "Retrieve the chat history for a session_key",
		Long: `Retrieve the chat history for a session_key:

smarter get chat-history [session_key]

The Smarter API will return the chat history for the session_key.`,
		APIKind: "ChatHistory",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "session_key", Shorthand: "s", Usage: "Chat session_key", Required: true},
	},
	{
		Use:   "chat-plugin-usage",
		Short: "Retrieve the chat plugin usage for a session_key",
		Long: `Retrieve the chat plugin usage for a session_key:

smarter get chat-plugin-usage [session_key]

The Smarter API will return the chat plugin usage for the session_key.`,
		APIKind: "ChatPluginUsage",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "session_key", Shorthand: "s", Usage: "Chat session_key", Required: true},
	},
	{
		Use:   "chats",
		Short: "Retrieve a list of Chats",
		Long: `Retrieves a list of Chats:

smarter get chats [flags]

The Smarter API will return a list of Chats.`,
		APIKind: "Chat",
		Flags: []cmd.FlagSpec{
			{Name: "chatbot", Shorthand: "c", Kind: cmd.FlagString, Usage: "Name of the chatbot"},
			{Name: "session_key", Shorthand: "s", Kind: cmd.FlagString, Usage: "Chat session_key"},
			{Name: "today", Kind: cmd.FlagBool, Usage: "Filter for today"},
			{Name: "yesterday", Kind: cmd.FlagBool, Usage: "Filter for yesterday"},
			{Name: "this-week", Kind: cmd.FlagBool, Usage: "Filter for this week"},
			{Name: "last-week", Kind: cmd.FlagBool, Usage: "Filter for last week"},
			{Name: "this-month", Kind: cmd.FlagBool, Usage: "Filter for this month"},
			{Name: "last-month", Kind: cmd.FlagBool, Usage: "Filter for last month"},
		},
	},
	{
		Use:   "chat-tool-calls",
		Short: "Retrieve the chat tool calls for a session_key",
		Long: `Retrieve the chat tool calls for a session_key:

smarter get chat-tool-calls [session_key]

The Smarter API will return the chat tool calls for the session_key.`,
		APIKind: "ChatToolCall",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "session_key", Shorthand: "s", Usage: "Chat session_key", Required: true},
	},
	{
		Use:   "plugins",
		Short: "Retrieve a list of Plugins",
		Long: `Retrieves a list of Plugins:

smarter get plugins [flags]


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
		APIKind: "Plugin",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "name", Shorthand: "n", Usage: "Name of the plugin"},
		Flags: []cmd.FlagSpec{
			{Name: "class", Shorthand: "c", Kind: cmd.FlagString, Usage: "Plugin class: static, sql, api", Choices: []string{"static", "sql", "api"}},
		},
	},
	{
		Use:   "users",
		Short: "Retrieve a list of Users",
		Long: `Retrieves a list of Users:

smarter get users [flags]

The Smarter API will return a list of Users in the specified format,
or a manifest for a specific User.`,
		APIKind: "User",
		NameArg: cmd.NameArgSpec{Mode: cmd.NameArgFlag, Kwarg: "username", Shorthand: "u", Usage: "Smarter username"},
	},
}

func init() {
	cmd.RegisterResources(getCmd, legacySpecs, APIRequest, ConsoleOutput, ErrorOutput)

	specs := make([]cmd.ResourceSpec, len(cmd.V14ResourceKinds))
	for i, k := range cmd.V14ResourceKinds {
		specs[i] = k.GetSpec()
	}
	cmd.RegisterResources(getCmd, specs, APIRequest, ConsoleOutput, ErrorOutput)
}
