package definitions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// FitnessAgeEndpoints defines all fitness age-related endpoints.
var FitnessAgeEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetFitnessAgeStats",
		Service:    "FitnessAge",
		Cassette:   "fitnessage",
		Path:       "/fitnessage-service/stats/daily/{startDate}/{endDate}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for fitness age data (default: last 7 days, max 28 days)"},
		},
		CLICommand:    "fitnessage",
		CLISubcommand: "stats",
		MCPTool:       "get_fitness_age_stats",
		Short:         "Get fitness age statistics",
		Long:          "Get daily fitness age statistics including fitness age, achievable fitness age, RHR, BMI, and vigorous activity days. Note: date range must be 28 days or less.",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}

			// Default: last 7 days
			end := time.Now()
			start := end.AddDate(0, 0, -7)
			if args.HasParam("end") {
				end = args.Date("end")
			}
			if args.HasParam("start") {
				start = args.Date("start")
			}

			// Validate: max 28 days
			days := int(end.Sub(start).Hours() / 24)
			if days > 28 {
				return nil, fmt.Errorf("date range must be 28 days or less (got %d days)", days)
			}
			if days < 1 {
				return nil, errors.New("date range must be at least 1 day")
			}

			return client.FitnessAge.GetStatsDaily(ctx, start, end)
		},
	},
}
