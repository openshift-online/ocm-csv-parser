package region

import (
	"fmt"

	"github.com/openshift-online/ocm-csv-parser/pkg/edit"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var args struct {
	id              string
	env             string
	enabled         bool
	ccsOnly         bool
	govcloud        bool
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
			"your local app-interface directly",
	)

	flags.StringVar(
		&args.env,
		"environment",
		"",
		"Environment of which constraintmap you want to edit (integration/staging/production -- required)")

	flags.StringVar(
		&args.id,
		"id",
		"",
		"ID of the region you want to edit (required)",
	)

	flags.BoolVar(
		&args.enabled,
		"enabled",
		false,
		"Whether or not the region you want to edit should be enabled/disabled (bool, use --enabled=<bool_value>)",
	)

	flags.BoolVar(
		&args.ccsOnly,
		"ccs-only",
		false,
		"Whether or not the region you want to edit should be ccs-only or not (bool, use --ccs-only=<bool_value>)",
	)

	flags.BoolVar(
		&args.govcloud,
		"govcloud",
		false,
		"Whether or not the region you want to edit is used for govcloud/fedramp (bool, use --govcloud=<bool_value>)",
	)

	Cmd.MarkFlagRequired("environment")
	Cmd.MarkFlagRequired("app-interface")
	Cmd.MarkFlagRequired("id")
}

var Cmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"r", "cloud-region", "cloudregion"},
	Short:   "Edits a region in a constraintmap in app-interface",
	Long: "Edits a single cloud region entry in app-interface at a time for a constraintmap in a specific env." +
		" Example: 'ocm-csv-parser edit constraint region --app-interface \"../app-interface\" --environment stage" +
		" --id us-west-1 --govcloud=true'",
	RunE: run,
}

func run(cmd *cobra.Command, _ []string) (err error) {
	constraintMap, err := edit.NewConstraintMap(args.env, args.appInterfaceDir)
	if err != nil {
		return errors.Errorf("Yaml constraintmap error: '%v'", err)
	}

	var ccsOnly *bool = nil
	var enabled *bool = nil
	var govcloud *bool = nil
	if cmd.Flags().Changed("ccs-only") {
		ccsOnly = &args.ccsOnly
	}
	if cmd.Flags().Changed("enabled") {
		enabled = &args.enabled
	}
	if cmd.Flags().Changed("govcloud") {
		govcloud = &args.govcloud
	}

	err = constraintMap.EditConstraint(args.id, ccsOnly, enabled, govcloud)
	if err != nil {
		return errors.Errorf("Edit constraint error: '%v'", err)
	}

	fmt.Printf("Finished, updated region constraint with ID '%s' in app-interface: '%s'\n\n", args.id,
		constraintMap.Path)
	return nil
}
