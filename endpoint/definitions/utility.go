package definitions

import (
	"context"
	"time"

	"github.com/llehouerou/go-garmin/endpoint"
)

// UtilityEndpoints defines utility endpoints that don't call the Garmin API.
var UtilityEndpoints = []endpoint.Endpoint{
	{
		Name:    "GetCurrentDate",
		Service: "Utility",
		MCPTool: "get_current_date",
		Short:   "Get current date",
		Long:    "Get the current date including year, month, day, and weekday. Useful for determining what date to use for other API calls.",
		Handler: func(_ context.Context, _ any, _ *endpoint.HandlerArgs) (any, error) {
			now := time.Now()
			return map[string]any{
				"date":       now.Format("2006-01-02"),
				"year":       now.Year(),
				"month":      int(now.Month()),
				"month_name": now.Month().String(),
				"day":        now.Day(),
				"weekday":    now.Weekday().String(),
				"iso8601":    now.Format(time.RFC3339),
			}, nil
		},
	},
}
