package edit

import (
	"github.com/openshift-online/ocm-csv-parser/cmd/edit/constraint"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edits a resource",
	Long:    "Edits a specific resource",
}

func init() {
	Cmd.AddCommand(constraint.Cmd)
}
