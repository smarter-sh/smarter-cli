/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import "fmt"

// ResourceKind is the canonical identity of a resource kind that follows the
// standard "describe/get/manifest/deploy <use> <name>" template shape. It's
// the single source of truth for that resource's command name, API kind
// string, and display text, so describe/get/manifest/deploy don't each carry
// their own copy that can drift out of sync. Resources whose Short/Long text
// is organic (pre-dates this template, e.g. "account" or "plugin") are kept
// as literal ResourceSpec entries in their verb package instead.
type ResourceKind struct {
	Singular      string // command name, e.g. "apiconnection"
	Plural        string // command name, e.g. "apiconnections"
	APIKind       string // the "kind" path segment sent to the Smarter API
	Display       string // human-readable singular name, e.g. "ApiConnection"
	DisplayPlural string // human-readable plural name, e.g. "ApiConnections"
	Deployable    bool   // whether this kind gets "deploy"/"undeploy" commands
}

// DescribeSpec builds this kind's "describe <name>" ResourceSpec.
func (k ResourceKind) DescribeSpec() ResourceSpec {
	a := Article(k.Display)
	return ResourceSpec{
		Use:   k.Singular + " <name>",
		Short: fmt.Sprintf("Retrieve %s %s manifest by name", a, k.Display),
		Long: fmt.Sprintf(`Retrieves a manifest for %s %s. For example:

	smarter describe %s <name> > my-plugin.yaml

This will generate a manifest for %s %s named <name> and write it to my-plugin.yaml in the current working directory.`, a, k.Display, k.Singular, a, k.Display),
		APIKind: k.APIKind,
		NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"},
	}
}

// ManifestSpec builds this kind's "manifest [flags]" ResourceSpec.
func (k ResourceKind) ManifestSpec() ResourceSpec {
	a := Article(k.Display)
	return ResourceSpec{
		Use:   k.Singular + " [flags]",
		Short: fmt.Sprintf("Generate an example manifest for %s %s.", a, k.Display),
		Long: fmt.Sprintf(`Generates an example manifest for %s %s. For example:

	smarter manifest %s [flags] > my-plugin.yaml

This will generate an example manifest for %s %s and write it to my-plugin.yaml in the current working directory.`, a, k.Display, k.Singular, a, k.Display),
		APIKind: k.APIKind,
	}
}

// GetSpec builds this kind's "get <plural> [flags]" list ResourceSpec.
func (k ResourceKind) GetSpec() ResourceSpec {
	return ResourceSpec{
		Use:   k.Plural,
		Short: fmt.Sprintf("Retrieve a list of %s", k.DisplayPlural),
		Long: fmt.Sprintf(`Retrieve a list of %s:

smarter get %s [flags]

The Smarter API will return a list of %s in the specified format,
or a manifest for a specific %s.`, k.DisplayPlural, k.Plural, k.DisplayPlural, k.Display),
		APIKind: k.APIKind,
		NameArg: NameArgSpec{Mode: NameArgFlag, Kwarg: "name", Shorthand: "n", Usage: fmt.Sprintf("Name of the %s", k.Display)},
	}
}

// DeploySpec builds this kind's "deploy <name>" ResourceSpec.
func (k ResourceKind) DeploySpec() ResourceSpec {
	a := Article(k.Display)
	return ResourceSpec{
		Use:   k.Singular + " <name>",
		Short: fmt.Sprintf("Deploy %s %s", a, k.Display),
		Long: fmt.Sprintf(`Deploys %s %s:

smarter deploy %s <name> [flags]

The Smarter API will deploy the %s.`, a, k.Display, k.Singular, k.Display),
		APIKind: k.APIKind,
		NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"},
	}
}

// UndeploySpec builds this kind's "undeploy <name>" ResourceSpec.
func (k ResourceKind) UndeploySpec() ResourceSpec {
	a := Article(k.Display)
	return ResourceSpec{
		Use:   k.Singular + " <name>",
		Short: fmt.Sprintf("Undo %s %s deployment", a, k.Display),
		Long: fmt.Sprintf(`Undo %s %s deployment. For example:

smarter undeploy %s <name>

This will reverse the effect of having deployed the %s.`, a, k.Display, k.Singular, k.Display),
		APIKind: k.APIKind,
		NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"},
	}
}

// ByUse returns the kind in kinds whose Singular field matches use.
func ByUse(kinds []ResourceKind, use string) (ResourceKind, bool) {
	for _, k := range kinds {
		if k.Singular == use {
			return k, true
		}
	}
	return ResourceKind{}, false
}
