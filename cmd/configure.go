/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
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

var (
	accountNumberPattern = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}$`)
	apiKeyPattern        = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	usernamePattern      = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	domainPattern        = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)+$`)
)

func validateAccountNumber(v string) error {
	if !accountNumberPattern.MatchString(v) {
		return fmt.Errorf("invalid account number. Smarter account numbers use the format, 1234-5678-9012")
	}
	return nil
}

func validateApiKey(v string) error {
	if !apiKeyPattern.MatchString(v) {
		return fmt.Errorf("invalid API key. API keys should be 64 hexadecimal characters")
	}
	return nil
}

func validateUsername(v string) error {
	if !usernamePattern.MatchString(v) {
		return fmt.Errorf("invalid username. Usernames should only contain alphanumeric characters and underscores")
	}
	return nil
}

func validateRootDomain(v string) error {
	if !domainPattern.MatchString(v) {
		return fmt.Errorf("invalid root domain. Expected a domain name, e.g. platform.smarter.sh")
	}
	return nil
}

func validateOutputFormat(v string) error {
	switch strings.ToLower(v) {
	case "json", "yaml":
		return nil
	default:
		return fmt.Errorf("invalid output format. Output format should be either 'json' or 'yaml'")
	}
}

// configField describes one configurable value: how to validate it, which
// viper key it lives at, and whether it's sensitive enough to hide from
// terminal echo.
type configField struct {
	label       string
	flagName    string
	key         string // viper key that is read from and written to
	fallbackKey string // secondary viper key to check when key is unset (optional)
	mask        bool
	validate    func(string) error
}

func (f configField) currentValue() string {
	if v := viper.GetString(f.key); v != "" {
		return v
	}
	if f.fallbackKey != "" {
		return viper.GetString(f.fallbackKey)
	}
	return ""
}

func (f configField) display(value string) string {
	if f.mask {
		return "********"
	}
	return value
}

// configFields lists the values the configure command manages, in prompt
// order. api_key is stored per-environment (e.g. "alpha.api_key") to match
// how fetchAPIKey() and initConfig() read it in api.go and root.go.
func configFields() []configField {
	environment := viper.GetString("environment")
	return []configField{
		{
			label:    "root_domain",
			flagName: "root_domain",
			key:      "config.root_domain",
			validate: validateRootDomain,
		},
		{
			label:    "account_number",
			flagName: "account_number",
			key:      "config.account_number",
			validate: validateAccountNumber,
		},
		{
			label:       "api_key",
			flagName:    "api_key",
			key:         fmt.Sprintf("%s.api_key", environment),
			fallbackKey: "api_key",
			mask:        true,
			validate:    validateApiKey,
		},
		{
			label:    "username",
			flagName: "username",
			key:      "config.username",
			validate: validateUsername,
		},
		{
			label:    "output_format",
			flagName: "output_format",
			key:      "output_format",
			validate: validateOutputFormat,
		},
	}
}

// promptField interactively prompts for a field's value, re-prompting on
// invalid input. Pressing enter accepts the current value shown in
// parentheses; if there's no current value, it keeps asking. On EOF (e.g.
// stdin is closed or non-interactive), it falls back to the current value
// if one exists, or returns an error rather than looping forever.
func promptField(reader *bufio.Reader, f configField) (string, error) {
	current := f.currentValue()

	for {
		if current == "" {
			fmt.Printf("%s: ", f.label)
		} else {
			fmt.Printf("%s (%s): ", f.label, f.display(current))
		}

		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			if current != "" {
				fmt.Println()
				return current, nil
			}
			return "", fmt.Errorf("no %s provided", f.label)
		}

		if input == "" {
			if current != "" {
				return current, nil
			}
			continue
		}

		if err := f.validate(input); err != nil {
			fmt.Println(err)
			continue
		}

		if input != current {
			viper.Set(f.key, input)
			fmt.Printf("%s set to %s\n", f.label, f.display(input))
		}
		return input, nil
	}
}

// applyField validates and stores a value supplied via command-line flag.
func applyField(f configField, value string) error {
	if err := f.validate(value); err != nil {
		return err
	}
	viper.Set(f.key, value)
	fmt.Printf("%s set to %s\n", f.label, f.display(value))
	return nil
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the smarter command-line interface",
	Long: `Configure the smarter command-line interface:

smarter configure

Set your account_number, username, api_key and application options.`,
	Run: func(cmd *cobra.Command, args []string) {
		fields := configFields()

		if cmd.Flags().NFlag() > 0 {
			for _, field := range fields {
				value, _ := cmd.Flags().GetString(field.flagName)
				if value == "" {
					continue
				}
				if err := applyField(field, value); err != nil {
					fmt.Println(err)
				}
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			for _, field := range fields {
				if _, err := promptField(reader, field); err != nil {
					log.Fatalf("configure: %v", err)
				}
			}
		}

		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error writing config: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Flags
	configureCmd.Flags().StringP("root_domain", "r", "", "Smarter platform root domain (e.g. platform.smarter.sh)")
	configureCmd.Flags().StringP("account_number", "a", "", "Smarter account number")
	configureCmd.Flags().StringP("api_key", "k", "", "Smarter cli secret key (64-character hash)")
	configureCmd.Flags().StringP("username", "u", "", "username (how you login to the Smarter web console)")
	configureCmd.Flags().StringP("output_format", "o", "", "Output format (json, yaml)")
}
