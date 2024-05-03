/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pluginsCmd represents the plugins command
var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Retrieve a list of Plugins or a manifest for a specific Plugin by name",
	Long: `Retrieve a list of Plugins,
	or a manifest for a specific Plugin:

smarter get plugins --name --json --yaml --csv --xml -n 10 --asc --desc


The Smarter API will return a list of Plugins in the specified format,
or a manifest for a specific Plugin.`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := GetAPI("apply")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response:", string(body))
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
}
