package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

	// Weight
	s.AddTool(
		mcp.NewTool("get_weight",
			mcp.WithDescription("Get weight data for a date or date range"),
			mcp.WithString("date",
				mcp.Description("Single date in YYYY-MM-DD format (defaults to today)"),
			),
			mcp.WithString("start",
				mcp.Description("Start date for range query"),
			),
			mcp.WithString("end",
				mcp.Description("End date for range query"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start, _ := request.RequireString("start")
			end, _ := request.RequireString("end")
			if start != "" && end != "" {
				startDate, err := time.Parse("2006-01-02", start)
				if err != nil {
					return errorResult(fmt.Errorf("invalid start date: %w", err)), nil
				}
				endDate, err := time.Parse("2006-01-02", end)
				if err != nil {
					return errorResult(fmt.Errorf("invalid end date: %w", err)), nil
				}
				data, err := client.Weight.GetRange(ctx, startDate, endDate)
				if err != nil {
					return errorResult(err), nil
				}
				return jsonResult(data), nil
			}
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Weight.GetDaily(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// HRV
	s.AddTool(
		mcp.NewTool("get_hrv",
			mcp.WithDescription("Get heart rate variability data for a date or date range"),
			mcp.WithString("date",
				mcp.Description("Single date in YYYY-MM-DD format (defaults to today)"),
			),
			mcp.WithString("start",
				mcp.Description("Start date for range query"),
			),
			mcp.WithString("end",
				mcp.Description("End date for range query"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start, _ := request.RequireString("start")
			end, _ := request.RequireString("end")
			if start != "" && end != "" {
				startDate, err := time.Parse("2006-01-02", start)
				if err != nil {
					return errorResult(fmt.Errorf("invalid start date: %w", err)), nil
				}
				endDate, err := time.Parse("2006-01-02", end)
				if err != nil {
					return errorResult(fmt.Errorf("invalid end date: %w", err)), nil
				}
				data, err := client.HRV.GetRange(ctx, startDate, endDate)
				if err != nil {
					return errorResult(err), nil
				}
				return jsonResult(data), nil
			}
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.HRV.GetDaily(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Devices - List
	s.AddTool(
		mcp.NewTool("list_devices",
			mcp.WithDescription("List all registered Garmin devices"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := client.Devices.GetDevices(ctx)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Devices - Get Settings
	s.AddTool(
		mcp.NewTool("get_device_settings",
			mcp.WithDescription("Get settings for a specific device"),
			mcp.WithString("device_id",
				mcp.Required(),
				mcp.Description("The device ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, err := request.RequireString("device_id")
			if err != nil {
				return errorResult(err), nil
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorResult(fmt.Errorf("invalid device_id: %w", err)), nil
			}
			data, err := client.Devices.GetSettings(ctx, id)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Training Readiness
	s.AddTool(
		mcp.NewTool("get_training_readiness",
			mcp.WithDescription("Get daily training readiness score"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetTrainingReadiness(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Training Status
	s.AddTool(
		mcp.NewTool("get_training_status",
			mcp.WithDescription("Get current training status"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetTrainingStatusAggregated(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - VO2 Max
	s.AddTool(
		mcp.NewTool("get_vo2max",
			mcp.WithDescription("Get VO2 max estimates for running and cycling"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetMaxMetLatest(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Endurance Score
	s.AddTool(
		mcp.NewTool("get_endurance_score",
			mcp.WithDescription("Get endurance score"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetEnduranceScore(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Hill Score
	s.AddTool(
		mcp.NewTool("get_hill_score",
			mcp.WithDescription("Get hill/climb performance score"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetHillScore(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Training Load Balance
	s.AddTool(
		mcp.NewTool("get_training_load",
			mcp.WithDescription("Get training load balance (acute vs chronic load)"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetTrainingLoadBalance(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Metrics - Heat/Altitude Acclimation
	s.AddTool(
		mcp.NewTool("get_acclimation",
			mcp.WithDescription("Get heat and altitude acclimation status"),
			mcp.WithString("date",
				mcp.Description("Date in YYYY-MM-DD format (defaults to today)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			date, err := parseOptionalDate(request, "date")
			if err != nil {
				return errorResult(err), nil
			}
			data, err := client.Metrics.GetHeatAltitudeAcclimation(ctx, date)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Biometric - Heart Rate Zones
	s.AddTool(
		mcp.NewTool("get_heart_rate_zones",
			mcp.WithDescription("Get heart rate zone definitions"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := client.Biometric.GetHeartRateZones(ctx)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Biometric - FTP
	s.AddTool(
		mcp.NewTool("get_ftp",
			mcp.WithDescription("Get functional threshold power (FTP) for cycling"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := client.Biometric.GetCyclingFTP(ctx)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Workout - List
	s.AddTool(
		mcp.NewTool("list_workouts",
			mcp.WithDescription("List saved workouts"),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of workouts to return (default 20)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			limit := 20
			if l, err := request.RequireFloat("limit"); err == nil {
				limit = int(l)
			}
			data, err := client.Workouts.List(ctx, 0, limit)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Workout - Get
	s.AddTool(
		mcp.NewTool("get_workout",
			mcp.WithDescription("Get details of a specific workout"),
			mcp.WithString("workout_id",
				mcp.Required(),
				mcp.Description("The workout ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, err := request.RequireString("workout_id")
			if err != nil {
				return errorResult(err), nil
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorResult(fmt.Errorf("invalid workout_id: %w", err)), nil
			}
			data, err := client.Workouts.Get(ctx, id)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)

	// Workout - Delete
	s.AddTool(
		mcp.NewTool("delete_workout",
			mcp.WithDescription("Delete a workout"),
			mcp.WithString("workout_id",
				mcp.Required(),
				mcp.Description("The workout ID to delete"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, err := request.RequireString("workout_id")
			if err != nil {
				return errorResult(err), nil
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return errorResult(fmt.Errorf("invalid workout_id: %w", err)), nil
			}
			err = client.Workouts.Delete(ctx, id)
			if err != nil {
				return errorResult(err), nil
			}
			return mcp.NewToolResultText("Workout deleted successfully"), nil
		},
	)

	// User Profile
	s.AddTool(
		mcp.NewTool("get_profile",
			mcp.WithDescription("Get user profile information"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := client.UserProfile.GetSocialProfile(ctx)
			if err != nil {
				return errorResult(err), nil
			}
			return jsonResult(data), nil
		},
	)
}
