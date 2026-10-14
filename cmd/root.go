/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Version is the cli's version, set by main from its -ldflags value.
var Version = "local.dev"

var RootCmd = &cobra.Command{
	Use:   "smarter",
	Short: "A command-line interface for working with Smarter resources",
	Long: `A command-line interface for working with Smarter resources.
Using the smarter cli, you can create Smarter plugins, add these to a ChatBot,
and deploy the ChatBot to a custom URL. You can interact with the ChatBot
on the command line, view chat log data, and manage your Smarter account.
Support: https://smarter.sh and support@smarter.sh.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version string) {
	Version = version
	if err := RootCmd.Execute(); err != nil {
		exitFunc(1)
	}
}

var environment string
var validEnvironments = []string{"local", "alpha", "beta", "next", "prod"}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smarter/config.yaml)")

	// Add the --environment flag
	RootCmd.PersistentFlags().StringVar(&environment, "environment", "", "environment to use: local, alpha, beta, next, prod. Default is prod")
	mustBindPFlag("environment", RootCmd.PersistentFlags().Lookup("environment"))

	// Add the --api_key flag
	RootCmd.PersistentFlags().String("api_key", "", "Smarter API key to use")
	mustBindPFlag("api_key", RootCmd.PersistentFlags().Lookup("api_key"))

	// Add the --verbose toggle
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	mustBindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	// Add the --output_format flag
	RootCmd.PersistentFlags().StringP("output_format", "o", "json", "output format: json, yaml, tabular")
	mustBindPFlag("output_format", RootCmd.PersistentFlags().Lookup("output_format"))

	// Add the --timeout flag
	RootCmd.PersistentFlags().Duration("timeout", DefaultTimeout, "how long to wait for the Smarter API, e.g. 30s or 5m")
	mustBindPFlag("timeout", RootCmd.PersistentFlags().Lookup("timeout"))

	// Add the --max_retries flag
	RootCmd.PersistentFlags().Int("max_retries", DefaultMaxRetries, "how many times to resend a request that failed for a transient reason")
	mustBindPFlag("max_retries", RootCmd.PersistentFlags().Lookup("max_retries"))

	// Bind the flag value validators
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateEnvironmentFlag(); err != nil {
			return err
		}
		if err := validateOutputToggles(); err != nil {
			return err
		}
		if err := validateTabular(cmd); err != nil {
			return err
		}
		return validateSmarterEnvironment()
	}
}

func validateOutputToggles() error {
	outputFormat := viper.GetString("output_format")

	// tabular is used by the get commands
	for _, format := range []string{"json", "yaml", "tabular"} {
		if outputFormat == format {
			return nil
		}
	}
	return fmt.Errorf("invalid output format: %v. Valid formats are 'json', 'yaml' or 'tabular'", outputFormat)
}

// validateTabular rejects -o tabular outside the get commands: only get
// responses have the titles and items of a table.
func validateTabular(cmd *cobra.Command) error {
	if viper.GetString("output_format") != "tabular" {
		return nil
	}
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == "get" && c.Parent() == RootCmd {
			return nil
		}
	}
	return fmt.Errorf("the tabular output format is only available for get commands; use json or yaml")
}

func validateEnvironmentFlag() error {
	if environment == "" {
		return nil
	}
	for _, validEnvironment := range validEnvironments {
		if environment == validEnvironment {
			return nil
		}
	}
	return errors.New("invalid environment. Valid values: " + strings.Join(validEnvironments, ", "))
}

// initConfig reads in config file and ENV variables if set.
//
// It never exits: a missing api_key is reported by the first command that
// needs one (verifyApiKey), so that "smarter configure" and "smarter version"
// work on a fresh install. The SMARTER_* variables that it needs are applied
// first (applySmarterEnvironment).
func initConfig() {
	applySmarterEnvironment()
	configFile, err := configFilePath()
	cobra.CheckErr(err)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in .smarter directory in home directory with name "config" (without extension).
		viper.SetConfigFile(configFile)
		viper.SetConfigType("yaml")
		viper.SetEnvPrefix("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Set defaults unconditionally so that they backfill missing keys in an
	// existing config file as well as seed a brand-new one.
	viper.SetDefault("config", map[string]interface{}{
		"account_number": "",
		"environment":    "prod",
		"output_format":  "yaml",
		"root_domain":    "platform.smarter.sh",
		"username":       "",
	})
	for _, env := range validEnvironments {
		viper.SetDefault(env, map[string]interface{}{"api_key": ""})
	}

	// If a config file is found, read it in. Otherwise, create one.
	if err := viper.ReadInConfig(); err != nil {
		if _, statErr := os.Stat(configFile); !os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "Failed to read config file %s: %v\n", configFile, err)
		} else {
			fmt.Fprintln(os.Stderr, strings.Repeat("*", 80))
			fmt.Fprintln(os.Stderr, "Welcome to the Smarter CLI!")
			fmt.Fprintln(os.Stderr, "Please note your smarter configuration file:", configFile)
			fmt.Fprintln(os.Stderr, "Set your credentials with: smarter configure")
			fmt.Fprintln(os.Stderr, strings.Repeat("*", 80))
			if err := writeConfigFile(configFile); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write default config file:", err)
			}
		}
	} else if viper.GetBool("verbose") {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set the environment from the config file if not set by flag
	if viper.GetString("environment") == "" {
		environment = viper.GetString("config.environment")
		if environment == "" {
			environment = "prod"
		}
		viper.Set("environment", environment)
	}

	if viper.GetBool("verbose") {
		log.Printf("Environment set to: %s", viper.GetString("environment"))
		if apiKey := fetchAPIKey(); apiKey != "" {
			log.Printf("API key set to: %s", maskSecret(apiKey))
		}
	}
}
