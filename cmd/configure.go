/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Account Number
//
// This section contains code related to account_number configuration.
func validateAccountNumber(accountNumber string) error {
	regex, err := regexp.Compile(`^\d{4}-\d{4}-\d{4}$`)
	if err != nil {
		return err
	}

	if !regex.MatchString(accountNumber) {
		return fmt.Errorf("invalid account number. Smarter account numbers use the format, 1234-5678-9012")
	}

	return nil
}

func getAccountNumber() string {
	fmt.Println("getAccountNumber()")
	accountNumber := viper.Get("account_number").(string)
	reader := bufio.NewReader(os.Stdin)
	valid := false

	for !valid {
		if accountNumber == "" {
			fmt.Print("account_number: ")
		} else {
			fmt.Printf("account_number (%s): ", accountNumber)
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			err := validateAccountNumber(input)
			if err != nil {
				fmt.Println(err)
			} else {
				if accountNumber != input {
					accountNumber = input
					viper.Set("config.account_number", accountNumber)
					fmt.Println("Account number set to", accountNumber)
				}
				valid = true
			}
		}
	}
	return accountNumber
}

func setAccountNumber(accountNumber string) {
	if accountNumber != "" {
		err := validateAccountNumber(accountNumber)
		if err != nil {
			fmt.Println(err)
		} else {
			viper.Set("config.account_number", accountNumber)
			fmt.Println("Account number set to", accountNumber)
		}
	}
}

// Api Key
//
// This section contains code related to api_key configuration.
func validateApiKey(apiKey string) error {
	regex, err := regexp.Compile(`^[a-fA-F0-9]{64}$`)
	if err != nil {
		return err
	}

	if !regex.MatchString(apiKey) {
		return fmt.Errorf("invalid API key. API keys should be 64 hexadecimal characters")
	}

	return nil
}

func getApiKey() string {
	fmt.Println("getApiKey()")
	apiKey := viper.Get("config.api_key").(string)
	if apiKey == "" {
		apiKey = viper.Get("api_key").(string)
	}
	reader := bufio.NewReader(os.Stdin)
	valid := false

	for !valid {
		if apiKey == "" {
			fmt.Print("api_key: ")
		} else {
			fmt.Printf("api_key (%s): ", apiKey)
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			err := validateApiKey(input)
			if err != nil {
				fmt.Println("Invalid API key. API keys should be a 64-character hash string.")
			} else {
				if apiKey != input {
					apiKey = input
					viper.Set("config.api_key", apiKey)
					fmt.Println("API key set to", apiKey)
				}
				valid = true
			}
		}
	}
	return apiKey
}

func setApiKey(apiKey string) {
	if apiKey != "" {
		err := validateApiKey(apiKey)
		if err != nil {
			fmt.Println(err)
		} else {
			viper.Set("config.api_key", apiKey)
			fmt.Println("API key set to", apiKey)
		}
	}
}

// Username
//
// This section contains code related to username configuration.
func validateUsername(username string) error {
	regex, err := regexp.Compile(`^[a-zA-Z0-9_]+$`)
	if err != nil {
		return err
	}

	if !regex.MatchString(username) {
		return fmt.Errorf("invalid username. Usernames should only contain alphanumeric characters and underscores")
	}

	return nil
}

func getUsername() string {
	fmt.Println("getUsername()")
	username := viper.Get("username").(string)
	reader := bufio.NewReader(os.Stdin)
	if username == "" {
		fmt.Print("username: ")
	} else {
		fmt.Printf("username (%s): ", username)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		err := validateUsername(input)
		if err != nil {
			fmt.Println("Invalid username. Please input your username for the Smarter web console https://platform.smarter.sh. Please try again.")
			return getUsername()
		}
		if username != input {
			username = input
			viper.Set("config.username", username)
			fmt.Println("Username set to", username)
		}
	}
	return username
}

func setUsername(username string) {
	err := validateUsername(username)
	if err != nil {
		fmt.Println(err)
	} else {
		viper.Set("config.username", username)
		fmt.Println("Username set to", username)
	}
}

// Output Format
//
// This section contains code related to output_format configuration.
func validateOutputFormat(format string) error {
	lowerFormat := strings.ToLower(format)

	if lowerFormat != "json" && lowerFormat != "yaml" {
		return fmt.Errorf("invalid output format. Output format should be either 'json' or 'yaml'")
	}

	return nil
}

func getOutputFormat() string {
	fmt.Println("getOutputFormat()")
	outputFormat := viper.Get("output_format").(string)
	reader := bufio.NewReader(os.Stdin)
	valid := false

	for !valid {
		if outputFormat == "" {
			fmt.Print("output_format: ")
		} else {
			fmt.Printf("output_format (%s): ", outputFormat)
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			err := validateOutputFormat(input)
			if err != nil {
				fmt.Println("Invalid output format. Allowed values are 'json' and 'yaml'.")
			} else {
				if outputFormat != input {
					outputFormat = input
					viper.Set("config.output_format", outputFormat)
					fmt.Println("Output format set to", outputFormat)
				}
				valid = true
			}
		}
	}
	return outputFormat
}

func setOutputFormat(outputFormat string) {
	err := validateOutputFormat(outputFormat)
	if err != nil {
		fmt.Println(err)
	} else {
		viper.Set("config.output_format", outputFormat)
		fmt.Println("Output format set to", outputFormat)
	}
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the smarter command-line interface",
	Long: `Configure the smarter command-line interface:

smarter configure

Set your account_number, username, api_key and application options.`,
	Run: func(cmd *cobra.Command, args []string) {
		accountNumber, _ := cmd.Flags().GetString("account_number")
		apiKey, _ := cmd.Flags().GetString("api_key")
		username, _ := cmd.Flags().GetString("username")
		outputFormat, _ := cmd.Flags().GetString("output_format")

		if cmd.Flags().NFlag() > 0 {
			if accountNumber != "" {
				setAccountNumber(accountNumber)
			}
			if apiKey != "" {
				setApiKey(apiKey)
			}
			if username != "" {
				setUsername(username)
			}
			if outputFormat != "" {
				setOutputFormat(outputFormat)
			}
		} else {
			if accountNumber == "" {
				getAccountNumber()
			}
			if apiKey == "" {
				getApiKey()
			}
			if username == "" {
				getUsername()
			}
			if outputFormat == "" {
				getOutputFormat()
			}
		}

		err := viper.WriteConfig()
		if err != nil {
			log.Fatalf("Error writing config: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Flags
	configureCmd.Flags().StringP("account_number", "a", "", "Smarter account number")
	configureCmd.Flags().StringP("api_key", "k", "", "Smarter cli secret key (64-character hash)")
	configureCmd.Flags().StringP("username", "u", "", "username (how you login to the Smarter web console)")
	configureCmd.Flags().StringP("output_format", "o", "", "Output format (json, yaml)")
}
