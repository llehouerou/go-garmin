package definitions

import (
	"context"
	"fmt"

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
			{Name: "range", Type: endpoint.ParamTypeDateRange, Required: false, Description: "Date range for fitness age data"},
		},
		CLICommand:    "fitnessage",
		CLISubcommand: "stats",
		MCPTool:       "get_fitness_age_stats",
		Short:         "Get fitness age statistics",
		Long:          "Get daily fitness age statistics including fitness age, achievable fitness age, RHR, BMI, and vigorous activity days",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Date("start")
			end := args.Date("end")
			return client.FitnessAge.GetStatsDaily(ctx, start, end)
		},
	},
}
