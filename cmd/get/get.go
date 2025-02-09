/*
Copyright © 2024 Lawrence McDaniel <Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	i := viper.GetInt("i")
	asc := viper.GetBool("asc")
	desc := viper.GetBool("desc")
	common_kwargs := map[string]string{
		"i":    strconv.Itoa(i),
		"asc":  strconv.FormatBool(asc),
		"desc": strconv.FormatBool(desc),
	}
	for key, value := range common_kwargs {
		kwargs[key] = value
	}

	// en route to /api/v1/cli/get/<str:kind>/
	return cmd.APIRequest("get/"+kind+"/", kwargs)

}

func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "tabular")
	}
	cmd.ConsoleOutput(bodyJson)
}

func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
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
	if err := viper.BindPFlag("i", getCmd.PersistentFlags().Lookup("i")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	getCmd.PersistentFlags().Bool("asc", false, "Sort results in ascending order")
	if err := viper.BindPFlag("asc", getCmd.PersistentFlags().Lookup("asc")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	getCmd.PersistentFlags().Bool("desc", false, "Sort results in descending order")
	if err := viper.BindPFlag("desc", getCmd.PersistentFlags().Lookup("desc")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
