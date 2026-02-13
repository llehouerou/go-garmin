package definitions

import (
	"context"
	"fmt"

	garmin "github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// CourseEndpoints defines all course-related API endpoints.
var CourseEndpoints = []endpoint.Endpoint{
	{
		Name:       "ListOwnerCourses",
		Service:    "Courses",
		Cassette:   "courses",
		Path:       "/web-gateway/course/owner",
		HTTPMethod: "GET",

		CLICommand:    "courses",
		CLISubcommand: "list",
		MCPTool:       "list_courses",
		Short:         "List owner courses",
		Long:          "List all courses/routes owned by the authenticated user, including distance, elevation, and activity type",

		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Courses.ListOwner(ctx)
		},
	},
	{
		Name:       "GetCourse",
		Service:    "Courses",
		Cassette:   "courses",
		Path:       "/course-service/course/{course_id}",
		HTTPMethod: "GET",

		Params: []endpoint.Param{
			{
				Name:        "course_id",
				Type:        endpoint.ParamTypeInt,
				Required:    true,
				Description: "Course ID to get details for",
			},
		},

		CLICommand:    "courses",
		CLISubcommand: "get",
		MCPTool:       "get_course",
		Short:         "Get course details",
		Long:          "Get detailed information about a specific course/route including distance, elevation, coordinates, and activity type",

		DependsOn: "ListOwnerCourses",
		ArgProvider: func(result any) map[string]any {
			resp, ok := result.(*garmin.CoursesForUserResponse)
			if !ok || len(resp.CoursesForUser) == 0 {
				return nil
			}
			return map[string]any{"course_id": resp.CoursesForUser[0].CourseID}
		},

		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Courses.Get(ctx, int64(args.Int("course_id")))
		},
	},
}
