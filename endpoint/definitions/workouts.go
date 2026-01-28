package definitions

import (
	"context"
	"fmt"
	"reflect"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// workoutBodyConfig provides documentation for the workout JSON structure.
var workoutBodyConfig = &endpoint.BodyConfig{
	Type: reflect.TypeFor[garmin.Workout](),
	Description: `JSON object representing a workout. Required fields: workoutName, sportType, workoutSegments.

Structure:
- workoutName (string, required): Name of the workout
- description (string): Workout description
- sportType (object, required): Sport type with sportTypeId and sportTypeKey
  - sportTypeId: 1=running, 2=cycling, 5=swimming, etc.
  - sportTypeKey: "running", "cycling", "lap_swimming", etc.
- workoutSegments (array, required): Array of workout segments, each containing:
  - segmentOrder (int): Order of the segment (1-based)
  - sportType (object): Same as above
  - workoutSteps (array): Array of workout steps

WorkoutStep types:
- "ExecutableStepDTO": A single exercise step
- "RepeatGroupDTO": A repeat group containing nested steps

Step fields:
- stepOrder (int): Order within segment (1-based)
- stepType: warmup, interval, recovery, rest, cooldown, other
- endCondition: time, distance, calories, heart.rate, iterations, lap.button
- endConditionValue: Value for the condition (seconds for time, meters for distance)
- targetType: no.target, heart.rate.zone, cadence.zone, speed.zone, pace.zone, power.zone
- zoneNumber: Zone number (1-5) when using zone targets

For repeat groups (type="RepeatGroupDTO"):
- numberOfIterations: Number of times to repeat
- workoutSteps: Nested array of steps to repeat`,
	Example: `{
  "workoutName": "Easy 30min Run",
  "description": "Easy aerobic run in Zone 2",
  "sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
  "workoutSegments": [{
    "segmentOrder": 1,
    "sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
    "workoutSteps": [
      {
        "type": "ExecutableStepDTO",
        "stepOrder": 1,
        "stepType": {"stepTypeId": 1, "stepTypeKey": "warmup"},
        "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
        "endConditionValue": 300,
        "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"}
      },
      {
        "type": "ExecutableStepDTO",
        "stepOrder": 2,
        "stepType": {"stepTypeId": 3, "stepTypeKey": "interval"},
        "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
        "endConditionValue": 1200,
        "targetType": {"workoutTargetTypeId": 4, "workoutTargetTypeKey": "heart.rate.zone"},
        "zoneNumber": 2
      },
      {
        "type": "ExecutableStepDTO",
        "stepOrder": 3,
        "stepType": {"stepTypeId": 2, "stepTypeKey": "cooldown"},
        "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
        "endConditionValue": 300,
        "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"}
      }
    ]
  }]
}`,
}

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
	{
		Name:       "CreateWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/workout",
		HTTPMethod: "POST",
		Body:       workoutBodyConfig,
		MCPTool:    "create_workout",
		Short:      "Create a new workout",
		Long:       "Create a new workout with segments and steps. The workout will be saved to your Garmin account.",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			workout, ok := args.Body.(*garmin.Workout)
			if !ok {
				return nil, fmt.Errorf("invalid workout body type: %T", args.Body)
			}
			return client.Workouts.Create(ctx, workout)
		},
	},
	{
		Name:       "UpdateWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/workout/{workoutId}",
		HTTPMethod: "PUT",
		Params: []endpoint.Param{
			{Name: "workout_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The workout ID to update"},
		},
		Body:    workoutBodyConfig,
		MCPTool: "update_workout",
		Short:   "Update an existing workout",
		Long:    "Update an existing workout. You must provide the complete workout structure including all segments and steps.",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			workout, ok := args.Body.(*garmin.Workout)
			if !ok {
				return nil, fmt.Errorf("invalid workout body type: %T", args.Body)
			}
			return client.Workouts.Update(ctx, int64(args.Int("workout_id")), workout)
		},
	},
	{
		Name:       "DeleteWorkout",
		Service:    "Workouts",
		Cassette:   "workouts",
		Path:       "/workout-service/workout/{workoutId}",
		HTTPMethod: "DELETE",
		Params: []endpoint.Param{
			{Name: "workout_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The workout ID to delete"},
		},
		MCPTool: "delete_workout",
		Short:   "Delete a workout",
		Long:    "Permanently delete a workout from your Garmin account",
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			err := client.Workouts.Delete(ctx, int64(args.Int("workout_id")))
			if err != nil {
				return nil, err
			}
			return map[string]string{"status": "success"}, nil
		},
	},
}
