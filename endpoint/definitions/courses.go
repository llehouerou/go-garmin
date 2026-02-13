package definitions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	{
		Name:       "DownloadCourseGPX",
		Service:    "Courses",
		Cassette:   "courses_download",
		Path:       "/course-service/course/gpx/{course_id}",
		HTTPMethod: "GET",
		RawOutput:  true,

		Params: []endpoint.Param{
			{
				Name:        "course_id",
				Type:        endpoint.ParamTypeInt,
				Required:    true,
				Description: "Course ID to download",
			},
		},

		CLICommand:    "courses",
		CLISubcommand: "download-gpx",
		Short:         "Download course as GPX",
		Long:          "Download a course/route as a GPX file. Output goes to stdout by default, use -o to write to a file.",

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
			return client.Courses.DownloadGPX(ctx, int64(args.Int("course_id")))
		},
	},
	{
		Name:       "DownloadCourseFIT",
		Service:    "Courses",
		Cassette:   "courses_download",
		Path:       "/course-service/course/fit/{course_id}/0",
		HTTPMethod: "GET",
		RawOutput:  true,

		Params: []endpoint.Param{
			{
				Name:        "course_id",
				Type:        endpoint.ParamTypeInt,
				Required:    true,
				Description: "Course ID to download",
			},
		},

		CLICommand:    "courses",
		CLISubcommand: "download-fit",
		Short:         "Download course as FIT",
		Long:          "Download a course/route as a FIT file. Output goes to stdout by default, use -o to write to a file.",

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
			return client.Courses.DownloadFIT(ctx, int64(args.Int("course_id")))
		},
	},
	{
		Name:       "ImportCourseGPX",
		Service:    "Courses",
		Cassette:   "none",
		Path:       "/course-service/course/import",
		HTTPMethod: "POST",

		Params: []endpoint.Param{
			{
				Name:        "file",
				Type:        endpoint.ParamTypeString,
				Required:    true,
				Description: "Path to GPX file to import",
			},
			{
				Name:        "activity-type",
				Type:        endpoint.ParamTypeInt,
				Required:    false,
				Description: "Activity type ID (e.g. 1=running, 3=hiking, 5=cycling)",
			},
			{
				Name:        "privacy",
				Type:        endpoint.ParamTypeInt,
				Required:    false,
				Description: "Privacy rule: 1=Public, 2=Private (default), 4=Group",
			},
		},

		CLICommand:    "courses",
		CLISubcommand: "import",
		Short:         "Import a GPX course",
		Long:          "Import a course/route from a GPX file to Garmin Connect",

		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			filePath := args.String("file")
			f, err := os.Open(filePath)
			if err != nil {
				return nil, fmt.Errorf("open file: %w", err)
			}
			defer f.Close()
			return client.Courses.Import(ctx, filepath.Base(filePath), f, args.Int("activity-type"), args.Int("privacy"))
		},
	},
}
