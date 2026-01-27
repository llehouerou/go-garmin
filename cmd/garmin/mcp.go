package main

import (
	"context"
	"fmt"
	"strconv"

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

	// Wellness - Stress
	s.AddTool(
		mcp.NewTool("get_stress",
			mcp.WithDescription("Get stress levels throughout the day"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetDailyStress(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Wellness - Body Battery
	s.AddTool(
		mcp.NewTool("get_body_battery",
			mcp.WithDescription("Get body battery drain and charge events throughout the day"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetBodyBatteryEvents(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Wellness - Heart Rate
	s.AddTool(
		mcp.NewTool("get_heart_rate",
			mcp.WithDescription("Get heart rate data for a day including resting HR and time in zones"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetDailyHeartRate(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Wellness - SpO2
	s.AddTool(
		mcp.NewTool("get_spo2",
			mcp.WithDescription("Get blood oxygen (SpO2) readings for a day"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetDailySpO2(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Wellness - Respiration
	s.AddTool(
		mcp.NewTool("get_respiration",
			mcp.WithDescription("Get respiration rate data for a day"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetDailyRespiration(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Wellness - Intensity Minutes
	s.AddTool(
		mcp.NewTool("get_intensity_minutes",
			mcp.WithDescription("Get weekly intensity minutes (moderate and vigorous activity)"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Wellness.GetDailyIntensityMinutes(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Activity - List
	s.AddTool(
		mcp.NewTool("list_activities",
			mcp.WithDescription("List activities with optional filters"),
			mcp.WithNumber("start",
				mcp.Description("Starting index (0-based, default 0)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of activities to return (default 20)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			opts := &garmin.ListOptions{}
			if start, err := request.RequireFloat("start"); err == nil {
				opts.Start = int(start)
			}
			if limit, err := request.RequireFloat("limit"); err == nil {
				opts.Limit = int(limit)
			}
			data, err := client.Activities.List(ctx, opts)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Activity - Get
	s.AddTool(
		mcp.NewTool("get_activity",
			mcp.WithDescription("Get detailed information about a specific activity"),
			mcp.WithString("activity_id",
				mcp.Required(),
				mcp.Description("The activity ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, err := request.RequireString("activity_id")
			if err != nil {
				return errorResult(err), nil
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorResult(fmt.Errorf("invalid activity_id: %w", err)), nil
			}
			data, err := client.Activities.Get(ctx, id)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Activity - Get Splits
	s.AddTool(
		mcp.NewTool("get_activity_splits",
			mcp.WithDescription("Get splits/laps for an activity"),
			mcp.WithString("activity_id",
				mcp.Required(),
				mcp.Description("The activity ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, err := request.RequireString("activity_id")
			if err != nil {
				return errorResult(err), nil
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorResult(fmt.Errorf("invalid activity_id: %w", err)), nil
			}
			data, err := client.Activities.GetSplits(ctx, id)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)
}
