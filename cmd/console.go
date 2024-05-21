package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

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

	fmt.Print(bodyJson)

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

func JsonOutput(bodyJson []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, bodyJson, "", "\t")
	if err != nil {
		fmt.Println("JSON parse error: ", err)
		return
	}
	fmt.Println("default")
	fmt.Println(prettyJSON.String())
}

func YamlOutput(bodyJson []byte) {
	bodyYaml, err := yaml.JSONToYAML(bodyJson)
	if err != nil {
		ErrorOutput(err)
	} else {
		fmt.Println(string(bodyYaml))
	}
}

func ConsoleOutput(bodyJson []byte) {
	jsonFlagValue := viper.GetBool("json")
	yamlFlagValue := viper.GetBool("yaml")

	switch {
	case jsonFlagValue:
		JsonOutput(bodyJson)
	case yamlFlagValue:
		YamlOutput(bodyJson)
	default:
		JsonOutput(bodyJson)
	}
}

func ErrorOutput(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	os.Exit(1)
}
