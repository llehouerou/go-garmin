package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var weightCmd = &cobra.Command{
	Use:   "weight",
	Short: "Weight data (daily, range)",
}

var weightDailyCmd = &cobra.Command{
	Use:   "daily [date]",
	Short: "Get daily weight data",
	Long:  "Get daily weight data for the specified date (defaults to today). Date format: YYYY-MM-DD",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWeightDaily,
}

var weightRangeCmd = &cobra.Command{
	Use:   "range <start> <end>",
	Short: "Get weight data for a date range",
	Long:  "Get weight data for the specified date range. Date format: YYYY-MM-DD",
	Args:  cobra.ExactArgs(2),
	RunE:  runWeightRange,
}

func init() {
	weightCmd.AddCommand(weightDailyCmd)
	weightCmd.AddCommand(weightRangeCmd)
}

func runWeightDaily(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Weight.GetDaily(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runWeightRange(cmd *cobra.Command, args []string) error {
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

	data, err := client.Weight.GetRange(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}
