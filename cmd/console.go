package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

var jsonFlagValue bool
var yamlFlagValue bool

func init() {
	jsonFlagValue = viper.GetBool("json")
	yamlFlagValue = viper.GetBool("yaml")
}

type Column struct {
	Title string      `json:"title"`
	Type  interface{} `json:"type"`
}

type Table struct {
	Titles []Column                 `json:"titles"`
	Data   []map[string]interface{} `json:"data"`
}

func TableOutput(bodyJson []byte) {
	var table Table

	err := json.Unmarshal([]byte(bodyJson), &table)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	// print column titles
	titles := make([]string, len(table.Titles))
	for i, title := range table.Titles {
		titles[i] = title.Title
	}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	// print data rows
	for _, row := range table.Data {
		values := make([]string, len(table.Titles))
		for i, title := range table.Titles {
			values[i] = fmt.Sprint(row[title.Title])
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	w.Flush()

}

func ConsoleOutput(bodyJson []byte) {
	switch {
	case jsonFlagValue:
		fmt.Println(string(bodyJson))
	case yamlFlagValue:
		bodyYaml, err := yaml.JSONToYAML(bodyJson)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(string(bodyYaml))
		}
	default:
		TableOutput(bodyJson)
	}
}
