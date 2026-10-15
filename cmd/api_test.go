package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFetchAPIKey(t *testing.T) {
	t.Run("prefers the environment-specific key", func(t *testing.T) {
		withViperValue(t, "environment", "alpha")
		withViperValue(t, "alpha.api_key", "alpha-key")
		withViperValue(t, "api_key", "generic-key")

		if got := fetchAPIKey(); got != "alpha-key" {
			t.Errorf("fetchAPIKey() = %q, want %q", got, "alpha-key")
		}
	})

	t.Run("falls back to the generic api_key", func(t *testing.T) {
		withViperValue(t, "environment", "beta")
		withViperValue(t, "beta.api_key", "")
		withViperValue(t, "api_key", "generic-key")

		if got := fetchAPIKey(); got != "generic-key" {
			t.Errorf("fetchAPIKey() = %q, want %q", got, "generic-key")
		}
	})
}

func TestFetchAPIKeyPrefersThePassedFlag(t *testing.T) {
	withViperValue(t, "environment", "alpha")
	withViperValue(t, "alpha.api_key", "alpha-key")
	flag := RootCmd.PersistentFlags().Lookup("api_key")
	if err := flag.Value.Set("flag-key"); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		_ = flag.Value.Set("")
		flag.Changed = false
	})

	if got := fetchAPIKey(); got != "flag-key" {
		t.Errorf("fetchAPIKey() = %q, want --api_key's value", got)
	}
}

func TestVerifyApiKey(t *testing.T) {
	t.Run("no error when a key is configured", func(t *testing.T) {
		withViperValue(t, "environment", "alpha")
		withViperValue(t, "alpha.api_key", "some-key")

		if err := verifyApiKey(); err != nil {
			t.Errorf("verifyApiKey() = %v, want nil", err)
		}
	})

	t.Run("returns an error, does not exit, when no key is configured", func(t *testing.T) {
		withViperValue(t, "environment", "alpha")
		withViperValue(t, "alpha.api_key", "")
		withViperValue(t, "api_key", "")

		if err := verifyApiKey(); err == nil {
			t.Error("verifyApiKey() = nil, want an error")
		}
	})
}

func TestGetAPIHost(t *testing.T) {
	t.Run("override takes precedence", func(t *testing.T) {
		SetAPIHostOverride("http://example.invalid")
		defer SetAPIHostOverride("")

		if got, err := getAPIHost(); err != nil || got != "http://example.invalid" {
			t.Errorf("getAPIHost() = %q, want override", got)
		}
	})

	cases := []struct {
		environment string
		want        string
	}{
		{"local", "http://localhost:9357"},
		{"alpha", "https://alpha.platform.smarter.sh"},
		{"beta", "https://beta.platform.smarter.sh"},
		{"next", "https://next.platform.smarter.sh"},
		{"prod", "https://platform.smarter.sh"},
	}
	for _, tc := range cases {
		t.Run(tc.environment, func(t *testing.T) {
			withViperValue(t, "environment", tc.environment)
			withViperValue(t, "config.root_domain", "platform.smarter.sh")

			if got, err := getAPIHost(); err != nil || got != tc.want {
				t.Errorf("getAPIHost() for %s = %q, want %q", tc.environment, got, tc.want)
			}
		})
	}

	t.Run("returns an error on an invalid environment", func(t *testing.T) {
		withViperValue(t, "environment", "nonsense")
		withViperValue(t, "config.root_domain", "platform.smarter.sh")

		if _, err := getAPIHost(); err == nil || !strings.Contains(err.Error(), "nonsense") {
			t.Errorf("getAPIHost() error = %v, want one naming the invalid environment", err)
		}
	})
}

func TestAPIRequestSuccess(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotQuery, gotBody string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotQuery = r.URL.RawQuery
		bodyBytes, _ := io.ReadAll(r.Body)
		gotBody = string(bodyBytes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	})
	defer srv.Close()

	body, err := APIRequest("widgets/", map[string]string{"name": "foo"}, `{"payload":true}`)
	if err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/api/v1/cli/widgets/" {
		t.Errorf("path = %q, want /api/v1/cli/widgets/", gotPath)
	}
	if gotAuth != "Token "+testAPIKey {
		t.Errorf("Authorization header = %q", gotAuth)
	}
	if !strings.Contains(gotQuery, "name=foo") {
		t.Errorf("query = %q, want to contain name=foo", gotQuery)
	}
	if gotBody != `{"payload":true}` {
		t.Errorf("request body = %q", gotBody)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("response body did not round-trip as JSON: %v", err)
	}
}

func TestAPIRequestNormalizesSlugToLowercaseWithTrailingSlash(t *testing.T) {
	var gotPath string
	srv := newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	if _, err := APIRequest("Some/MixedCase", map[string]string{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if gotPath != "/api/v1/cli/some/mixedcase/" {
		t.Errorf("path = %q, want /api/v1/cli/some/mixedcase/", gotPath)
	}
}

func TestAPIRequestMissingApiKeyReturnsErrorWithoutExiting(t *testing.T) {
	withViperValue(t, "environment", "alpha")
	withViperValue(t, "alpha.api_key", "")
	withViperValue(t, "api_key", "")

	_, err := APIRequest("whoami", map[string]string{})
	if err == nil {
		t.Fatal("APIRequest() error = nil, want an error for a missing api_key")
	}
}

func TestAPIRequestNon200ReturnsThePlatformsError(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"description":"boom","context":"test-context","stacktrace":"trace-here"}`))
	})

	_, err := APIRequest("whoami", map[string]string{})
	if err == nil {
		t.Fatal("APIRequest() error = nil, want the platform's error")
	}
	for _, want := range []string{"500", "boom", "test-context"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want it to contain %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "trace-here") {
		t.Errorf("error = %q, want no stack trace without --verbose", err)
	}

	withViperValue(t, "verbose", true)
	_, err = APIRequest("whoami", map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "trace-here") {
		t.Errorf("error = %v, want the stack trace with --verbose", err)
	}
}

func TestAPIRequestReadsTheJournalsErrorObject(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"api":"smarter.sh/v1","error":{"errorClass":"SAMBrokerErrorNotFound","description":"Guardrail x not found","context":"thing=Guardrail, command=describe","stacktrace":"trace-here"}}`))
	})

	_, err := APIRequest("describe/guardrail", map[string]string{"name": "x"})
	if err == nil || !strings.Contains(err.Error(), "Guardrail x not found") || !strings.Contains(err.Error(), "command=describe") {
		t.Errorf("error = %v, want the journal's description and context", err)
	}

	withViperValue(t, "verbose", true)
	_, err = APIRequest("describe/guardrail", map[string]string{"name": "x"})
	if err == nil || !strings.Contains(err.Error(), "SAMBrokerErrorNotFound") || !strings.Contains(err.Error(), "trace-here") {
		t.Errorf("error = %v, want the error class and stack trace with --verbose", err)
	}
}

func TestAPIRequestReadsDRFDetail(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"Invalid token."}`))
	})
	_, err := APIRequest("whoami", map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "Invalid token.") {
		t.Errorf("error = %v, want DRF's detail", err)
	}
}

func TestAPIRequestNonJSONErrorBodyIsQuoted(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>" + strings.Repeat("x", 300) + "</html>"))
	})

	_, err := APIRequest("whoami", map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "<html>") {
		t.Fatalf("error = %v, want the status and the start of the body", err)
	}
	if !strings.HasSuffix(err.Error(), "...") {
		t.Errorf("error = %q, want a long body shortened", err)
	}
}

func TestAPIRequestInvalidEnvironmentReturnsError(t *testing.T) {
	withViperValue(t, "environment", "nonsense")
	withViperValue(t, "nonsense.api_key", testAPIKey)

	if _, err := APIRequest("whoami", map[string]string{}); err == nil {
		t.Fatal("APIRequest() error = nil, want an error for an invalid environment")
	}
}

func TestAPIRequestUnreachableHostReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()
	SetAPIHostOverride(url)
	t.Cleanup(func() { SetAPIHostOverride("") })
	withViperValue(t, "api_key", testAPIKey)

	if _, err := APIRequest("whoami", map[string]string{}); err == nil {
		t.Fatal("APIRequest() error = nil, want a connection error")
	}
}

func TestAPIRequestTimesOut(t *testing.T) {
	release := make(chan struct{})
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		<-release
	})
	t.Cleanup(func() { close(release) })
	withViperValue(t, "timeout", 50*time.Millisecond)

	_, err := APIRequest("whoami", map[string]string{})
	if err == nil {
		t.Fatal("APIRequest() error = nil, want a timeout")
	}
}

func TestRequestTimeoutDefault(t *testing.T) {
	withViperValue(t, "timeout", time.Duration(0))
	if got := requestTimeout(); got != DefaultTimeout {
		t.Errorf("requestTimeout() = %v, want %v", got, DefaultTimeout)
	}
}

func TestAPIRequestSendsIdentifyingHeadersNotKwargs(t *testing.T) {
	var header http.Header
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		header = r.Header.Clone()
		_, _ = w.Write([]byte(`{}`))
	})

	if _, err := APIRequest("whoami", map[string]string{"name": "foo"}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	if header.Get("X-Smarter-Client") != ClientName {
		t.Errorf("X-Smarter-Client = %q, want %q", header.Get("X-Smarter-Client"), ClientName)
	}
	if !strings.HasPrefix(header.Get("User-Agent"), ClientName+"/") {
		t.Errorf("User-Agent = %q", header.Get("User-Agent"))
	}
	if header.Get("name") != "" {
		t.Error("kwargs must go in the query string only, not in headers")
	}
}

func TestAPIRequestVerboseMasksTheKeyOnStderr(t *testing.T) {
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	})
	withViperValue(t, "verbose", true)

	var logged bytes.Buffer
	log.SetOutput(&logged)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	stdout := captureStdout(t, func() {
		if _, err := APIRequest("whoami", map[string]string{}); err != nil {
			t.Fatalf("APIRequest() error = %v", err)
		}
	})
	if stdout != "" {
		t.Errorf("stdout = %q, want diagnostics on stderr only", stdout)
	}
	if strings.Contains(logged.String(), testAPIKey) {
		t.Error("verbose output contains the full api key")
	}
	if !strings.Contains(logged.String(), maskSecret(testAPIKey)) {
		t.Errorf("verbose output = %q, want the masked key", logged.String())
	}
}

func TestMaskSecret(t *testing.T) {
	if got := maskSecret("abc"); got != "****" {
		t.Errorf("maskSecret(short) = %q", got)
	}
	if got := maskSecret("0123456789"); got != "****6789" {
		t.Errorf("maskSecret() = %q", got)
	}
}
