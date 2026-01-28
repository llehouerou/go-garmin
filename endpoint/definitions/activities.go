package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// ActivityEndpoints defines all activity-related endpoints.
var ActivityEndpoints = []endpoint.Endpoint{
	{
		Name:       "ListActivities",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activitylist-service/activities/search/activities",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "start", Type: endpoint.ParamTypeInt, Required: false, Description: "Starting index (0-based, defaults to 0)"},
			{Name: "limit", Type: endpoint.ParamTypeInt, Required: false, Description: "Maximum number of activities to return (defaults to 20)"},
		},
		CLICommand:    "activities",
		CLISubcommand: "list",
		MCPTool:       "list_activities",
		Short:         "List activities",
		Long:          "List activities with pagination including distance, duration, heart rate, and other metrics",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			opts := &garmin.ListOptions{
				Start: args.Int("start"),
				Limit: args.Int("limit"),
			}
			if opts.Limit == 0 {
				opts.Limit = 20
			}
			return client.Activities.List(ctx, opts)
		},
	},
	{
		Name:       "GetActivity",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "get",
		MCPTool:       "get_activity",
		Short:         "Get activity details",
		Long:          "Get detailed information about a specific activity including metadata, summary, and splits",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.Get(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivityWeather",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/weather",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "weather",
		MCPTool:       "get_activity_weather",
		Short:         "Get activity weather",
		Long:          "Get weather data for a specific activity including temperature, humidity, and wind",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetWeather(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivitySplits",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/splits",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "splits",
		MCPTool:       "get_activity_splits",
		Short:         "Get activity splits",
		Long:          "Get splits/laps data for a specific activity including pace, heart rate, and elevation per lap",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetSplits(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivityDetails",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/details",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "details",
		MCPTool:       "get_activity_details",
		Short:         "Get activity time-series details",
		Long:          "Get extended details with time-series metrics for an activity",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetDetails(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivityHRTimeInZones",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/hrTimeInZones",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "hr-zones",
		MCPTool:       "get_activity_hr_zones",
		Short:         "Get activity HR time in zones",
		Long:          "Get heart rate time in zones for an activity",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetHRTimeInZones(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivityPowerTimeInZones",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/powerTimeInZones",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "power-zones",
		MCPTool:       "get_activity_power_zones",
		Short:         "Get activity power time in zones",
		Long:          "Get power time in zones for an activity",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetPowerTimeInZones(ctx, int64(args.Int("activity_id")))
		},
	},
	{
		Name:       "GetActivityExerciseSets",
		Service:    "Activities",
		Cassette:   "activities",
		Path:       "/activity-service/activity/{activityId}/exerciseSets",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "activity_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The activity ID"},
		},
		CLICommand:    "activities",
		CLISubcommand: "exercise-sets",
		MCPTool:       "get_activity_exercise_sets",
		Short:         "Get activity exercise sets",
		Long:          "Get exercise sets for a strength workout activity",
		DependsOn:     "ListActivities",
		ArgProvider: func(result any) map[string]any {
			activities, ok := result.([]garmin.Activity)
			if !ok || len(activities) == 0 {
				return nil
			}
			return map[string]any{"activity_id": activities[0].ActivityID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Activities.GetExerciseSets(ctx, int64(args.Int("activity_id")))
		},
	},
}
