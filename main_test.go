package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSmarterBinaryVersionCommand is an end-to-end smoke test: it builds the
// real smarter binary and runs it as a subprocess (with a throwaway $HOME so
// it doesn't touch the developer's real ~/.smarter config), verifying the
// wiring in main.go - flag registration, cobra command dispatch, and the
// version command's non-verbose path - works together in the actual
// executable, not just in package-level unit tests.
func TestSmarterBinaryVersionCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping binary build/run smoke test in -short mode")
	}

	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "smarter")

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}

	home := t.TempDir()
	configDir := filepath.Join(home, ".smarter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	// initConfig() (cmd/root.go) log.Fatalf's if no api_key is configured
	// for the resolved environment (prod, by default), so the fixture needs
	// one even though "version" without --verbose never calls the API.
	const fixture = "config:\n  environment: prod\nprod:\n  api_key: \"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcd\"\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(fixture), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	run := exec.Command(binPath, "version")
	run.Env = append(os.Environ(), "HOME="+home)

	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("smarter version failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Local version:") {
		t.Errorf("output = %q, want it to contain 'Local version:'", out)
	}
}
