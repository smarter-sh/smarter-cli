/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

// ResourceKinds is the canonical, single-sourced list of v0.14 resource
// kinds that follow the standard template shape (see ResourceKind). Each
// verb package (describe/get/manifest/deploy) generates its own leaf
// ResourceSpecs from this list instead of maintaining a parallel copy of the
// Use/APIKind/display text.
//
// There should be a 1:1 mapping of these kinds against those in these
// two Python enum classes:
// - https://github.com/smarter-sh/smarter/blob/main/smarter/smarter/apps/api/v1/manifests/enum.py#L53
// - https://github.com/smarter-sh/smarter/blob/main/smarter/smarter/lib/journal/enum.py#L56

var ResourceKinds = []ResourceKind{
	{Singular: "apiconnection", Plural: "apiconnections", APIKind: "ApiConnection", Display: "ApiConnection", DisplayPlural: "ApiConnections"},
	{Singular: "sqlconnection", Plural: "sqlconnections", APIKind: "SqlConnection", Display: "SqlConnection", DisplayPlural: "SqlConnections"},
	{Singular: "apikey", Plural: "apikeys", APIKind: "SmarterAuthToken", Display: "SmarterAuthToken", DisplayPlural: "SmarterAuthTokens"},
	{Singular: "apiplugin", Plural: "apiplugins", APIKind: "ApiPlugin", Display: "ApiPlugin", DisplayPlural: "ApiPlugins"},
	{Singular: "sqlplugin", Plural: "sqlplugins", APIKind: "SqlPlugin", Display: "SqlPlugin", DisplayPlural: "SqlPlugins"},
	{Singular: "skillplugin", Plural: "skillplugins", APIKind: "SkillPlugin", Display: "SkillPlugin", DisplayPlural: "SkillPlugins"},
	{Singular: "guardrail", Plural: "guardrails", APIKind: "Guardrail", Display: "Guardrail", DisplayPlural: "Guardrails"},
	{Singular: "llmclient", Plural: "llmclients", APIKind: "LlmClient", Display: "LlmClient", DisplayPlural: "LlmClients"},
	{Singular: "llmhost", Plural: "llmhosts", APIKind: "LLMHost", Display: "LLMHost", DisplayPlural: "LLMHosts"},
	{Singular: "mcpclient", Plural: "mcpclients", APIKind: "MCPClient", Display: "MCPClient", DisplayPlural: "MCPClients"},
	{Singular: "orchestrator", Plural: "orchestrators", APIKind: "Orchestrator", Display: "Orchestrator", DisplayPlural: "Orchestrators"},
	{Singular: "prompt", Plural: "prompts", APIKind: "Prompt", Display: "Prompt", DisplayPlural: "Prompts", Deployable: true},
	{Singular: "promptconfig", Plural: "promptconfigs", APIKind: "PromptConfig", Display: "PromptConfig", DisplayPlural: "PromptConfigs"},
	{Singular: "provider", Plural: "providers", APIKind: "Provider", Display: "Provider", DisplayPlural: "Providers"},
	{Singular: "proxy", Plural: "proxies", APIKind: "Proxy", Display: "Proxy", DisplayPlural: "Proxies"},
	{Singular: "secret", Plural: "secrets", APIKind: "Secret", Display: "Secret", DisplayPlural: "Secrets"},
	{Singular: "vectorsearch", Plural: "vectorsearches", APIKind: "Vectorsearch", Display: "Vectorsearch", DisplayPlural: "Vectorsearches"},
	{Singular: "vectorstore", Plural: "vectorstores", APIKind: "Vectorstore", Display: "Vectorstore", DisplayPlural: "Vectorstores"},
}
