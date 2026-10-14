/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

// ResourceKinds is the canonical, single-sourced list of resource kinds that
// follow the standard template shape (see ResourceKind). Each verb package
// (get, describe, delete, logs, manifest, json-schema, deploy, undeploy)
// generates its own leaf ResourceSpecs from this list instead of maintaining
// a parallel copy of the Use/APIKind/display text.
//
// Together with the legacySpecs of the verb packages (account, plugin, user),
// it covers every SAMKind of the platform's OpenAPI document, which
// TestCommandsMatchTheOpenAPIDocument checks against the pinned copy in
// testdata/contracts/openapi.json: see doc/adr/022-pinned-openapi.md.
//
// Deployable is true only for kinds whose platform broker implements deploy,
// as the platform's resources command reports. CustomDomain's and
// LLMHostCompute's brokers deploy but don't undeploy, so their undeploy
// commands return the platform's "not implemented" error.

var ResourceKinds = []ResourceKind{
	{Singular: "apiconnection", Plural: "apiconnections", APIKind: "ApiConnection", Display: "ApiConnection", DisplayPlural: "ApiConnections"},
	{Singular: "apikey", Plural: "apikeys", APIKind: "SmarterAuthToken", Display: "SmarterAuthToken", DisplayPlural: "SmarterAuthTokens", Deployable: true},
	{Singular: "apiplugin", Plural: "apiplugins", APIKind: "ApiPlugin", Display: "ApiPlugin", DisplayPlural: "ApiPlugins"},
	{Singular: "budget", Plural: "budgets", APIKind: "Budget", Display: "Budget", DisplayPlural: "Budgets"},
	{Singular: "customdomain", Plural: "customdomains", APIKind: "CustomDomain", Display: "CustomDomain", DisplayPlural: "CustomDomains", Deployable: true},
	{Singular: "guardrail", Plural: "guardrails", APIKind: "Guardrail", Display: "Guardrail", DisplayPlural: "Guardrails"},
	{Singular: "imagesearchplugin", Plural: "imagesearchplugins", APIKind: "ImageSearchPlugin", Display: "ImageSearchPlugin", DisplayPlural: "ImageSearchPlugins"},
	{Singular: "llmclient", Plural: "llmclients", APIKind: "LLMClient", Display: "LLMClient", DisplayPlural: "LLMClients", Deployable: true},
	{Singular: "llmhost", Plural: "llmhosts", APIKind: "LLMHost", Display: "LLMHost", DisplayPlural: "LLMHosts", Deployable: true},
	{Singular: "llmhostcompute", Plural: "llmhostcompute", APIKind: "LLMHostCompute", Display: "LLMHostCompute", DisplayPlural: "LLMHostCompute", Deployable: true},
	{Singular: "mcpclient", Plural: "mcpclients", APIKind: "MCPClient", Display: "MCPClient", DisplayPlural: "MCPClients"},
	{Singular: "orchestrator", Plural: "orchestrators", APIKind: "Orchestrator", Display: "Orchestrator", DisplayPlural: "Orchestrators"},
	{Singular: "prompt", Plural: "prompts", APIKind: "Prompt", Display: "Prompt", DisplayPlural: "Prompts"},
	{Singular: "provider", Plural: "providers", APIKind: "Provider", Display: "Provider", DisplayPlural: "Providers"},
	{Singular: "proxy", Plural: "proxies", APIKind: "Proxy", Display: "Proxy", DisplayPlural: "Proxies"},
	{Singular: "secret", Plural: "secrets", APIKind: "Secret", Display: "Secret", DisplayPlural: "Secrets"},
	{Singular: "sqlconnection", Plural: "sqlconnections", APIKind: "SqlConnection", Display: "SqlConnection", DisplayPlural: "SqlConnections"},
	{Singular: "sqlplugin", Plural: "sqlplugins", APIKind: "SqlPlugin", Display: "SqlPlugin", DisplayPlural: "SqlPlugins"},
	{Singular: "skillplugin", Plural: "skillplugins", APIKind: "SkillPlugin", Display: "SkillPlugin", DisplayPlural: "SkillPlugins"},
	{Singular: "vectorsearch", Plural: "vectorsearches", APIKind: "Vectorsearch", Display: "VectorSearch", DisplayPlural: "VectorSearches"},
	{Singular: "vectorstore", Plural: "vectorstores", APIKind: "Vectorstore", Display: "VectorStore", DisplayPlural: "VectorStores", Deployable: true},
	{Singular: "websearchplugin", Plural: "websearchplugins", APIKind: "WebsearchPlugin", Display: "WebsearchPlugin", DisplayPlural: "WebsearchPlugins"},
}
