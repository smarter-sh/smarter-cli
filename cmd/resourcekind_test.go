package cmd

import (
	"strings"
	"testing"
)

func TestArticle(t *testing.T) {
	cases := []struct {
		display string
		want    string
	}{
		{"ApiConnection", "an"},
		{"Orchestrator", "an"},
		{"Guardrail", "a"},
		{"SqlConnection", "a"},
		{"", "a"},
	}
	for _, tc := range cases {
		if got := Article(tc.display); got != tc.want {
			t.Errorf("Article(%q) = %q, want %q", tc.display, got, tc.want)
		}
	}
}

func sampleKind() ResourceKind {
	return ResourceKind{
		Singular:      "widget",
		Plural:        "widgets",
		APIKind:       "Widget",
		Display:       "Widget",
		DisplayPlural: "Widgets",
		Deployable:    true,
	}
}

func TestDescribeSpec(t *testing.T) {
	spec := sampleKind().DescribeSpec()

	if spec.Use != "widget <name>" {
		t.Errorf("Use = %q", spec.Use)
	}
	if spec.APIKind != "Widget" {
		t.Errorf("APIKind = %q", spec.APIKind)
	}
	if spec.NameArg.Mode != NameArgPositional || spec.NameArg.Kwarg != "name" {
		t.Errorf("NameArg = %+v", spec.NameArg)
	}
	if !strings.Contains(spec.Short, "a Widget") {
		t.Errorf("Short = %q, want it to use the article for Widget", spec.Short)
	}
	if !strings.Contains(spec.Long, "smarter describe widget") {
		t.Errorf("Long = %q, want an example invocation", spec.Long)
	}
}

func TestManifestSpec(t *testing.T) {
	spec := sampleKind().ManifestSpec()

	if spec.Use != "widget [flags]" {
		t.Errorf("Use = %q", spec.Use)
	}
	if spec.APIKind != "Widget" {
		t.Errorf("APIKind = %q", spec.APIKind)
	}
	if spec.NameArg.Mode != NameArgNone {
		t.Errorf("NameArg.Mode = %v, want NameArgNone", spec.NameArg.Mode)
	}
}

func TestGetSpec(t *testing.T) {
	spec := sampleKind().GetSpec()

	if spec.Use != "widgets" {
		t.Errorf("Use = %q", spec.Use)
	}
	if spec.NameArg.Mode != NameArgFlag || spec.NameArg.Kwarg != "name" || spec.NameArg.Shorthand != "n" {
		t.Errorf("NameArg = %+v", spec.NameArg)
	}
	if !strings.Contains(spec.Short, "Widgets") {
		t.Errorf("Short = %q, want the plural display name", spec.Short)
	}
}

func TestDeploySpec(t *testing.T) {
	spec := sampleKind().DeploySpec()

	if spec.Use != "widget <name>" {
		t.Errorf("Use = %q", spec.Use)
	}
	if spec.NameArg.Mode != NameArgPositional || spec.NameArg.Kwarg != "name" {
		t.Errorf("NameArg = %+v", spec.NameArg)
	}
}

func TestUndeploySpec(t *testing.T) {
	spec := sampleKind().UndeploySpec()

	if spec.Use != "widget <name>" {
		t.Errorf("Use = %q", spec.Use)
	}
	if !strings.Contains(spec.Short, "Undo") {
		t.Errorf("Short = %q, want it to describe undoing the deployment", spec.Short)
	}
}

func TestByUse(t *testing.T) {
	kinds := []ResourceKind{sampleKind(), {Singular: "gadget", APIKind: "Gadget"}}

	got, ok := ByUse(kinds, "gadget")
	if !ok || got.APIKind != "Gadget" {
		t.Errorf("ByUse(gadget) = %+v, %v", got, ok)
	}

	_, ok = ByUse(kinds, "missing")
	if ok {
		t.Error("ByUse(missing) = true, want false")
	}
}
