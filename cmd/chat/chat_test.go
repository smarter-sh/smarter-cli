package chat

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/viper"
)

func TestGetUniqueID(t *testing.T) {
	got := getUniqueID()
	if got == "" {
		t.Fatal("getUniqueID() returned an empty string")
	}
	if !strings.Contains(got, "-") {
		t.Errorf("getUniqueID() = %q, want host-macAddr shape", got)
	}
}

func TestAPIRequestRoutesUnderChatPrefix(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("widget/", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/chat/widget/" {
		t.Errorf("path = %q, want /api/v1/cli/chat/widget/", gotPath)
	}
}

func TestAPIRequestForwardsFileContents(t *testing.T) {
	var gotBody string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("widget/", map[string]string{}, `{"prompt":"hi"}`); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotBody != `{"prompt":"hi"}` {
		t.Errorf("body = %q", gotBody)
	}
}

func TestConsoleOutputPreservesAnAlreadySetFormat(t *testing.T) {
	withViperValue(t, "output_format", "yaml")

	// ConsoleOutput only forces "json" when output_format isn't set at all;
	// it must leave an explicitly configured format alone.
	ConsoleOutput([]byte(`{"a":1}`))

	if got := viper.GetString("output_format"); got != "yaml" {
		t.Errorf("output_format = %q, want it left as yaml", got)
	}
}

func TestChatCmdRegisteredWithFlags(t *testing.T) {
	if chatCmd.PersistentFlags().Lookup("chatbot") == nil {
		t.Error("chatCmd is missing the --chatbot persistent flag")
	}
	if chatCmd.PersistentFlags().Lookup("new_session") == nil {
		t.Error("chatCmd is missing the --new_session persistent flag")
	}

	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "chat" {
			found = true
		}
	}
	if !found {
		t.Error("chat command is not registered on RootCmd")
	}
}
