/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>

The SMARTER_* environment variables, which the Smarter SDKs read too: see
doc/adr/002-configuration.md. A flag that is passed wins over its variable,
and a variable wins over the configuration file.

SMARTER_ROOT_DOMAIN is deliberately not read: the platform's own settings use
that name for its apex domain (e.g. smarter.sh, its marketing site), and a
developer's shell often exports it, which would send the cli's requests there.
*/

package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// SmarterEnvironmentVariables are the SMARTER_* variables that the cli reads.
// Tests unset them, so that a developer's shell can't change their outcome.
var SmarterEnvironmentVariables = []string{
	"SMARTER_API_KEY",
	"SMARTER_BASE_URL",
	"SMARTER_CONFIG_FILE",
	"SMARTER_ENVIRONMENT",
	"SMARTER_LOG",
	"SMARTER_MAX_RETRIES",
	"SMARTER_TIMEOUT",
}

// DefaultMaxRetries is how many times a request that failed for a transient
// reason is sent again, as in the SDKs.
const DefaultMaxRetries = 2

// smarterEnv returns the value of a SMARTER_* environment variable, trimmed.
func smarterEnv(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

// flagChanged reports whether a global flag was passed on the command line.
func flagChanged(name string) bool {
	flag := RootCmd.PersistentFlags().Lookup(name)
	return flag != nil && flag.Changed
}

// validateSmarterEnvironment rejects SMARTER_* variables whose values can't be
// used, with the same messages as the SDKs.
func validateSmarterEnvironment() error {
	if value := smarterEnv("SMARTER_ENVIRONMENT"); value != "" {
		valid := false
		for _, env := range validEnvironments {
			valid = valid || value == env
		}
		if !valid {
			return fmt.Errorf("SMARTER_ENVIRONMENT must be one of %s, not %q", strings.Join(validEnvironments, ", "), value)
		}
	}
	if value := smarterEnv("SMARTER_BASE_URL"); value != "" {
		parsed, err := url.Parse(value)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return fmt.Errorf("SMARTER_BASE_URL must be an http or https url, not %q", value)
		}
	}
	if value := smarterEnv("SMARTER_TIMEOUT"); value != "" {
		if seconds, err := strconv.ParseFloat(value, 64); err != nil || seconds <= 0 {
			return fmt.Errorf("SMARTER_TIMEOUT must be a number of seconds, more than 0, not %q", value)
		}
	}
	if value := smarterEnv("SMARTER_MAX_RETRIES"); value != "" {
		if retries, err := strconv.Atoi(value); err != nil || retries < 0 {
			return fmt.Errorf("SMARTER_MAX_RETRIES must be a whole number, 0 or more, not %q", value)
		}
	}
	if value := smarterEnv("SMARTER_LOG"); value != "" && value != "debug" && value != "info" {
		return fmt.Errorf("SMARTER_LOG must be debug or info, not %q", value)
	}
	return nil
}

// applySmarterEnvironment applies the variables that initConfig needs before
// it resolves the configuration: SMARTER_CONFIG_FILE, SMARTER_ENVIRONMENT and
// SMARTER_LOG. Invalid values are ignored here, and reported by
// validateSmarterEnvironment before the command runs.
func applySmarterEnvironment() {
	if value := smarterEnv("SMARTER_CONFIG_FILE"); value != "" && cfgFile == "" {
		cfgFile = value
	}
	if value := smarterEnv("SMARTER_ENVIRONMENT"); value != "" && !flagChanged("environment") {
		for _, env := range validEnvironments {
			if value == env {
				environment = value
				viper.Set("environment", value)
			}
		}
	}
	if value := smarterEnv("SMARTER_LOG"); (value == "debug" || value == "info") && !flagChanged("verbose") {
		viper.Set("verbose", true)
	}
}

// envBaseURL returns SMARTER_BASE_URL without a trailing slash, unless
// --environment was passed, which then decides the host.
func envBaseURL() string {
	if flagChanged("environment") {
		return ""
	}
	return strings.TrimRight(smarterEnv("SMARTER_BASE_URL"), "/")
}

// requestTimeout returns --timeout when it was passed, else SMARTER_TIMEOUT
// (seconds), else the configured or default timeout.
func requestTimeout() time.Duration {
	if value := smarterEnv("SMARTER_TIMEOUT"); value != "" && !flagChanged("timeout") {
		if seconds, err := strconv.ParseFloat(value, 64); err == nil && seconds > 0 {
			return time.Duration(seconds * float64(time.Second))
		}
	}
	if timeout := viper.GetDuration("timeout"); timeout > 0 {
		return timeout
	}
	return DefaultTimeout
}

// maxRetries returns --max_retries when it was passed, else
// SMARTER_MAX_RETRIES, else the configured or default number of retries.
func maxRetries() int {
	if value := smarterEnv("SMARTER_MAX_RETRIES"); value != "" && !flagChanged("max_retries") {
		if retries, err := strconv.Atoi(value); err == nil && retries >= 0 {
			return retries
		}
	}
	if !viper.IsSet("max_retries") {
		return DefaultMaxRetries
	}
	if retries := viper.GetInt("max_retries"); retries >= 0 {
		return retries
	}
	return DefaultMaxRetries
}
