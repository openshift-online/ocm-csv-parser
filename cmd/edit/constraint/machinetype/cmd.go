package machinetype

import (
	"fmt"

	"github.com/openshift-online/ocm-csv-parser/pkg/edit"
	"github.com/spf13/cobra"
	errors "github.com/zgalor/weberr"
)

var args struct {
	id              string
	env             string
	enabled         bool
	ccsOnly         bool
	appInterfaceDir string
}

func init() {
	flags := Cmd.Flags()

	flags.StringVarP(
		&args.appInterfaceDir,
		"app-interface",
		"a",
		"",
		"This flag should be set to the filepath of your local app-interface. This allows the CLI to update "+
			"your local app-interface directly (required)",
	)

	flags.StringVar(
		&args.env,
		"environment",
		"",
		"Environment of which constraintmap you want to edit (integration/stage/production -- required)")

	flags.StringVar(
		&args.id,
		"id",
		"",
		"ID of the machine type you want to edit (required)",
	)

	flags.BoolVar(
		&args.enabled,
		"enabled",
		false,
		"Whether or not the machine type you want to edit should be enabled/disabled (bool, use --enabled=<bool_value>)",
	)

	flags.BoolVar(
		&args.ccsOnly,
		"ccs-only",
		false,
		"Whether or not the machine type you want to edit should be ccs-only or not (bool, use --ccs-only=<bool_value>)",
	)

	Cmd.MarkFlagRequired("environment")
	Cmd.MarkFlagRequired("app-interface")
	Cmd.MarkFlagRequired("id")
}

var Cmd = &cobra.Command{
	Use:     "machinetype",
	Aliases: []string{"m", "machine-type"},
	Short:   "Edits the machine types constraintmap in app-interface",
	Long: "Edits a single machine type entry in app-interface at a time for a constraintmap in a specific env." +
		" Example: 'ocm-csv-parser edit constraint machinetype --app-interface \"../app-interface\"" +
		" --environment integration --id m5.xlarge --enabled=false'",
	RunE: run,
}

func run(cmd *cobra.Command, _ []string) (err error) {
	constraintMap, err := edit.NewConstraintMap(args.env, args.appInterfaceDir)
	if err != nil {
		return errors.Errorf("Yaml constraintmap error: '%v'", err)
	}

	var ccsOnly *bool = nil
	var enabled *bool = nil
	if cmd.Flags().Changed("ccs-only") {
		ccsOnly = &args.ccsOnly
	}
	if cmd.Flags().Changed("enabled") {
		enabled = &args.enabled
	}

	err = constraintMap.EditConstraint(args.id, ccsOnly, enabled, nil)
	if err != nil {
		return errors.Errorf("Edit constraint error: '%v'", err)
	}

	fmt.Printf("Finished, updated machinetype constraint with ID '%s' in app-interface: '%s'\n\n", args.id,
		constraintMap.Path)
	return nil
}
