package get

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/pflag"
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

// TestDeprecatedFlagsAreNotSent checks that --i, --asc, --desc and --class,
// which the Smarter API ignores, are still accepted but never sent.
func TestDeprecatedFlagsAreNotSent(t *testing.T) {
	var gotQuery string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":{"data":{"titles":[],"items":[]}}}`))
	})
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"get", "plugins", "--i", "25", "--asc", "--class", "sql"})
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
	for _, unwanted := range []string{"i=", "asc=", "desc=", "class="} {
		if strings.Contains(gotQuery, unwanted) {
			t.Errorf("query = %q, want no %s", gotQuery, unwanted)
		}
	}
	for _, name := range []string{"i", "asc", "desc"} {
		if getCmd.PersistentFlags().Lookup(name).Deprecated == "" {
			t.Errorf("--%s isn't marked deprecated", name)
		}
	}
}

func TestGetUsersSendsName(t *testing.T) {
	var gotQuery string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":{"data":{"titles":[],"items":[]}}}`))
	})
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"get", "users", "-u", "jane"})
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
	if gotQuery != "name=jane" {
		t.Errorf("query = %q, want name=jane: the platform filters users by name", gotQuery)
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

	legacy := []string{"account", "plugins", "users"}
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

	cmd.RootCmd.SetArgs([]string{"get", "plugins"})
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/get/plugin/" {
		t.Errorf("path = %q, want /api/v1/cli/get/plugin/", gotPath)
	}
	if strings.Contains(gotQuery, "class=") {
		t.Errorf("query = %q, want no class", gotQuery)
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

func TestConsoleOutputExitsOnAnUnparsableBody(t *testing.T) {
	withViperValue(t, "output_format", "tabular")
	type exit struct{}
	cmd.SetExitFunc(func(int) { panic(exit{}) })
	t.Cleanup(func() { cmd.SetExitFunc(nil) })

	defer func() {
		if _, ok := recover().(exit); !ok {
			t.Error("ConsoleOutput did not exit on an unparsable body")
		}
	}()
	ConsoleOutput([]byte(`not json`))
}

func TestConsoleOutputDefaultsToTabular(t *testing.T) {
	// viper.IsSet reports a nil value as unset
	withViperValue(t, "output_format", nil)

	captureStdout(t, func() { ConsoleOutput([]byte(`{"data":{}}`)) })

	if got := viper.GetString("output_format"); got != "tabular" {
		t.Errorf("output_format = %q, want tabular when it wasn't set", got)
	}
}

// resetAnnotationFlag clears --annotation after a test that set it, because
// cobra keeps a flag's values between Executes.
func resetAnnotationFlag(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		flag := getCmd.PersistentFlags().Lookup(annotationFlag)
		if err := flag.Value.(pflag.SliceValue).Replace(nil); err != nil {
			t.Errorf("resetting --annotation: %v", err)
		}
		flag.Changed = false
	})
}

// TestGetSendsEachAnnotationRequirement checks that every --annotation is
// sent as its own annotation query parameter, unchanged and in order: the
// platform reads them with getlist, and a value may contain a comma.
func TestGetSendsEachAnnotationRequirement(t *testing.T) {
	var gotPath string
	var got url.Values
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, got = r.URL.Path, r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{"data":{"titles":[],"items":[]}}}`))
	})
	withViperValue(t, "output_format", "json")
	resetAnnotationFlag(t)

	want := []string{
		"acme.com/crm-account-id=A-10042",
		"!acme.com/archived",
		"acme.com/tier!=free",
		"acme.com/owner",
		`acme.com/note=one, two "three"`,
	}
	args := []string{"get", "llmclients", "--name", "stackademy"}
	for _, requirement := range want {
		args = append(args, "--annotation", requirement)
	}
	cmd.RootCmd.SetArgs(args)
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	if gotPath != "/api/v1/cli/get/llmclient/" {
		t.Errorf("path = %q, want /api/v1/cli/get/llmclient/", gotPath)
	}
	if !reflect.DeepEqual(got["annotation"], want) {
		t.Errorf("annotation = %q, want %q", got["annotation"], want)
	}
	if got.Get("name") != "stackademy" {
		t.Errorf("name = %q, want stackademy", got.Get("name"))
	}
}

func TestGetWithoutAnnotationSendsNone(t *testing.T) {
	var gotQuery string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":{"data":{"titles":[],"items":[]}}}`))
	})
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"get", "secrets"})
	_ = captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
	if strings.Contains(gotQuery, "annotation") {
		t.Errorf("query = %q, want no annotation", gotQuery)
	}
}

// TestGetReportsAnInvalidSelector checks that the platform's 400 for an
// invalid requirement reaches the user as its description.
func TestGetReportsAnInvalidSelector(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"errorClass":"SmarterAnnotationError","description":"Invalid annotation key 'Bad Key'."}}`))
	})
	withViperValue(t, annotationFlag, []string{"Bad Key=1"})

	_, err := APIRequest("LLMClient", map[string]string{})
	if err == nil {
		t.Fatal("APIRequest() error = nil, want the platform's 400")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "Invalid annotation key 'Bad Key'.") {
		t.Errorf("error = %q, want the 400 and its description", err)
	}
}

func TestAnnotationFlagIsRepeatable(t *testing.T) {
	flag := getCmd.PersistentFlags().Lookup(annotationFlag)
	if flag == nil {
		t.Fatal("get has no --annotation flag")
	}
	if flag.Value.Type() != "stringArray" {
		t.Errorf("--annotation is a %s, want a stringArray, which never splits a value at its commas", flag.Value.Type())
	}
}
