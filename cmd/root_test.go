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

// TestValidateOutputTogglesExitsOnInvalidFormat verifies the log.Fatalf path
// for an unrecognized output format. Runs in a subprocess since log.Fatalf
// calls os.Exit.
func TestValidateOutputTogglesExitsOnInvalidFormat(t *testing.T) {
	if os.Getenv("SMARTER_CLI_TEST_CRASHER") == "1" {
		viper.Set("output_format", "xml")
		_ = validateOutputToggles()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestValidateOutputTogglesExitsOnInvalidFormat")
	cmd.Env = append(os.Environ(), "SMARTER_CLI_TEST_CRASHER=1")
	output, err := cmd.CombinedOutput()

	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.Success() {
		t.Fatalf("process exited with err=%v, want non-zero status; output: %s", err, output)
	}
	if !strings.Contains(string(output), "Invalid output format") {
		t.Errorf("expected output to explain the invalid format, got: %s", output)
	}
}

// TestInitConfigBootstrapsFreshHome exercises initConfig against a throwaway
// $HOME that has no ~/.smarter directory yet, verifying it creates the
// directory and a default config.yaml, and resolves environment/api_key from
// it. This is a distinct scenario from TestMain's pre-seeded home (which
// exists precisely so *other* tests' incidental Execute() calls don't hit
// this bootstrap path).
// TestInitConfigBootstrapsFreshHome exercises initConfig against a throwaway
// $HOME that has no ~/.smarter directory yet, verifying it creates the
// directory and writes a default config.yaml. This runs in its own
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
		// Route around the log.Fatalf("No api_key found...") branch:
		// initConfig only reaches that check if viper's api_key is still
		// empty after reading the (freshly created, all-empty) config file.
		viper.Set("api_key", "preset-key")

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
