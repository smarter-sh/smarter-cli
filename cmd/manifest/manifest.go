/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"encoding/json"
	"log"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {

	return cmd.APIRequest("manifest/"+kind+"/", kwargs)

}
func ConsoleOutput(bodyJson []byte) {
	jsonFlagValue := viper.GetBool("json")
	yamlFlagValue := viper.GetBool("yaml")
	if !jsonFlagValue && !yamlFlagValue {
		viper.Set("yaml", true)
	}
	var data map[string]interface{}
	err := json.Unmarshal(bodyJson, &data)
	if err != nil {
		log.Fatalf("Error occurred during unmarshalling. %v", err)
	}

	value, ok := data["data"]
	if ok {
		bodyJson, err = json.Marshal(value)
		if err != nil {
			log.Fatalf("Error occurred during marshalling. %v", err)
		}
	}
	cmd.ConsoleOutput(bodyJson)
}
func ErrorOutput(err error) {
	cmd.ErrorOutput(err)
}

var manifestCmd = &cobra.Command{
	Use:   "manifest <kind> [flags]",
	Short: "Generate an example manifest for the resource kind",
	Long: `Generate an example manifest for the resource kind. For example:

	smarter manifest <kind> [flags] > my-plugin.yaml

This will generate an example manifest for the specified kind of resource and write it to my-plugin.yaml in the current working directory.`,
}

func init() {
	cmd.RootCmd.AddCommand(manifestCmd)
}
