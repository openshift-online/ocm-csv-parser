package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/openshift-online/ocm-csv-parser/cmd/edit"
	"github.com/openshift-online/ocm-csv-parser/cmd/parse"
	"github.com/openshift-online/ocm-csv-parser/pkg/color"
	"github.com/spf13/cobra"
)

var (
	resources []map[string]interface{}
)

var root = &cobra.Command{
	Use:   "ocm-csv-parser",
	Short: "Command line tool for parsing CSV files",
	Long: "Command line tool for downloading and parsing CSV files into yaml configmaps for use in " +
		"the Cloud Resource Service\n",
}

func init() {
	// Add the command line flags:
	color.AddFlag(root)

	root.AddCommand(parse.Cmd)
	root.AddCommand(edit.Cmd)
}

func main() {
	// Execute the root command:
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "Did you mean this?") {
			fmt.Fprintf(os.Stderr, "Failed to execute root command: %s\n", err)
		}
		os.Exit(1)
	}
}
