package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// testAPIKey is a syntactically valid (64 hex char) placeholder API key used
// throughout the test config fixture below.
const testAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcd"

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
// this package runs. This matters because cobra.OnInitialize(initConfig)
// fires on *every* Command.Execute() call, including ones made deep in this
// package's tests against leaf commands. Without a valid pre-seeded config,
// initConfig would try to create ~/.smarter or log.Fatalf on a missing
// api_key, which would corrupt the developer's real home directory or crash
// the test binary.
func TestMain(m *testing.M) {
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
