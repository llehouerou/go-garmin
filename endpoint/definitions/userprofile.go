package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// UserProfileEndpoints defines all user profile-related endpoints.
var UserProfileEndpoints = []endpoint.Endpoint{
	{
		Name:          "GetSocialProfile",
		Service:       "UserProfile",
		Cassette:      "userprofile",
		Path:          "/userprofile-service/socialProfile",
		HTTPMethod:    "GET",
		CLICommand:    "profile",
		CLISubcommand: "social",
		MCPTool:       "get_social_profile",
		Short:         "Get social profile",
		Long:          "Get the user's social profile including display name, bio, and visibility settings",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.UserProfile.GetSocialProfile(ctx)
		},
	},
	{
		Name:          "GetUserSettings",
		Service:       "UserProfile",
		Cassette:      "userprofile",
		Path:          "/userprofile-service/userprofile/user-settings",
		HTTPMethod:    "GET",
		CLICommand:    "profile",
		CLISubcommand: "settings",
		MCPTool:       "get_user_settings",
		Short:         "Get user settings",
		Long:          "Get the user's settings including personal data, sleep settings, and preferences",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.UserProfile.GetUserSettings(ctx)
		},
	},
	{
		Name:          "GetProfileSettings",
		Service:       "UserProfile",
		Cassette:      "userprofile",
		Path:          "/userprofile-service/userprofile/settings",
		HTTPMethod:    "GET",
		CLICommand:    "profile",
		CLISubcommand: "display",
		MCPTool:       "get_profile_settings",
		Short:         "Get profile display settings",
		Long:          "Get the user's profile display settings including locale, units, and format preferences",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.UserProfile.GetProfileSettings(ctx)
		},
	},
}
