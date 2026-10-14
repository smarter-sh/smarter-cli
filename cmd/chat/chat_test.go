package chat

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
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

func TestAPIRequestRoutesUnderPromptPrefix(t *testing.T) {
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
	if gotPath != "/api/v1/cli/prompt/widget/" {
		t.Errorf("path = %q, want /api/v1/cli/prompt/widget/", gotPath)
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
	if chatCmd.PersistentFlags().Lookup("llmclient") == nil {
		t.Error("chatCmd is missing the --llmclient persistent flag")
	}
	if f := chatCmd.PersistentFlags().Lookup("chatbot"); f == nil || f.Deprecated == "" {
		t.Error("chatCmd should keep --chatbot as a deprecated alias")
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

// exits replaces cmd.ErrorOutput's exit with a panic for the duration of
// the test, runs fn, and reports whether it exited.
func exits(t *testing.T, fn func()) (exited bool) {
	t.Helper()
	type exit struct{}
	cmd.SetExitFunc(func(int) { panic(exit{}) })
	t.Cleanup(func() { cmd.SetExitFunc(nil) })
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(exit); !ok {
				panic(r)
			}
			exited = true
		}
	}()
	fn()
	return false
}

// setFlags sets c's flags for the test, and restores their previous values.
func setFlags(t *testing.T, c *cobra.Command, values map[string]string) {
	t.Helper()
	// merge the parent's persistent flags into c.Flags(), as Execute does
	_ = c.InheritedFlags()
	for name, value := range values {
		f := c.Flags().Lookup(name)
		old := f.Value.String()
		if err := f.Value.Set(value); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = f.Value.Set(old) })
	}
}

func TestLLMClientNamePrefersTheNewFlag(t *testing.T) {
	setFlags(t, promptCmd, map[string]string{"llmclient": "new", "chatbot": "old"})
	if name, err := llmClientName(promptCmd); err != nil || name != "new" {
		t.Errorf("llmClientName() = %q, %v; want new", name, err)
	}
	setFlags(t, promptCmd, map[string]string{"llmclient": ""})
	if name, err := llmClientName(promptCmd); err != nil || name != "old" {
		t.Errorf("llmClientName() = %q, %v; want the deprecated alias's value", name, err)
	}
}

func TestChatCommandsExitWithoutAnLLMClient(t *testing.T) {
	setFlags(t, promptCmd, map[string]string{"llmclient": "", "chatbot": "", "prompt": "hi"})
	if !exits(t, func() { promptCmd.Run(promptCmd, nil) }) {
		t.Error("chat prompt without an LLMClient did not exit")
	}
	setFlags(t, configCmd, map[string]string{"llmclient": "", "chatbot": ""})
	if !exits(t, func() { configCmd.Run(configCmd, nil) }) {
		t.Error("chat config without an LLMClient did not exit")
	}
}

func TestChatPromptExitsWithoutAPrompt(t *testing.T) {
	setFlags(t, promptCmd, map[string]string{"llmclient": "x", "prompt": ""})
	if !exits(t, func() { promptCmd.Run(promptCmd, nil) }) {
		t.Error("chat prompt without a prompt did not exit")
	}
}

func TestChatCommandsExitOnAnAPIError(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"description":"not found"}`))
	})
	setFlags(t, promptCmd, map[string]string{"llmclient": "x", "prompt": "hi"})
	if !exits(t, func() { promptCmd.Run(promptCmd, nil) }) {
		t.Error("chat prompt did not exit on a 404")
	}
	setFlags(t, configCmd, map[string]string{"llmclient": "x"})
	if !exits(t, func() { configCmd.Run(configCmd, nil) }) {
		t.Error("chat config did not exit on a 404")
	}
}
