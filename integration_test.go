// integration_test.go
//
// Integration tests using recorded API interactions (cassettes).
// To record new cassettes:
//
//	go run ./cmd/record-fixtures -email=EMAIL -password=PASSWORD
//
// Tests will skip if cassettes don't exist.
package garmin

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/llehouerou/go-garmin/testutil"
)

func skipIfNoCassette(t *testing.T, name string) {
	t.Helper()
	cassettePath := filepath.Join(testutil.CassetteDir, name+".yaml")
	if _, err := os.Stat(cassettePath); os.IsNotExist(err) {
		t.Skipf("cassette %s not found, run record-fixtures first", name)
	}
}

// newTestClient creates a test client with a fake session loaded and VCR recorder attached.
func newTestClient(t *testing.T, rec *recorder.Recorder) *Client {
	t.Helper()

	client := New(Options{
		HTTPClient: testutil.HTTPClientWithRecorder(rec),
	})

	// Load fake session to make client "authenticated"
	if err := client.LoadSession(strings.NewReader(testutil.FakeSessionJSON())); err != nil {
		t.Fatalf("failed to load fake session: %v", err)
	}

	return client
}

func TestIntegration_Sleep_GetDaily(t *testing.T) {
	skipIfNoCassette(t, "sleep_daily")

	rec, err := testutil.NewRecorder("sleep_daily", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	sleep, err := client.Sleep.GetDaily(ctx, date)
	if err != nil {
		t.Fatalf("GetDaily failed: %v", err)
	}

	if sleep == nil {
		t.Fatal("expected sleep data, got nil")
	}

	// Verify we got actual data
	if sleep.DailySleepDTO.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
}

func TestIntegration_Wellness_GetDailyStress(t *testing.T) {
	skipIfNoCassette(t, "wellness_stress")

	rec, err := testutil.NewRecorder("wellness_stress", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	stress, err := client.Wellness.GetDailyStress(ctx, date)
	if err != nil {
		t.Fatalf("GetDailyStress failed: %v", err)
	}

	if stress == nil {
		t.Fatal("expected stress data, got nil")
	}

	if stress.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
}

func TestIntegration_Wellness_GetBodyBatteryEvents(t *testing.T) {
	skipIfNoCassette(t, "wellness_body_battery")

	rec, err := testutil.NewRecorder("wellness_body_battery", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	events, err := client.Wellness.GetBodyBatteryEvents(ctx, date)
	if err != nil {
		t.Fatalf("GetBodyBatteryEvents failed: %v", err)
	}

	if events == nil {
		t.Fatal("expected body battery events, got nil")
	}
}
