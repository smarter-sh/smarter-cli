package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

// TestConfigureCommandFlagMode exercises configureCmd.Run's non-interactive
// branch (cmd.Flags().NFlag() > 0), which applies flag values directly
// instead of prompting on stdin. Calling Run directly (rather than
// Execute()) sidesteps cobra.OnInitialize(initConfig) so this test isn't
// coupled to global config-file state; it still needs real pflag.Changed
// tracking for NFlag(), which Flags().Set provides.
func TestConfigureCommandFlagMode(t *testing.T) {
	t.Cleanup(func() {
		_ = configureCmd.Flags().Set("root_domain", "")
		_ = configureCmd.Flags().Set("account_number", "")
	})

	// configureCmd.Run ends with viper.WriteConfig(), which requires viper to
	// already have a config file associated (normally guaranteed in
	// production by cobra.OnInitialize(initConfig) running before any
	// command's Run). Since this test calls Run directly to avoid the
	// cobra/OnInitialize machinery, and depending on test execution order
	// nothing else may have triggered initConfig yet, call it explicitly so
	// viper.WriteConfig() has somewhere to write (TestMain's fixture home).
	initConfig()

	withViperValue(t, "config.root_domain", "")
	withViperValue(t, "config.account_number", "")

	if err := configureCmd.Flags().Set("root_domain", "example.com"); err != nil {
		t.Fatalf("Set(root_domain) error = %v", err)
	}
	if err := configureCmd.Flags().Set("account_number", "1234-5678-9012"); err != nil {
		t.Fatalf("Set(account_number) error = %v", err)
	}

	out := captureStdout(t, func() {
		configureCmd.Run(configureCmd, []string{})
	})

	if got := viper.GetString("config.root_domain"); got != "example.com" {
		t.Errorf("config.root_domain = %q, want example.com", got)
	}
	if got := viper.GetString("config.account_number"); got != "1234-5678-9012" {
		t.Errorf("config.account_number = %q, want 1234-5678-9012", got)
	}
	if out == "" {
		t.Error("expected configureCmd.Run to print confirmation of the applied values")
	}
}

func TestConfigureRegisteredOnRootCmd(t *testing.T) {
	found := false
	for _, c := range RootCmd.Commands() {
		if c.Name() == "configure" {
			found = true
		}
	}
	if !found {
		t.Error("configure command is not registered on RootCmd")
	}
}
