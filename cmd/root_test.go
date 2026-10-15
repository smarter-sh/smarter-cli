package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateEnvironmentFlag(t *testing.T) {
	oldEnv := environment
	t.Cleanup(func() { environment = oldEnv })

	environment = ""
	if err := validateEnvironmentFlag(); err != nil {
		t.Errorf("validateEnvironmentFlag() with empty environment = %v, want nil", err)
	}

	for _, valid := range validEnvironments {
		environment = valid
		if err := validateEnvironmentFlag(); err != nil {
			t.Errorf("validateEnvironmentFlag() for %q = %v, want nil", valid, err)
		}
	}

	environment = "not-a-real-environment"
	if err := validateEnvironmentFlag(); err == nil {
		t.Error("validateEnvironmentFlag() with an invalid environment = nil, want an error")
	}
}

func TestValidateOutputTogglesAcceptsValidFormats(t *testing.T) {
	for _, format := range []string{"json", "yaml", "tabular"} {
		withViperValue(t, "output_format", format)
		if err := validateOutputToggles(); err != nil {
			t.Errorf("validateOutputToggles() for %q = %v, want nil", format, err)
		}
	}
}

func TestValidateOutputTogglesRejectsInvalidFormat(t *testing.T) {
	withViperValue(t, "output_format", "xml")
	err := validateOutputToggles()
	if err == nil || !strings.Contains(err.Error(), "invalid output format") {
		t.Errorf("validateOutputToggles() = %v, want an error explaining the invalid format", err)
	}
}

// TestInitConfigBootstrapsFreshHome exercises initConfig against a throwaway
// $HOME that has no ~/.smarter directory yet, verifying it creates the
// directory and writes a default config.yaml, owner-only, holding only the
// configuration keys, without exiting for the missing api_key. This runs in its own
// subprocess rather than in-process: viper is a global singleton that caches
// the last successfully-read config file's values for the life of the
// process, and a failed re-read (which is what happens here, since the file
// doesn't exist yet) does not clear that cache. Any other test in this
// binary that has already run initConfig() against TestMain's fixture home
// (which cobra.OnInitialize triggers on *any* command's Execute(), not just
// this one) would otherwise leave stale values that shadow this test's own
// config. A subprocess guarantees this is the only initConfig call that has
// ever happened in that process.
func TestInitConfigBootstrapsFreshHome(t *testing.T) {
	if os.Getenv("SMARTER_CLI_TEST_FRESH_INIT") == "1" {
		freshHome := os.Getenv("SMARTER_CLI_TEST_FRESH_HOME_DIR")
		if err := os.Setenv("HOME", freshHome); err != nil {
			t.Fatalf("os.Setenv() error = %v", err)
		}
		cfgFile = filepath.Join(freshHome, ".smarter", "config.yaml")
		environment = ""

		// No api_key anywhere: initConfig must not exit, so that
		// "smarter configure" can run on a fresh install.
		initConfig()

		written, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Println("MISSING_FILE:", err)
			os.Exit(1)
		}
		content := string(written)
		if !strings.Contains(content, "root_domain: platform.smarter.sh") {
			fmt.Println("BAD_ROOT_DOMAIN:", content)
			os.Exit(1)
		}
		if !strings.Contains(content, "environment: prod") {
			fmt.Println("BAD_ENVIRONMENT:", content)
			os.Exit(1)
		}
		// Only configuration keys, not every flag that viper knows.
		for _, junk := range []string{"chatbot:", "asc:", "today:", "filename:"} {
			if strings.Contains(content, junk) {
				fmt.Println("FLAG_WRITTEN:", junk, content)
				os.Exit(1)
			}
		}
		info, err := os.Stat(cfgFile)
		if err != nil || info.Mode().Perm() != 0o600 {
			fmt.Println("BAD_MODE:", info.Mode().Perm(), err)
			os.Exit(1)
		}
		fmt.Println("OK")
		return
	}

	freshHome := t.TempDir()
	c := exec.Command(os.Args[0], "-test.run=^TestInitConfigBootstrapsFreshHome$")
	c.Env = append(os.Environ(),
		"SMARTER_CLI_TEST_FRESH_INIT=1",
		"SMARTER_CLI_TEST_FRESH_HOME_DIR="+freshHome,
	)
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("subprocess failed: %v\noutput:\n%s", err, out)
	}
	if !strings.Contains(string(out), "OK") {
		t.Fatalf("subprocess did not report success, output:\n%s", out)
	}
}

// useConfigFile points initConfig at path for the test, and afterwards
// re-reads TestMain's fixture config, because viper keeps the last config it
// read for the life of the process.
func useConfigFile(t *testing.T, path string) {
	t.Helper()
	old := cfgFile
	cfgFile = path
	t.Cleanup(func() {
		cfgFile = old
		initConfig()
	})
}

func TestInitConfigWritesAMissingConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "smarter", "config.yaml")
	useConfigFile(t, path)
	withViperValue(t, "environment", "")
	withViperValue(t, "verbose", true)

	initConfig()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("initConfig() didn't write %s: %v", path, err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", info.Mode().Perm())
	}
	if got := viper.GetString("environment"); got == "" {
		t.Error("initConfig() didn't resolve an environment")
	}
}

func TestInitConfigReportsAnUnreadableConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("config: [not: valid"), 0o600); err != nil {
		t.Fatal(err)
	}
	useConfigFile(t, path)

	initConfig()

	body, _ := os.ReadFile(path)
	if string(body) != "config: [not: valid" {
		t.Errorf("initConfig() overwrote an unreadable config file: %q", body)
	}
}

func TestExecuteRunsACommand(t *testing.T) {
	RootCmd.SetArgs([]string{"version"})
	out := captureStdout(t, func() { Execute("9.9.9") })
	if !strings.Contains(out, "9.9.9") {
		t.Errorf("output = %q, want the version", out)
	}
}

func TestExecuteExitsOnAnUnknownCommand(t *testing.T) {
	RootCmd.SetArgs([]string{"no-such-command"})
	t.Cleanup(func() { RootCmd.SetArgs(nil) })
	if !exits(t, func() {
		_ = captureStdout(t, func() { Execute("9.9.9") })
	}) {
		t.Error("Execute() did not exit on an unknown command")
	}
}
