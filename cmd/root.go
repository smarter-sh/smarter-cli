/*
Copyright © 2024 Lawrence McDaniel lawrence@querium.com
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Version string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "smarter",
	Short: "A command-line interface for working with Smarter resources",
	Long: `A command-line interface for working with Smarter resources.
Using the smarter cli, you can create Smarter plugins, add these to a ChatBot,
and deploy the ChatBot to a custom URL. You can interact with the ChatBot
on the command line, view chat log data, and manage your Smarter account.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version string) {
	Version = version
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var environment string
var validEnvironments = []string{"local", "alpha", "beta", "next", "prod"}

func init() {
	cobra.OnInitialize(initConfig)
	initConfig()

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smarter.yaml)")

	// Add the --environment flag
	// Set up a global --environment flag and bind this to viper.
	RootCmd.PersistentFlags().StringVar(&environment, "environment", "", "environment to use: local, alpha, beta, next, prod. Default is prod")
	if err := viper.BindPFlag("environment", RootCmd.PersistentFlags().Lookup("environment")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	// Add the --api_key flag
	RootCmd.PersistentFlags().String("api_key", "", "Smarter API key to use")
	if err := viper.BindPFlag("api_key", RootCmd.PersistentFlags().Lookup("api_key")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	// Add the --json toggle
	RootCmd.PersistentFlags().BoolP("json", "j", false, "output in JSON format")
	if err := viper.BindPFlag("json", RootCmd.PersistentFlags().Lookup("json")); err != nil {
		log.Fatalf("Error binding toggle: %v", err)
	}

	// Add the --yaml toggle
	RootCmd.PersistentFlags().BoolP("yaml", "y", false, "output in YAML format")
	if err := viper.BindPFlag("yaml", RootCmd.PersistentFlags().Lookup("yaml")); err != nil {
		log.Fatalf("Error binding toggle: %v", err)
	}

	// Bind the flag value validators
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateEnvironmentFlag(); err != nil {
			return err
		}
		if err := validateOutputToggles(); err != nil {
			return err
		}
		return nil
	}

}

func validateOutputToggles() error {
	jsonOutput := viper.GetBool("json")
	yamlOutput := viper.GetBool("yaml")
	output_format := viper.GetString("config.output_format")

	if jsonOutput && yamlOutput {
		return errors.New("cannot specify both --json and --yaml")
	}
	if !jsonOutput && !yamlOutput {
		// check the config file
		if output_format == "json" {
			viper.Set("json", true)
		} else {
			if output_format != "yaml" {
				viper.Set("yaml", true)
			}
		}
	}
	return nil

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
func initConfig() {

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configDir := filepath.Join(home, ".smarter")
	configFile := filepath.Join(configDir, "config.yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in .smarter directory in home directory with name "config" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in. Otherwise,
	// create a default config file.
	if err := viper.ReadInConfig(); err != nil {
		defaultConfig := map[string]interface{}{
			"account_number": "",
			"username":       "",
			"api_key":        "",
			"environment":    "",
			"output_format":  "",
		}
		viper.SetDefault("config", defaultConfig)

		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.Mkdir(configDir, 0755); err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(os.Stderr, strings.Repeat("*", 80))
			fmt.Fprintln(os.Stderr, "Welcome to the Smarter CLI!")
			fmt.Fprintln(os.Stderr, "Please note your smarter configuration path:", configDir)
			fmt.Fprintln(os.Stderr, strings.Repeat("*", 80))
		}

		err := viper.SafeWriteConfigAs(configFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write default config file:", err)
		}
	}
}
