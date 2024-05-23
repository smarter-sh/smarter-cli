package cmd

import (
	"bytes"
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

func TableOutput(bodyJson []byte) {
	var body Body

	err := json.Unmarshal(bodyJson, &body)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
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
			if title.Type == "DateTimeField" {
				t, err := time.Parse(time.RFC3339, item[title.Name].(string))
				if err != nil {
					log.Fatalf("Error parsing date: %v", err)
				}
				values[i] = t.Format("2006-Jan-02 15:04")
			} else {
				values[i] = fmt.Sprint(item[title.Name])
			}
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
	outputFormat := viper.GetString("output_format")

	switch {
	case outputFormat == "json":
		JsonOutput(bodyJson)
	case outputFormat == "yaml":
		YamlOutput(bodyJson)
	case outputFormat == "tabular":
		TableOutput(bodyJson)
	default:
		JsonOutput(bodyJson)
	}
}

func ErrorOutput(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	os.Exit(1)
}
