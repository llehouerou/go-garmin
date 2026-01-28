package main

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/llehouerou/go-garmin/endpoint"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server for LLM integration",
	Long:  "Start a Model Context Protocol server that exposes Garmin data to LLM assistants like Claude.",
	RunE:  runMCP,
}

func runMCP(_ *cobra.Command, _ []string) error {
	client, err := loadClient()
	if err != nil {
		return err
	}

	s := server.NewMCPServer(
		"garmin",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all tools from the endpoint registry
	mcpGen := endpoint.NewMCPGenerator(endpointRegistry, client)
	mcpGen.RegisterTools(s)

	return server.ServeStdio(s)
}
