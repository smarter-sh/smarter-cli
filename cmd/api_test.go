package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/viper"
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

		if got := getAPIHost(); got != "http://example.invalid" {
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

			if got := getAPIHost(); got != tc.want {
				t.Errorf("getAPIHost() for %s = %q, want %q", tc.environment, got, tc.want)
			}
		})
	}

	t.Run("panics on an invalid environment", func(t *testing.T) {
		withViperValue(t, "environment", "nonsense")
		withViperValue(t, "config.root_domain", "platform.smarter.sh")

		defer func() {
			if recover() == nil {
				t.Error("getAPIHost() did not panic on an invalid environment")
			}
		}()
		getAPIHost()
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

// TestAPIRequestNon200ExitsProcess verifies that a non-200 response from the
// API causes the process to exit with status 1 via ErrorOutput. Since that
// path calls os.Exit, it must be exercised in a subprocess rather than in
// this test binary directly.
func TestAPIRequestNon200ExitsProcess(t *testing.T) {
	if os.Getenv("SMARTER_CLI_TEST_CRASHER") == "1" {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"description":"boom","context":"test-context"}`))
		}))
		defer srv.Close()

		SetAPIHostOverride(srv.URL)
		viper.Set("api_key", testAPIKey)
		viper.Set("environment", "alpha")

		_, _ = APIRequest("whoami", map[string]string{})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestAPIRequestNon200ExitsProcess")
	cmd.Env = append(os.Environ(), "SMARTER_CLI_TEST_CRASHER=1")
	output, err := cmd.CombinedOutput()

	if exitErr, ok := err.(*exec.ExitError); ok && !exitErr.Success() {
		if !strings.Contains(string(output), "boom") {
			t.Errorf("expected error output to mention the server's description, got: %s", output)
		}
		return
	}
	t.Fatalf("process exited with err=%v, want a non-zero exit status; output: %s", err, output)
}
