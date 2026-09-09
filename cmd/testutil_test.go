package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

// withViperValue sets key to value for the duration of the test and restores
// its previous value afterwards. viper is a global singleton, and several
// package-level cobra commands bind flags to it in init(), so tests must not
// call viper.Reset() (it would destroy those bindings); overriding individual
// keys and restoring them is the safe way to get per-test isolation.
func withViperValue(t *testing.T, key string, value interface{}) {
	t.Helper()
	old := viper.Get(key)
	viper.Set(key, value)
	t.Cleanup(func() {
		viper.Set(key, old)
	})
}

// newMockAPIServer starts an httptest.Server and points APIRequest at it for
// the duration of the test via SetAPIHostOverride, restoring the previous
// override on cleanup. It also ensures a valid api_key is configured so
// verifyApiKey() doesn't reject the request before it reaches the server.
func newMockAPIServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	SetAPIHostOverride(srv.URL)
	t.Cleanup(func() { SetAPIHostOverride("") })

	withViperValue(t, "api_key", testAPIKey)
	withViperValue(t, "environment", "alpha")

	return srv
}
