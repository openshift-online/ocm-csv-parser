package parse

import (
	"github.com/openshift-online/ocm-csv-parser/cmd/parse/machinetypes"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "parse",
	Aliases: []string{"p"},
	Short:   "Parses resources from a CSV file into a configmap",
	Long:    "Parses resources from a CSV file into a yaml configmap for use with the Cloud Resource Service",
}

func init() {
	Cmd.AddCommand(machinetypes.Cmd)
}
