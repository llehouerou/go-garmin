package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const weightUsage = `Usage: garmin weight <command> [arguments]

Commands:
    daily [date]              Get daily weight data (default: today)
    range <start> <end>       Get weight data for a date range

Date format: YYYY-MM-DD

Examples:
    garmin weight daily
    garmin weight daily 2026-01-27
    garmin weight range 2025-12-28 2026-01-27
`

func weightCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, weightUsage)
		os.Exit(1)
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	ctx := context.Background()

	switch args[0] {
	case "daily":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Weight.GetDaily(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "range":
		if len(args) < 3 {
			printError(errors.New("range requires start and end dates"))
			fmt.Fprint(os.Stderr, weightUsage)
			os.Exit(1)
		}

		startDate, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			printError(fmt.Errorf("invalid start date: %s", args[1]))
			os.Exit(1)
		}

		endDate, err := time.Parse("2006-01-02", args[2])
		if err != nil {
			printError(fmt.Errorf("invalid end date: %s", args[2]))
			os.Exit(1)
		}

		data, err := client.Weight.GetRange(ctx, startDate, endDate)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(weightUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown weight command: %s\n\n%s", args[0], weightUsage)
		os.Exit(1)
	}
}
