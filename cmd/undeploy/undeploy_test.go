package undeploy

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

func TestAPIRequestRoutesUnderUndeployPrefix(t *testing.T) {
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
	if gotPath != "/api/v1/cli/undeploy/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/undeploy/chatbot/", gotPath)
	}
}

func TestConsoleOutputPrintsUndeployed(t *testing.T) {
	out := captureStdout(t, ConsoleOutput)
	if !strings.Contains(out, "undeployed.") {
		t.Errorf("output = %q, want it to mention undeployed", out)
	}
}

func TestUndeployOnlyRegistersDeployableKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range undeployCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}

	if !got["chatbot"] {
		t.Error("expected legacy 'chatbot' undeploy command to be registered")
	}

	for _, k := range cmd.ResourceKinds {
		if k.Deployable && !got[k.Singular] {
			t.Errorf("expected deployable kind %q to have an undeploy command", k.Singular)
		}
		if !k.Deployable && got[k.Singular] {
			t.Errorf("non-deployable kind %q should not have an undeploy command", k.Singular)
		}
	}
}

func TestUndeployChatbotEndToEnd(t *testing.T) {
	var gotPath, gotQuery string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	cmd.RootCmd.SetArgs([]string{"undeploy", "chatbot", "my-bot"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/undeploy/chatbot/" {
		t.Errorf("path = %q, want /api/v1/cli/undeploy/chatbot/", gotPath)
	}
	if !strings.Contains(gotQuery, "name=my-bot") {
		t.Errorf("query = %q, want it to contain name=my-bot", gotQuery)
	}
	if !strings.Contains(out, "undeployed.") {
		t.Errorf("output = %q, want it to mention undeployed", out)
	}
}

func TestUndeployRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "undeploy" {
			found = true
		}
	}
	if !found {
		t.Error("undeploy command is not registered on RootCmd")
	}
}
