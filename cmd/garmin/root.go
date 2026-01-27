package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "garmin",
	Short: "Garmin Connect CLI",
	Long:  "A command-line interface for interacting with Garmin Connect API.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(sleepCmd)
	rootCmd.AddCommand(wellnessCmd)
	rootCmd.AddCommand(activitiesCmd)
	rootCmd.AddCommand(devicesCmd)
	rootCmd.AddCommand(hrvCmd)
	rootCmd.AddCommand(weightCmd)
	rootCmd.AddCommand(metricsCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(biometricCmd)
	rootCmd.AddCommand(workoutsCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(mcpCmd)
}
