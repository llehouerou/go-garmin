package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Metrics data (training readiness, VO2 max, scores)",
}

var metricsReadinessCmd = &cobra.Command{
	Use:   "readiness [date]",
	Short: "Get training readiness",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsReadiness,
}

var metricsEnduranceCmd = &cobra.Command{
	Use:   "endurance [date]",
	Short: "Get endurance score",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsEndurance,
}

var metricsHillCmd = &cobra.Command{
	Use:   "hill [date]",
	Short: "Get hill score",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsHill,
}

var metricsVO2MaxCmd = &cobra.Command{
	Use:   "vo2max [date]",
	Short: "Get latest VO2 max / MET",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsVO2Max,
}

var metricsVO2MaxRangeCmd = &cobra.Command{
	Use:   "vo2max-range <start> <end>",
	Short: "Get VO2 max / MET for a date range",
	Args:  cobra.ExactArgs(2),
	RunE:  runMetricsVO2MaxRange,
}

var metricsStatusCmd = &cobra.Command{
	Use:   "status [date]",
	Short: "Get daily training status",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsStatus,
}

var metricsStatusAggCmd = &cobra.Command{
	Use:   "status-agg [date]",
	Short: "Get aggregated training status",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsStatusAgg,
}

var metricsLoadBalanceCmd = &cobra.Command{
	Use:   "load-balance [date]",
	Short: "Get training load balance",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsLoadBalance,
}

var metricsAcclimationCmd = &cobra.Command{
	Use:   "acclimation [date]",
	Short: "Get heat/altitude acclimation",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMetricsAcclimation,
}

func init() {
	metricsCmd.AddCommand(metricsReadinessCmd)
	metricsCmd.AddCommand(metricsEnduranceCmd)
	metricsCmd.AddCommand(metricsHillCmd)
	metricsCmd.AddCommand(metricsVO2MaxCmd)
	metricsCmd.AddCommand(metricsVO2MaxRangeCmd)
	metricsCmd.AddCommand(metricsStatusCmd)
	metricsCmd.AddCommand(metricsStatusAggCmd)
	metricsCmd.AddCommand(metricsLoadBalanceCmd)
	metricsCmd.AddCommand(metricsAcclimationCmd)
}

func runMetricsReadiness(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetTrainingReadiness(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsEndurance(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetEnduranceScore(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsHill(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetHillScore(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsVO2Max(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetMaxMetLatest(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsVO2MaxRange(cmd *cobra.Command, args []string) error {
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

	data, err := client.Metrics.GetMaxMetDaily(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsStatus(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetTrainingStatusDaily(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsStatusAgg(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetTrainingStatusAggregated(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsLoadBalance(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetTrainingLoadBalance(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runMetricsAcclimation(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Metrics.GetHeatAltitudeAcclimation(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}
