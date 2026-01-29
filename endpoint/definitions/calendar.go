package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// CalendarEndpoints defines all calendar-related endpoints.
var CalendarEndpoints = []endpoint.Endpoint{
	{
		Name:       "GetCalendar",
		Service:    "Calendar",
		Cassette:   "calendar",
		Path:       "/calendar-service/year/{year}/month/{month}/day/{day}/start/{start}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "year", Type: endpoint.ParamTypeInt, Required: true, Description: "Year (e.g., 2026)"},
			{Name: "month", Type: endpoint.ParamTypeInt, Required: false, Description: "Month (0-11, January=0)"},
			{Name: "day", Type: endpoint.ParamTypeInt, Required: false, Description: "Day of month (requires month and start)"},
			{Name: "start", Type: endpoint.ParamTypeInt, Required: false, Description: "Week start day, 1=Monday (required when day is provided)"},
		},
		CLICommand:    "calendar",
		CLISubcommand: "get",
		MCPTool:       "get_calendar",
		Short:         "Get calendar",
		Long:          "Get calendar items including activities, workouts, and weight entries. Parameters are hierarchical: month requires year, day requires both month and start.",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			var opts *garmin.CalendarOptions
			if args.HasParam("month") {
				opts = &garmin.CalendarOptions{}
				month := args.Int("month")
				opts.Month = &month
				if args.HasParam("day") {
					day := args.Int("day")
					opts.Day = &day
					if args.HasParam("start") {
						start := args.Int("start")
						opts.Start = &start
					}
				}
			}
			return client.Calendar.Get(ctx, args.Int("year"), opts)
		},
	},
}
