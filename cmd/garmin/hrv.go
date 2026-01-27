package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var hrvCmd = &cobra.Command{
	Use:   "hrv",
	Short: "HRV data (daily, range)",
}

var hrvDailyCmd = &cobra.Command{
	Use:   "daily [date]",
	Short: "Get daily HRV data",
	Long:  "Get daily HRV data for the specified date (defaults to today). Date format: YYYY-MM-DD",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHRVDaily,
}

var hrvRangeCmd = &cobra.Command{
	Use:   "range <start> <end>",
	Short: "Get HRV data for a date range",
	Long:  "Get HRV data for the specified date range. Date format: YYYY-MM-DD",
	Args:  cobra.ExactArgs(2),
	RunE:  runHRVRange,
}

func init() {
	hrvCmd.AddCommand(hrvDailyCmd)
	hrvCmd.AddCommand(hrvRangeCmd)
}

func runHRVDaily(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.HRV.GetDaily(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runHRVRange(cmd *cobra.Command, args []string) error {
	startDate, err := time.Parse("2006-01-02", args[0])
	if err != nil {
		return fmt.Errorf("invalid start date: %s", args[0])
	}

	endDate, err := time.Parse("2006-01-02", args[1])
	if err != nil {
		return fmt.Errorf("invalid end date: %s", args[1])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.HRV.GetRange(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}
