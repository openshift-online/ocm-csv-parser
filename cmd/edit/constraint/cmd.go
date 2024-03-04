package constraint

import (
	"github.com/openshift-online/ocm-csv-parser/cmd/edit/constraint/machinetype"
	"github.com/openshift-online/ocm-csv-parser/cmd/edit/constraint/region"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "constraint",
	Aliases: []string{"p"},
	Short:   "",
	Long:    "Edits a constraint's values for use with the Cloud Resource Service",
}

func init() {
	Cmd.AddCommand(machinetype.Cmd)
	Cmd.AddCommand(region.Cmd)
}
