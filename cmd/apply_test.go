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

func TestApplyCommandSendsManifestAndDryRunFlag(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte("kind: Plugin\nname: foo\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var gotPath, gotQuery, gotBody string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	withViperValue(t, "filename", path)
	withViperValue(t, "dry_run", true)

	out := captureStdout(t, func() {
		applyCmd.Run(applyCmd, []string{})
	})

	if gotPath != "/api/v1/cli/apply/" {
		t.Errorf("path = %q, want /api/v1/cli/apply/", gotPath)
	}
	if !strings.Contains(gotQuery, "dry_run=true") {
		t.Errorf("query = %q, want it to contain dry_run=true", gotQuery)
	}
	if gotBody != "kind: Plugin\nname: foo\n" {
		t.Errorf("request body = %q, want the manifest contents", gotBody)
	}
	if !strings.Contains(out, "dry run") {
		t.Errorf("output = %q, want it to mention the dry run", out)
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
