package deploy

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

func TestAPIRequestRoutesUnderDeployPrefix(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("chatbot", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/deploy/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/deploy/chatbot/", gotPath)
	}
}

func TestConsoleOutputPrintsDeployed(t *testing.T) {
	out := captureStdout(t, ConsoleOutput)
	if !strings.Contains(out, "deployed.") {
		t.Errorf("output = %q, want it to mention deployed", out)
	}
}

// TestDeployOnlyRegistersDeployableKinds guards the v0.14 rule described in
// deploy/resources.go: of the generated ResourceKind specs, only kinds
// flagged Deployable (currently just "prompt") get a deploy leaf command,
// on top of the legacy "chatbot" leaf.
func TestDeployOnlyRegistersDeployableKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range deployCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}

	if !got["chatbot"] {
		t.Error("expected legacy 'chatbot' deploy command to be registered")
	}

	for _, k := range cmd.ResourceKinds {
		if k.Deployable && !got[k.Singular] {
			t.Errorf("expected deployable kind %q to have a deploy command", k.Singular)
		}
		if !k.Deployable && got[k.Singular] {
			t.Errorf("non-deployable kind %q should not have a deploy command", k.Singular)
		}
	}
}

func TestDeployChatbotEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	cmd.RootCmd.SetArgs([]string{"deploy", "chatbot", "my-bot"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/deploy/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/deploy/chatbot/", gotPath)
	}
	if !strings.Contains(gotQuery, "name=my-bot") {
		t.Errorf("query = %q, want it to contain name=my-bot", gotQuery)
	}
	if !strings.Contains(out, "deployed.") {
		t.Errorf("output = %q, want it to mention deployed", out)
	}
}

func TestDeployRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "deploy" {
			found = true
		}
	}
	if !found {
		t.Error("deploy command is not registered on RootCmd")
	}
}
