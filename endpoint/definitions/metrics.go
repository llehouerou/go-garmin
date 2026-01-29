package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// MetricsEndpoints defines all metrics-related endpoints.
var MetricsEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetTrainingReadiness",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/trainingreadiness/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get training readiness for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "readiness",
		MCPTool:       "get_training_readiness",
		Short:         "Get training readiness",
		Long:          "Get training readiness data including score, sleep, recovery time, and HRV factors",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetTrainingReadiness(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetEnduranceScore",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/endurancescore",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get endurance score for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "endurance",
		MCPTool:       "get_endurance_score",
		Short:         "Get endurance score",
		Long:          "Get endurance score data including overall score, classification, and contributors",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetEnduranceScore(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetHillScore",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/hillscore",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get hill score for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "hill",
		MCPTool:       "get_hill_score",
		Short:         "Get hill score",
		Long:          "Get hill score data including strength, endurance, and VO2 max",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetHillScore(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetMaxMetLatest",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/maxmet/latest/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get VO2 max for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "vo2max",
		MCPTool:       "get_vo2max",
		Short:         "Get latest VO2 max",
		Long:          "Get the latest VO2 max / MET data including generic and cycling values",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetMaxMetLatest(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetMaxMetDaily",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/maxmet/daily/{startDate}/{endDate}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for VO2 max data"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "vo2max-range",
		Short:         "Get VO2 max range",
		Long:          "Get VO2 max / MET data for a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.Metrics.GetMaxMetDaily(ctx, start, end)
		},
	},
	{
		Name:       "GetTrainingStatusAggregated",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/trainingstatus/aggregated/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get training status for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "training-status",
		MCPTool:       "get_training_status",
		Short:         "Get aggregated training status",
		Long:          "Get aggregated training status including VO2 max, load balance, and heat/altitude acclimation",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetTrainingStatusAggregated(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetTrainingStatusDaily",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/trainingstatus/daily/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get daily training status for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "training-status-daily",
		Short:         "Get daily training status",
		Long:          "Get daily training status data including weekly load and acute chronic workload ratio",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetTrainingStatusDaily(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetTrainingLoadBalance",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/trainingloadbalance/latest/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get training load balance for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "load-balance",
		MCPTool:       "get_training_load_balance",
		Short:         "Get training load balance",
		Long:          "Get training load balance data including aerobic and anaerobic load targets",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetTrainingLoadBalance(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetHeatAltitudeAcclimation",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/heataltitudeacclimation/latest/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get heat/altitude acclimation for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "acclimation",
		MCPTool:       "get_heat_altitude_acclimation",
		Short:         "Get heat/altitude acclimation",
		Long:          "Get heat and altitude acclimation data including percentages and trends",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Metrics.GetHeatAltitudeAcclimation(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetRacePredictions",
		Service:    "Metrics",
		Cassette:   "metrics",
		Path:       "/metrics-service/metrics/racepredictions/latest/{displayName}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "display_name", Type: endpoint.ParamTypeString, Required: false, Description: "User display name (defaults to current user)"},
		},
		CLICommand:    "metrics",
		CLISubcommand: "race-predictions",
		MCPTool:       "get_race_predictions",
		Short:         "Get race predictions",
		Long:          "Get predicted race times for 5K, 10K, half marathon, and marathon based on current fitness",
		DependsOn:     "GetSocialProfile",
		ArgProvider: func(result any) map[string]any {
			profile, ok := result.(*garmin.SocialProfile)
			if !ok || profile == nil {
				return nil
			}
			return map[string]any{"display_name": profile.DisplayName}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			displayName := args.String("display_name")
			if displayName == "" {
				// Auto-fetch display name from current user's social profile
				profile, err := client.UserProfile.GetSocialProfile(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to get display name: %w", err)
				}
				displayName = profile.DisplayName
			}
			return client.Metrics.GetRacePredictionsLatest(ctx, displayName)
		},
	},
}
