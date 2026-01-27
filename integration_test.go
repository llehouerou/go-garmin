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

func TestIntegration_Activity_List(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	activities, err := client.Activities.List(ctx, &ListOptions{Start: 0, Limit: 5})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(activities) == 0 {
		t.Fatal("expected activities, got none")
	}

	// Verify first activity has expected fields
	first := activities[0]
	if first.ActivityID == 0 {
		t.Error("expected ActivityID to be set")
	}
	if first.ActivityName == "" {
		t.Error("expected ActivityName to be set")
	}
	if first.ActivityType.TypeKey == "" {
		t.Error("expected ActivityType.TypeKey to be set")
	}
	if first.Distance == 0 {
		t.Error("expected Distance to be set")
	}
	if first.Duration == 0 {
		t.Error("expected Duration to be set")
	}

	// Verify RawJSON is available
	if first.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_Get(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21661023200)

	detail, err := client.Activities.Get(ctx, activityID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if detail == nil {
		t.Fatal("expected activity detail, got nil")
	}

	if detail.ActivityID != activityID {
		t.Errorf("ActivityID = %d, want %d", detail.ActivityID, activityID)
	}
	if detail.ActivityName == "" {
		t.Error("expected ActivityName to be set")
	}
	if detail.ActivityTypeDTO.TypeKey == "" {
		t.Error("expected ActivityTypeDTO.TypeKey to be set")
	}
	if detail.SummaryDTO.Distance == 0 {
		t.Error("expected SummaryDTO.Distance to be set")
	}
	if detail.SummaryDTO.Duration == 0 {
		t.Error("expected SummaryDTO.Duration to be set")
	}

	// Verify RawJSON is available
	if detail.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}
