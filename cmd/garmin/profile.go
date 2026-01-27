package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

const profileUsage = `Usage: garmin profile <command>

Commands:
    social      Get social profile
    settings    Get profile settings
    user        Get user settings (detailed)

Examples:
    garmin profile social
    garmin profile settings
    garmin profile user
`

func profileCmd(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, profileUsage)
		os.Exit(1)
	}

	client, err := loadClient()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	ctx := context.Background()

	switch args[0] {
	case "social":
		data, err := client.UserProfile.GetSocialProfile(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "settings":
		data, err := client.UserProfile.GetProfileSettings(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "user":
		data, err := client.UserProfile.GetUserSettings(ctx)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(data)

	case "-h", "--help", "help": //nolint:goconst // CLI help flags
		fmt.Print(profileUsage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown profile command: %s\n\n%s", args[0], profileUsage)
		os.Exit(1)
	}
}
