package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/llehouerou/go-garmin/endpoint"
)

var (
	rootCmd = &cobra.Command{
		Use:   "garmin",
		Short: "Garmin Connect CLI",
		Long:  "A command-line interface for interacting with Garmin Connect API.",
	}

	cliGenerator *endpoint.CLIGenerator
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Administrative commands (no client needed for help)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(mcpCmd)

	// Generate data commands from endpoint registry
	cliGenerator = endpoint.NewCLIGenerator(endpointRegistry)
	for _, cmd := range cliGenerator.GenerateCommands() {
		// Add PersistentPreRunE to load client before running data commands
		cmd.PersistentPreRunE = loadClientForCLI
		rootCmd.AddCommand(cmd)
	}
}

// loadClientForCLI loads the Garmin client and sets it on the CLI generator.
func loadClientForCLI(_ *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}
	cliGenerator.SetClient(client)
	return nil
}
