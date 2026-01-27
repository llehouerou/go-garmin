// endpoint/definitions/sleep.go
package definitions

import (
	"context"
	"fmt"

	garmin "github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// SleepEndpoints defines all sleep-related API endpoints.
var SleepEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetDailySleep",
		Service:    "Sleep",
		Cassette:   "sleep_daily",
		Path:       "/sleep-service/sleep/dailySleepData",
		HTTPMethod: "GET",

		Params: []endpoint.Param{
			{
				Name:        "date",
				Type:        endpoint.ParamTypeDate,
				Required:    false,
				Description: "Date to get sleep data for (YYYY-MM-DD, defaults to today)",
			},
		},

		CLICommand: "sleep",
		MCPTool:    "get_sleep",
		Short:      "Get sleep data for a date",
		Long:       "Get sleep data including duration, stages (deep, light, REM), and sleep score",

		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Sleep.GetDaily(ctx, args.Date("date"))
		},
	},
}
