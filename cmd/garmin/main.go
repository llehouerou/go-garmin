package main

import (
	"fmt"
	"os"
)

const usage = `Usage: garmin <command> [arguments]

Commands:
    login       Authenticate with Garmin Connect
    logout      Remove saved session
    sleep       Sleep data
    wellness    Wellness data (stress, body battery)

Run 'garmin <command> -h' for command-specific help.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		loginCmd(os.Args[2:])
	case "logout":
		logoutCmd(os.Args[2:])
	case "sleep":
		sleepCmd(os.Args[2:])
	case "wellness":
		wellnessCmd(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n%s", os.Args[1], usage)
		os.Exit(1)
	}
}
