package chat

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
)

func TestChatPromptCommandSuccess(t *testing.T) {
	var gotPath, gotBody string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"response":"hello back"}}`))
	})
	defer srv.Close()
	withViperValue(t, "output_format", "json")

	// See config_test.go: Execute() always redirects to the real root, so
	// args must be set there with the full command path.
	cmd.RootCmd.SetArgs([]string{"chat", "prompt", "--chatbot", "my-bot", "--prompt", "hi there"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.HasPrefix(gotPath, "/api/v1/cli/chat/my-bot/") {
		t.Errorf("path = %q, want it to start with /api/v1/cli/chat/my-bot/", gotPath)
	}
	if !strings.Contains(gotBody, `"prompt":"hi there"`) {
		t.Errorf("body = %q, want it to contain the prompt", gotBody)
	}
}

func TestChatPromptRegisteredOnChatCmd(t *testing.T) {
	found := false
	for _, c := range chatCmd.Commands() {
		if c.Name() == "prompt" {
			found = true
		}
	}
	if !found {
		t.Error("prompt command is not registered on chatCmd")
	}
}
