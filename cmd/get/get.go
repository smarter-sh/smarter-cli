/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/QueriumCorp/smarter-cli/cmd"

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
	return cmd.APIRequest("get/"+kind+"/", kwargs)

}

func ConsoleOutput(bodyJson []byte) {
	cmd.ConsoleOutput(bodyJson)
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get <kind> --json --yaml -n <10> --asc --desc",
	Short: "Generate a list of Smarter resources",
	Long: `Generate a list of Smarter resources:

smarter get <kind> --json --yaml -n <10> --asc --desc

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.`,
}

func init() {
	cmd.RootCmd.AddCommand(GetCmd)

	GetCmd.PersistentFlags().Int("i", 10, "Number of resources to retrieve")
	if err := viper.BindPFlag("i", GetCmd.PersistentFlags().Lookup("i")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	GetCmd.PersistentFlags().Bool("asc", false, "Sort results in ascending order")
	if err := viper.BindPFlag("asc", GetCmd.PersistentFlags().Lookup("asc")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	GetCmd.PersistentFlags().Bool("desc", false, "Sort results in descending order")
	if err := viper.BindPFlag("desc", GetCmd.PersistentFlags().Lookup("desc")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
