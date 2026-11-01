/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	ApiBasePath = "/api/v1/cli/"

	// ClientName is sent as the X-Smarter-Client header, which identifies the
	// cli in the platform's logs and journal, as the SDKs identify themselves.
	ClientName = "smarter-cli"

	// DefaultTimeout bounds every request. Deploying an LLMHost returns when
	// the cluster accepts its resources, and a prompt waits for the LLM, so it
	// is generous; --timeout changes it.
	DefaultTimeout = 2 * time.Minute
)

// fetchAPIKey returns the api key: --api_key when it was passed, else
// SMARTER_API_KEY, else the configured environment's <environment>.api_key,
// else any other value of the api_key key. See doc/adr/002-configuration.md.
func fetchAPIKey() string {
	if flag := RootCmd.PersistentFlags().Lookup("api_key"); flag != nil && flag.Changed {
		return flag.Value.String()
	}
	if apiKey := smarterEnv("SMARTER_API_KEY"); apiKey != "" {
		return apiKey
	}
	environment := viper.GetString("environment")
	if apiKey := viper.GetString(fmt.Sprintf("%s.api_key", environment)); apiKey != "" {
		return apiKey
	}
	return viper.GetString("api_key")
}

func verifyApiKey() error {
	if fetchAPIKey() == "" {
		return fmt.Errorf(`no api_key is configured for the %q environment. Set it with:

smarter configure --api_key <api_key>`, viper.GetString("environment"))
	}
	return nil
}

// maskSecret returns secret with all but its last four characters hidden, for
// logs and terminal echo. Short secrets are hidden entirely.
func maskSecret(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return "****" + secret[len(secret)-4:]
}

// apiHostOverride, when non-empty, is returned by getAPIHost instead of the
// environment-derived host. It exists so tests can point APIRequest at an
// httptest.Server without a live network dependency.
var apiHostOverride string

// SetAPIHostOverride overrides the host used by APIRequest. Pass "" to clear
// the override and resume deriving the host from the environment/root_domain
// config. For use by tests only.
func SetAPIHostOverride(host string) {
	apiHostOverride = host
}

// getAPIHost returns SMARTER_BASE_URL, or the platform url of the configured
// environment: see doc/adr/002-configuration.md. The SDKs derive it the same way.
func getAPIHost() (string, error) {
	if apiHostOverride != "" {
		return apiHostOverride, nil
	}
	if baseURL := envBaseURL(); baseURL != "" {
		return baseURL, nil
	}
	environment := viper.GetString("environment")
	rootDomain := viper.GetString("config.root_domain")
	if viper.GetBool("verbose") {
		log.Printf("Environment: %s", environment)
	}
	switch environment {
	case "local":
		return "http://localhost:9357", nil
	case "alpha", "beta", "next":
		return fmt.Sprintf("https://%s.%s", environment, rootDomain), nil
	case "prod":
		return fmt.Sprintf("https://%s", rootDomain), nil
	default:
		return "", fmt.Errorf("invalid environment %q. Valid values: %s", environment, strings.Join(validEnvironments, ", "))
	}
}

// maxRetryAfter is the longest Retry-After that is waited for.
const maxRetryAfter = 60 * time.Second

// retryableStatuses are the statuses of a failed request that may succeed when
// it is sent again, as in the SDKs. A 503 is decided by isRetryable.
var retryableStatuses = map[int]bool{
	http.StatusRequestTimeout:  true,
	http.StatusTooManyRequests: true,
	http.StatusBadGateway:      true,
	http.StatusGatewayTimeout:  true,
}

// isRetryable reports whether a failed response may succeed when the request
// is sent again. A 503 is retried only when it isn't a journaled response: the
// platform itself answers 503 for a resource that isn't ready, which a retry
// won't change, whereas a 503 from a load balancer means the platform is
// briefly unavailable.
func isRetryable(status int, body []byte) bool {
	if status == http.StatusServiceUnavailable {
		var document map[string]interface{}
		if json.Unmarshal(body, &document) == nil {
			if _, journaled := document["error"].(map[string]interface{}); journaled {
				return false
			}
		}
		return true
	}
	return retryableStatuses[status]
}

// retryDelay returns how long to wait before retry number attempt (0 first):
// the response's Retry-After, up to a minute, or exponential backoff from half
// a second, up to eight, with jitter. The SDKs wait the same way.
func retryDelay(attempt int, header http.Header) time.Duration {
	if header != nil {
		if value := strings.TrimSpace(header.Get("Retry-After")); value != "" {
			if seconds, err := strconv.ParseFloat(value, 64); err == nil && seconds >= 0 {
				if delay := time.Duration(seconds * float64(time.Second)); delay <= maxRetryAfter {
					return delay
				}
			}
		}
	}
	backoff := min(0.5*float64(int(1)<<min(attempt, 10)), 8.0)
	jitter := 1 - 0.25*rand.Float64() //nolint:gosec // G404: jitter, not security
	return time.Duration(backoff * jitter * float64(time.Second))
}

// sleep waits between retries; tests replace it.
var sleep = time.Sleep

// APIRequest POSTs to the platform's cli api, /api/v1/cli/<slug>/, with kwargs
// as the query string and the optional body, and returns the body of a 200
// response. Any other outcome is returned as an error, for the command to
// report with ErrorOutput: see doc/adr/006-error-handling.md and 008-http-transport.md.
func APIRequest(slug string, kwargs map[string]string, fileContents ...string) ([]byte, error) {
	params := url.Values{}
	for key, value := range kwargs {
		params.Add(key, value)
	}
	return APIRequestWithParams(slug, params, fileContents...)
}

// APIRequestWithParams is APIRequest for a query string whose parameters may
// repeat, e.g. get's annotation selector, one requirement per parameter.
func APIRequestWithParams(slug string, params url.Values, fileContents ...string) ([]byte, error) {
	verbose := viper.GetBool("verbose")

	if err := verifyApiKey(); err != nil {
		return nil, err
	}
	apiKey := fetchAPIKey()
	apiHost, err := getAPIHost()
	if err != nil {
		return nil, err
	}
	urlPath := path.Clean("/" + ApiBasePath + slug)
	urlOrig := strings.ToLower(apiHost + urlPath)
	if !strings.HasSuffix(urlOrig, "/") {
		urlOrig += "/"
	}

	var textData string
	if len(fileContents) > 0 {
		textData = fileContents[0]
	}

	reqURL := urlOrig
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}
	if verbose {
		log.Printf("HTTP Request: POST %s", reqURL)
	}
	client := &http.Client{Timeout: requestTimeout()}
	retries := maxRetries()
	var resp *http.Response
	var bodyBytes []byte
	for attempt := 0; ; attempt++ {
		resp, bodyBytes, err = sendRequest(client, reqURL, textData, apiKey, verbose)
		if err != nil {
			if attempt < retries {
				delay := retryDelay(attempt, nil)
				if verbose {
					log.Printf("Retrying in %s after: %v", delay, err)
				}
				sleep(delay)
				continue
			}
			return nil, fmt.Errorf("request to %s failed: %w", urlOrig, err)
		}
		if resp.StatusCode != http.StatusOK && attempt < retries && isRetryable(resp.StatusCode, bodyBytes) {
			delay := retryDelay(attempt, resp.Header)
			if verbose {
				log.Printf("Retrying in %s after an http response %d", delay, resp.StatusCode)
			}
			sleep(delay)
			continue
		}
		break
	}

	if resp.StatusCode != http.StatusOK {
		return nil, responseError(resp.StatusCode, urlOrig, bodyBytes, verbose)
	}
	return bodyBytes, nil
}

// sendRequest sends one attempt of a request, and reads its response.
func sendRequest(client *http.Client, reqURL, body, apiKey string, verbose bool) (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed creating the request: %w", err)
	}

	// see https://jazzband.github.io/django-rest-knox/auth/
	req.Header.Set("User-Agent", ClientName+"/"+Version)
	req.Header.Set("X-Smarter-Client", ClientName)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	if verbose {
		if reqDump, err := httputil.DumpRequestOut(req, true); err == nil {
			log.Printf("Request: %s", strings.ReplaceAll(string(reqDump), apiKey, maskSecret(apiKey)))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if verbose {
		if respDump, err := httputil.DumpResponse(resp, false); err == nil {
			log.Printf("Response: %s", respDump)
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading the response: %w", err)
	}
	return resp, bodyBytes, nil
}

// responseError describes a non-200 response with the platform's description
// and context, and with --verbose its stack trace. The platform's journal
// nests them in an "error" object (SAMJournalApiErrorResponse in its OpenAPI
// document); older platforms put them at the top level, and authentication
// failures come from Django REST Framework as "detail". A body that isn't
// json (e.g. a proxy's html error page) is quoted, shortened.
func responseError(status int, url string, body []byte, verbose bool) error {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return fmt.Errorf("received an http response %d from %s: %s", status, url, snippet)
	}
	if nested, ok := result["error"].(map[string]interface{}); ok {
		result = nested
	}
	description, ok := result["description"]
	if !ok {
		if description, ok = result["detail"]; !ok {
			description = "unknown error"
		}
	}
	message := fmt.Sprint(description)
	if context, ok := result["context"]; ok && context != "" {
		message += fmt.Sprintf(" (%v)", context)
	}
	if verbose {
		if errorClass, ok := result["errorClass"]; ok {
			message += fmt.Sprintf("\nError class: %v", errorClass)
		}
		stackTrace, ok := result["stacktrace"]
		if !ok {
			stackTrace = "unknown stack trace"
		}
		message += fmt.Sprintf("\nStack trace: %v", stackTrace)
	}
	return fmt.Errorf("received an http response %d from %s: %s", status, url, message)
}
