package cmd

import (
	"bufio"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateAccountNumber(t *testing.T) {
	valid := []string{"1234-5678-9012"}
	invalid := []string{"", "12345678901", "1234-5678-901", "abcd-efgh-ijkl", "1234 5678 9012"}
	for _, v := range valid {
		if err := validateAccountNumber(v); err != nil {
			t.Errorf("validateAccountNumber(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range invalid {
		if err := validateAccountNumber(v); err == nil {
			t.Errorf("validateAccountNumber(%q) = nil, want an error", v)
		}
	}
}

func TestValidateApiKey(t *testing.T) {
	valid := strings.Repeat("a1", 32) // 64 hex chars
	if err := validateApiKey(valid); err != nil {
		t.Errorf("validateApiKey(valid) = %v, want nil", err)
	}
	invalid := []string{"", "short", strings.Repeat("z", 64), strings.Repeat("a", 63)}
	for _, v := range invalid {
		if err := validateApiKey(v); err == nil {
			t.Errorf("validateApiKey(%q) = nil, want an error", v)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	valid := []string{"alice", "alice_42", "ALICE"}
	invalid := []string{"", "alice smith", "alice-smith", "alice@example.com"}
	for _, v := range valid {
		if err := validateUsername(v); err != nil {
			t.Errorf("validateUsername(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range invalid {
		if err := validateUsername(v); err == nil {
			t.Errorf("validateUsername(%q) = nil, want an error", v)
		}
	}
}

func TestValidateRootDomain(t *testing.T) {
	valid := []string{"platform.smarter.sh", "example.com", "sub.example.co.uk"}
	invalid := []string{"", "no-dot", "-badstart.com", ".com"}
	for _, v := range valid {
		if err := validateRootDomain(v); err != nil {
			t.Errorf("validateRootDomain(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range invalid {
		if err := validateRootDomain(v); err == nil {
			t.Errorf("validateRootDomain(%q) = nil, want an error", v)
		}
	}
}

func TestValidateOutputFormat(t *testing.T) {
	valid := []string{"json", "yaml", "JSON", "Yaml"}
	invalid := []string{"", "tabular", "xml"}
	for _, v := range valid {
		if err := validateOutputFormat(v); err != nil {
			t.Errorf("validateOutputFormat(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range invalid {
		if err := validateOutputFormat(v); err == nil {
			t.Errorf("validateOutputFormat(%q) = nil, want an error", v)
		}
	}
}

func TestConfigFieldCurrentValue(t *testing.T) {
	withViperValue(t, "config.root_domain", "example.com")
	withViperValue(t, "fallback.key", "fallback-value")

	f := configField{key: "config.root_domain"}
	if got := f.currentValue(); got != "example.com" {
		t.Errorf("currentValue() = %q, want example.com", got)
	}

	f2 := configField{key: "missing.key", fallbackKey: "fallback.key"}
	if got := f2.currentValue(); got != "fallback-value" {
		t.Errorf("currentValue() with fallback = %q, want fallback-value", got)
	}

	f3 := configField{key: "missing.key"}
	if got := f3.currentValue(); got != "" {
		t.Errorf("currentValue() with no fallback = %q, want empty", got)
	}
}

func TestConfigFieldDisplay(t *testing.T) {
	masked := configField{mask: true}
	if got := masked.display("secret"); got != "********" {
		t.Errorf("display() = %q, want masked", got)
	}
	plain := configField{mask: false}
	if got := plain.display("visible"); got != "visible" {
		t.Errorf("display() = %q, want visible", got)
	}
}

func TestConfigFieldsUsesCurrentEnvironmentForApiKey(t *testing.T) {
	withViperValue(t, "environment", "beta")

	fields := configFields()
	var apiKeyField *configField
	for i := range fields {
		if fields[i].label == "api_key" {
			apiKeyField = &fields[i]
		}
	}
	if apiKeyField == nil {
		t.Fatal("configFields() did not include an api_key field")
	}
	if apiKeyField.key != "beta.api_key" {
		t.Errorf("api_key field key = %q, want beta.api_key", apiKeyField.key)
	}
	if !apiKeyField.mask {
		t.Error("api_key field should be masked")
	}
}

func TestApplyFieldValidAndInvalid(t *testing.T) {
	withViperValue(t, "config.root_domain", "")
	f := configField{label: "root_domain", key: "config.root_domain", validate: validateRootDomain}

	if err := applyField(f, "example.com"); err != nil {
		t.Fatalf("applyField() error = %v", err)
	}
	if got := viper.GetString("config.root_domain"); got != "example.com" {
		t.Errorf("viper value = %q, want example.com", got)
	}

	if err := applyField(f, "not a domain"); err == nil {
		t.Error("applyField() error = nil, want an error for an invalid value")
	}
}

func TestPromptFieldAcceptsValidInput(t *testing.T) {
	withViperValue(t, "config.root_domain", "")
	f := configField{label: "root_domain", key: "config.root_domain", validate: validateRootDomain}

	reader := bufio.NewReader(strings.NewReader("example.com\n"))
	got, err := promptField(reader, f)
	if err != nil {
		t.Fatalf("promptField() error = %v", err)
	}
	if got != "example.com" {
		t.Errorf("promptField() = %q, want example.com", got)
	}
}

func TestPromptFieldReprompts(t *testing.T) {
	withViperValue(t, "config.account_number", "")
	f := configField{label: "account_number", key: "config.account_number", validate: validateAccountNumber}

	// First line is invalid, second is valid; promptField should re-prompt
	// and use the second line.
	reader := bufio.NewReader(strings.NewReader("not-valid\n1234-5678-9012\n"))
	got, err := promptField(reader, f)
	if err != nil {
		t.Fatalf("promptField() error = %v", err)
	}
	if got != "1234-5678-9012" {
		t.Errorf("promptField() = %q, want 1234-5678-9012", got)
	}
}

func TestPromptFieldKeepsCurrentValueOnEmptyInput(t *testing.T) {
	withViperValue(t, "config.root_domain", "existing.com")
	f := configField{label: "root_domain", key: "config.root_domain", validate: validateRootDomain}

	reader := bufio.NewReader(strings.NewReader("\n"))
	got, err := promptField(reader, f)
	if err != nil {
		t.Fatalf("promptField() error = %v", err)
	}
	if got != "existing.com" {
		t.Errorf("promptField() = %q, want existing.com (unchanged)", got)
	}
}

func TestPromptFieldEOFWithNoCurrentValueReturnsError(t *testing.T) {
	withViperValue(t, "config.root_domain", "")
	f := configField{label: "root_domain", key: "config.root_domain", validate: validateRootDomain}

	reader := bufio.NewReader(strings.NewReader(""))
	if _, err := promptField(reader, f); err == nil {
		t.Error("promptField() error = nil, want an error on EOF with no current value")
	}
}

func TestPromptFieldEOFWithCurrentValueFallsBack(t *testing.T) {
	withViperValue(t, "config.root_domain", "existing.com")
	f := configField{label: "root_domain", key: "config.root_domain", validate: validateRootDomain}

	reader := bufio.NewReader(strings.NewReader(""))
	got, err := promptField(reader, f)
	if err != nil {
		t.Fatalf("promptField() error = %v", err)
	}
	if got != "existing.com" {
		t.Errorf("promptField() = %q, want existing.com", got)
	}
}
