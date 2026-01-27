# go-garmin

A Go client library and CLI for the Garmin Connect API.

## Installation

```bash
go install github.com/llehouerou/go-garmin/cmd/garmin@latest
```

## CLI Usage

### Authentication

```bash
# Login to Garmin Connect
garmin login -email=your@email.com -password=yourpassword

# Logout
garmin logout
```

### Commands

```bash
# Sleep data
garmin sleep [date]

# Wellness data
garmin wellness stress [date]
garmin wellness body-battery [date]
garmin wellness heart-rate [date]

# Activities
garmin activities list [--limit=20]
garmin activities get <activity-id>

# Weight and HRV
garmin weight [date]
garmin hrv [date]

# Training metrics
garmin metrics readiness [date]
garmin metrics vo2max [date]
garmin metrics endurance [date]

# Devices
garmin devices

# User profile
garmin profile

# Workouts
garmin workouts list
garmin workouts get <workout-id>
```

All commands output JSON for easy parsing.

## MCP Server (LLM Integration)

The CLI includes an MCP (Model Context Protocol) server that lets LLM assistants like Claude access your Garmin data.

### Setup with Claude Code

Add to your Claude Code MCP settings (`~/.claude.json` or project `.claude/settings.json`):

```json
{
  "mcpServers": {
    "garmin": {
      "command": "garmin",
      "args": ["mcp"]
    }
  }
}
```

### Prerequisites

1. Install the CLI: `go install github.com/llehouerou/go-garmin/cmd/garmin@latest`
2. Login once: `garmin login -email=your@email.com -password=yourpassword`

The MCP server reuses your CLI session, so you only need to login once.

### Available Tools

Once configured, you can ask Claude questions like:

- "How did I sleep last night?"
- "What's my training readiness today?"
- "Show me my recent activities"
- "What's my current VO2 max?"
- "How's my stress level today?"

The MCP server exposes 27 tools:

| Category | Tools |
|----------|-------|
| Sleep | `get_sleep` |
| Wellness | `get_stress`, `get_body_battery`, `get_heart_rate`, `get_spo2`, `get_respiration`, `get_intensity_minutes` |
| Activity | `list_activities`, `get_activity`, `get_activity_splits` |
| Weight/HRV | `get_weight`, `get_hrv` |
| Device | `list_devices`, `get_device_settings` |
| Metrics | `get_training_readiness`, `get_training_status`, `get_vo2max`, `get_endurance_score`, `get_hill_score`, `get_training_load`, `get_acclimation` |
| Biometric | `get_heart_rate_zones`, `get_ftp` |
| Workout | `list_workouts`, `get_workout`, `delete_workout` |
| Profile | `get_profile` |

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/llehouerou/go-garmin"
)

func main() {
    client := garmin.New(garmin.Options{})

    // Login
    err := client.Login(context.Background(), "email", "password")
    if err != nil {
        panic(err)
    }

    // Get today's sleep data
    sleep, err := client.Sleep.GetDaily(context.Background(), time.Now())
    if err != nil {
        panic(err)
    }

    fmt.Printf("Sleep score: %d\n", sleep.SleepScores.Overall.Value)
}
```

## License

MIT
