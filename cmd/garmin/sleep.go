package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const sleepUsage = `Usage: garmin sleep [date]

Get sleep data for the specified date.

Date format: YYYY-MM-DD (defaults to today)
`

func sleepCmd(args []string) {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		fmt.Print(sleepUsage)
		return
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	// Parse date (default to today)
	date := time.Now()
	if len(args) > 0 {
		date, err = time.Parse("2006-01-02", args[0])
		if err != nil {
			printError(errors.New("invalid date format, use YYYY-MM-DD"))
			os.Exit(1)
		}
	}

	ctx := context.Background()

	data, err := client.Sleep.GetDaily(ctx, date)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(data)
}
