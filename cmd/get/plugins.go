/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pluginsCmd represents the plugins command
var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Retrieve a list of Plugins",
	Long: `Retrieves a list of Plugins:

smarter get plugins [flags]


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the class value
		class := viper.GetString("class")

		// Define allowed classes
		allowedClasses := []string{"static", "sql", "api"}

		// Check if the class is allowed
		isValidClass := false
		for _, allowedClass := range allowedClasses {
			if class == allowedClass {
				isValidClass = true
				break
			}
		}

		// If the class is not allowed, log an error and exit
		if !isValidClass {
			log.Fatalf("Invalid class '%s'. Allowed classes are: %v", class, allowedClasses)
		}

		name := viper.GetString("name")
		plugin_class := viper.GetString("class")

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
