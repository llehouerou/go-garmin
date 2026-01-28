package definitions

import (
	"context"
	"fmt"

	"github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/endpoint"
)

// DeviceEndpoints defines all device-related endpoints.
var DeviceEndpoints = []endpoint.Endpoint{
	{
		Name:          "ListDevices",
		Service:       "Devices",
		Cassette:      "devices",
		Path:          "/device-service/deviceregistration/devices",
		HTTPMethod:    "GET",
		CLICommand:    "devices",
		CLISubcommand: "list",
		MCPTool:       "list_devices",
		Short:         "List all registered devices",
		Long:          "List all registered Garmin devices with their capabilities and status",
		Handler: func(ctx context.Context, c any, _ *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Devices.GetDevices(ctx)
		},
	},
	{
		Name:       "GetDeviceSettings",
		Service:    "Devices",
		Cassette:   "devices",
		Path:       "/device-service/deviceservice/device-info/settings",
		HTTPMethod: "GET",
		Params: []endpoint.Param{
			{Name: "device_id", Type: endpoint.ParamTypeInt, Required: true, Description: "The device ID"},
		},
		CLICommand:    "devices",
		CLISubcommand: "settings",
		MCPTool:       "get_device_settings",
		Short:         "Get device settings",
		Long:          "Get settings for a specific device including alarms, activity tracking, and preferences",
		DependsOn:     "ListDevices",
		ArgProvider: func(result any) map[string]any {
			devices, ok := result.([]garmin.Device)
			if !ok || len(devices) == 0 {
				return nil
			}
			return map[string]any{"device_id": devices[0].DeviceID}
		},
		Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
			client, ok := c.(*garmin.Client)
			if !ok {
				return nil, fmt.Errorf("handler received invalid client type: %T, expected *garmin.Client", c)
			}
			return client.Devices.GetSettings(ctx, int64(args.Int("device_id")))
		},
	},
}
