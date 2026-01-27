package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	garmin "github.com/llehouerou/go-garmin"
)

var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "Activities data (list, details, weather, splits, download)",
}

var activitiesListCmd = &cobra.Command{
	Use:   "list [limit]",
	Short: "List recent activities",
	Long:  "List recent activities (default: 10)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runActivitiesList,
}

var activitiesGetCmd = &cobra.Command{
	Use:   "get <activity-id>",
	Short: "Get detailed activity information",
	Args:  cobra.ExactArgs(1),
	RunE:  runActivitiesGet,
}

var activitiesWeatherCmd = &cobra.Command{
	Use:   "weather <activity-id>",
	Short: "Get weather data for an activity",
	Args:  cobra.ExactArgs(1),
	RunE:  runActivitiesWeather,
}

var activitiesSplitsCmd = &cobra.Command{
	Use:   "splits <activity-id>",
	Short: "Get splits/laps data for an activity",
	Args:  cobra.ExactArgs(1),
	RunE:  runActivitiesSplits,
}

var activitiesDownloadCmd = &cobra.Command{
	Use:       "download <activity-id> <format>",
	Short:     "Download activity file",
	Long:      "Download activity file in the specified format (fit, tcx, gpx, kml, csv)",
	Args:      cobra.ExactArgs(2),
	ValidArgs: []string{"fit", "tcx", "gpx", "kml", "csv"},
	RunE:      runActivitiesDownload,
}

func init() {
	activitiesCmd.AddCommand(activitiesListCmd)
	activitiesCmd.AddCommand(activitiesGetCmd)
	activitiesCmd.AddCommand(activitiesWeatherCmd)
	activitiesCmd.AddCommand(activitiesSplitsCmd)
	activitiesCmd.AddCommand(activitiesDownloadCmd)
}

func runActivitiesList(cmd *cobra.Command, args []string) error {
	limit := 10
	if len(args) > 0 {
		var err error
		limit, err = strconv.Atoi(args[0])
		if err != nil || limit < 1 {
			return fmt.Errorf("invalid limit: %s", args[0])
		}
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	activities, err := client.Activities.List(cmd.Context(), &garmin.ListOptions{Limit: limit})
	if err != nil {
		return err
	}

	return printJSON(activities)
}

func runActivitiesGet(cmd *cobra.Command, args []string) error {
	activityID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid activity ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	activity, err := client.Activities.Get(cmd.Context(), activityID)
	if err != nil {
		return err
	}

	return printJSON(activity)
}

func runActivitiesWeather(cmd *cobra.Command, args []string) error {
	activityID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid activity ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	weather, err := client.Activities.GetWeather(cmd.Context(), activityID)
	if err != nil {
		return err
	}

	return printJSON(weather)
}

func runActivitiesSplits(cmd *cobra.Command, args []string) error {
	activityID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid activity ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	splits, err := client.Activities.GetSplits(cmd.Context(), activityID)
	if err != nil {
		return err
	}

	return printJSON(splits)
}

func runActivitiesDownload(cmd *cobra.Command, args []string) error {
	activityID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid activity ID: %s", args[0])
	}

	format := args[1]

	client, err := loadClient()
	if err != nil {
		return err
	}

	var data []byte
	ctx := cmd.Context()

	switch format {
	case "fit":
		data, err = client.Activities.DownloadFIT(ctx, activityID)
	case "tcx":
		data, err = client.Activities.DownloadTCX(ctx, activityID)
	case "gpx":
		data, err = client.Activities.DownloadGPX(ctx, activityID)
	case "kml":
		data, err = client.Activities.DownloadKML(ctx, activityID)
	case "csv":
		data, err = client.Activities.DownloadCSV(ctx, activityID)
	default:
		return fmt.Errorf("unknown format: %s (use fit, tcx, gpx, kml, or csv)", format)
	}

	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%d.%s", activityID, format)
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s (%d bytes)\n", filename, len(data))
	return nil
}
