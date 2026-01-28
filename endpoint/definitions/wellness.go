package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// WellnessEndpoints defines all wellness-related endpoints.
var WellnessEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetDailyStress",
		Service:    "Wellness",
		Cassette:   "wellness_stress",
		Path:       "/wellness-service/wellness/dailyStress",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get stress data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "stress",
		MCPTool:       "get_stress",
		Short:         "Get stress levels for a date",
		Long:          "Get stress levels throughout the day including max, average, and stress chart values",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetDailyStress(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetBodyBatteryEvents",
		Service:    "Wellness",
		Cassette:   "wellness_body_battery",
		Path:       "/wellness-service/wellness/bodyBattery/events",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get body battery events for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "body-battery",
		MCPTool:       "get_body_battery",
		Short:         "Get body battery events for a date",
		Long:          "Get body battery drain and charge events throughout the day including sleep, activity, and stress impacts",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetBodyBatteryEvents(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetDailyHeartRate",
		Service:    "Wellness",
		Cassette:   "wellness_heart_rate",
		Path:       "/wellness-service/wellness/dailyHeartRate",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get heart rate data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "heart-rate",
		MCPTool:       "get_heart_rate",
		Short:         "Get heart rate data for a date",
		Long:          "Get heart rate data for a day including resting HR, max HR, and time in zones",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetDailyHeartRate(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetDailySpO2",
		Service:    "Wellness",
		Cassette:   "wellness_extended",
		Path:       "/wellness-service/wellness/daily/spo2",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get SpO2 data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "spo2",
		MCPTool:       "get_spo2",
		Short:         "Get blood oxygen (SpO2) for a date",
		Long:          "Get blood oxygen (SpO2) readings for a day including average, lowest, and sleep SpO2",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetDailySpO2(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetDailyRespiration",
		Service:    "Wellness",
		Cassette:   "wellness_extended",
		Path:       "/wellness-service/wellness/daily/respiration",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get respiration data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "respiration",
		MCPTool:       "get_respiration",
		Short:         "Get respiration data for a date",
		Long:          "Get respiration rate data for a day including waking and sleep respiration averages",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetDailyRespiration(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetDailyIntensityMinutes",
		Service:    "Wellness",
		Cassette:   "wellness_extended",
		Path:       "/wellness-service/wellness/daily/im",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get intensity minutes for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "wellness",
		CLISubcommand: "intensity-minutes",
		MCPTool:       "get_intensity_minutes",
		Short:         "Get intensity minutes for a date",
		Long:          "Get weekly intensity minutes (moderate and vigorous activity) and progress toward weekly goal",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Wellness.GetDailyIntensityMinutes(ctx, args.Date("date"))
		},
	},
}
