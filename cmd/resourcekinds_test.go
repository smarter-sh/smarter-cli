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

func TestResourceKindsDeployableSubset(t *testing.T) {
	found := false
	for _, k := range ResourceKinds {
		if k.Deployable {
			found = true
			if k.Singular != "prompt" {
				t.Errorf("unexpected deployable kind %q; deploy/undeploy tests assume only 'prompt' is deployable", k.Singular)
			}
		}
	}
	if !found {
		t.Error("expected at least one ResourceKind to be marked Deployable")
	}
}
