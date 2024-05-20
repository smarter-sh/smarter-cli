/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pluginsCmd represents the plugins command
var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Retrieve a list of Plugins or a manifest for a specific Plugin by name",
	Long: `Retrieve a list of Plugins,
	or a manifest for a specific Plugin:

smarter get plugins --name --json --yaml --csv --xml -n 10 --asc --desc --class


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFlagValue := viper.GetBool("json")
		yamlFlagValue := viper.GetBool("yaml")

		name := viper.GetString("name")
		plugin_class := viper.GetString("class")
		n := viper.GetInt("n")

		kwargs := map[string]string{
			"name":  name,
			"class": plugin_class,
			"n":     strconv.Itoa(n),
		}

		bodyJson, err := GetAPIResponse("plugins", kwargs)
		if err != nil {
			panic(err)
		} else {
			ConsoleOutput(bodyJson, jsonFlagValue, yamlFlagValue)
		}

	},
}

func init() {
	GetCmd.AddCommand(pluginsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pluginsCmd.Flags().String("name", "", "Name of the plugin")
	pluginsCmd.Flags().String("class", "", "Plugin class: static, sql, api")
	pluginsCmd.Flags().Int("n", 10, "Number of plugins to retrieve")

	if err := viper.BindPFlag("name", pluginsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

	if err := viper.BindPFlag("class", pluginsCmd.Flags().Lookup("class")); err != nil {
		log.Fatalf("Error binding flag 'class': %v", err)
	}

	if err := viper.BindPFlag("n", pluginsCmd.Flags().Lookup("n")); err != nil {
		log.Fatalf("Error binding flag 'n': %v", err)
	}
}
