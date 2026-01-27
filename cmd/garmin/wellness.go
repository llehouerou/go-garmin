package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const wellnessUsage = `Usage: garmin wellness <command> [date]

Commands:
    stress        Get stress data
    body-battery  Get body battery data
    heart-rate    Get heart rate data

Date format: YYYY-MM-DD (defaults to today)
`

func wellnessCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, wellnessUsage)
		os.Exit(1)
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	// Parse date (default to today)
	date := time.Now()
	if len(args) > 1 {
		date, err = time.Parse("2006-01-02", args[1])
		if err != nil {
			printError(errors.New("invalid date format, use YYYY-MM-DD"))
			os.Exit(1)
		}
	}

	ctx := context.Background()

	switch args[0] {
	case "stress":
		data, err := client.Wellness.GetDailyStress(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "body-battery":
		data, err := client.Wellness.GetBodyBatteryEvents(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data.Events)

	case "heart-rate":
		data, err := client.Wellness.GetDailyHeartRate(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(wellnessUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown wellness command: %s\n\n%s", args[0], wellnessUsage)
		os.Exit(1)
	}
}
