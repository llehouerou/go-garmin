package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// WeightEndpoints defines all weight-related endpoints.
var WeightEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetDailyWeight",
		Service:    "Weight",
		Cassette:   "weight",
		Path:       "/weight-service/weight/dayview",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date to get weight data for (YYYY-MM-DD, defaults to today)"},
		},
		CLICommand:    "weight",
		CLISubcommand: "daily",
		MCPTool:       "get_weight",
		Short:         "Get weight data for a date",
		Long:          "Get weight data for a date including BMI, body fat, muscle mass, and other body composition metrics",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Weight.GetDaily(ctx, args.Date("date"))
		},
	},
	{
		Name:       "GetWeightRange",
		Service:    "Weight",
		Cassette:   "weight",
		Path:       "/weight-service/weight/range",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for weight data"},
		},
		CLICommand:    "weight",
		CLISubcommand: "range",
		Short:         "Get weight data for a date range",
		Long:          "Get weight data summaries for a date range including averages",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.Weight.GetRange(ctx, start, end)
		},
	},
}
