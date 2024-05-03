package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetAPIResponse(slug string) (string, error) {
	apiHost := os.Getenv("API_HOST")
	root_url := apiHost + "/api/v0/cli"
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
