package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigFilePath(t *testing.T) {
	old := cfgFile
	t.Cleanup(func() { cfgFile = old })

	cfgFile = "/tmp/elsewhere.yaml"
	if got, err := configFilePath(); err != nil || got != "/tmp/elsewhere.yaml" {
		t.Errorf("configFilePath() = %q, %v; want --config", got, err)
	}
	cfgFile = ""
	if got, err := configFilePath(); err != nil || !strings.HasSuffix(got, filepath.Join(".smarter", "config.yaml")) {
		t.Errorf("configFilePath() = %q, %v; want ~/.smarter/config.yaml", got, err)
	}
}

func TestConfigDocumentHoldsOnlyConfigurationKeys(t *testing.T) {
	withViperValue(t, "config.root_domain", "example.com")
	withViperValue(t, "alpha.api_key", testAPIKey)
	withViperValue(t, "chatbot", "should-not-be-written")
	old := outputFormatConfigured
	t.Cleanup(func() { outputFormatConfigured = old })
	outputFormatConfigured = false

	document := configDocument()
	config := document["config"].(map[string]interface{})
	if config["root_domain"] != "example.com" {
		t.Errorf("config.root_domain = %v", config["root_domain"])
	}
	for _, env := range validEnvironments {
		if _, ok := document[env]; !ok {
			t.Errorf("document has no %s section", env)
		}
	}
	if document["alpha"].(map[string]interface{})["api_key"] != testAPIKey {
		t.Error("alpha.api_key was not written")
	}
	for _, key := range []string{"chatbot", "output_format", "name", "i"} {
		if _, ok := document[key]; ok {
			t.Errorf("document has the key %q, which isn't configuration", key)
		}
	}

	outputFormatConfigured = true
	withViperValue(t, "output_format", "yaml")
	if configDocument()["output_format"] != "yaml" {
		t.Error("a configured output_format was not written")
	}
}

func TestWriteConfigFileIsOwnerOnly(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "new-dir")
	path := filepath.Join(dir, "config.yaml")

	if err := writeConfigFile(path); err != nil {
		t.Fatalf("writeConfigFile() error = %v", err)
	}
	if info, err := os.Stat(dir); err != nil || info.Mode().Perm() != 0o700 {
		t.Errorf("directory mode = %v, %v; want 0700", info.Mode().Perm(), err)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %v, %v; want 0600", info.Mode().Perm(), err)
	}
}

func TestWriteConfigFileRestrictsAnExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("old: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeConfigFile(path); err != nil {
		t.Fatalf("writeConfigFile() error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %v, %v; want 0600", info.Mode().Perm(), err)
	}
	body, _ := os.ReadFile(path)
	if strings.Contains(string(body), "old:") {
		t.Errorf("file = %q, want only configuration keys", body)
	}
}

func TestWriteConfigFileReportsErrors(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "a-file")
	if err := os.WriteFile(blocker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeConfigFile(filepath.Join(blocker, "config.yaml")); err == nil {
		t.Error("writeConfigFile() under a file = nil, want an error")
	}
	dir := t.TempDir()
	if err := writeConfigFile(dir); err == nil {
		t.Error("writeConfigFile() onto a directory = nil, want an error")
	}
}

func TestSetFieldNotesAConfiguredOutputFormat(t *testing.T) {
	old := outputFormatConfigured
	t.Cleanup(func() { outputFormatConfigured = old })
	outputFormatConfigured = false
	withViperValue(t, "output_format", "json")

	setField(configField{key: "output_format"}, "yaml")
	if !outputFormatConfigured {
		t.Error("setField(output_format) didn't note the configured output format")
	}
}

// exits replaces ErrorOutput's exit with a panic for the duration of the
// test, runs fn, and reports whether it exited.
func exits(t *testing.T, fn func()) (exited bool) {
	t.Helper()
	type exit struct{}
	SetExitFunc(func(int) { panic(exit{}) })
	t.Cleanup(func() { SetExitFunc(nil) })
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(exit); !ok {
				panic(r)
			}
			exited = true
		}
	}()
	fn()
	return false
}

func TestConfigureExitsOnAnInvalidFlagValue(t *testing.T) {
	if err := configureCmd.Flags().Set("account_number", "not-a-number"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = configureCmd.Flags().Set("account_number", "") })

	if !exits(t, func() {
		_ = captureStdout(t, func() { configureCmd.Run(configureCmd, nil) })
	}) {
		t.Error("configure with an invalid account number did not exit")
	}
}
