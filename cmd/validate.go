/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validation is the data of a validate response: see the platform's
// CliApiValidation in its OpenAPI document.
type validation struct {
	Data struct {
		Valid  bool `json:"valid"`
		Errors []struct {
			Loc     []interface{} `json:"loc"`
			Message string        `json:"message"`
			Type    string        `json:"type"`
		} `json:"errors"`
	} `json:"data"`
}

// reportValidation prints the outcome of a validate response, and returns an
// error when the manifest is invalid, so that the command exits 1.
func reportValidation(bodyJson []byte) error {
	var result validation
	if err := json.Unmarshal(bodyJson, &result); err != nil {
		return fmt.Errorf("failed parsing the validation: %w", err)
	}
	if viper.IsSet("output_format") {
		if err := ConsoleOutput(bodyJson); err != nil {
			return err
		}
	}
	if result.Data.Valid {
		fmt.Println("manifest is valid.")
		return nil
	}
	for _, e := range result.Data.Errors {
		path := make([]string, len(e.Loc))
		for i, part := range e.Loc {
			path[i] = fmt.Sprint(part)
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", strings.Join(path, "."), e.Message)
	}
	return fmt.Errorf("manifest has %d error(s)", len(result.Data.Errors))
}

var validateCmd = &cobra.Command{
	Use:   "validate -f <manifest.yaml>",
	Short: "Validate a Smarter manifest without applying it",
	Long: `Validate a Smarter manifest without applying it:

smarter validate -f <manifest.yaml>

The Smarter API will check the manifest against the model of its kind, as
apply does first, and report each error with its path in the manifest. The
command exits 1 when the manifest is invalid, so it can be used in CI.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flag directly: "filename" is apply's viper key.
		filename, _ := cmd.Flags().GetString("filename")
		fileContents, err := readManifest(filename)
		if err != nil {
			ErrorOutput(err)
			return
		}

		// this request goes to /api/v1/cli/validate/
		bodyJson, err := APIRequest("validate", map[string]string{}, fileContents)
		if err == nil {
			err = reportValidation(bodyJson)
		}
		if err != nil {
			ErrorOutput(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("filename", "f", "", "Path and filename of the manifest to validate")
	mustMarkRequired(validateCmd, "filename")
}
