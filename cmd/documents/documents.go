/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/

// Package documents is the documents verb: list a Vectorstore's documents, add
// documents to it, or remove one, with the platform's cli api,
// /api/v1/cli/documents/<kind>/. Only a Vectorstore has documents.
package documents

import (
	"errors"
	"log"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vectorstoreKind is the kind of the resource whose documents the command manages.
const vectorstoreKind = "Vectorstore"

// fileField is the multipart field of the files to add: the platform reads "files".
const fileField = "files"

// Options are the documents command's flags.
type Options struct {
	Add      []string // paths of files to add
	URL      string   // a public url of a page or file to add
	Metadata string   // a JSON object, stored with every chunk of every document added
	Remove   string   // the id or name of a document to remove
}

// adding reports whether the options add documents.
func (o Options) adding() bool {
	return len(o.Add) > 0 || o.URL != ""
}

// APIRequest sends the documents command for the resource of kind named name:
// a list without options, a multipart upload with Add or URL, or a removal
// with Remove.
func APIRequest(kind, name string, options Options) ([]byte, error) {
	if options.adding() && options.Remove != "" {
		return nil, errors.New("add documents, or remove one, but not both")
	}
	if options.Metadata != "" && !options.adding() {
		return nil, errors.New("--metadata needs documents to add, with --add or --url")
	}
	kwargs := map[string]string{"name": name}
	if options.Remove != "" {
		kwargs["remove"] = options.Remove
	}
	// en route to /api/v1/cli/documents/<str:kind>/
	slug := "documents/" + kind + "/"
	if !options.adding() {
		return cmd.APIRequest(slug, kwargs)
	}
	form := cmd.MultipartForm{
		Files:     options.Add,
		FileField: fileField,
		Fields:    map[string]string{"url": options.URL, "metadata": options.Metadata},
	}
	return cmd.APIRequestMultipart(slug, kwargs, form)
}

// ConsoleOutput adapts cmd.ConsoleOutput's error return.
func ConsoleOutput(bodyJson []byte) {
	if err := cmd.ConsoleOutput(bodyJson); err != nil {
		cmd.ErrorOutput(err)
	}
}

var documentsCmd = &cobra.Command{
	Use:   "documents <kind> <name>",
	Short: "List, add or remove a resource's documents",
	Long: `Lists, adds or removes the documents of a resource that holds them, a Vectorstore:

smarter documents vectorstore <name>
smarter documents vectorstore <name> --add handbook.md --add policies.pdf
smarter documents vectorstore <name> --url https://example.com/faq.html
smarter documents vectorstore <name> --remove handbook.md`,
}

var vectorstoreCmd = &cobra.Command{
	Use:   "vectorstore <name>",
	Short: "List, add or remove a Vectorstore's documents",
	Long: `Lists, adds or removes a Vectorstore's documents:

smarter documents vectorstore <name>
	Lists its documents, newest first, with each one's status and number of chunks.

smarter documents vectorstore <name> --add <file> [--add <file> ...] [--url <url>] [--metadata <json>]
	Uploads files (PDF, text, Markdown, CSV, JSON, YAML or HTML, up to 50 MiB each),
	and the page or file at a public url. Their text is extracted now, and they are
	loaded in the background, once the Vectorstore is deployed and ready. --metadata,
	a JSON object, is stored with every chunk of every document, e.g. '{"source": "faq"}'.
	A document whose content the Vectorstore already has is not added again.

smarter documents vectorstore <name> --remove <document>
	Removes a document, by its id or its name, and its chunks, in the background.

Only the Vectorstore's owner, an account admin, may add or remove its documents.`,
	Args: cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		add, _ := c.Flags().GetStringArray("add")
		url, _ := c.Flags().GetString("url")
		metadata, _ := c.Flags().GetString("metadata")
		remove, _ := c.Flags().GetString("remove")
		options := Options{Add: add, URL: url, Metadata: metadata, Remove: remove}
		bodyJson, err := APIRequest(vectorstoreKind, args[0], options)
		if err != nil {
			cmd.ErrorOutput(err)
			return
		}
		ConsoleOutput(bodyJson)
	},
}

// bindFlag binds a flag to viper, as every flag of the cli is, under a key of
// this command's own, so that it never shares one with another command's flag.
func bindFlag(name string) {
	if err := viper.BindPFlag("documents_"+name, vectorstoreCmd.Flags().Lookup(name)); err != nil {
		log.Fatalf("Error binding flag '%s': %v", name, err)
	}
}

func init() {
	cmd.RootCmd.AddCommand(documentsCmd)
	documentsCmd.AddCommand(vectorstoreCmd)

	vectorstoreCmd.Flags().StringArrayP("add", "a", nil, "a file to add; repeat it to add several")
	vectorstoreCmd.Flags().String("url", "", "a public url of a page or file to add")
	vectorstoreCmd.Flags().String("metadata", "", "a JSON object stored with every chunk of the documents added")
	vectorstoreCmd.Flags().StringP("remove", "r", "", "the id or name of a document to remove")
	for _, name := range []string{"add", "url", "metadata", "remove"} {
		bindFlag(name)
	}
}
