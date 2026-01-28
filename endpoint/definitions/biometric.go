package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// BiometricEndpoints defines all biometric-related endpoints.
var BiometricEndpoints = []endpoint.Endpoint{
	{
		Name:          "GetLatestLactateThreshold",
		Service:       "Biometric",
		Cassette:      "biometric",
		Path:          "/biometric-service/biometric/latestLactateThreshold",
		HTTPMethod:    "GET",
		CLICommand:    "biometric",
		CLISubcommand: "lactate-threshold",
		MCPTool:       "get_lactate_threshold",
		Short:         "Get latest lactate threshold",
		Long:          "Get the latest lactate threshold data including speed and heart rate",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Biometric.GetLatestLactateThreshold(ctx)
		},
	},
	{
		Name:          "GetCyclingFTP",
		Service:       "Biometric",
		Cassette:      "biometric",
		Path:          "/biometric-service/biometric/latestFunctionalThresholdPower/CYCLING",
		HTTPMethod:    "GET",
		CLICommand:    "biometric",
		CLISubcommand: "ftp",
		MCPTool:       "get_cycling_ftp",
		Short:         "Get cycling FTP",
		Long:          "Get the latest cycling Functional Threshold Power (FTP) in watts",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Biometric.GetCyclingFTP(ctx)
		},
	},
	{
		Name:       "GetPowerToWeight",
		Service:    "Biometric",
		Cassette:   "biometric",
		Path:       "/biometric-service/biometric/powerToWeight/latest/{date}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get power-to-weight for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "biometric",
		CLISubcommand: "power-weight",
		MCPTool:       "get_power_to_weight",
		Short:         "Get power-to-weight ratio",
		Long:          "Get the power-to-weight ratio for running on a specific date",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Biometric.GetPowerToWeight(ctx, args.Date("date"))
		},
	},
	{
		Name:          "GetHeartRateZones",
		Service:       "Biometric",
		Cassette:      "biometric",
		Path:          "/biometric-service/heartRateZones/",
		HTTPMethod:    "GET",
		CLICommand:    "biometric",
		CLISubcommand: "hr-zones",
		MCPTool:       "get_heart_rate_zones",
		Short:         "Get heart rate zones",
		Long:          "Get heart rate zone configurations for all sports",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Biometric.GetHeartRateZones(ctx)
		},
	},
	{
		Name:       "GetLactateThresholdSpeedRange",
		Service:    "Biometric",
		Cassette:   "biometric",
		Path:       "/biometric-service/stats/lactateThresholdSpeed/range/{startDate}/{endDate}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for lactate threshold speed stats"},
		},
		CLICommand:    "biometric",
		CLISubcommand: "lt-speed-range",
		Short:         "Get lactate threshold speed range",
		Long:          "Get lactate threshold speed statistics for a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.Biometric.GetLactateThresholdSpeedRange(ctx, start, end)
		},
	},
	{
		Name:       "GetLactateThresholdHRRange",
		Service:    "Biometric",
		Cassette:   "biometric",
		Path:       "/biometric-service/stats/lactateThresholdHeartRate/range/{startDate}/{endDate}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for lactate threshold HR stats"},
		},
		CLICommand:    "biometric",
		CLISubcommand: "lt-hr-range",
		Short:         "Get lactate threshold HR range",
		Long:          "Get lactate threshold heart rate statistics for a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.Biometric.GetLactateThresholdHRRange(ctx, start, end)
		},
	},
	{
		Name:       "GetFTPRange",
		Service:    "Biometric",
		Cassette:   "biometric",
		Path:       "/biometric-service/stats/functionalThresholdPower/range/{startDate}/{endDate}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for FTP stats"},
		},
		CLICommand:    "biometric",
		CLISubcommand: "ftp-range",
		Short:         "Get FTP range",
		Long:          "Get Functional Threshold Power statistics for a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.Biometric.GetFTPRange(ctx, start, end)
		},
	},
}
