package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGetAPIHost(t *testing.T) {
	testCases := []struct {
		environment string
		expected    string
	}{
		{"local", "http://localhost:8000"},
		{"alpha", "https://alpha.platform.smarter.sh"},
		{"beta", "https://beta.platform.smarter.sh"},
		{"next", "https://next.platform.smarter.sh"},
		{"prod", "https://platform.smarter.sh"},
	}

	for _, tc := range testCases {
		viper.Set("config.environment", tc.environment)
		result := getAPIHost()
		if result != tc.expected {
			t.Errorf("For environment %s, expected %s but got %s", tc.environment, tc.expected, result)
		}
	}
}

func TestGetAPIResponseResponse(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test that the method is POST
		if req.Method != http.MethodPost {
			t.Errorf("Expected 'POST' request, got '%s'", req.Method)
		}
		// Test that the URL is correct
		if req.URL.String() != "/test/" {
			t.Errorf("Expected request to '/test/', got '%s'", req.URL.String())
		}
		// Send response to be tested
		_, err := rw.Write([]byte(`OK`))
		if err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	// Now we can run our function with the mock server's URL
	_, err := APIRequest("test", map[string]string{}, server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

type TestCase struct {
	Name           string            `json:"name"`
	InputURL       string            `json:"inputURL"`
	InputParams    map[string]string `json:"inputParams"`
	ExpectedStatus int               `json:"expectedStatus"`
	ExpectedBody   string            `json:"expectedBody"`
}

func TestGetAPIResponseResponseWithRealRequests(t *testing.T) {
	file, err := os.Open("testdata.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var testCases []TestCase
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&testCases); err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bodyBytes, err := APIRequest(tc.InputURL, tc.InputParams)
			if err != nil {
				if tc.ExpectedStatus != 0 {
					t.Errorf("Expected no error, got %v", err)
				}
				return
			}

			var result map[string]interface{}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				t.Errorf("Expected a valid JSON dictionary, but got an error: %v", err)
			}
		})
	}
}
