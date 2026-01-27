package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const metricsUsage = `Usage: garmin metrics <command> [arguments]

Commands:
    readiness [date]           Get training readiness (default: today)
    endurance [date]           Get endurance score (default: today)
    hill [date]                Get hill score (default: today)
    vo2max [date]              Get latest VO2 max / MET (default: today)
    vo2max-range <start> <end> Get VO2 max / MET for a date range
    status [date]              Get daily training status (default: today)
    status-agg [date]          Get aggregated training status (default: today)
    load-balance [date]        Get training load balance (default: today)
    acclimation [date]         Get heat/altitude acclimation (default: today)

Date format: YYYY-MM-DD

Examples:
    garmin metrics readiness
    garmin metrics endurance 2026-01-27
    garmin metrics vo2max-range 2025-12-28 2026-01-27
`

func metricsCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, metricsUsage)
		os.Exit(1)
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	ctx := context.Background()

	switch args[0] {
	case "readiness":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetTrainingReadiness(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "endurance":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetEnduranceScore(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "hill":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetHillScore(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "vo2max":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetMaxMetLatest(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "vo2max-range":
		if len(args) < 3 {
			printError(errors.New("vo2max-range requires start and end dates"))
			fmt.Fprint(os.Stderr, metricsUsage)
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

		data, err := client.Metrics.GetMaxMetDaily(ctx, startDate, endDate)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "status":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetTrainingStatusDaily(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "status-agg":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetTrainingStatusAggregated(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "load-balance":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetTrainingLoadBalance(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "acclimation":
		date := time.Now()
		if len(args) > 1 {
			date, err = time.Parse("2006-01-02", args[1])
			if err != nil {
				printError(errors.New("invalid date format, use YYYY-MM-DD"))
				os.Exit(1)
			}
		}

		data, err := client.Metrics.GetHeatAltitudeAcclimation(ctx, date)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(metricsUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown metrics command: %s\n\n%s", args[0], metricsUsage)
		os.Exit(1)
	}
}
