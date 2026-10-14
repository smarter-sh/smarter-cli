/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	// en route to /api/v1/cli/get/<str:kind>/
	return cmd.APIRequest("get/"+kind+"/", kwargs)

}

func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "tabular")
	}
	if err := cmd.ConsoleOutput(bodyJson); err != nil {
		cmd.ErrorOutput(err)
	}
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Generate a list of Smarter resources",
	Long: `Generate a list of Smarter resources:

smarter get [kind] [flags]

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().Int("i", 10, "Number of resources to retrieve")
	getCmd.PersistentFlags().Bool("asc", false, "Sort results in ascending order")
	getCmd.PersistentFlags().Bool("desc", false, "Sort results in descending order")
	// The Smarter API ignores these: each kind returns its resources in a fixed
	// order, up to a fixed number. They stay accepted, so that old scripts run.
	for _, name := range []string{"i", "asc", "desc"} {
		if err := getCmd.PersistentFlags().MarkDeprecated(name, "the Smarter API ignores it"); err != nil {
			log.Fatalf("Error deprecating flag: %v", err)
		}
	}
}
