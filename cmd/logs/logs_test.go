package logs

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

func TestAPIRequestRoutesUnderLogsPrefix(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("LLMClient", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/logs/llmclient/" {
		t.Errorf("path = %q, want /api/v1/cli/logs/llmclient/", gotPath)
	}
}

func TestConsoleOutputRendersBody(t *testing.T) {
	withViperValue(t, "output_format", "json")
	out := captureStdout(t, func() {
		ConsoleOutput([]byte(`{"lines":"log output"}`))
	})
	if !strings.Contains(out, "log output") {
		t.Errorf("output = %q, want it to contain the body", out)
	}
}

func TestLogsRegistersExpectedLeafCommands(t *testing.T) {
	want := []string{"plugin"}
	for _, k := range cmd.ResourceKinds {
		want = append(want, k.Singular)
	}
	got := map[string]bool{}
	for _, c := range logsCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("expected logs command %q to be registered", name)
		}
	}
	if len(logsCmd.Commands()) != len(want) {
		t.Errorf("logs has %d subcommands, want %d", len(logsCmd.Commands()), len(want))
	}
}

func TestLogsLLMClientEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"lines":"hello"}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"logs", "llmclient", "my-bot"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/logs/llmclient/" {
		t.Errorf("path = %q, want /api/v1/cli/logs/llmclient/", gotPath)
	}
	if !strings.Contains(gotQuery, "name=my-bot") {
		t.Errorf("query = %q, want it to contain name=my-bot", gotQuery)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("output = %q, want it to contain the response body", out)
	}
}

func TestLogsRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "logs" {
			found = true
		}
	}
	if !found {
		t.Error("logs command is not registered on RootCmd")
	}
}
