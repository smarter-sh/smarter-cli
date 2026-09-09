package describe

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

func TestAPIRequestRoutesUnderDescribePrefix(t *testing.T) {
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
	if gotPath != "/api/v1/cli/describe/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/describe/plugin/", gotPath)
	}
}

func TestConsoleOutputRendersBody(t *testing.T) {
	withViperValue(t, "output_format", "json")
	out := captureStdout(t, func() {
		ConsoleOutput([]byte(`{"name":"foo"}`))
	})
	if !strings.Contains(out, "foo") {
		t.Errorf("output = %q, want it to contain the body", out)
	}
}

// TestDescribeRegistersLegacyAndGeneratedKinds checks that describeCmd's
// subcommands are the union of the hand-written legacySpecs and one
// generated DescribeSpec() per entry in the shared cmd.ResourceKinds table.
func TestDescribeRegistersLegacyAndGeneratedKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range describeCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}

	for _, legacyName := range []string{"account", "chatbot", "plugin", "user"} {
		if !got[legacyName] {
			t.Errorf("expected legacy describe command %q to be registered", legacyName)
		}
	}
	for _, k := range cmd.ResourceKinds {
		if !got[k.Singular] {
			t.Errorf("expected generated describe command %q to be registered", k.Singular)
		}
	}

	wantCount := 4 + len(cmd.ResourceKinds)
	if len(describeCmd.Commands()) != wantCount {
		t.Errorf("describe has %d subcommands, want %d", len(describeCmd.Commands()), wantCount)
	}
}

func TestDescribePluginEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"kind":"Plugin"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"describe", "plugin", "my-plugin"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/describe/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/describe/plugin/", gotPath)
	}
	if !strings.Contains(gotQuery, "name=my-plugin") {
		t.Errorf("query = %q, want it to contain name=my-plugin", gotQuery)
	}
	if !strings.Contains(out, "Plugin") {
		t.Errorf("output = %q, want it to contain the response body", out)
	}
}

func TestDescribeRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "describe" {
			found = true
		}
	}
	if !found {
		t.Error("describe command is not registered on RootCmd")
	}
}
