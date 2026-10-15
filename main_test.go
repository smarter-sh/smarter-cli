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

// TestSmarterBinaryConfiguresAFreshInstall builds the binary and runs it on an
// empty $HOME: version and configure must work before any api_key exists,
// and configure must write an owner-only file holding the key and only the
// configuration keys. See doc/adr/002-configuration.md.
func TestSmarterBinaryConfiguresAFreshInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping binary build/run smoke test in -short mode")
	}

	binPath := filepath.Join(t.TempDir(), "smarter")
	if out, err := exec.Command("go", "build", "-o", binPath, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	home := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		c := exec.Command(binPath, args...)
		c.Env = append(os.Environ(), "HOME="+home)
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("smarter %v failed: %v\n%s", args, err, out)
		}
		return string(out)
	}

	if out := run("version"); !strings.Contains(out, "Local version:") {
		t.Errorf("version output = %q", out)
	}
	run("configure", "--api_key", testAPIKey, "--account_number", "1234-5678-9012")

	path := filepath.Join(home, ".smarter", "config.yaml")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("config.yaml mode = %v, want 0600", info.Mode().Perm())
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(body)
	if !strings.Contains(content, testAPIKey) || !strings.Contains(content, "1234-5678-9012") {
		t.Errorf("config.yaml = %q, want the configured key and account number", content)
	}
	for _, flag := range []string{"chatbot:", "asc:", "today:", "filename:"} {
		if strings.Contains(content, flag) {
			t.Errorf("config.yaml contains the flag %q", flag)
		}
	}
}
