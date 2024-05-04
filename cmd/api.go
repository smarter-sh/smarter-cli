package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

const (
	RootDomain  = "platform.smarter.sh"
	ApiBasePath = "/api/v0/cli"
)

func getAPIHost() string {
	environment := viper.GetString("config.environment")
	if viper.IsSet("environment") {
		environment = viper.GetString("environment")
	}

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

func GetAPIResponse(slug string) (map[string]interface{}, error) {

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
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return result, nil
}
