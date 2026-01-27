package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	garmin "github.com/llehouerou/go-garmin"
)

const activitiesUsage = `Usage: garmin activities <command> [arguments]

Commands:
    list [limit]         List recent activities (default: 10)
    get <activity-id>    Get detailed activity information
    weather <activity-id> Get weather data for an activity
    splits <activity-id>  Get splits/laps data for an activity

Examples:
    garmin activities list
    garmin activities list 20
    garmin activities get 21661023200
    garmin activities weather 21661023200
    garmin activities splits 21661023200
`

func activitiesCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, activitiesUsage)
		os.Exit(1)
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	ctx := context.Background()

	switch args[0] {
	case "list":
		limit := 10
		if len(args) > 1 {
			limit, err = strconv.Atoi(args[1])
			if err != nil || limit < 1 {
				printError(fmt.Errorf("invalid limit: %s", args[1]))
				os.Exit(1)
			}
		}

		activities, err := client.Activities.List(ctx, &garmin.ListOptions{Limit: limit})
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(activities)

	case "get":
		if len(args) < 2 {
			printError(errors.New("missing activity ID"))
			fmt.Fprint(os.Stderr, activitiesUsage)
			os.Exit(1)
		}

		activityID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			printError(fmt.Errorf("invalid activity ID: %s", args[1]))
			os.Exit(1)
		}

		activity, err := client.Activities.Get(ctx, activityID)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(activity)

	case "weather":
		if len(args) < 2 {
			printError(errors.New("missing activity ID"))
			fmt.Fprint(os.Stderr, activitiesUsage)
			os.Exit(1)
		}

		activityID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			printError(fmt.Errorf("invalid activity ID: %s", args[1]))
			os.Exit(1)
		}

		weather, err := client.Activities.GetWeather(ctx, activityID)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(weather)

	case "splits":
		if len(args) < 2 {
			printError(errors.New("missing activity ID"))
			fmt.Fprint(os.Stderr, activitiesUsage)
			os.Exit(1)
		}

		activityID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			printError(fmt.Errorf("invalid activity ID: %s", args[1]))
			os.Exit(1)
		}

		splits, err := client.Activities.GetSplits(ctx, activityID)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(splits)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(activitiesUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown activities command: %s\n\n%s", args[0], activitiesUsage)
		os.Exit(1)
	}
}
