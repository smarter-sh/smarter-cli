package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/spf13/viper"
)

const (
	RootDomain  = "platform.smarter.sh"
	ApiBasePath = "/api/v1/cli/"
)

func fetchAPIKey() string {
	apiKey := viper.Get("config.api_key").(string)
	if apiKey == "" {
		apiKey = viper.Get("api_key").(string)
	}
	return apiKey
}

func verifyApiKey() error {
	apiKey := fetchAPIKey()
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
	apiKey := fetchAPIKey()
	apiHost := getAPIHost()
	url_path := path.Clean("/" + ApiBasePath + slug)
	url := apiHost + url_path
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	var textData string
	if len(fileContents) > 0 {
		textData = fileContents[0]
	} else {
		textData = ""
	}

	fmt.Println("HTTP Request:", url, textData)
	req, err := http.NewRequest("POST", url, strings.NewReader(textData))
	if err != nil {
		ErrorOutput(err)
	}

	// Set headers from kwargs
	for key, value := range kwargs {
		req.Header.Set(key, value)
	}

	// see https://jazzband.github.io/django-rest-knox/auth/
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("HTTP Response Status:", resp.StatusCode)
	if err != nil {
		ErrorOutput(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		ErrorOutput(fmt.Errorf("HTTP Response Status: %d", resp.StatusCode))
	}

	return bodyBytes, nil
}
