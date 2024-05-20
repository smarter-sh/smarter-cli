/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a Smarter manifest",
	Long: `Apply a Smarter manifest:

smarter apply -f <manifest.yaml> --dry-run

The Smarter API will apply the manifest to the Smarter account,
migrating the resource to the new state. The --json and --yaml
flags will output the manifest in the specified format. The
--dry-run flag will simulate the apply without making any changes.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		filename := viper.GetString("f")
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Failed opening file: %s", err)
		}
		defer file.Close()
		byteValue, _ := io.ReadAll(file)
		fileContents := string(byteValue)

		kwargs := map[string]string{}
		bodyJson, err := GetAPIResponse("apply", kwargs, fileContents)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.Flags().String("f", "", "Path and filename of the manifest to apply")
	if err := viper.BindPFlag("f", applyCmd.Flags().Lookup("f")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
