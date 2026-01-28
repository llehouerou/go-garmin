package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// HRVEndpoints defines all HRV-related endpoints.
var HRVEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetDailyHRV",
		Service:    "HRV",
		Cassette:   "hrv",
		Path:       "/hrv-service/hrv",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get HRV data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "hrv",
		CLISubcommand: "daily",
		MCPTool:       "get_hrv",
		Short:         "Get HRV data for a date",
		Long:          "Get heart rate variability data for a date including weekly average, last night average, and baseline",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.HRV.GetDaily(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetHRVRange",
		Service:    "HRV",
		Cassette:   "hrv",
		Path:       "/hrv-service/hrv/daily",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for HRV data"},
		},
		CLICommand:    "hrv",
		CLISubcommand: "range",
		Short:         "Get HRV data for a date range",
		Long:          "Get heart rate variability summaries for a date range",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.HRV.GetRange(ctx, start, end)
		},
	},
}
