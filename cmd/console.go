/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>

Shared functions for generated formatted console output in any of json, yaml,
or tabular formats.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

type Title struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Item map[string]interface{}

type InnerData struct {
	Titles []Title `json:"titles"`
	Items  []Item  `json:"items"`
}

type Data struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Metadata   struct {
		Count int `json:"count"`
	} `json:"metadata"`
	Kwargs map[string]interface{} `json:"kwargs"`
	Data   InnerData              `json:"data"`
}

type Body struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

func TableOutput(bodyJson []byte) error {
	var body Body

	if viper.GetBool("verbose") {
		log.Printf("TableOutput()")
	}

	if err := json.Unmarshal(bodyJson, &body); err != nil {
		return fmt.Errorf("failed parsing JSON: %w", err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	// print column titles
	titles := body.Data.Data.Titles
	titleNames := make([]string, len(titles))
	for i, title := range titles {
		titleNames[i] = title.Name
	}
	fmt.Fprintln(w, strings.Join(titleNames, "\t"))

	// print dashed line
	dashes := make([]string, len(titles))
	for i, title := range titles {
		dashes[i] = strings.Repeat("-", len(title.Name))
	}
	fmt.Fprintln(w, strings.Join(dashes, "\t"))

	// print data rows
	for _, item := range body.Data.Data.Items {
		values := make([]string, len(titles))
		for i, title := range titles {
			value, ok := item[title.Name]
			if !ok || value == nil {
				values[i] = ""
				continue
			}
			if title.Type == "DateTimeField" {
				t, err := time.Parse(time.RFC3339, value.(string))
				if err != nil {
					return fmt.Errorf("failed parsing date: %w", err)
				}
				values[i] = t.Format("2006-Jan-02 15:04")
			} else {
				values[i] = fmt.Sprint(value)
			}
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	w.Flush()
	return nil
}

func JsonOutput(bodyJson []byte) error {
	var jsonData map[string]interface{}

	if err := json.Unmarshal(bodyJson, &jsonData); err != nil {
		return fmt.Errorf("failed unmarshalling JSON: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed formatting JSON: %w", err)
	}

	fmt.Println(string(prettyJSON))
	return nil
}

func YamlOutput(bodyJson []byte) error {
	var jsonData map[string]interface{}

	if viper.GetBool("verbose") {
		log.Printf("YamlOutput()")
	}
	if err := json.Unmarshal(bodyJson, &jsonData); err != nil {
		return fmt.Errorf("failed unmarshalling JSON: %w", err)
	}

	if data, ok := jsonData["data"]; ok {
		newData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed marshalling data: %w", err)
		}
		bodyJson = newData
	}

	bodyYaml, err := yaml.JSONToYAML(bodyJson)
	if err != nil {
		return fmt.Errorf("failed converting JSON to YAML: %w", err)
	}
	fmt.Println(string(bodyYaml))
	return nil
}

func ConsoleOutput(bodyJson []byte) error {
	outputFormat := viper.GetString("output_format")

	switch outputFormat {
	case "json":
		return JsonOutput(bodyJson)
	case "yaml":
		return YamlOutput(bodyJson)
	case "tabular":
		return TableOutput(bodyJson)
	default:
		return JsonOutput(bodyJson)
	}
}

func ErrorOutput(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	os.Exit(1)
}
