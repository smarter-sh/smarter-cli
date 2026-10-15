package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadManifestSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := readManifest(path)
	if err != nil {
		t.Fatalf("readManifest() error = %v", err)
	}
	if got != "kind: Plugin\n" {
		t.Errorf("readManifest() = %q", got)
	}
}

func TestReadManifestMissingFile(t *testing.T) {
	if _, err := readManifest(filepath.Join(t.TempDir(), "does-not-exist.yaml")); err == nil {
		t.Error("readManifest() error = nil, want an error for a missing file")
	}
}

// TestApplyDryRunValidatesInsteadOfApplying checks that --dry-run never
// reaches apply/: the platform's apply has no dry run, so the cli validates.
func TestApplyDryRunValidatesInsteadOfApplying(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\nname: foo\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var gotPaths []string
	var gotQuery, gotBody string
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		gotQuery = r.URL.RawQuery
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(`{"data":{"valid":true,"errors":[]}}`))
	})

	withViperValue(t, "filename", path)
	withViperValue(t, "dry_run", true)
	withViperValue(t, "output_format", nil)

	out := captureStdout(t, func() {
		applyCmd.Run(applyCmd, []string{})
	})

	if len(gotPaths) != 1 || gotPaths[0] != "/api/v1/cli/validate/" {
		t.Errorf("paths = %v, want only /api/v1/cli/validate/", gotPaths)
	}
	if gotQuery != "" {
		t.Errorf("query = %q, want none", gotQuery)
	}
	if gotBody != "kind: Plugin\nname: foo\n" {
		t.Errorf("request body = %q, want the manifest contents", gotBody)
	}
	if !strings.Contains(out, "manifest is valid.") || !strings.Contains(out, "nothing was applied") {
		t.Errorf("output = %q, want the validation and the dry run's note", out)
	}
}

func TestApplyDryRunExitsOnAnInvalidManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\n"), 0644); err != nil {
		t.Fatal(err)
	}
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"valid":false,"errors":[{"loc":["metadata"],"message":"missing","type":"x"}]}}`))
	})
	withViperValue(t, "filename", path)
	withViperValue(t, "dry_run", true)
	withViperValue(t, "output_format", nil)

	if !exits(t, func() {
		_ = captureStdout(t, func() { applyCmd.Run(applyCmd, []string{}) })
	}) {
		t.Error("apply --dry-run did not exit for an invalid manifest")
	}
}

func TestApplyExitsOnAnAPIError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\n"), 0644); err != nil {
		t.Fatal(err)
	}
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"description":"bad manifest"}}`))
	})
	withViperValue(t, "filename", path)
	withViperValue(t, "dry_run", false)

	if !exits(t, func() { applyCmd.Run(applyCmd, []string{}) }) {
		t.Error("apply did not exit on a 400")
	}
}

func TestApplyCommandNonDryRunMessage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	withViperValue(t, "filename", path)
	withViperValue(t, "dry_run", false)

	out := captureStdout(t, func() {
		applyCmd.Run(applyCmd, []string{})
	})

	if !strings.Contains(out, "manifest applied.") || strings.Contains(out, "dry run") {
		t.Errorf("output = %q, want a plain applied message", out)
	}
}

func TestApplyRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "apply" {
			found = true
		}
	}
	if !found {
		t.Error("apply command is not registered on RootCmd")
	}
}
