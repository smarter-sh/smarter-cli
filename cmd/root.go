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
func Execute() {
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smarter.yaml)")

	// Add the --environment flag
	RootCmd.PersistentFlags().StringVar(&environment, "environment", "", "environment to use: local, alpha, beta, next, prod. Default is prod")
	RootCmd.PersistentPreRunE = validateEnvironmentFlag

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Add this function to validate the environment flag
func validateEnvironmentFlag(cmd *cobra.Command, args []string) error {
	for _, validEnvironment := range validEnvironments {
		if environment == validEnvironment {
			return nil
		}
	}
	return errors.New("invalid environment. Valid environments are development, staging, production")
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
			"api_key":     "",
			"environment": "",
			"output":      "json",
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
	} else {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
