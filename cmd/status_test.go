package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestStatusCommandCallsStatusEndpoint(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"status":"ok"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	out := captureStdout(t, func() {
		statusCmd.Run(statusCmd, []string{})
	})

	if gotPath != "/api/v1/cli/status/" {
		t.Errorf("path = %q, want /api/v1/cli/status/", gotPath)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("output = %q, want it to contain the response body", out)
	}
}

func TestStatusRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "status" {
			found = true
		}
	}
	if !found {
		t.Error("status command is not registered on RootCmd")
	}
}
