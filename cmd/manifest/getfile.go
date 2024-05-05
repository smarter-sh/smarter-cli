/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

func getLocalFileContents(kind string) (string, error) {
	filePath := fmt.Sprintf("./data/manifests/%s.yaml", kind)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func formatOutput(contents string) (string, error) {
	jsonFlagValue := viper.GetBool("json")
	yamlFlagValue := viper.GetBool("yaml")
	outputFormat := viper.GetString("config.output_format")

	// if neither flag is set, then first check output_format.
	// if output_format is set, then use that. it is either "json" or "yaml"
	// otherwise, default to yaml
	if !jsonFlagValue && !yamlFlagValue {
		if outputFormat == "json" {
			jsonFlagValue = true
		} else {
			yamlFlagValue = true
		}
	}
	if yamlFlagValue {
		return contents, nil
	}

	// IF jsonFlagValue is true, we need to convert the YAML to JSON
	if jsonFlagValue {
		bodyYaml := []byte(contents)
		bodyJson, err := yaml.YAMLToJSON(bodyYaml)
		if err != nil {
			return "", err
		} else {
			return string(bodyJson), nil
		}
	}
	return "", nil
}

func GetAndPrintManifest(url string, kind string) (string, error) {
	environment := viper.GetString("environment")

	// If we are in a local environment, we can just read the file from the
	// repository
	if environment == "local" {
		localContents, err := getLocalFileContents(kind)
		if err != nil {
			return "", err
		}
		return formatOutput(localContents)
	}

	// Otherwise, we need to make an HTTP request to the environment
	// cdn to get the deployed file contents. This is a 2-step process:
	// 1. Get the URL of the file from the smarter API
	// 2. Get the contents of the file from the environment CDN
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// example url: https://cdn.platform.smarter.sh/cli/example-manifests/plugin.yaml
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	remoteContents := string(bodyBytes)
	return formatOutput(remoteContents)

}
