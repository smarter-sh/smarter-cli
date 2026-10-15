package logs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/viper"
)

// withViperValue sets key to value for the duration of the test and restores
// its previous value afterwards. viper is a global singleton shared across
// this package and cmd, so tests must not call viper.Reset().
func withViperValue(t *testing.T, key string, value interface{}) {
	t.Helper()
	old := viper.Get(key)
	viper.Set(key, value)
	t.Cleanup(func() {
		viper.Set(key, old)
	})
}

// newMockAPIServer starts an httptest.Server and points cmd.APIRequest at it
// for the duration of the test.
func newMockAPIServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cmd.SetAPIHostOverride(srv.URL)
	t.Cleanup(func() { cmd.SetAPIHostOverride("") })

	withViperValue(t, "api_key", testAPIKey)
	withViperValue(t, "environment", "alpha")

	return srv
}

// exits replaces cmd.ErrorOutput's exit with a panic for the duration of
// the test, runs fn, and reports whether it exited.
func exits(t *testing.T, fn func()) (exited bool) {
	t.Helper()
	type exit struct{}
	cmd.SetExitFunc(func(int) { panic(exit{}) })
	t.Cleanup(func() { cmd.SetExitFunc(nil) })
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(exit); !ok {
				panic(r)
			}
			exited = true
		}
	}()
	fn()
	return false
}
