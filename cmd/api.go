package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

const (
	RootDomain  = "platform.smarter.sh"
	ApiBasePath = "/api/v1/cli/"
)

func verifyApiKey() error {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		errMsg := `api_key is missing. Please set the API key using the command:
smarter configure --api_key <api_key string>
Contact support@querium.com if you need help finding your API key`
		return errors.New(errMsg)
	}
	return nil
}

func getAPIHost() string {
	environment := viper.GetString("environment")
	baseURL := fmt.Sprintf("https://%%s.%s", RootDomain)

	switch environment {
	case "local":
		return "http://localhost:8000"
	case "alpha":
		return fmt.Sprintf(baseURL, "alpha")
	case "beta":
		return fmt.Sprintf(baseURL, "beta")
	case "next":
		return fmt.Sprintf(baseURL, "next")
	case "prod":
		return fmt.Sprintf("https://%s", RootDomain)
	default:
		panic(fmt.Sprintf("invalid environment: %s", environment))
	}
}

func APIRequest(slug string, kwargs map[string]string, fileContents ...string) ([]byte, error) {

	checkApiKey := verifyApiKey()
	if checkApiKey != nil {
		return []byte{}, checkApiKey
	}
	apiHost := getAPIHost()
	root_url := apiHost + ApiBasePath
	url := fmt.Sprintf("%s/%s/", root_url, slug)

	var textData string
	if len(fileContents) > 0 {
		textData = fileContents[0]
	} else {
		textData = ""
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(textData))
	if err != nil {
		panic(err)
	}
	// Set headers from kwargs
	for key, value := range kwargs {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return bodyBytes, nil
}
