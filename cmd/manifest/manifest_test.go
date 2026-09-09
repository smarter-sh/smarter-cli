package manifest

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestAPIRequestRoutesUnderExampleManifestPrefix(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("plugin", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/example_manifest/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/example_manifest/plugin/", gotPath)
	}
}

func TestConsoleOutputUnwrapsDataEnvelope(t *testing.T) {
	withViperValue(t, "output_format", "json")
	out := captureStdout(t, func() {
		ConsoleOutput([]byte(`{"data":{"kind":"Plugin"},"message":"ok"}`))
	})
	if !strings.Contains(out, "Plugin") {
		t.Errorf("output = %q, want it to contain the unwrapped data", out)
	}
	if strings.Contains(out, "\"message\"") {
		t.Errorf("output = %q, want the envelope's message field stripped", out)
	}
}

func TestConsoleOutputWithoutDataEnvelope(t *testing.T) {
	withViperValue(t, "output_format", "json")
	out := captureStdout(t, func() {
		ConsoleOutput([]byte(`{"kind":"Plugin"}`))
	})
	if !strings.Contains(out, "Plugin") {
		t.Errorf("output = %q, want it to contain the body", out)
	}
}

// TestManifestRegistersLegacyAndGeneratedKinds checks manifestCmd's
// subcommands are the union of the hand-written legacySpecs and one
// generated ManifestSpec() per entry in the shared cmd.ResourceKinds table.
func TestManifestRegistersLegacyAndGeneratedKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range manifestCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}

	for _, legacyName := range []string{"account", "chat", "chatbot", "plugin", "user"} {
		if !got[legacyName] {
			t.Errorf("expected legacy manifest command %q to be registered", legacyName)
		}
	}
	for _, k := range cmd.ResourceKinds {
		if !got[k.Singular] {
			t.Errorf("expected generated manifest command %q to be registered", k.Singular)
		}
	}

	wantCount := 5 + len(cmd.ResourceKinds)
	if len(manifestCmd.Commands()) != wantCount {
		t.Errorf("manifest has %d subcommands, want %d", len(manifestCmd.Commands()), wantCount)
	}
}

func TestManifestPluginEndToEnd(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"kind":"Plugin","name":"example"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"manifest", "plugin"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/example_manifest/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/example_manifest/plugin/", gotPath)
	}
	if !strings.Contains(out, "example") {
		t.Errorf("output = %q, want it to contain the response body", out)
	}
}

func TestManifestRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "manifest" {
			found = true
		}
	}
	if !found {
		t.Error("manifest command is not registered on RootCmd")
	}
}
