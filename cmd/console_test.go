package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// whatever was written to it. Several functions under test (TableOutput,
// JsonOutput, YamlOutput) print directly to os.Stdout rather than accepting
// an io.Writer, so this is the only way to observe their output.
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
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	return buf.String()
}

const sampleListBody = `{
	"data": {
		"apiVersion": "v1",
		"kind": "Plugin",
		"name": "plugins",
		"metadata": {"count": 2},
		"data": {
			"titles": [
				{"name": "name", "type": "CharField"},
				{"name": "created", "type": "DateTimeField"}
			],
			"items": [
				{"name": "foo", "created": "2024-01-02T15:04:05Z"},
				{"name": "bar", "created": null}
			]
		}
	},
	"message": "ok"
}`

func TestTableOutput(t *testing.T) {
	out := captureStdout(t, func() {
		if err := TableOutput([]byte(sampleListBody)); err != nil {
			t.Fatalf("TableOutput() error = %v", err)
		}
	})

	if !strings.Contains(out, "name") || !strings.Contains(out, "created") {
		t.Errorf("output missing column titles: %q", out)
	}
	if !strings.Contains(out, "foo") || !strings.Contains(out, "bar") {
		t.Errorf("output missing row data: %q", out)
	}
	if !strings.Contains(out, "2024-Jan-02") {
		t.Errorf("output did not format the DateTimeField, got: %q", out)
	}
}

func TestTableOutputInvalidJSON(t *testing.T) {
	if err := TableOutput([]byte("not json")); err == nil {
		t.Error("TableOutput() error = nil, want an error for invalid JSON")
	}
}

func TestTableOutputInvalidDate(t *testing.T) {
	body := `{"data":{"data":{"titles":[{"name":"created","type":"DateTimeField"}],"items":[{"created":"not-a-date"}]}}}`
	if err := TableOutput([]byte(body)); err == nil {
		t.Error("TableOutput() error = nil, want an error for an unparseable date")
	}
}

func TestJsonOutput(t *testing.T) {
	out := captureStdout(t, func() {
		if err := JsonOutput([]byte(`{"b":2,"a":1}`)); err != nil {
			t.Fatalf("JsonOutput() error = %v", err)
		}
	})

	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, out)
	}
	if got["a"] != 1.0 || got["b"] != 2.0 {
		t.Errorf("got %v, want a=1 b=2", got)
	}
}

func TestJsonOutputInvalidJSON(t *testing.T) {
	if err := JsonOutput([]byte("not json")); err == nil {
		t.Error("JsonOutput() error = nil, want an error for invalid JSON")
	}
}

func TestYamlOutput(t *testing.T) {
	out := captureStdout(t, func() {
		if err := YamlOutput([]byte(`{"data":{"name":"foo"},"message":"ok"}`)); err != nil {
			t.Fatalf("YamlOutput() error = %v", err)
		}
	})

	// YamlOutput unwraps the "data" envelope before converting to YAML.
	if !strings.Contains(out, "name: foo") {
		t.Errorf("output = %q, want it to contain the unwrapped data", out)
	}
	if strings.Contains(out, "message") {
		t.Errorf("output = %q, want the envelope's message field stripped", out)
	}
}

func TestYamlOutputWithoutDataEnvelope(t *testing.T) {
	out := captureStdout(t, func() {
		if err := YamlOutput([]byte(`{"name":"foo"}`)); err != nil {
			t.Fatalf("YamlOutput() error = %v", err)
		}
	})
	if !strings.Contains(out, "name: foo") {
		t.Errorf("output = %q, want it to contain name: foo", out)
	}
}

func TestYamlOutputInvalidJSON(t *testing.T) {
	if err := YamlOutput([]byte("not json")); err == nil {
		t.Error("YamlOutput() error = nil, want an error for invalid JSON")
	}
}

// TestErrorOutputExitsProcess verifies ErrorOutput prints to stderr and exits
// with status 1. Since os.Exit ends the process, this must run in a
// subprocess rather than in this test binary directly.
func TestErrorOutputExitsProcess(t *testing.T) {
	if os.Getenv("SMARTER_CLI_TEST_CRASHER") == "1" {
		ErrorOutput(errors.New("boom"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestErrorOutputExitsProcess")
	cmd.Env = append(os.Environ(), "SMARTER_CLI_TEST_CRASHER=1")
	output, err := cmd.CombinedOutput()

	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.Success() {
		t.Fatalf("process exited with err=%v, want a non-zero exit status; output: %s", err, output)
	}
	if !strings.Contains(string(output), "boom") {
		t.Errorf("expected output to contain the error message, got: %s", output)
	}
}

func TestConsoleOutputDispatchesOnFormat(t *testing.T) {
	body := []byte(`{"data":{"name":"foo"}}`)

	cases := []struct {
		format string
		want   string
	}{
		{"json", `"name"`},
		{"yaml", "name: foo"},
		{"tabular", ""}, // no titles/items in this body; just must not error
		{"unrecognized", `"name"`},
	}

	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			withViperValue(t, "output_format", tc.format)
			out := captureStdout(t, func() {
				if err := ConsoleOutput(body); err != nil {
					t.Fatalf("ConsoleOutput() error = %v", err)
				}
			})
			if tc.want != "" && !strings.Contains(out, tc.want) {
				t.Errorf("output = %q, want it to contain %q", out, tc.want)
			}
		})
	}
}
