package get

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/viper"
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

func TestAPIRequestMergesPagingAndSortFlags(t *testing.T) {
	withViperValue(t, "i", 25)
	withViperValue(t, "asc", true)
	withViperValue(t, "desc", false)

	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("Plugin", map[string]string{"name": "foo"}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/get/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/get/plugin/", gotPath)
	}
	for _, want := range []string{"name=foo", "i=25", "asc=true", "desc=false"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query = %q, want it to contain %q", gotQuery, want)
		}
	}
}

func TestConsoleOutputPreservesAnAlreadySetFormat(t *testing.T) {
	withViperValue(t, "output_format", "json")
	ConsoleOutput([]byte(`{"data":{}}`))
	if got := viper.GetString("output_format"); got != "json" {
		t.Errorf("output_format = %q, want it left as json", got)
	}
}

// TestGetRegistersLegacyAndGeneratedKinds checks getCmd's subcommands are
// the union of the hand-written legacySpecs and one generated GetSpec() per
// entry in the shared cmd.ResourceKinds table.
func TestGetRegistersLegacyAndGeneratedKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range getCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}

	legacy := []string{"account", "chatbots", "chat-history", "chat-plugin-usage", "chats", "chat-tool-calls", "plugins", "users"}
	for _, name := range legacy {
		if !got[name] {
			t.Errorf("expected legacy get command %q to be registered", name)
		}
	}
	for _, k := range cmd.ResourceKinds {
		if !got[k.Plural] {
			t.Errorf("expected generated get command %q to be registered", k.Plural)
		}
	}

	wantCount := len(legacy) + len(cmd.ResourceKinds)
	if len(getCmd.Commands()) != wantCount {
		t.Errorf("get has %d subcommands, want %d", len(getCmd.Commands()), wantCount)
	}
}

func TestGetPluginsEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"titles":[],"items":[]}}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"get", "plugins", "--class", "sql"})
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/get/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/get/plugin/", gotPath)
	}
	if !strings.Contains(gotQuery, "class=sql") {
		t.Errorf("query = %q, want it to contain class=sql", gotQuery)
	}
}

func TestGetRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "get" {
			found = true
		}
	}
	if !found {
		t.Error("get command is not registered on RootCmd")
	}
}
