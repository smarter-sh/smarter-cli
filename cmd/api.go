package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

const (
	RootDomain  = "platform.smarter.sh"
	ApiBasePath = "/api/v0/cli"
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

func GetAPIResponse(slug string) ([]byte, error) {

	checkApiKey := verifyApiKey()
	if checkApiKey != nil {
		return []byte{}, checkApiKey
	}
	apiHost := getAPIHost()
	root_url := apiHost + ApiBasePath
	url := fmt.Sprintf("%s/%s/", root_url, slug)
	jsonData := []byte(`{}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return json.Marshal(result)
}
