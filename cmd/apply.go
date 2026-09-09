/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

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
the resource to the new state. The --dry-run flag simulates the apply
without making any changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileContents, err := readManifest(viper.GetString("filename"))
		if err != nil {
			ErrorOutput(err)
			return
		}

		dryRun := viper.GetBool("dry_run")
		kwargs := map[string]string{
			"dry_run": strconv.FormatBool(dryRun),
		}

		// this request goes to /api/v1/cli/apply/
		_, err = APIRequest("apply", kwargs, fileContents)
		if err != nil {
			ErrorOutput(err)
		} else if dryRun {
			fmt.Println("manifest applied. (dry run, no changes were made)")
		} else {
			fmt.Println("manifest applied.")
		}
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
