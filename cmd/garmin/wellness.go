package main

import (
	"github.com/spf13/cobra"
)

var wellnessCmd = &cobra.Command{
	Use:   "wellness",
	Short: "Wellness data (stress, body battery, heart rate, SpO2, respiration, intensity)",
}

var wellnessStressCmd = &cobra.Command{
	Use:   "stress [date]",
	Short: "Get stress data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessStress,
}

var wellnessBodyBatteryCmd = &cobra.Command{
	Use:   "body-battery [date]",
	Short: "Get body battery data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessBodyBattery,
}

var wellnessHeartRateCmd = &cobra.Command{
	Use:   "heart-rate [date]",
	Short: "Get heart rate data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessHeartRate,
}

var wellnessSpO2Cmd = &cobra.Command{
	Use:   "spo2 [date]",
	Short: "Get blood oxygen (SpO2) data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessSpO2,
}

var wellnessRespirationCmd = &cobra.Command{
	Use:   "respiration [date]",
	Short: "Get respiration data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessRespiration,
}

var wellnessIntensityCmd = &cobra.Command{
	Use:   "intensity [date]",
	Short: "Get intensity minutes data",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWellnessIntensity,
}

func init() {
	wellnessCmd.AddCommand(wellnessStressCmd)
	wellnessCmd.AddCommand(wellnessBodyBatteryCmd)
	wellnessCmd.AddCommand(wellnessHeartRateCmd)
	wellnessCmd.AddCommand(wellnessSpO2Cmd)
	wellnessCmd.AddCommand(wellnessRespirationCmd)
	wellnessCmd.AddCommand(wellnessIntensityCmd)
}

func runWellnessStress(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetDailyStress(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runWellnessBodyBattery(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetBodyBatteryEvents(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data.Events)
}

func runWellnessHeartRate(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetDailyHeartRate(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runWellnessSpO2(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetDailySpO2(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runWellnessRespiration(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetDailyRespiration(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}

func runWellnessIntensity(cmd *cobra.Command, args []string) error {
	date, err := parseDate(args)
	if err != nil {
		return err
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Wellness.GetDailyIntensityMinutes(cmd.Context(), date)
	if err != nil {
		return err
	}

	return printJSON(data)
}
