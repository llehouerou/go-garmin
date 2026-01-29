# go-garmin

A Go client library and CLI for the Garmin Connect API.

## Installation

### Prerequisites

- Go 1.25 or later ([install Go](https://go.dev/doc/install))
- `$GOPATH/bin` (usually `~/go/bin`) in your PATH

### Install CLI

```bash
go install github.com/llehouerou/go-garmin/cmd/garmin@latest
```

Verify installation:

```bash
garmin --version
```

### Build from Source

```bash
git clone https://github.com/llehouerou/go-garmin.git
cd go-garmin
go build -o garmin ./cmd/garmin
./garmin --help
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
garmin wellness spo2 [date]
garmin wellness respiration [date]
garmin wellness intensity-minutes [date]

# Activities
garmin activities list [--start=0] [--limit=20]
garmin activities get <activity-id>
garmin activities types
garmin activities splits <activity-id>
garmin activities weather <activity-id>
garmin activities details <activity-id>
garmin activities hr-zones <activity-id>
garmin activities power-zones <activity-id>
garmin activities exercise-sets <activity-id>

# Weight and HRV
garmin weight daily [date]
garmin weight range --start=YYYY-MM-DD --end=YYYY-MM-DD
garmin hrv daily [date]
garmin hrv range --start=YYYY-MM-DD --end=YYYY-MM-DD

# Training metrics
garmin metrics readiness [date]
garmin metrics vo2max [date]
garmin metrics endurance [date]
garmin metrics hill [date]
garmin metrics training-status [date]
garmin metrics load-balance [date]
garmin metrics acclimation [date]
garmin metrics race-predictions [display-name]

# Fitness age
garmin fitnessage stats --start=YYYY-MM-DD --end=YYYY-MM-DD

# Fitness stats
garmin fitnessstats get [--start=YYYY-MM-DD] [--end=YYYY-MM-DD] [--aggregation=weekly] [--metrics=calories,distance,duration]
garmin fitnessstats activities [--start=YYYY-MM-DD] [--end=YYYY-MM-DD] [--activity_type=running] [--metrics=name,startLocal,activityType]

# Biometric data
garmin biometric lactate-threshold
garmin biometric ftp
garmin biometric hr-zones
garmin biometric power-weight [date]

# Devices
garmin devices list
garmin devices settings <device-id>

# User profile
garmin profile social
garmin profile settings
garmin profile display

# Workouts
garmin workouts list [--start=0] [--limit=20]
garmin workouts get <workout-id>
garmin workouts create --file=workout.json
garmin workouts create --json='{"workoutName": "..."}'
cat workout.json | garmin workouts create
garmin workouts update <workout-id> --file=workout.json
garmin workouts delete <workout-id>
garmin workouts schedule <workout-id> <date>
garmin workouts unschedule <schedule-id>

# Calendar (month is 0-indexed: January=0)
garmin calendar get --year=2026 [--month=0] [--day=28] [--start=1]
```

All commands output JSON for easy parsing.

## MCP Server (LLM Integration)

The CLI includes an MCP (Model Context Protocol) server that lets LLM assistants like Claude access your Garmin data.

### Prerequisites

1. Install the CLI (see [Installation](#installation))
2. Login once to create a session:
   ```bash
   garmin login -email=your@email.com -password=yourpassword
   ```

The MCP server reuses your CLI session, so you only need to login once.

### Claude Code

Add to `~/.claude.json` (global) or `.claude/settings.json` (project):

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

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

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

### Cursor

Add to Cursor settings (Settings > MCP Servers):

```json
{
  "garmin": {
    "command": "garmin",
    "args": ["mcp"]
  }
}
```

### Troubleshooting

If the MCP server isn't working:

1. Verify the CLI is in your PATH: `which garmin`
2. Check you're logged in: `garmin sleep` should return data
3. Use the full path to the binary if needed:
   ```json
   {
     "mcpServers": {
       "garmin": {
         "command": "/full/path/to/garmin",
         "args": ["mcp"]
       }
     }
   }
   ```

### Available Tools

Once configured, you can ask Claude questions like:

- "How did I sleep last night?"
- "What's my training readiness today?"
- "Show me my recent activities"
- "What's my current VO2 max?"
- "How's my stress level today?"

The MCP server exposes 47 tools across these categories:

| Category | Tools |
|----------|-------|
| Sleep | `get_sleep` |
| Wellness | `get_stress`, `get_body_battery`, `get_heart_rate`, `get_spo2`, `get_respiration`, `get_intensity_minutes` |
| Activity | `list_activities`, `get_activity`, `get_activity_types`, `get_activity_splits`, `get_activity_weather`, `get_activity_details`, `get_activity_hr_zones`, `get_activity_power_zones`, `get_activity_exercise_sets` |
| Weight | `get_weight` |
| HRV | `get_hrv` |
| Device | `list_devices`, `get_device_settings` |
| Metrics | `get_training_readiness`, `get_training_status`, `get_vo2max`, `get_endurance_score`, `get_hill_score`, `get_training_load_balance`, `get_heat_altitude_acclimation`, `get_race_predictions` |
| Fitness Age | `get_fitness_age_stats` |
| Fitness Stats | `get_fitness_stats`, `get_fitness_stats_activities` |
| Biometric | `get_lactate_threshold`, `get_cycling_ftp`, `get_heart_rate_zones`, `get_power_to_weight` |
| Workout | `list_workouts`, `get_workout`, `create_workout`, `update_workout`, `delete_workout`, `schedule_workout`, `unschedule_workout` |
| Calendar | `get_calendar` |
| Profile | `get_social_profile`, `get_user_settings`, `get_profile_settings` |
| Utility | `get_current_date` |

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

## Architecture

This project uses a **Declarative Endpoint Registry** system. All endpoints are defined once in `endpoint/definitions/` and automatically generate:

- CLI commands (via `endpoint.CLIGenerator`)
- MCP tools (via `endpoint.MCPGenerator`)
- Validation rules (via `endpoint.Validator`)

See [CLAUDE.md](CLAUDE.md) for instructions on adding new endpoints.

## License

MIT
