package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/llehouerou/go-garmin"
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

	registerTools(s, client)

	return server.ServeStdio(s)
}

func registerTools(s *server.MCPServer, client *garmin.Client) {
	// Sleep
	s.AddTool(
		mcp.NewTool("get_sleep",
			mcp.WithDescription("Get sleep data for a specific date including duration, stages, and sleep score"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Sleep.GetDaily(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)
}
