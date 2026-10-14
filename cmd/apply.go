/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// readManifest loads the manifest file named by filename and returns its
// contents as a string, or an error if the file can't be opened or read.
func readManifest(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed opening manifest file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed reading manifest file: %w", err)
	}
	return string(byteValue), nil
}

var applyCmd = &cobra.Command{
	Use:   "apply -f <manifest.yaml> [--dry-run]",
	Short: "Apply a Smarter manifest",
	Long: `Apply a Smarter manifest:

smarter apply -f <manifest.yaml> [--dry-run]

The Smarter API will apply the manifest to the Smarter account, migrating
the resource to the new state. The --dry-run flag validates the manifest
instead, as smarter validate does, and changes nothing.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileContents, err := readManifest(viper.GetString("filename"))
		if err != nil {
			ErrorOutput(err)
			return
		}

		// The platform's apply has no dry run, so --dry-run validates instead:
		// /api/v1/cli/validate/ checks the manifest as apply does first,
		// without saving it.
		if viper.GetBool("dry_run") {
			bodyJson, err := APIRequest("validate", map[string]string{}, fileContents)
			if err == nil {
				err = reportValidation(bodyJson)
			}
			if err != nil {
				ErrorOutput(err)
				return
			}
			fmt.Println("(dry run: the manifest was validated, and nothing was applied)")
			return
		}

		// this request goes to /api/v1/cli/apply/
		if _, err := APIRequest("apply", map[string]string{}, fileContents); err != nil {
			ErrorOutput(err)
			return
		}
		fmt.Println("manifest applied.")
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringP("filename", "f", "", "Path and filename of the manifest to apply")
	mustBindPFlag("filename", applyCmd.Flags().Lookup("filename"))
	mustMarkRequired(applyCmd, "filename")

	applyCmd.Flags().Bool("dry-run", false, "Simulate the apply without making any changes")
	mustBindPFlag("dry_run", applyCmd.Flags().Lookup("dry-run"))
}
