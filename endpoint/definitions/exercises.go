package definitions

import (
	"context"
	"errors"
	"fmt"

	"github.com/llehouerou/go-garmin/endpoint"
	"github.com/llehouerou/go-garmin/exercises"
)

// ExerciseEndpoints defines endpoints for the exercise library.
// These are static endpoints that don't require authentication.
var ExerciseEndpoints = []endpoint.Endpoint{
	{
		Name:          "ListExerciseCategories",
		Service:       "Exercises",
		Cassette:      "none",
		Path:          "/exercises/categories",
		HTTPMethod:    "GET",
		CLICommand:    "exercises",
		CLISubcommand: "categories",
		MCPTool:       "list_exercise_categories",
		Short:         "List exercise categories",
		Long:          "List all available exercise categories (e.g., BENCH_PRESS, CURL, DEADLIFT)",
		Handler: func(_ context.Context, _ any, _ *endpoint.HandlerArgs) (any, error) {
			return exercises.Get().Categories(), nil
		},
	},
	{
		Name:          "ListMuscleGroups",
		Service:       "Exercises",
		Cassette:      "none",
		Path:          "/exercises/muscles",
		HTTPMethod:    "GET",
		CLICommand:    "exercises",
		CLISubcommand: "muscles",
		MCPTool:       "list_muscle_groups",
		Short:         "List muscle groups",
		Long:          "List all available muscle groups (e.g., CHEST, BICEPS, QUADS)",
		Handler: func(_ context.Context, _ any, _ *endpoint.HandlerArgs) (any, error) {
			return exercises.Get().Muscles(), nil
		},
	},
	{
		Name:          "ListEquipmentTypes",
		Service:       "Exercises",
		Cassette:      "none",
		Path:          "/exercises/equipment",
		HTTPMethod:    "GET",
		CLICommand:    "exercises",
		CLISubcommand: "equipment",
		MCPTool:       "list_equipment_types",
		Short:         "List equipment types",
		Long:          "List all available equipment types (e.g., DUMBBELL, BARBELL, MACHINE)",
		Handler: func(_ context.Context, _ any, _ *endpoint.HandlerArgs) (any, error) {
			return exercises.Get().Equipment(), nil
		},
	},
	{
		Name:       "ListExercises",
		Service:    "Exercises",
		Cassette:   "none",
		Path:       "/exercises",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "category", Type: endpoint.ParamTypeString, Required: false, Description: "Filter by category (e.g., BENCH_PRESS)"},
			{Name: "muscle", Type: endpoint.ParamTypeString, Required: false, Description: "Filter by muscle group (e.g., CHEST)"},
			{Name: "equipment", Type: endpoint.ParamTypeString, Required: false, Description: "Filter by equipment (e.g., DUMBBELL)"},
			{Name: "search", Type: endpoint.ParamTypeString, Required: false, Description: "Search exercise names"},
		},
		CLICommand:    "exercises",
		CLISubcommand: "list",
		MCPTool:       "list_exercises",
		Short:         "List/search exercises",
		Long:          "List exercises with optional filters. All filters are combined with AND logic.",
		Handler: func(_ context.Context, _ any, args *endpoint.HandlerArgs) (any, error) {
			category := args.String("category")
			muscle := args.String("muscle")
			equipment := args.String("equipment")
			search := args.String("search")
			return exercises.Get().Filter(category, muscle, equipment, search), nil
		},
	},
	{
		Name:       "GetExercise",
		Service:    "Exercises",
		Cassette:   "none",
		Path:       "/exercises/:key",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "key", Type: endpoint.ParamTypeString, Required: true, Description: "Exercise key (e.g., BARBELL_BENCH_PRESS)"},
		},
		CLICommand:    "exercises",
		CLISubcommand: "get",
		MCPTool:       "get_exercise",
		Short:         "Get exercise by key",
		Long:          "Get exercise details by key. May return multiple exercises if the key exists in multiple categories.",
		Handler: func(_ context.Context, _ any, args *endpoint.HandlerArgs) (any, error) {
			key := args.String("key")
			if key == "" {
				return nil, errors.New("key is required")
			}
			result := exercises.Get().ByKey(key)
			if len(result) == 0 {
				return nil, fmt.Errorf("exercise not found: %s", key)
			}
			return result, nil
		},
	},
}
