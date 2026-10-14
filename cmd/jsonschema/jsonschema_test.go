package jsonschema

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

func TestAPIRequestRoutesUnderJSONSchemaPrefix(t *testing.T) {
	var gotPath string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := APIRequest("Guardrail", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/json-schema/guardrail/" {
		t.Errorf("path = %q, want /api/v1/cli/json-schema/guardrail/", gotPath)
	}
}

// TestConsoleOutputPrintsTheSchemaAsJSON checks that the schema itself, not
// the response envelope, is printed, as json by default.
func TestConsoleOutputPrintsTheSchemaAsJSON(t *testing.T) {
	old := viper.Get("output_format")
	t.Cleanup(func() { viper.Set("output_format", old) })
	viper.Set("output_format", nil)

	out := captureStdout(t, func() {
		ConsoleOutput([]byte(`{"data":{"title":"Guardrail","type":"object"},"message":"ok"}`))
	})
	if !strings.Contains(out, `"title": "Guardrail"`) || strings.Contains(out, `"message"`) {
		t.Errorf("output = %q, want the schema alone, as json", out)
	}
}

func TestJSONSchemaRegistersLegacyAndGeneratedKinds(t *testing.T) {
	got := map[string]bool{}
	for _, c := range jsonSchemaCmd.Commands() {
		got[strings.Fields(c.Use)[0]] = true
	}
	for _, name := range []string{"account", "plugin", "user"} {
		if !got[name] {
			t.Errorf("expected legacy json-schema command %q to be registered", name)
		}
	}
	for _, k := range cmd.ResourceKinds {
		if !got[k.Singular] {
			t.Errorf("expected generated json-schema command %q to be registered", k.Singular)
		}
	}
	if want := 3 + len(cmd.ResourceKinds); len(jsonSchemaCmd.Commands()) != want {
		t.Errorf("json-schema has %d subcommands, want %d", len(jsonSchemaCmd.Commands()), want)
	}
}

func TestJSONSchemaGuardrailEndToEnd(t *testing.T) {
	var gotPath string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"title":"Guardrail"}}`))
	})
	withViperValue(t, "output_format", "json")

	cmd.RootCmd.SetArgs([]string{"json-schema", "guardrail"})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
	if gotPath != "/api/v1/cli/json-schema/guardrail/" {
		t.Errorf("path = %q", gotPath)
	}
	if !strings.Contains(out, "Guardrail") {
		t.Errorf("output = %q, want the schema", out)
	}
}

func TestJSONSchemaRegisteredOnRootCmd(t *testing.T) {
	for _, c := range cmd.RootCmd.Commands() {
		if c.Name() == "json-schema" {
			return
		}
	}
	t.Error("json-schema command is not registered on RootCmd")
}
