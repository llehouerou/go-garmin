package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var biometricCmd = &cobra.Command{
	Use:   "biometric",
	Short: "Biometric data (lactate threshold, FTP, power-to-weight)",
}

var biometricLactateCmd = &cobra.Command{
	Use:   "lactate",
	Short: "Get latest lactate threshold",
	Args:  cobra.NoArgs,
	RunE:  runBiometricLactate,
}

var biometricFTPCmd = &cobra.Command{
	Use:   "ftp",
	Short: "Get latest cycling FTP",
	Args:  cobra.NoArgs,
	RunE:  runBiometricFTP,
}

var biometricPowerToWeightCmd = &cobra.Command{
	Use:   "power-to-weight [date]",
	Short: "Get power-to-weight ratio for running",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runBiometricPowerToWeight,
}

var biometricLactateSpeedRangeCmd = &cobra.Command{
	Use:   "lactate-speed-range <start> <end>",
	Short: "Get lactate threshold speed for a date range",
	Args:  cobra.ExactArgs(2),
	RunE:  runBiometricLactateSpeedRange,
}

var biometricLactateHRRangeCmd = &cobra.Command{
	Use:   "lactate-hr-range <start> <end>",
	Short: "Get lactate threshold heart rate for a date range",
	Args:  cobra.ExactArgs(2),
	RunE:  runBiometricLactateHRRange,
}

var biometricFTPRangeCmd = &cobra.Command{
	Use:   "ftp-range <start> <end>",
	Short: "Get FTP for a date range",
	Args:  cobra.ExactArgs(2),
	RunE:  runBiometricFTPRange,
}

func init() {
	biometricCmd.AddCommand(biometricLactateCmd)
	biometricCmd.AddCommand(biometricFTPCmd)
	biometricCmd.AddCommand(biometricPowerToWeightCmd)
	biometricCmd.AddCommand(biometricLactateSpeedRangeCmd)
	biometricCmd.AddCommand(biometricLactateHRRangeCmd)
	biometricCmd.AddCommand(biometricFTPRangeCmd)
}

func runBiometricLactate(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Biometric.GetLatestLactateThreshold(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runBiometricFTP(cmd *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Biometric.GetCyclingFTP(cmd.Context())
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runBiometricPowerToWeight(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Biometric.GetPowerToWeight(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runBiometricLactateSpeedRange(cmd *cobra.Command, args []string) error {
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

	data, err := client.Biometric.GetLactateThresholdSpeedRange(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runBiometricLactateHRRange(cmd *cobra.Command, args []string) error {
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

	data, err := client.Biometric.GetLactateThresholdHRRange(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runBiometricFTPRange(cmd *cobra.Command, args []string) error {
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

	data, err := client.Biometric.GetFTPRange(cmd.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	return printJSON(data)
}
