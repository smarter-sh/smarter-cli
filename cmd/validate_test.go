package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeManifest(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "manifest.yaml")
	if err := os.WriteFile(path, []byte("apiVersion: smarter.sh/v1\nkind: Guardrail\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runValidate(t *testing.T, filename string) (exited bool, stdout string) {
	t.Helper()
	if err := validateCmd.Flags().Set("filename", filename); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = validateCmd.Flags().Set("filename", "") })
	exited = exits(t, func() {
		stdout = captureStdout(t, func() { validateCmd.Run(validateCmd, nil) })
	})
	return exited, stdout
}

func TestValidateAValidManifest(t *testing.T) {
	var gotPath, gotBody string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(`{"data":{"valid":true,"errors":[]},"message":"Manifest is valid"}`))
	})
	withViperValue(t, "output_format", nil)

	exited, out := runValidate(t, writeManifest(t))
	if exited {
		t.Fatal("validate exited for a valid manifest")
	}
	if gotPath != "/api/v1/cli/validate/" {
		t.Errorf("path = %q, want /api/v1/cli/validate/", gotPath)
	}
	if !strings.Contains(gotBody, "kind: Guardrail") {
		t.Errorf("body = %q, want the manifest", gotBody)
	}
	if !strings.Contains(out, "manifest is valid.") {
		t.Errorf("output = %q", out)
	}
}

func TestValidateAnInvalidManifestExits(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"valid":false,"errors":[{"loc":["spec","config",0],"message":"bad","type":"x"}]}}`))
	})
	withViperValue(t, "output_format", nil)

	if exited, _ := runValidate(t, writeManifest(t)); !exited {
		t.Error("validate did not exit for an invalid manifest")
	}
}

func TestReportValidation(t *testing.T) {
	err := reportValidation([]byte(`{"data":{"valid":false,"errors":[{"loc":["spec","config",0],"message":"bad","type":"x"},{"loc":[],"message":"worse","type":"y"}]}}`))
	if err == nil || !strings.Contains(err.Error(), "2 error(s)") {
		t.Errorf("reportValidation() = %v, want the number of errors", err)
	}
	if err := reportValidation([]byte(`not json`)); err == nil {
		t.Error("reportValidation(not json) = nil, want an error")
	}
}

func TestValidatePrintsTheResponseWithAnOutputFormat(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"valid":true,"errors":[]},"message":"Manifest is valid"}`))
	})
	withViperValue(t, "output_format", "json")

	if _, out := runValidate(t, writeManifest(t)); !strings.Contains(out, `"valid": true`) {
		t.Errorf("output = %q, want the response", out)
	}
}

func TestValidateAMissingFileExits(t *testing.T) {
	if exited, _ := runValidate(t, filepath.Join(t.TempDir(), "missing.yaml")); !exited {
		t.Error("validate did not exit for a missing file")
	}
}

func TestValidateRegisteredOnRootCmd(t *testing.T) {
	for _, c := range RootCmd.Commands() {
		if c.Name() == "validate" {
			return
		}
	}
	t.Error("validate command is not registered on RootCmd")
}
