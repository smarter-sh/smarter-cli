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

func getYamlFileContents(kind string) (string, error) {
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

	if environment == "local" {
		return getYamlFileContents(kind)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

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
