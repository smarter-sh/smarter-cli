/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>

Reading and writing the configuration file, ~/.smarter/config.yaml. Its format
is a contract with the Smarter SDKs, which read it: see doc/adr/002-configuration.md.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

// configSectionKeys are the keys of the configuration file's "config" section.
var configSectionKeys = []string{"account_number", "environment", "output_format", "root_domain", "username"}

// outputFormatConfigured records that "smarter configure" set the top-level
// output_format key in this run, so that writeConfigFile keeps it.
var outputFormatConfigured bool

// configFilePath returns --config, or ~/.smarter/config.yaml.
func configFilePath() (string, error) {
	if cfgFile != "" {
		return cfgFile, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".smarter", "config.yaml"), nil
}

// configDocument returns the configuration as it is written to the file: the
// "config" section, an api_key per environment, and output_format when the
// user set it. Flags that viper also knows (name, i, chatbot, ...) are never
// written, unlike viper.WriteConfig, which writes every key.
func configDocument() map[string]interface{} {
	config := map[string]interface{}{}
	for _, key := range configSectionKeys {
		config[key] = viper.GetString("config." + key)
	}
	document := map[string]interface{}{"config": config}
	for _, env := range validEnvironments {
		document[env] = map[string]interface{}{"api_key": viper.GetString(env + ".api_key")}
	}
	if outputFormatConfigured || viper.InConfig("output_format") {
		document["output_format"] = viper.GetString("output_format")
	}
	return document
}

// writeConfigFile writes configDocument to path, readable only by its owner,
// because it holds api keys: see doc/adr/016-security.md.
func writeConfigFile(path string) error {
	body, err := yaml.Marshal(configDocument())
	if err != nil {
		return fmt.Errorf("failed formatting the configuration: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("failed creating %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return fmt.Errorf("failed writing %s: %w", path, err)
	}
	// os.WriteFile keeps the mode of an existing file, which older versions
	// of the cli wrote world-readable.
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("failed restricting %s: %w", path, err)
	}
	return nil
}
