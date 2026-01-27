package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	garmin "github.com/llehouerou/go-garmin"
)

var workoutsCmd = &cobra.Command{
	Use:   "workouts",
	Short: "Workouts data (list, get, create, update, delete, download, schedule)",
}

var workoutsListCmd = &cobra.Command{
	Use:   "list [limit]",
	Short: "List workouts",
	Long:  "List workouts (default: 10)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWorkoutsList,
}

var workoutsGetCmd = &cobra.Command{
	Use:   "get <workout-id>",
	Short: "Get workout by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkoutsGet,
}

var workoutsDownloadCmd = &cobra.Command{
	Use:   "download <workout-id>",
	Short: "Download workout as FIT file",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkoutsDownload,
}

var workoutsDeleteCmd = &cobra.Command{
	Use:   "delete <workout-id>",
	Short: "Delete a workout",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkoutsDelete,
}

var workoutsCreateCmd = &cobra.Command{
	Use:   "create <json-file>",
	Short: "Create a workout from JSON file",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkoutsCreate,
}

var workoutsUpdateCmd = &cobra.Command{
	Use:   "update <workout-id> <json-file>",
	Short: "Update a workout from JSON file",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkoutsUpdate,
}

var workoutsScheduleCmd = &cobra.Command{
	Use:   "schedule <workout-id> <date>",
	Short: "Schedule a workout for a date (YYYY-MM-DD)",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorkoutsSchedule,
}

var workoutsGetScheduledCmd = &cobra.Command{
	Use:   "get-scheduled <schedule-id>",
	Short: "Get a scheduled workout by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkoutsGetScheduled,
}

func init() {
	workoutsCmd.AddCommand(workoutsListCmd)
	workoutsCmd.AddCommand(workoutsGetCmd)
	workoutsCmd.AddCommand(workoutsDownloadCmd)
	workoutsCmd.AddCommand(workoutsDeleteCmd)
	workoutsCmd.AddCommand(workoutsCreateCmd)
	workoutsCmd.AddCommand(workoutsUpdateCmd)
	workoutsCmd.AddCommand(workoutsScheduleCmd)
	workoutsCmd.AddCommand(workoutsGetScheduledCmd)
}

func runWorkoutsList(cmd *cobra.Command, args []string) error {
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

	workouts, err := client.Workouts.List(cmd.Context(), 0, limit)
	if err != nil {
		return err
	}

	return printJSON(workouts)
}

func runWorkoutsGet(cmd *cobra.Command, args []string) error {
	workoutID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	workout, err := client.Workouts.Get(cmd.Context(), workoutID)
	if err != nil {
		return err
	}

	return printJSON(workout)
}

func runWorkoutsDownload(cmd *cobra.Command, args []string) error {
	workoutID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	data, err := client.Workouts.DownloadFIT(cmd.Context(), workoutID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("workout_%d.fit", workoutID)
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s (%d bytes)\n", filename, len(data))
	return nil
}

func runWorkoutsDelete(cmd *cobra.Command, args []string) error {
	workoutID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	if err := client.Workouts.Delete(cmd.Context(), workoutID); err != nil {
		return err
	}

	fmt.Printf("Deleted workout %d\n", workoutID)
	return nil
}

func runWorkoutsCreate(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var workout garmin.Workout
	if err := json.Unmarshal(data, &workout); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	created, err := client.Workouts.Create(cmd.Context(), &workout)
	if err != nil {
		return err
	}

	return printJSON(created)
}

func runWorkoutsUpdate(cmd *cobra.Command, args []string) error {
	workoutID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %s", args[0])
	}

	data, err := os.ReadFile(args[1])
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var workout garmin.Workout
	if err := json.Unmarshal(data, &workout); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	updated, err := client.Workouts.Update(cmd.Context(), workoutID, &workout)
	if err != nil {
		return err
	}

	return printJSON(updated)
}

func runWorkoutsSchedule(cmd *cobra.Command, args []string) error {
	workoutID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %s", args[0])
	}

	date, err := time.Parse("2006-01-02", args[1])
	if err != nil {
		return fmt.Errorf("invalid date format (use YYYY-MM-DD): %s", args[1])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	scheduled, err := client.Workouts.Schedule(cmd.Context(), workoutID, date)
	if err != nil {
		return err
	}

	return printJSON(scheduled)
}

func runWorkoutsGetScheduled(cmd *cobra.Command, args []string) error {
	scheduleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid schedule ID: %s", args[0])
	}

	client, err := loadClient()
	if err != nil {
		return err
	}

	scheduled, err := client.Workouts.GetScheduled(cmd.Context(), scheduleID)
	if err != nil {
		return err
	}

	return printJSON(scheduled)
}
