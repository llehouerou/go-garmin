package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const devicesUsage = `Usage: garmin devices <command> [arguments]

Commands:
    list                   List registered devices
    settings <device-id>   Get settings for a specific device
    messages               Get device messages
    primary                Get primary training device info

Examples:
    garmin devices list
    garmin devices settings 12345678
    garmin devices messages
    garmin devices primary
`

func devicesCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, devicesUsage)
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
		devices, err := client.Devices.GetDevices(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(devices)

	case "settings":
		if len(args) < 2 {
			printError(errors.New("missing device ID"))
			fmt.Fprint(os.Stderr, devicesUsage)
			os.Exit(1)
		}

		deviceID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			printError(fmt.Errorf("invalid device ID: %s", args[1]))
			os.Exit(1)
		}

		settings, err := client.Devices.GetSettings(ctx, deviceID)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(settings)

	case "messages":
		messages, err := client.Devices.GetMessages(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(messages)

	case "primary":
		info, err := client.Devices.GetPrimaryTrainingDevice(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(info)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(devicesUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown devices command: %s\n\n%s", args[0], devicesUsage)
		os.Exit(1)
	}
}
