package cmd

import "testing"

// TestResourceKindsAreWellFormed guards against copy-paste mistakes in the
// hand-maintained ResourceKinds table: every entry must have all of its
// identity fields populated, and singular/plural command names must be
// unique so RegisterResourceCmd never registers two commands with the same
// name on a parent.
func TestResourceKindsAreWellFormed(t *testing.T) {
	seenSingular := map[string]bool{}
	seenPlural := map[string]bool{}

	for _, k := range ResourceKinds {
		if k.Singular == "" || k.Plural == "" || k.APIKind == "" || k.Display == "" || k.DisplayPlural == "" {
			t.Errorf("ResourceKind %+v has an empty required field", k)
		}
		if seenSingular[k.Singular] {
			t.Errorf("duplicate Singular %q in ResourceKinds", k.Singular)
		}
		seenSingular[k.Singular] = true

		if seenPlural[k.Plural] {
			t.Errorf("duplicate Plural %q in ResourceKinds", k.Plural)
		}
		seenPlural[k.Plural] = true
	}
}

// TestResourceKindsDeployableSubset pins the kinds whose platform broker
// implements deploy, as the platform's resources command reports them (User,
// also deployable, is a legacy spec). Prompt's broker raises "not
// implemented", so Prompt is not among them.
func TestResourceKindsDeployableSubset(t *testing.T) {
	want := map[string]bool{
		"apikey": true, "customdomain": true, "llmclient": true,
		"llmhost": true, "llmhostcompute": true, "vectorstore": true,
	}
	for _, k := range ResourceKinds {
		if k.Deployable != want[k.Singular] {
			t.Errorf("ResourceKind %q has Deployable = %v, want %v", k.Singular, k.Deployable, want[k.Singular])
		}
	}
}
