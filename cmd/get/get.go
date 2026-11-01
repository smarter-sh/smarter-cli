/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package get

import (
	"log"
	"net/url"

	"github.com/smarter-sh/smarter-cli/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// annotationFlag is get's repeatable annotation selector flag, and its viper key.
const annotationFlag = "annotation"

// APIRequest requests a list of resources of a kind, narrowed by kwargs (e.g.
// name) and by every --annotation requirement, each sent as its own
// annotation query parameter: requirement values may contain commas, so they
// can't be joined into one.
func APIRequest(kind string, kwargs map[string]string) ([]byte, error) {
	params := url.Values{}
	for key, value := range kwargs {
		params.Add(key, value)
	}
	for _, requirement := range viper.GetStringSlice(annotationFlag) {
		params.Add(annotationFlag, requirement)
	}

	// en route to /api/v1/cli/get/<str:kind>/
	return cmd.APIRequestWithParams("get/"+kind+"/", params)
}

func ConsoleOutput(bodyJson []byte) {
	if !viper.IsSet("output_format") {
		viper.Set("output_format", "tabular")
	}
	if err := cmd.ConsoleOutput(bodyJson); err != nil {
		cmd.ErrorOutput(err)
	}
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Generate a list of Smarter resources",
	Long: `Generate a list of Smarter resources:

smarter get [kind] [flags]

The Smarter API will return a list of resources in the specified format,
or a manifest for a specific resource.

Narrow the list by the resources' annotations with --annotation, one
requirement per flag, in the syntax of Kubernetes label selectors. A resource
must match every requirement:

  key          the resource has the annotation
  !key         the resource doesn't have it
  key=value    the annotation has this value
  key!=value   the annotation doesn't have this value, or is missing

smarter get llmclients --annotation acme.com/crm-account-id=A-10042
smarter get secrets --annotation acme.com/owner --annotation '!acme.com/archived'`,
}

func init() {
	cmd.RootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringArray(annotationFlag, nil, "an annotation selector requirement that every resource must match: key, !key, key=value or key!=value; repeat it for several")
	if err := viper.BindPFlag(annotationFlag, getCmd.PersistentFlags().Lookup(annotationFlag)); err != nil {
		log.Fatalf("Error binding flag '%s': %v", annotationFlag, err)
	}

	getCmd.PersistentFlags().Int("i", 10, "Number of resources to retrieve")
	getCmd.PersistentFlags().Bool("asc", false, "Sort results in ascending order")
	getCmd.PersistentFlags().Bool("desc", false, "Sort results in descending order")
	// The Smarter API ignores these: each kind returns its resources in a fixed
	// order, up to a fixed number. They stay accepted, so that old scripts run.
	for _, name := range []string{"i", "asc", "desc"} {
		if err := getCmd.PersistentFlags().MarkDeprecated(name, "the Smarter API ignores it"); err != nil {
			log.Fatalf("Error deprecating flag: %v", err)
		}
	}
}
