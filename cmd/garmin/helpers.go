package main

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// parseDate parses a date from args.
// Returns today's date if args is empty.
func parseDate(args []string) (time.Time, error) {
	if len(args) == 0 {
		return time.Now(), nil
	}
	date, err := time.Parse("2006-01-02", args[0])
	if err != nil {
		return time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
	}
	return date, nil
}

// printJSON encodes the given value as JSON to stdout.
func printJSON(v any) error {
	return json.NewEncoder(os.Stdout).Encode(v)
}
