package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestVersionCommandNonVerbosePrintsLocalVersionOnly(t *testing.T) {
	oldVersion := Version
	t.Cleanup(func() { Version = oldVersion })
	Version = "9.9.9"

	withViperValue(t, "verbose", false)

	out := captureStdout(t, func() {
		versionCmd.Run(versionCmd, []string{})
	})

	if !strings.Contains(out, "9.9.9") {
		t.Errorf("output = %q, want it to contain the local version", out)
	}
}

func TestVersionCommandVerboseMergesRemoteVersion(t *testing.T) {
	oldVersion := Version
	t.Cleanup(func() { Version = oldVersion })
	Version = "9.9.9"

	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"remote_version":"1.0.0"}`))
	})
	defer srv.Close()

	withViperValue(t, "verbose", true)
	withViperValue(t, "output_format", "json")

	out := captureStdout(t, func() {
		versionCmd.Run(versionCmd, []string{})
	})

	if gotPath != "/api/v1/cli/version/" {
		t.Errorf("path = %q, want /api/v1/cli/version/", gotPath)
	}
	if !strings.Contains(out, "9.9.9") || !strings.Contains(out, "1.0.0") {
		t.Errorf("output = %q, want it to contain both local and remote versions", out)
	}
}

func TestVersionRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "version" {
			found = true
		}
	}
	if !found {
		t.Error("version command is not registered on RootCmd")
	}
}
