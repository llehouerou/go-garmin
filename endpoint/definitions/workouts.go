package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// WorkoutEndpoints defines all workout-related endpoints.
var WorkoutEndpoints = []endpoint.Endpoint{
	{
		Name:       "ListWorkouts",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/workouts",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "start", Type: endpoint.ParamTypeInt, Required: false, Description: "Starting index (0-based, defaults to 0)"},
			{Name: "limit", Type: endpoint.ParamTypeInt, Required: false, Description: "Maximum number of workouts to return (defaults to 20)"},
		},
		CLICommand:    "workouts",
		CLISubcommand: "list",
		MCPTool:       "list_workouts",
		Short:         "List workouts",
		Long:          "List workouts with pagination including name, sport type, and estimated duration",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			start := args.Int("start")
			limit := args.Int("limit")
			if limit == 0 {
				limit = 20
			}
			return client.Workouts.List(ctx, start, limit)
		},
	},
	{
		Name:       "GetWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/workout/{workoutId}",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "workout_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The workout ID"},
		},
		CLICommand:    "workouts",
		CLISubcommand: "get",
		MCPTool:       "get_workout",
		Short:         "Get workout details",
		Long:          "Get detailed information about a specific workout including segments and steps",
		DependsOn:     "ListWorkouts",
		ArgProvider: func(result any) map[string]any {
			list, ok := result.(*garmin.WorkoutList)
			if !ok || len(list.Workouts) == 0 {
				return nil
			}
			return map[string]any{"workout_id": list.Workouts[0].WorkoutID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Workouts.Get(ctx, int64(args.Int("workout_id")))
		},
	},
	{
		Name:       "ScheduleWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/schedule/{workoutId}",
		HTTPMethod: "POST",
		Params: []endpoint.Param{
			{Name: "workout_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The workout ID to schedule"},
			{Name: "date", Type: endpoint.ParamTypeDate, Required: true, Description: "Date to schedule the workout (YYYY-MM-DD)"},
		},
		CLICommand:    "workouts",
		CLISubcommand: "schedule",
		MCPTool:       "schedule_workout",
		Short:         "Schedule a workout",
		Long:          "Schedule a workout for a specific date",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Workouts.Schedule(ctx, int64(args.Int("workout_id")), args.Date("date"))
		},
	},
	{
		Name:       "UnscheduleWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/schedule/{scheduleId}",
		HTTPMethod: "DELETE",
		Params: []endpoint.Param{
			{Name: "schedule_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The schedule ID to remove"},
		},
		CLICommand:    "workouts",
		CLISubcommand: "unschedule",
		MCPTool:       "unschedule_workout",
		Short:         "Unschedule a workout",
		Long:          "Remove a scheduled workout by its schedule ID",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			err := client.Workouts.Unschedule(ctx, int64(args.Int("schedule_id")))
			if err != nil {
				return nil, err
			}
			return map[string]string{"status": "success"}, nil
		},
	},
}
