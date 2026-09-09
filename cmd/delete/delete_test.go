package delete

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// whatever was written to it.
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

func TestAPIRequestRoutesUnderDeletePrefix(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("chatbot", map[string]string{"name": "foo"}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/delete/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/delete/chatbot/", gotPath)
	}
}

func TestConsoleOutputPrintsDeleted(t *testing.T) {
	out := captureStdout(t, ConsoleOutput)
	if !strings.Contains(out, "deleted.") {
		t.Errorf("output = %q, want it to mention deleted", out)
	}
}

func TestDeleteCmdRegistersExpectedLeafCommands(t *testing.T) {
	want := map[string]string{
		"apiconnection": "PluginDataApiConnection",
		"apikey":        "SmarterAuthToken",
		"chat":          "chat",
		"chatbot":       "chatbot",
		"plugin":        "plugin",
		"sqlconnection": "PluginDataSqlConnection",
		"user":          "user",
	}

	got := map[string]bool{}
	for _, c := range deleteCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}
	for name := range want {
		if !got[name] {
			t.Errorf("expected delete subcommand %q to be registered, got commands: %v", name, got)
		}
	}
	if len(deleteCmd.Commands()) != len(want) {
		t.Errorf("delete has %d subcommands, want %d", len(deleteCmd.Commands()), len(want))
	}
}

func TestDeleteChatbotEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	cmd.RootCmd.SetArgs([]string{"delete", "chatbot", "my-bot"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/delete/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/delete/chatbot/", gotPath)
	}
	if !strings.Contains(gotQuery, "name=my-bot") {
		t.Errorf("query = %q, want it to contain name=my-bot", gotQuery)
	}
	if !strings.Contains(out, "deleted.") {
		t.Errorf("output = %q, want it to mention deleted", out)
	}
}

func TestDeleteRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "delete" {
			found = true
		}
	}
	if !found {
		t.Error("delete command is not registered on RootCmd")
	}
}
