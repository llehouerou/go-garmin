# Claude Code Instructions

## Verification

Always run `make check` after making changes to verify the code compiles, passes linting, and tests pass.

## Adding New Endpoints

Follow these steps when implementing a new Garmin API endpoint:

### 1. Update Record Fixtures Tool

Add the new endpoint to `cmd/record-fixtures/main.go`:

1. Create a new `record<ServiceName>()` function that makes the API call
2. Call it from `recordFixtures()`
3. Use `testutil.NewRecordingRecorder("<cassette_name>")` for the recorder

```go
func recordNewEndpoint(ctx context.Context, session []byte, date time.Time) error {
    rec, err := testutil.NewRecordingRecorder("cassette_name")
    if err != nil {
        return err
    }
    defer func() { _ = stopRecorder(rec) }()

    client, err := loadSession(rec, session)
    if err != nil {
        return err
    }

    _, err = client.ServiceName.Method(ctx, date)
    if err != nil {
        fmt.Printf("  Warning: %v\n", err)
    }
    return nil
}
```

### 2. Record Cassette

Run the fixture recorder to capture real API responses:

```bash
go run ./cmd/record-fixtures -email=USER -password=PASS [-date=YYYY-MM-DD]
```

### 3. Check Sensitive Information

Review the recorded cassette in `testdata/cassettes/` for personal data:

- User IDs (`ownerId`, `userProfileId`, `userProfilePk`)
- Names (`ownerFullName`, `fullname`, `displayName`, `ownerDisplayName`)
- Profile image URLs
- Email addresses
- Any other PII

If new patterns are found, update `testutil/vcr.go`:

1. Add new regex patterns to the `var` block
2. Add replacement logic in `anonymizeBody()` function
3. Re-record the cassette

### 4. Create Service Implementation

Create `service_<name>.go` with:

1. **Type definitions** - Complete structs for all response fields
2. **Helper methods** - Convenience methods like `StartTime()`, `Duration()`, `RawJSON()`
3. **Service methods** - API methods that call `s.client.doAPI()`

```go
// Type with all fields from API response
type ResponseType struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
    // ... all fields

    raw json.RawMessage
}

// RawJSON returns the original JSON response (fallback for API changes)
func (r *ResponseType) RawJSON() json.RawMessage {
    return r.raw
}

// Service method
func (s *ServiceName) GetData(ctx context.Context, date time.Time) (*ResponseType, error) {
    path := fmt.Sprintf("/service-path/endpoint/%s", date.Format("2006-01-02"))

    resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    raw, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result ResponseType
    if err := json.Unmarshal(raw, &result); err != nil {
        return nil, err
    }
    result.raw = raw

    return &result, nil
}
```

### 5. Add Unit Tests

Create `service_<name>_test.go` with:

- JSON unmarshaling tests
- Helper method tests (conversions, calculations)
- Edge case tests (nil values, empty responses)

### 6. Add Integration Tests

Add to `integration_test.go`:

```go
func TestIntegration_Service_Method(t *testing.T) {
    skipIfNoCassette(t, "cassette_name")

    rec, err := testutil.NewRecorder("cassette_name", recorder.ModeReplayOnly)
    if err != nil {
        t.Fatalf("failed to create recorder: %v", err)
    }
    defer func() { _ = rec.Stop() }()

    client := newTestClient(t, rec)
    ctx := context.Background()

    result, err := client.ServiceName.Method(ctx, args)
    if err != nil {
        t.Fatalf("Method failed: %v", err)
    }

    // Verify expected fields
    if result.Field == "" {
        t.Error("expected Field to be set")
    }
}
```

### 7. Add CLI Command

Create or update `cmd/garmin/<service>.go`:

1. Define usage string with commands and examples
2. Implement command handler with subcommands
3. Output JSON to stdout

Update `cmd/garmin/main.go`:
1. Add command to usage string
2. Add case to switch statement

### 8. Update Documentation

Mark endpoints as implemented in `ENDPOINTS.md`:

```markdown
| [x] | GET | `/service/endpoint` | Description |
```

### 9. Verify and Commit

```bash
make check
git add -A
git commit -m "feat: add ServiceName with endpoint methods"
```

## File Structure

```
├── service_<name>.go       # Service implementation and types
├── service_<name>_test.go  # Unit tests
├── integration_test.go     # Integration tests (append to existing)
├── cmd/garmin/<name>.go    # CLI commands
├── cmd/record-fixtures/    # Update main.go for new cassettes
├── testdata/cassettes/     # VCR cassettes
├── testutil/vcr.go         # Sanitization patterns
└── ENDPOINTS.md            # Implementation status
```

## Code Style

- Use pointers for optional fields (`*int`, `*float64`, `*string`)
- Always include `raw json.RawMessage` and `RawJSON()` method for API change fallback
- Use `time.Time` helpers for timestamp conversions
- Follow existing naming conventions (e.g., `GetDaily`, `List`, `Get`)
- Keep CLI output as JSON for easy piping/parsing
