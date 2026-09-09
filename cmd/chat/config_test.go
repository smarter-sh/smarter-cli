package chat

import (
	"net/http"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
)

func TestChatConfigCommandSuccess(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"apiUrl":"https://example.com"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	// cobra's Execute() always redirects to the true root command
	// (Command.Root().ExecuteC()) whenever the receiver has a parent, so
	// args must be set on cmd.RootCmd - configCmd.SetArgs would be silently
	// ignored in favor of whatever RootCmd's own args are.
	cmd.RootCmd.SetArgs([]string{"chat", "config", "--chatbot", "my-bot"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.HasPrefix(gotPath, "/api/v1/cli/chat/config/my-bot/") {
		t.Errorf("path = %q, want it to start with /api/v1/cli/chat/config/my-bot/", gotPath)
	}
}

func TestChatConfigRegisteredOnChatCmd(t *testing.T) {
	found := false
	for _, c := range chatCmd.Commands() {
		if c.Name() == "config" {
			found = true
		}
	}
	if !found {
		t.Error("config command is not registered on chatCmd")
	}
}
