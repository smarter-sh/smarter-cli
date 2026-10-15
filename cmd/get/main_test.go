package get

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
)

// testAPIKey is a syntactically valid (64 hex char) placeholder API key.
const testAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

const testConfigYAML = `config:
  account_number: "1234-5678-9012"
  environment: alpha
  output_format: yaml
  root_domain: platform.smarter.sh
local:
  api_key: "` + testAPIKey + `"
alpha:
  api_key: "` + testAPIKey + `"
beta:
  api_key: "` + testAPIKey + `"
prod:
  api_key: "` + testAPIKey + `"
`

// TestMain seeds a fake $HOME with a complete config.yaml before any test in
// this package runs. cobra.OnInitialize(cmd.initConfig) fires on *every*
// Command.Execute() call in the process (this package imports cmd, which
// registers that hook in its own init()), so without a valid pre-seeded
// config it would try to create ~/.smarter or log.Fatalf on a missing
// api_key - corrupting the developer's real home directory or crashing the
// test binary.
func TestMain(m *testing.M) {
	// a developer's SMARTER_* variables must not change the tests' outcome.
	for _, name := range cmd.SmarterEnvironmentVariables {
		if err := os.Unsetenv(name); err != nil {
			panic(err)
		}
	}
	dir, err := os.MkdirTemp("", "smarter-cli-test-home")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	configDir := filepath.Join(dir, ".smarter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(testConfigYAML), 0644); err != nil {
		panic(err)
	}

	if err := os.Setenv("HOME", dir); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
