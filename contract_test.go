package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// The contract test checks the cli against the platform's OpenAPI document,
// pinned in testdata/contracts/openapi.json (make contracts): see
// doc/adr/022-pinned-openapi.md. The unit tests' mock servers answer any
// path, so without this test the cli can call a route that the platform
// renamed and still pass.

// testAPIKey is a syntactically valid (64 hex char) placeholder API key.
const testAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// TestMain seeds a fake $HOME with a complete config.yaml, because
// cobra.OnInitialize(initConfig) runs on every Execute, and must never touch
// the developer's real ~/.smarter.
func TestMain(m *testing.M) {
	// a developer's SMARTER_* variables must not change the tests' outcome.
	for _, name := range cmd.SmarterEnvironmentVariables {
		if err := os.Unsetenv(name); err != nil {
			panic(err)
		}
	}
	dir, err := os.MkdirTemp("", "smarter-cli-test-home")
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(dir, ".smarter")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		panic(err)
	}
	config := "config:\n  environment: alpha\n  root_domain: platform.smarter.sh\nalpha:\n  api_key: \"" + testAPIKey + "\"\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(config), 0o600); err != nil {
		panic(err)
	}
	if err := os.Setenv("HOME", dir); err != nil {
		panic(err)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// operation is one operation of the OpenAPI document.
type operation struct {
	id       string
	method   string
	template string // e.g. /cli/get/{kind}/
	pattern  *regexp.Regexp
	query    map[string]bool // the documented query parameters
}

type openAPIDocument struct {
	Paths map[string]map[string]struct {
		OperationID string `json:"operationId"`
		Parameters  []struct {
			Name string `json:"name"`
			In   string `json:"in"`
		} `json:"parameters"`
	} `json:"paths"`
	Components struct {
		Schemas struct {
			SAMKind struct {
				Enum []string `json:"enum"`
			} `json:"SAMKind"`
		} `json:"schemas"`
	} `json:"components"`
}

// loadOpenAPI returns the document's /cli/ operations, and its SAMKinds.
func loadOpenAPI(t *testing.T) ([]operation, []string) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", "contracts", "openapi.json"))
	if err != nil {
		t.Fatalf("reading the pinned OpenAPI document (make contracts): %v", err)
	}
	var document openAPIDocument
	if err := json.Unmarshal(body, &document); err != nil {
		t.Fatalf("parsing the pinned OpenAPI document: %v", err)
	}
	var operations []operation
	for template, item := range document.Paths {
		if !strings.HasPrefix(template, "/cli/") {
			continue
		}
		parts := strings.Split(template, "/")
		for i, part := range parts {
			if strings.HasPrefix(part, "{") {
				parts[i] = "[^/]+"
			} else {
				parts[i] = regexp.QuoteMeta(part)
			}
		}
		pattern := regexp.MustCompile("(?i)^/api/v1" + strings.Join(parts, "/") + "$")
		for method, op := range item {
			query := map[string]bool{}
			for _, parameter := range op.Parameters {
				if parameter.In == "query" {
					query[parameter.Name] = true
				}
			}
			operations = append(operations, operation{id: op.OperationID, method: strings.ToUpper(method), template: template, pattern: pattern, query: query})
		}
	}
	if len(operations) == 0 || len(document.Components.Schemas.SAMKind.Enum) == 0 {
		t.Fatal("the pinned OpenAPI document has no /cli/ operations or no SAMKind enum")
	}
	return operations, document.Components.Schemas.SAMKind.Enum
}

// request is one request that a command sent to the mock platform.
type request struct {
	command string
	method  string
	path    string
	query   []string // the query parameters' names
}

// commandArgs are the arguments, beyond the command path and a name, that a
// command needs to send its request.
func commandArgs(t *testing.T, path string) []string {
	manifest := filepath.Join(t.TempDir(), "manifest.yaml")
	if err := os.WriteFile(manifest, []byte("apiVersion: smarter.sh/v1\nkind: Secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	switch path {
	case "smarter apply", "smarter validate":
		return []string{"-f", manifest}
	case "smarter version":
		return []string{"-v"}
	case "smarter chat prompt":
		return []string{"--llmclient", "example", "--prompt", "hello"}
	case "smarter chat config":
		return []string{"--llmclient", "example"}
	}
	return nil
}

// noRequest are the commands that don't call the platform.
var noRequest = map[string]bool{
	"smarter configure": true,
}

// leaves returns every runnable command under c.
func leaves(c *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, child := range c.Commands() {
		if child.Name() == "help" || child.Name() == "completion" || child.Hidden {
			continue
		}
		if child.Runnable() {
			out = append(out, child)
		}
		out = append(out, leaves(child)...)
	}
	return out
}

// resetFlags restores every flag of the command tree to its default, because
// cobra keeps flag values between Execute calls in one process.
func resetFlags(c *cobra.Command) {
	reset := func(f *pflag.Flag) {
		// Set appends to a slice flag, and its DefValue reads "[]", so a
		// slice flag (every one has an empty default) is emptied instead.
		if slice, ok := f.Value.(pflag.SliceValue); ok {
			_ = slice.Replace(nil)
		} else {
			_ = f.Value.Set(f.DefValue)
		}
		f.Changed = false
	}
	c.Flags().VisitAll(reset)
	c.PersistentFlags().VisitAll(reset)
	for _, child := range c.Commands() {
		resetFlags(child)
	}
}

// runEveryCommand executes every command against a mock platform, and returns
// the requests it received.
func runEveryCommand(t *testing.T) []request {
	t.Helper()
	var mu sync.Mutex
	var requests []request
	var current string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		var query []string
		for name := range r.URL.Query() {
			query = append(query, name)
		}
		requests = append(requests, request{command: current, method: r.Method, path: r.URL.Path, query: query})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"titles":[],"items":[],"valid":true,"errors":[]},"message":"ok"}`))
	}))
	t.Cleanup(srv.Close)
	cmd.SetAPIHostOverride(srv.URL)
	t.Cleanup(func() { cmd.SetAPIHostOverride("") })

	type exited struct{ code int }
	cmd.SetExitFunc(func(code int) { panic(exited{code}) })
	t.Cleanup(func() { cmd.SetExitFunc(nil) })

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = devNull.Close() })
	if os.Getenv("SMARTER_CLI_CONTRACT_DEBUG") != "1" {
		log.SetOutput(devNull)
		t.Cleanup(func() { log.SetOutput(os.Stderr) })
	}

	type invocation struct {
		path string
		args []string
	}
	var invocations []invocation
	for _, leaf := range leaves(cmd.RootCmd) {
		path := leaf.CommandPath()
		if noRequest[path] {
			continue
		}
		args := strings.Fields(path)[1:]
		if strings.Contains(leaf.Use, "<") {
			args = append(args, "example")
		}
		invocations = append(invocations, invocation{path, append(args, commandArgs(t, path)...)})
	}
	// options that send a different request
	invocations = append(invocations, invocation{"smarter apply --dry-run", append([]string{"apply", "--dry-run"}, commandArgs(t, "smarter apply")...)})

	for _, inv := range invocations {
		path, args := inv.path, inv.args

		mu.Lock()
		current = path
		before := len(requests)
		mu.Unlock()

		func() {
			stdout, stderr := os.Stdout, os.Stderr
			os.Stdout, os.Stderr = devNull, devNull
			defer func() {
				os.Stdout, os.Stderr = stdout, stderr
				if r := recover(); r != nil {
					t.Errorf("%s failed: %v (rerun it with SMARTER_CLI_CONTRACT_DEBUG=1 to see its output)", path, r)
				}
			}()
			if os.Getenv("SMARTER_CLI_CONTRACT_DEBUG") == "1" {
				os.Stdout, os.Stderr = stdout, stderr
			}
			cmd.RootCmd.SetArgs(args)
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Errorf("%s: %v", path, err)
			}
		}()
		resetFlags(cmd.RootCmd)
		for _, key := range []string{"output_format", "verbose"} {
			viper.Set(key, nil)
		}

		mu.Lock()
		if len(requests) == before {
			t.Errorf("%s sent no request", path)
		}
		mu.Unlock()
	}
	return requests
}

// match returns the operation that a request is, or nil.
func match(operations []operation, r request) *operation {
	for i, op := range operations {
		if op.method == r.method && op.pattern.MatchString(r.path) {
			return &operations[i]
		}
	}
	return nil
}

// kindOf returns the {kind} segment of a request, and its verb, e.g. ("get", "secret").
func kindOf(op *operation, path string) (string, string, bool) {
	templateParts := strings.Split(op.template, "/")
	pathParts := strings.Split(strings.TrimPrefix(path, "/api/v1"), "/")
	for i, part := range templateParts {
		if part == "{kind}" && i < len(pathParts) {
			return templateParts[i-1], strings.ToLower(pathParts[i]), true
		}
	}
	return "", "", false
}

// kindVerbs are the verbs that every SAMKind has a command for, apart from
// kindExceptions.
var kindVerbs = []string{"get", "describe", "delete", "logs", "example-manifest", "json-schema"}

// kindExceptions are the verbs that a kind deliberately has no command for.
var kindExceptions = map[string][]string{
	"account": {"delete", "logs"}, // an account isn't deleted from the cli, and has no name to read logs by
	"user":    {"logs"},           // user commands take a username; the logs api reads a name
}

func TestCommandsMatchTheOpenAPIDocument(t *testing.T) {
	operations, samKinds := loadOpenAPI(t)
	requests := runEveryCommand(t)

	called := map[string]bool{}
	kindsByVerb := map[string]map[string]bool{}
	known := map[string]bool{}
	for _, kind := range samKinds {
		known[strings.ToLower(kind)] = true
	}

	for _, r := range requests {
		op := match(operations, r)
		if op == nil {
			t.Errorf("%s sent %s %s, which is not an operation of the platform's OpenAPI document", r.command, r.method, r.path)
			continue
		}
		called[op.id] = true
		for _, name := range r.query {
			if !op.query[name] {
				t.Errorf("%s sent the query parameter %q, which %s %s doesn't document", r.command, name, op.method, op.template)
			}
		}
		if verb, kind, ok := kindOf(op, r.path); ok {
			if !known[kind] {
				t.Errorf("%s sent the kind %q, which is not a SAMKind of the platform's OpenAPI document", r.command, kind)
			}
			if kindsByVerb[verb] == nil {
				kindsByVerb[verb] = map[string]bool{}
			}
			kindsByVerb[verb][kind] = true
		}
	}

	var uncalled []string
	for _, op := range operations {
		if !called[op.id] {
			uncalled = append(uncalled, fmt.Sprintf("%s %s", op.method, op.template))
		}
	}
	sort.Strings(uncalled)
	if len(uncalled) > 0 {
		t.Errorf("the platform publishes /cli/ operations that no command calls:\n  %s", strings.Join(uncalled, "\n  "))
	}

	for _, kind := range samKinds {
		lower := strings.ToLower(kind)
		for _, verb := range kindVerbs {
			exempt := false
			for _, e := range kindExceptions[lower] {
				exempt = exempt || e == verb
			}
			if !exempt && !kindsByVerb[verb][lower] {
				t.Errorf("the SAMKind %s has no %q command", kind, verb)
			}
		}
	}
}
