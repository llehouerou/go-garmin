# Testing

This project uses VCR-style testing to record and replay API interactions.

## Overview

- **Unit tests**: Test data structures and helpers without API calls
- **Integration tests**: Test real API interactions using recorded cassettes

## Recording Fixtures

To record new API interactions for testing:

```bash
# Build the record-fixtures tool
make record-fixtures

# Record cassettes with your Garmin credentials
./record-fixtures -email=your@email.com -password=yourpassword

# Optionally specify a date
./record-fixtures -email=your@email.com -password=yourpassword -date=2026-01-27
```

This creates cassette files in `testdata/cassettes/`:
- `auth.yaml` - Authentication flow (recorded once, reused for session)
- `sleep_daily.yaml` - Sleep service endpoints
- `wellness_stress.yaml` - Stress data endpoints
- `wellness_body_battery.yaml` - Body battery endpoints

## Running Tests

```bash
# Run all tests (integration tests skip if no cassettes)
make test

# Run with verbose output
go test -v ./...

# Run only integration tests
go test -v -run Integration ./...
```

## How It Works

1. **Recording**: The `record-fixtures` command authenticates once and records the auth flow to `auth.yaml`. It then reuses the session to record each API endpoint to separate cassettes.

2. **Sanitization**: Sensitive data is automatically anonymized before saving:
   - Headers: Authorization, Cookie, Set-Cookie are redacted
   - URLs: OAuth tickets are redacted
   - Bodies: Passwords are redacted
   - Personal info: `userProfilePK`, names, emails are replaced with anonymous values

3. **Replay**: Integration tests load a fake session (to satisfy the client's auth check) and replay the API cassettes without making real API calls.

## Adding New Tests

1. Add a new recording function in `cmd/record-fixtures/main.go`
2. Call it from `recordFixtures()` with the session
3. Add corresponding test in `integration_test.go`
4. Re-run `record-fixtures` to generate the cassette

Example:

```go
// In cmd/record-fixtures/main.go
func recordNewEndpoint(ctx context.Context, session []byte, date time.Time) error {
    rec, err := testutil.NewRecordingRecorder("new_endpoint")
    if err != nil {
        return err
    }
    defer func() { _ = stopRecorder(rec) }()

    client, err := loadSession(rec, session)
    if err != nil {
        return err
    }

    fmt.Printf("  Getting new endpoint data for %s...\n", date.Format("2006-01-02"))
    _, err = client.NewService.GetData(ctx, date)
    if err != nil {
        fmt.Printf("  Warning: %v\n", err)
    }

    return nil
}
```

```go
// In integration_test.go
func TestIntegration_NewService_GetData(t *testing.T) {
    skipIfNoCassette(t, "new_endpoint")

    rec, err := testutil.NewRecorder("new_endpoint", recorder.ModeReplayOnly)
    if err != nil {
        t.Fatalf("failed to create recorder: %v", err)
    }
    defer func() { _ = rec.Stop() }()

    client := newTestClient(t, rec)
    ctx := context.Background()
    date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

    data, err := client.NewService.GetData(ctx, date)
    if err != nil {
        t.Fatalf("GetData failed: %v", err)
    }

    // Add assertions...
}
```

## Security Notes

- Never commit real credentials
- Cassettes are sanitized but review them before committing
- Personal information (userProfilePK, names, emails) is automatically anonymized
- The `testdata/cassettes/` directory can be committed as cassettes contain only anonymized data
