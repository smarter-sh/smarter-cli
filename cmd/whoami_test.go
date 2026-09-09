package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestWhoamiCommandCallsWhoamiEndpoint(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"username":"alice"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	out := captureStdout(t, func() {
		whoamiCmd.Run(whoamiCmd, []string{})
	})

	if gotPath != "/api/v1/cli/whoami/" {
		t.Errorf("path = %q, want /api/v1/cli/whoami/", gotPath)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("output = %q, want it to contain the response body", out)
	}
}

func TestWhoamiRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "whoami" {
			found = true
		}
	}
	if !found {
		t.Error("whoami command is not registered on RootCmd")
	}
}
