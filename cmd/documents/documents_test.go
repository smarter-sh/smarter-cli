package documents

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smarter-sh/smarter-cli/cmd"
	"github.com/spf13/pflag"
)

// documentsBody is the platform's answer: a journaled response whose data is CliApiDocuments.
const documentsBody = `{"data": {"vectorstore": "handbook", "documents": [], "added": [], "duplicates": []}, "message": "Vectorstore handbook has 0 documents."}`

// request is what the mock platform received.
type request struct {
	path        string
	query       map[string]string
	contentType string
	files       map[string]string // file name: content
	fileTypes   map[string]string // file name: content type
	fields      map[string]string
}

// recordRequests starts a mock platform that records each request and answers with documentsBody.
func recordRequests(t *testing.T) *[]request {
	t.Helper()
	var requests []request
	newMockAPIServer(t, func(w http.ResponseWriter, r *http.Request) {
		got := request{
			path:        r.URL.Path,
			query:       map[string]string{},
			contentType: r.Header.Get("Content-Type"),
			files:       map[string]string{},
			fileTypes:   map[string]string{},
			fields:      map[string]string{},
		}
		for key := range r.URL.Query() {
			got.query[key] = r.URL.Query().Get(key)
		}
		if strings.HasPrefix(got.contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Errorf("ParseMultipartForm() error = %v", err)
			}
			for key, values := range r.MultipartForm.Value {
				got.fields[key] = values[0]
			}
			for _, header := range r.MultipartForm.File[fileField] {
				file, _ := header.Open()
				content, _ := io.ReadAll(file)
				_ = file.Close()
				got.files[header.Filename] = string(content)
				got.fileTypes[header.Filename] = header.Header.Get("Content-Type")
			}
		}
		requests = append(requests, got)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(documentsBody))
	})
	return &requests
}

// writeFile writes a file in the test's temporary directory, and returns its path.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// resetFlags restores the command's flags after a test runs it: cobra keeps
// their values between runs, and a string array flag accumulates them.
func resetFlags(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		cmd.RootCmd.SetArgs(nil)
		vectorstoreCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if slice, ok := flag.Value.(pflag.SliceValue); ok {
				_ = slice.Replace(nil)
			} else {
				_ = flag.Value.Set(flag.DefValue)
			}
			flag.Changed = false
		})
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestListSendsTheNameOnly(t *testing.T) {
	requests := recordRequests(t)
	if _, err := APIRequest(vectorstoreKind, "handbook", Options{}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	got := (*requests)[0]
	if got.path != "/api/v1/cli/documents/vectorstore/" {
		t.Errorf("path = %q, want /api/v1/cli/documents/vectorstore/", got.path)
	}
	if len(got.query) != 1 || got.query["name"] != "handbook" {
		t.Errorf("query = %v, want only name=handbook", got.query)
	}
	if got.contentType != "application/json" {
		t.Errorf("content type = %q, want application/json", got.contentType)
	}
}

func TestAddUploadsFilesUrlAndMetadata(t *testing.T) {
	requests := recordRequests(t)
	handbook := writeFile(t, "handbook.md", "# Refunds")
	policies := writeFile(t, "policies.txt", "Be kind.")
	options := Options{
		Add:      []string{handbook, policies},
		URL:      "https://example.com/faq.html",
		Metadata: `{"source": "handbook"}`,
	}
	if _, err := APIRequest(vectorstoreKind, "handbook", options); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	got := (*requests)[0]
	if !strings.HasPrefix(got.contentType, "multipart/form-data; boundary=") {
		t.Errorf("content type = %q, want multipart/form-data", got.contentType)
	}
	if got.files["handbook.md"] != "# Refunds" || got.files["policies.txt"] != "Be kind." {
		t.Errorf("files = %v, want both files, by their base names", got.files)
	}
	if !strings.HasPrefix(got.fileTypes["handbook.md"], "text/markdown") {
		t.Errorf("handbook.md's content type = %q, want text/markdown", got.fileTypes["handbook.md"])
	}
	if got.fields["url"] != "https://example.com/faq.html" || got.fields["metadata"] != `{"source": "handbook"}` {
		t.Errorf("fields = %v, want url and metadata", got.fields)
	}
	if _, ok := got.query["remove"]; ok {
		t.Errorf("query = %v, want no remove", got.query)
	}
}

func TestAddOnlyAUrl(t *testing.T) {
	requests := recordRequests(t)
	if _, err := APIRequest(vectorstoreKind, "handbook", Options{URL: "https://example.com/faq.html"}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	got := (*requests)[0]
	if len(got.files) != 0 || got.fields["url"] != "https://example.com/faq.html" {
		t.Errorf("files = %v, fields = %v, want only the url", got.files, got.fields)
	}
	if _, ok := got.fields["metadata"]; ok {
		t.Errorf("fields = %v, want no empty metadata", got.fields)
	}
}

func TestRemoveSendsItsQueryParameter(t *testing.T) {
	requests := recordRequests(t)
	if _, err := APIRequest(vectorstoreKind, "handbook", Options{Remove: "handbook.md"}); err != nil {
		t.Fatalf("APIRequest() error = %v", err)
	}
	got := (*requests)[0]
	if got.query["remove"] != "handbook.md" || got.query["name"] != "handbook" {
		t.Errorf("query = %v, want name and remove", got.query)
	}
}

func TestRefusals(t *testing.T) {
	requests := recordRequests(t)
	handbook := writeFile(t, "handbook.md", "# Refunds")
	cases := []struct {
		name    string
		options Options
		want    string
	}{
		{"add and remove", Options{Add: []string{handbook}, Remove: "x"}, "not both"},
		{"metadata alone", Options{Metadata: "{}"}, "--metadata needs documents to add"},
		{"a missing file", Options{Add: []string{filepath.Join(t.TempDir(), "missing.md")}}, "could not read"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := APIRequest(vectorstoreKind, "handbook", tc.options)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error = %v, want it to contain %q", err, tc.want)
			}
		})
	}
	if len(*requests) != 0 {
		t.Errorf("sent %d requests, want none", len(*requests))
	}
}

func TestTheCommandAddsAndPrintsTheResponse(t *testing.T) {
	requests := recordRequests(t)
	withViperValue(t, "output_format", "json")
	resetFlags(t)
	handbook := writeFile(t, "handbook.md", "# Refunds")
	cmd.RootCmd.SetArgs([]string{"documents", "vectorstore", "handbook", "--add", handbook, "-a", handbook})
	out := captureStdout(t, func() {
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
	if !strings.Contains(out, "Vectorstore handbook has 0 documents.") {
		t.Errorf("output = %q, want the response", out)
	}
	if len(*requests) != 1 || len((*requests)[0].files) != 1 {
		t.Errorf("requests = %v, want one upload", *requests)
	}
}

func TestTheCommandReportsAnError(t *testing.T) {
	recordRequests(t)
	resetFlags(t)
	cmd.RootCmd.SetArgs([]string{"documents", "vectorstore", "handbook", "--metadata", "{}"})
	if !exits(t, func() { _ = cmd.RootCmd.Execute() }) {
		t.Error("the command didn't exit with an error")
	}
}

func TestTheVerbHasAVectorstoreCommand(t *testing.T) {
	var uses []string
	for _, c := range documentsCmd.Commands() {
		uses = append(uses, c.Use)
	}
	if len(uses) != 1 || uses[0] != "vectorstore <name>" {
		t.Errorf("subcommands = %v, want only vectorstore <name>", uses)
	}
}
