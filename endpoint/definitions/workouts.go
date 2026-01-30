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

SPORT TYPES (sportTypeId / sportTypeKey):
  1=running, 2=cycling, 3=other, 4=swimming, 5=strength_training,
  6=cardio_training, 7=yoga, 8=pilates, 9=hiit, 10=multi_sport, 11=mobility

STEP TYPES (stepTypeId / stepTypeKey):
  1=warmup, 2=cooldown, 3=interval, 4=recovery, 5=rest, 6=repeat, 7=other, 8=main

END CONDITIONS (conditionTypeId / conditionTypeKey):
  1=lap.button, 2=time (seconds), 3=distance (meters), 4=calories, 5=power,
  6=heart.rate, 7=iterations (for repeats), 8=fixed.rest, 9=fixed.repetition,
  10=reps (for strength), 17=velocity.loss, 24=power.loss

TARGET TYPES (workoutTargetTypeId / workoutTargetTypeKey):
  1=no.target, 2=power.zone, 3=cadence, 4=heart.rate.zone, 5=speed.zone,
  6=pace.zone, 7=grade, 14=swim.stroke, 15=resistance, 19=instruction

Structure:
- workoutName (string, required): Name of the workout
- description (string): Workout description
- sportType (object, required): {"sportTypeId": N, "sportTypeKey": "key"}
- workoutSegments (array, required): Array of workout segments

WorkoutStep types:
- "ExecutableStepDTO": A single exercise step
- "RepeatGroupDTO": A repeat group containing nested steps

ExecutableStepDTO fields:
- stepOrder (int): Order within segment (1-based)
- stepType: {"stepTypeId": N, "stepTypeKey": "key"}
- description (string): Optional note/description
- endCondition: {"conditionTypeId": N, "conditionTypeKey": "key"}
- endConditionValue: Value for condition (seconds/meters/reps)
- targetType: {"workoutTargetTypeId": N, "workoutTargetTypeKey": "key"}
- zoneNumber: Zone number (1-5) for zone targets
- targetValueOne/targetValueTwo: Custom target range

STRENGTH TRAINING specific fields (sportTypeId=5):
- category: Exercise category (e.g., "BENCH_PRESS", "CURL", "DEADLIFT")
- exerciseName: Exercise key (e.g., "BARBELL_BENCH_PRESS", "DUMBBELL_CURL")
- weightValue: Weight amount (optional)
- weightUnit: Weight unit info (optional)
Use list_exercise_categories, list_exercises, get_exercise tools to find valid values.

SWIMMING specific fields (sportTypeId=4):
- strokeType: {"strokeTypeId": N, "strokeTypeKey": "key"}
  (1=any_stroke, 2=backstroke, 3=breaststroke, 4=drill, 5=fly, 6=free, 7=individual_medley)
- equipmentType: {"equipmentTypeId": N, "equipmentTypeKey": "key"}
  (1=fins, 2=kickboard, 3=paddles, 4=pull_buoy, 5=snorkel)

RepeatGroupDTO fields:
- numberOfIterations: Number of times to repeat
- workoutSteps: Nested array of steps to repeat
- smartRepeat (bool): Enable smart repeat
- skipLastRestStep (bool): Skip last recovery in group`,
	Example: `RUNNING EXAMPLE:
{
  "workoutName": "Easy 30min Run",
  "sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
  "workoutSegments": [{
    "segmentOrder": 1,
    "sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
    "workoutSteps": [
      {"type": "ExecutableStepDTO", "stepOrder": 1,
       "stepType": {"stepTypeId": 1, "stepTypeKey": "warmup"},
       "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
       "endConditionValue": 300,
       "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"}},
      {"type": "ExecutableStepDTO", "stepOrder": 2,
       "stepType": {"stepTypeId": 3, "stepTypeKey": "interval"},
       "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
       "endConditionValue": 1500,
       "targetType": {"workoutTargetTypeId": 4, "workoutTargetTypeKey": "heart.rate.zone"},
       "zoneNumber": 2}
    ]
  }]
}

STRENGTH TRAINING EXAMPLE:
{
  "workoutName": "Upper Body Workout",
  "sportType": {"sportTypeId": 5, "sportTypeKey": "strength_training"},
  "workoutSegments": [{
    "segmentOrder": 1,
    "sportType": {"sportTypeId": 5, "sportTypeKey": "strength_training"},
    "workoutSteps": [
      {"type": "ExecutableStepDTO", "stepOrder": 1,
       "stepType": {"stepTypeId": 1, "stepTypeKey": "warmup"},
       "endCondition": {"conditionTypeId": 10, "conditionTypeKey": "reps"},
       "endConditionValue": 10,
       "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"},
       "category": "ROW", "exerciseName": "BANDED_FACE_PULLS"},
      {"type": "RepeatGroupDTO", "stepOrder": 2,
       "stepType": {"stepTypeId": 6, "stepTypeKey": "repeat"},
       "numberOfIterations": 3,
       "endCondition": {"conditionTypeId": 7, "conditionTypeKey": "iterations"},
       "workoutSteps": [
         {"type": "ExecutableStepDTO", "stepOrder": 3,
          "stepType": {"stepTypeId": 3, "stepTypeKey": "interval"},
          "childStepId": 1,
          "endCondition": {"conditionTypeId": 10, "conditionTypeKey": "reps"},
          "endConditionValue": 8,
          "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"},
          "category": "BENCH_PRESS", "exerciseName": "BARBELL_BENCH_PRESS"},
         {"type": "ExecutableStepDTO", "stepOrder": 4,
          "stepType": {"stepTypeId": 4, "stepTypeKey": "recovery"},
          "childStepId": 1,
          "endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
          "endConditionValue": 90,
          "targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"}}
       ]}
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
		Name:          "CreateWorkout",
		Service:       "Workouts",
		Cassette:      "workouts",
		Path:          "/workout-service/workout",
		HTTPMethod:    "POST",
		Body:          workoutBodyConfig,
		CLICommand:    "workouts",
		CLISubcommand: "create",
		MCPTool:       "create_workout",
		Short:         "Create a new workout",
		Long:          "Create a new workout with segments and steps. Use --file to read from a file, --json to pass inline JSON, or pipe JSON to stdin.",
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
		Body:          workoutBodyConfig,
		CLICommand:    "workouts",
		CLISubcommand: "update",
		MCPTool:       "update_workout",
		Short:         "Update an existing workout",
		Long:          "Update an existing workout. Use --file to read from a file, --json to pass inline JSON, or pipe JSON to stdin.",
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
		CLICommand:    "workouts",
		CLISubcommand: "delete",
		MCPTool:       "delete_workout",
		Short:         "Delete a workout",
		Long:          "Permanently delete a workout from your Garmin account",
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
