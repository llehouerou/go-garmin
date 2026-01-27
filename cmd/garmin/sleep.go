package main

import (
	"github.com/spf13/cobra"
)

var sleepCmd = &cobra.Command{
	Use:   "sleep [date]",
	Short: "Get sleep data",
	Long:  "Get sleep data for the specified date. Date format: YYYY-MM-DD (defaults to today).",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSleep,
}

func runSleep(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Sleep.GetDaily(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}
