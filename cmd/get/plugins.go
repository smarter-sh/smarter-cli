/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateClass(class string) bool {
	allowedClasses := []string{"static", "sql", "api"}

	for _, allowedClass := range allowedClasses {
		if class == allowedClass {
			return true
		}
	}

	log.Fatalf("Invalid class '%s'. Allowed classes are: %v", class, allowedClasses)
	return false
}

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Retrieve a list of Plugins",
	Long: `Retrieves a list of Plugins:

smarter get plugins [flags]


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("name")

		plugin_class := viper.GetString("class")
		if plugin_class != "" {
			validateClass(plugin_class)
		}

		kwargs := map[string]string{
			"name":  name,
			"class": plugin_class,
		}

		bodyJson, err := APIRequest("Plugin", kwargs)
		if err != nil {
			ErrorOutput(err)
		} else {
			ConsoleOutput(bodyJson)
		}

	},
}

func init() {
	GetCmd.AddCommand(pluginsCmd)

	pluginsCmd.Flags().StringP("name", "n", "", "Name of the plugin")
	if err := viper.BindPFlag("name", pluginsCmd.Flags().Lookup("name")); err != nil {
		log.Fatalf("Error binding flag 'name': %v", err)
	}

	pluginsCmd.Flags().StringP("class", "c", "", "Plugin class: static, sql, api")
	if err := viper.BindPFlag("class", pluginsCmd.Flags().Lookup("class")); err != nil {
		log.Fatalf("Error binding flag 'class': %v", err)
	}
}
