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
)

func getLocalFileContents(kind string) (string, error) {
	filePath := fmt.Sprintf("./data/manifests/%s.yaml", kind)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func GetAndPrintYAMLResponse(url string, kind string) (string, error) {
	var environment string
	if viper.IsSet("environment") {
		environment = viper.GetString("environment")
	} else {
		environment = "prod"
	}

	// If we are in a local environment, we can just read the file from the
	// repository
	if environment == "local" {
		return getLocalFileContents(kind)
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
	return string(bodyBytes), nil

}
