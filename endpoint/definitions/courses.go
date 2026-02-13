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

		CLICommand: "courses",
		MCPTool:    "list_courses",
		Short:      "List owner courses",
		Long:       "List all courses/routes owned by the authenticated user, including distance, elevation, and activity type",

		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Courses.ListOwner(ctx)
		},
	},
}
