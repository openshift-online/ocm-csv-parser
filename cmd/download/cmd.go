package download

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"google.golang.org/api/sheets/v4"
)

var Cmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads the CSV file from the Google Sheet",
	Long:  "Downloads the CSV file from the Google Sheet and places in /tmp/",
	RunE:  run,
}

func run(cmd *cobra.Command, _ []string) (err error) {
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx)
	_ = sheetsService
	fmt.Println("Command under development, please download manually for now and use '--file / -f' when parsing")
	return nil
}
