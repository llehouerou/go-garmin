package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const hrvUsage = `Usage: garmin hrv <command> [arguments]

Commands:
    daily [date]              Get daily HRV data (default: today)
    range <start> <end>       Get HRV data for a date range

Date format: YYYY-MM-DD

Examples:
    garmin hrv daily
    garmin hrv daily 2026-01-27
    garmin hrv range 2026-01-20 2026-01-27
`

func hrvCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, hrvUsage)
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

		data, err := client.HRV.GetDaily(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "range":
		if len(args) < 3 {
			printError(errors.New("range requires start and end dates"))
			fmt.Fprint(os.Stderr, hrvUsage)
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

		data, err := client.HRV.GetRange(ctx, startDate, endDate)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(hrvUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown hrv command: %s\n\n%s", args[0], hrvUsage)
		os.Exit(1)
	}
}
