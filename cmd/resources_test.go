package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestResourcesCommand(t *testing.T) {
	var gotPath string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[{"Singular":"Secret","Plural":"Secrets"}]}`))
	})
	withViperValue(t, "output_format", "json")

	out := captureStdout(t, func() { resourcesCmd.Run(resourcesCmd, nil) })
	if gotPath != "/api/v1/cli/resources/" {
		t.Errorf("path = %q, want /api/v1/cli/resources/", gotPath)
	}
	if !strings.Contains(out, "Secrets") {
		t.Errorf("output = %q, want the kinds", out)
	}
}

func TestResourcesCommandExitsOnAnError(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"description":"no"}`))
	})
	if !exits(t, func() { resourcesCmd.Run(resourcesCmd, nil) }) {
		t.Error("resources did not exit on a 403")
	}
}
