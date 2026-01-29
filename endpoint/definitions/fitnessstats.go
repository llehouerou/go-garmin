package definitions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// parseMetrics parses a comma-separated string of metrics into a slice.
func parseMetrics(metricsStr string) ([]garmin.FitnessMetric, error) {
	var metrics []garmin.FitnessMetric
	for m := range strings.SplitSeq(metricsStr, ",") {
		m = strings.TrimSpace(m)
		switch m {
		case "calories":
			metrics = append(metrics, garmin.MetricCalories)
		case "distance":
			metrics = append(metrics, garmin.MetricDistance)
		case "duration":
			metrics = append(metrics, garmin.MetricDuration)
		case "avgSpeed":
			metrics = append(metrics, garmin.MetricAvgSpeed)
		case "maxHr":
			metrics = append(metrics, garmin.MetricMaxHR)
		case "avgHr":
			metrics = append(metrics, garmin.MetricAvgHR)
		case "elevationGain":
			metrics = append(metrics, garmin.MetricElevationGain)
		case "avgRunCadence":
			metrics = append(metrics, garmin.MetricAvgRunCadence)
		case "avgGroundContactBalance":
			metrics = append(metrics, garmin.MetricAvgGroundContactBalance)
		case "avgStrideLength":
			metrics = append(metrics, garmin.MetricAvgStrideLength)
		case "avgVerticalOscillation":
			metrics = append(metrics, garmin.MetricAvgVerticalOscillation)
		case "avgVerticalRatio":
			metrics = append(metrics, garmin.MetricAvgVerticalRatio)
		case "avgGroundContactTime":
			metrics = append(metrics, garmin.MetricAvgGroundContactTime)
		case "startLocal":
			metrics = append(metrics, garmin.MetricStartLocal)
		case "activityType":
			metrics = append(metrics, garmin.MetricActivityType)
		case "activitySubType":
			metrics = append(metrics, garmin.MetricActivitySubType)
		case "name":
			metrics = append(metrics, garmin.MetricName)
		case "aerobicTrainingEffect":
			metrics = append(metrics, garmin.MetricAerobicTrainingEffect)
		case "anaerobicTrainingEffect":
			metrics = append(metrics, garmin.MetricAnaerobicTrainingEffect)
		default:
			return nil, fmt.Errorf("invalid metric: %s", m)
		}
	}
	return metrics, nil
}

// FitnessStatsEndpoints defines all fitness stats-related endpoints.
var FitnessStatsEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetFitnessStats",
		Service:    "FitnessStats",
		Cassette:   "fitnessstats",
		Path:       "/fitnessstats-service/activity",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for fitness stats (default: last 3 months)"},
			{Name: "aggregation", Type: endpoint.ParamTypeString, Required: false, Description: "Aggregation period: daily, weekly, monthly, yearly (default: weekly)"},
			{Name: "metrics", Type: endpoint.ParamTypeString, Required: false, Description: "Comma-separated metrics: calories, distance, duration, avgSpeed, maxHr, avgHr, elevationGain, avgRunCadence, avgGroundContactBalance, avgStrideLength, avgVerticalOscillation, avgVerticalRatio, avgGroundContactTime, aerobicTrainingEffect, anaerobicTrainingEffect (default: calories,distance,duration)"},
			{Name: "group_by_activity_type", Type: endpoint.ParamTypeBool, Required: false, Description: "Group stats by activity type (e.g., running, hiking)"},
			{Name: "standardized_units", Type: endpoint.ParamTypeBool, Required: false, Description: "Use standardized units in response"},
		},
		CLICommand:    "fitnessstats",
		CLISubcommand: "get",
		MCPTool:       "get_fitness_stats",
		Short:         "Get fitness statistics",
		Long:          "Get aggregated fitness statistics for activities including calories, distance, and duration over a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}

			// Parse aggregation
			aggregation := garmin.AggregationWeekly
			if agg := args.String("aggregation"); agg != "" {
				switch agg {
				case "daily":
					aggregation = garmin.AggregationDaily
				case "weekly":
					aggregation = garmin.AggregationWeekly
				case "monthly":
					aggregation = garmin.AggregationMonthly
				case "yearly":
					aggregation = garmin.AggregationYearly
				default:
					return nil, fmt.Errorf("invalid aggregation: %s (valid: daily, weekly, monthly, yearly)", agg)
				}
			}

			// Parse metrics
			metricsStr := args.String("metrics")
			if metricsStr == "" {
				metricsStr = "calories,distance,duration"
			}
			metrics, err := parseMetrics(metricsStr)
			if err != nil {
				return nil, err
			}

			// Default date range: last 3 months
			endDate := time.Now()
			startDate := endDate.AddDate(0, -3, 0)
			if args.HasParam("end") {
				endDate = args.Date("end")
			}
			if args.HasParam("start") {
				startDate = args.Date("start")
			}

			opts := &garmin.FitnessStatsOptions{
				Aggregation:         aggregation,
				StartDate:           startDate,
				EndDate:             endDate,
				Metrics:             metrics,
				GroupByActivityType: args.Bool("group_by_activity_type"),
				StandardizedUnits:   args.Bool("standardized_units"),
			}

			return client.FitnessStats.GetActivityStats(ctx, opts)
		},
	},
	{
		Name:       "GetFitnessStatsActivities",
		Service:    "FitnessStats",
		Cassette:   "fitnessstats",
		Path:       "/fitnessstats-service/activity/all",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for activities (default: last 3 months)"},
			{Name: "activity_type", Type: endpoint.ParamTypeString, Required: false, Description: "Filter by activity type (e.g., running, hiking, cycling)"},
			{Name: "metrics", Type: endpoint.ParamTypeString, Required: false, Description: "Comma-separated metrics: calories, distance, duration, avgSpeed, maxHr, avgHr, elevationGain, avgRunCadence, avgGroundContactBalance, avgStrideLength, avgVerticalOscillation, avgVerticalRatio, avgGroundContactTime, startLocal, activityType, activitySubType, name, aerobicTrainingEffect, anaerobicTrainingEffect (default: name,startLocal,activityType,duration,distance,calories)"},
		},
		CLICommand:    "fitnessstats",
		CLISubcommand: "activities",
		MCPTool:       "get_fitness_stats_activities",
		Short:         "Get individual activity data",
		Long:          "Get individual activity data without aggregation, including activity names, types, and training effects",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}

			// Parse metrics
			metricsStr := args.String("metrics")
			if metricsStr == "" {
				metricsStr = "name,startLocal,activityType,duration,distance,calories"
			}
			metrics, err := parseMetrics(metricsStr)
			if err != nil {
				return nil, err
			}

			// Default date range: last 3 months
			endDate := time.Now()
			startDate := endDate.AddDate(0, -3, 0)
			if args.HasParam("end") {
				endDate = args.Date("end")
			}
			if args.HasParam("start") {
				startDate = args.Date("start")
			}

			opts := &garmin.FitnessStatsOptions{
				StartDate:    startDate,
				EndDate:      endDate,
				Metrics:      metrics,
				ActivityType: args.String("activity_type"),
			}

			return client.FitnessStats.GetAllActivities(ctx, opts)
		},
	},
}
