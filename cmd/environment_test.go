package cmd

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// setFlag passes a global flag for the test, and restores it afterwards.
func setFlag(t *testing.T, name, value string) {
	t.Helper()
	flag := RootCmd.PersistentFlags().Lookup(name)
	old := flag.Value.String()
	if err := flag.Value.Set(value); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		_ = flag.Value.Set(old)
		flag.Changed = false
	})
}

func TestValidateSmarterEnvironment(t *testing.T) {
	if err := validateSmarterEnvironment(); err != nil {
		t.Fatalf("validateSmarterEnvironment() with nothing set = %v", err)
	}
	valid := map[string]string{
		"SMARTER_ENVIRONMENT": "local",
		"SMARTER_BASE_URL":    "http://localhost:9357",
		"SMARTER_TIMEOUT":     "2.5",
		"SMARTER_MAX_RETRIES": "0",
		"SMARTER_LOG":         "debug",
	}
	for name, value := range valid {
		t.Setenv(name, value)
	}
	if err := validateSmarterEnvironment(); err != nil {
		t.Errorf("validateSmarterEnvironment() with valid values = %v", err)
	}

	for name, value := range map[string]string{
		"SMARTER_ENVIRONMENT": "staging",
		"SMARTER_BASE_URL":    "ftp://example.com",
		"SMARTER_TIMEOUT":     "-1",
		"SMARTER_MAX_RETRIES": "two",
		"SMARTER_LOG":         "trace",
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(name, value)
			err := validateSmarterEnvironment()
			if err == nil || !strings.Contains(err.Error(), name) {
				t.Errorf("validateSmarterEnvironment() = %v, want an error naming %s", err, name)
			}
		})
	}
}

func TestApplySmarterEnvironment(t *testing.T) {
	oldCfg, oldEnv := cfgFile, environment
	t.Cleanup(func() { cfgFile, environment = oldCfg, oldEnv })
	withViperValue(t, "environment", "alpha")
	withViperValue(t, "verbose", false)
	cfgFile = ""
	t.Setenv("SMARTER_CONFIG_FILE", "/tmp/smarter-config.yaml")
	t.Setenv("SMARTER_ENVIRONMENT", "local")
	t.Setenv("SMARTER_LOG", "info")

	applySmarterEnvironment()

	if cfgFile != "/tmp/smarter-config.yaml" {
		t.Errorf("cfgFile = %q, want SMARTER_CONFIG_FILE", cfgFile)
	}
	if viper.GetString("environment") != "local" || environment != "local" {
		t.Errorf("environment = %q, want SMARTER_ENVIRONMENT", viper.GetString("environment"))
	}
	if !viper.GetBool("verbose") {
		t.Error("SMARTER_LOG=info didn't turn on verbose output")
	}
}

func TestApplySmarterEnvironmentYieldsToFlags(t *testing.T) {
	oldEnv := environment
	t.Cleanup(func() { environment = oldEnv })
	withViperValue(t, "environment", "beta")
	setFlag(t, "environment", "beta")
	t.Setenv("SMARTER_ENVIRONMENT", "local")

	applySmarterEnvironment()

	if viper.GetString("environment") != "beta" {
		t.Errorf("environment = %q, want the --environment flag's", viper.GetString("environment"))
	}
}

func TestFetchAPIKeyReadsSmarterAPIKey(t *testing.T) {
	withViperValue(t, "environment", "alpha")
	withViperValue(t, "alpha.api_key", "config-key")
	t.Setenv("SMARTER_API_KEY", "env-key")

	if got := fetchAPIKey(); got != "env-key" {
		t.Errorf("fetchAPIKey() = %q, want SMARTER_API_KEY over the configuration file", got)
	}
	setFlag(t, "api_key", "flag-key")
	if got := fetchAPIKey(); got != "flag-key" {
		t.Errorf("fetchAPIKey() = %q, want --api_key over SMARTER_API_KEY", got)
	}
}

func TestGetAPIHostReadsSmarterBaseURL(t *testing.T) {
	withViperValue(t, "environment", "prod")
	t.Setenv("SMARTER_BASE_URL", "http://localhost:9357/")

	if got, err := getAPIHost(); err != nil || got != "http://localhost:9357" {
		t.Errorf("getAPIHost() = %q, %v; want SMARTER_BASE_URL without its trailing slash", got, err)
	}
	setFlag(t, "environment", "prod")
	if got, _ := getAPIHost(); got != "https://platform.smarter.sh" {
		t.Errorf("getAPIHost() = %q, want --environment over SMARTER_BASE_URL", got)
	}
}

// TestGetAPIHostIgnoresThePlatformsRootDomain guards against reading the
// platform's own SMARTER_ROOT_DOMAIN (its apex domain), which developers'
// shells export.
func TestGetAPIHostIgnoresThePlatformsRootDomain(t *testing.T) {
	withViperValue(t, "environment", "prod")
	withViperValue(t, "config.root_domain", "platform.smarter.sh")
	t.Setenv("SMARTER_ROOT_DOMAIN", "smarter.sh")

	if got, _ := getAPIHost(); got != "https://platform.smarter.sh" {
		t.Errorf("getAPIHost() = %q, want the configured root domain", got)
	}
}

func TestRequestTimeoutAndMaxRetries(t *testing.T) {
	withViperValue(t, "timeout", time.Minute)
	if got := requestTimeout(); got != time.Minute {
		t.Errorf("requestTimeout() = %v, want the configured timeout", got)
	}
	t.Setenv("SMARTER_TIMEOUT", "1.5")
	if got := requestTimeout(); got != 1500*time.Millisecond {
		t.Errorf("requestTimeout() = %v, want SMARTER_TIMEOUT", got)
	}
	setFlag(t, "timeout", "3s")
	withViperValue(t, "timeout", 3*time.Second)
	if got := requestTimeout(); got != 3*time.Second {
		t.Errorf("requestTimeout() = %v, want --timeout over SMARTER_TIMEOUT", got)
	}

	withViperValue(t, "max_retries", nil)
	if got := maxRetries(); got != DefaultMaxRetries {
		t.Errorf("maxRetries() = %d, want the default", got)
	}
	withViperValue(t, "max_retries", -1)
	if got := maxRetries(); got != DefaultMaxRetries {
		t.Errorf("maxRetries() = %d, want the default for a negative value", got)
	}
	withViperValue(t, "max_retries", 5)
	if got := maxRetries(); got != 5 {
		t.Errorf("maxRetries() = %d, want the configured value", got)
	}
	t.Setenv("SMARTER_MAX_RETRIES", "0")
	if got := maxRetries(); got != 0 {
		t.Errorf("maxRetries() = %d, want SMARTER_MAX_RETRIES", got)
	}
}

func TestIsRetryable(t *testing.T) {
	cases := []struct {
		status int
		body   string
		want   bool
	}{
		{http.StatusTooManyRequests, ``, true},
		{http.StatusBadGateway, ``, true},
		{http.StatusGatewayTimeout, ``, true},
		{http.StatusRequestTimeout, ``, true},
		{http.StatusServiceUnavailable, `<html>maintenance</html>`, true},
		{http.StatusServiceUnavailable, `{"error":{"errorClass":"SAMBrokerErrorNotReady"}}`, false},
		{http.StatusInternalServerError, ``, false},
		{http.StatusNotFound, ``, false},
	}
	for _, tc := range cases {
		if got := isRetryable(tc.status, []byte(tc.body)); got != tc.want {
			t.Errorf("isRetryable(%d, %q) = %v, want %v", tc.status, tc.body, got, tc.want)
		}
	}
}

func TestRetryDelay(t *testing.T) {
	header := http.Header{}
	header.Set("Retry-After", "3")
	if got := retryDelay(0, header); got != 3*time.Second {
		t.Errorf("retryDelay() = %v, want Retry-After", got)
	}
	header.Set("Retry-After", "600")
	if got := retryDelay(0, header); got > 500*time.Millisecond {
		t.Errorf("retryDelay() = %v, want backoff when Retry-After is over a minute", got)
	}
	if got := retryDelay(10, nil); got > 8*time.Second || got < 6*time.Second {
		t.Errorf("retryDelay(10) = %v, want backoff capped at 8s, less jitter", got)
	}
}

func TestAPIRequestRetries(t *testing.T) {
	var delays []time.Duration
	old := sleep
	sleep = func(d time.Duration) { delays = append(delays, d) }
	t.Cleanup(func() { sleep = old })

	attempts := 0
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	withViperValue(t, "max_retries", 2)

	if _, err := APIRequest("whoami", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v, want success on the third attempt", err)
	}
	if attempts != 3 || len(delays) != 2 {
		t.Errorf("attempts = %d, delays = %v; want 3 attempts and 2 waits", attempts, delays)
	}

	attempts = -10
	withViperValue(t, "max_retries", 1)
	if _, err := APIRequest("whoami", map[string]string{}); err == nil || !strings.Contains(err.Error(), "503") {
		t.Errorf("APIRequest() error = %v, want the last 503 after the retries", err)
	}
}

func TestAPIRequestDoesNotRetryAJournaledError(t *testing.T) {
	attempts := 0
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"errorClass":"SAMBrokerErrorNotReady","description":"not ready"}}`))
	})

	if _, err := APIRequest("whoami", map[string]string{}); err == nil {
		t.Fatal("APIRequest() error = nil")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1: the platform's own 503 isn't retried", attempts)
	}
}

func TestValidateTabular(t *testing.T) {
	withViperValue(t, "output_format", "tabular")
	if err := validateTabular(whoamiCmd); err == nil {
		t.Error("validateTabular(whoami) = nil, want an error")
	}
	getCmd, _, err := RootCmd.Find([]string{"get"})
	if err != nil || getCmd == RootCmd {
		t.Skip("the get verb isn't registered in this package's tests")
	}
	if err := validateTabular(getCmd); err != nil {
		t.Errorf("validateTabular(get) = %v", err)
	}
}
