package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetAPIResponse(slug string) (string, error) {
	apiHost := os.Getenv("API_HOST")
	root_url := apiHost + "/api/v0/cli"
	url := fmt.Sprintf("%s/%s", root_url, slug)

	resp, err := http.Get(url)
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
