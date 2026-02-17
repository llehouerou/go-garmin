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
	activityID := int64(21680374805)

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

func TestIntegration_Activity_GetWeather(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	weather, err := client.Activities.GetWeather(ctx, activityID)
	if err != nil {
		t.Fatalf("GetWeather failed: %v", err)
	}

	if weather == nil {
		t.Fatal("expected weather data, got nil")
	}

	// Verify we got actual weather data
	if weather.IssueDate == "" {
		t.Error("expected IssueDate to be set")
	}
	if weather.WindDirectionCompassPoint == "" {
		t.Error("expected WindDirectionCompassPoint to be set")
	}
	if weather.WeatherStationDTO.ID == "" {
		t.Error("expected WeatherStationDTO.ID to be set")
	}

	// Verify conversion methods work
	tempC := weather.TempCelsius()
	if tempC < -50 || tempC > 60 {
		t.Errorf("TempCelsius() = %v, seems unreasonable", tempC)
	}

	// Verify RawJSON is available
	if weather.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_GetSplits(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	splits, err := client.Activities.GetSplits(ctx, activityID)
	if err != nil {
		t.Fatalf("GetSplits failed: %v", err)
	}

	if splits == nil {
		t.Fatal("expected splits data, got nil")
	}

	if splits.ActivityID != activityID {
		t.Errorf("ActivityID = %d, want %d", splits.ActivityID, activityID)
	}

	// Verify we got lap data
	if len(splits.LapDTOs) == 0 {
		t.Error("expected LapDTOs to have at least one lap")
	}

	// Verify first lap has expected fields
	if len(splits.LapDTOs) > 0 {
		firstLap := splits.LapDTOs[0]
		if firstLap.Duration == 0 {
			t.Error("expected lap Duration to be set")
		}
		if firstLap.Distance == 0 {
			t.Error("expected lap Distance to be set")
		}

		// Verify conversion methods work
		dur := firstLap.DurationTime()
		if dur <= 0 {
			t.Errorf("DurationTime() = %v, expected positive", dur)
		}
		distKm := firstLap.DistanceKm()
		if distKm <= 0 {
			t.Errorf("DistanceKm() = %v, expected positive", distKm)
		}
	}

	// Verify RawJSON is available
	if splits.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Wellness_GetDailyHeartRate(t *testing.T) {
	skipIfNoCassette(t, "wellness_heart_rate")

	rec, err := testutil.NewRecorder("wellness_heart_rate", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	hr, err := client.Wellness.GetDailyHeartRate(ctx, date)
	if err != nil {
		t.Fatalf("GetDailyHeartRate failed: %v", err)
	}

	if hr == nil {
		t.Fatal("expected heart rate data, got nil")
	}

	if hr.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if hr.MaxHeartRate == 0 {
		t.Error("expected MaxHeartRate to be set")
	}
	if hr.MinHeartRate == 0 {
		t.Error("expected MinHeartRate to be set")
	}
	if hr.RestingHeartRate == 0 {
		t.Error("expected RestingHeartRate to be set")
	}
	if len(hr.HeartRateValues) == 0 {
		t.Error("expected HeartRateValues to have data")
	}

	// Verify RawJSON is available
	if hr.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_HRV_GetDaily(t *testing.T) {
	skipIfNoCassette(t, "hrv")

	rec, err := testutil.NewRecorder("hrv", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	hrv, err := client.HRV.GetDaily(ctx, date)
	if err != nil {
		t.Fatalf("GetDaily failed: %v", err)
	}

	if hrv == nil {
		t.Fatal("expected HRV data, got nil")
	}

	if hrv.HRVSummary.CalendarDate == "" {
		t.Error("expected HRVSummary.CalendarDate to be set")
	}
	if hrv.HRVSummary.Status == "" {
		t.Error("expected HRVSummary.Status to be set")
	}
	if hrv.HRVSummary.WeeklyAvg == 0 {
		t.Error("expected HRVSummary.WeeklyAvg to be set")
	}
	if len(hrv.HRVReadings) == 0 {
		t.Error("expected HRVReadings to have data")
	}

	// Verify RawJSON is available
	if hrv.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_HRV_GetRange(t *testing.T) {
	skipIfNoCassette(t, "hrv")

	rec, err := testutil.NewRecorder("hrv", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	endDate := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, 0, -7)

	hrvRange, err := client.HRV.GetRange(ctx, startDate, endDate)
	if err != nil {
		t.Fatalf("GetRange failed: %v", err)
	}

	if hrvRange == nil {
		t.Fatal("expected HRV range data, got nil")
	}

	if len(hrvRange.HRVSummaries) == 0 {
		t.Error("expected HRVSummaries to have data")
	}

	// Verify each summary has expected fields
	for i, summary := range hrvRange.HRVSummaries {
		if summary.CalendarDate == "" {
			t.Errorf("HRVSummaries[%d].CalendarDate is empty", i)
		}
		if summary.Status == "" {
			t.Errorf("HRVSummaries[%d].Status is empty", i)
		}
	}

	// Verify RawJSON is available
	if hrvRange.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Weight_GetDaily(t *testing.T) {
	skipIfNoCassette(t, "weight")

	rec, err := testutil.NewRecorder("weight", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	weight, err := client.Weight.GetDaily(ctx, date)
	if err != nil {
		t.Fatalf("GetDaily failed: %v", err)
	}

	if weight == nil {
		t.Fatal("expected weight data, got nil")
	}

	if weight.StartDate == "" {
		t.Error("expected StartDate to be set")
	}
	if weight.EndDate == "" {
		t.Error("expected EndDate to be set")
	}

	// Verify we have weight entries
	if len(weight.DateWeightList) == 0 {
		t.Error("expected DateWeightList to have data")
	}

	// Verify first entry has expected fields
	if len(weight.DateWeightList) > 0 {
		entry := weight.DateWeightList[0]
		if entry.CalendarDate == "" {
			t.Error("expected entry CalendarDate to be set")
		}
		if entry.Weight == nil {
			t.Error("expected entry Weight to be set")
		}

		// Verify conversion methods work
		kg := entry.WeightKg()
		if kg <= 0 {
			t.Errorf("WeightKg() = %v, expected positive", kg)
		}
	}

	// Verify RawJSON is available
	if weight.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Weight_GetRange(t *testing.T) {
	skipIfNoCassette(t, "weight")

	rec, err := testutil.NewRecorder("weight", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	endDate := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, 0, -30)

	weightRange, err := client.Weight.GetRange(ctx, startDate, endDate)
	if err != nil {
		t.Fatalf("GetRange failed: %v", err)
	}

	if weightRange == nil {
		t.Fatal("expected weight range data, got nil")
	}

	// Verify we have summaries
	if len(weightRange.DailyWeightSummaries) == 0 {
		t.Error("expected DailyWeightSummaries to have data")
	}

	// Verify first summary has expected fields
	if len(weightRange.DailyWeightSummaries) > 0 {
		summary := weightRange.DailyWeightSummaries[0]
		if summary.SummaryDate == "" {
			t.Error("expected summary SummaryDate to be set")
		}
		if summary.MinWeight == 0 {
			t.Error("expected summary MinWeight to be set")
		}
	}

	// Verify RawJSON is available
	if weightRange.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetTrainingReadiness(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	readiness, err := client.Metrics.GetTrainingReadiness(ctx, date)
	if err != nil {
		t.Fatalf("GetTrainingReadiness failed: %v", err)
	}

	if readiness == nil {
		t.Fatal("expected training readiness data, got nil")
	}

	if len(readiness.Entries) == 0 {
		t.Error("expected Entries to have data")
	}

	// Verify first entry has expected fields
	if len(readiness.Entries) > 0 {
		entry := readiness.Entries[0]
		if entry.CalendarDate == "" {
			t.Error("expected entry CalendarDate to be set")
		}
		if entry.Level == "" {
			t.Error("expected entry Level to be set")
		}
		if entry.Score == 0 {
			t.Error("expected entry Score to be set")
		}
	}

	if readiness.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetEnduranceScore(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	score, err := client.Metrics.GetEnduranceScore(ctx, date)
	if err != nil {
		t.Fatalf("GetEnduranceScore failed: %v", err)
	}

	if score == nil {
		t.Fatal("expected endurance score data, got nil")
	}

	if score.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if score.OverallScore == 0 {
		t.Error("expected OverallScore to be set")
	}

	if score.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetEnduranceScoreStats(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	endDate := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, 0, -84)

	stats, err := client.Metrics.GetEnduranceScoreStats(ctx, startDate, endDate, AggregationWeekly)
	if err != nil {
		t.Fatalf("GetEnduranceScoreStats failed: %v", err)
	}

	if stats == nil {
		t.Fatal("expected endurance score stats, got nil")
	}

	if stats.StartDate == "" {
		t.Error("expected StartDate to be set")
	}
	if stats.EndDate == "" {
		t.Error("expected EndDate to be set")
	}
	if stats.Avg == 0 {
		t.Error("expected Avg to be set")
	}
	if stats.Max == 0 {
		t.Error("expected Max to be set")
	}
	if len(stats.GroupMap) == 0 {
		t.Error("expected GroupMap to have entries")
	}

	if stats.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetHillScore(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	score, err := client.Metrics.GetHillScore(ctx, date)
	if err != nil {
		t.Fatalf("GetHillScore failed: %v", err)
	}

	if score == nil {
		t.Fatal("expected hill score data, got nil")
	}

	if score.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if score.VO2Max == 0 {
		t.Error("expected VO2Max to be set")
	}

	if score.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetMaxMetLatest(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	maxMet, err := client.Metrics.GetMaxMetLatest(ctx, date)
	if err != nil {
		t.Fatalf("GetMaxMetLatest failed: %v", err)
	}

	if maxMet == nil {
		t.Fatal("expected max met data, got nil")
	}

	if maxMet.Generic == nil {
		t.Error("expected Generic to be set")
	}
	if maxMet.Generic != nil && maxMet.Generic.VO2MaxValue == 0 {
		t.Error("expected Generic.VO2MaxValue to be set")
	}

	if maxMet.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetMaxMetDaily(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	endDate := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, 0, -30)

	maxMet, err := client.Metrics.GetMaxMetDaily(ctx, startDate, endDate)
	if err != nil {
		t.Fatalf("GetMaxMetDaily failed: %v", err)
	}

	if maxMet == nil {
		t.Fatal("expected max met daily data, got nil")
	}

	if len(maxMet.Entries) == 0 {
		t.Error("expected Entries to have data")
	}

	if maxMet.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetTrainingStatusAggregated(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	status, err := client.Metrics.GetTrainingStatusAggregated(ctx, date)
	if err != nil {
		t.Fatalf("GetTrainingStatusAggregated failed: %v", err)
	}

	if status == nil {
		t.Fatal("expected training status aggregated data, got nil")
	}

	if status.MostRecentVO2Max == nil {
		t.Error("expected MostRecentVO2Max to be set")
	}
	if status.MostRecentTrainingStatus == nil {
		t.Error("expected MostRecentTrainingStatus to be set")
	}

	if status.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetTrainingStatusDaily(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	status, err := client.Metrics.GetTrainingStatusDaily(ctx, date)
	if err != nil {
		t.Fatalf("GetTrainingStatusDaily failed: %v", err)
	}

	if status == nil {
		t.Fatal("expected training status daily data, got nil")
	}

	if len(status.LatestTrainingStatusData) == 0 {
		t.Error("expected LatestTrainingStatusData to have data")
	}

	if status.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetTrainingLoadBalance(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	balance, err := client.Metrics.GetTrainingLoadBalance(ctx, date)
	if err != nil {
		t.Fatalf("GetTrainingLoadBalance failed: %v", err)
	}

	if balance == nil {
		t.Fatal("expected training load balance data, got nil")
	}

	if len(balance.MetricsTrainingLoadBalanceDTOMap) == 0 {
		t.Error("expected MetricsTrainingLoadBalanceDTOMap to have data")
	}

	if balance.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Metrics_GetHeatAltitudeAcclimation(t *testing.T) {
	skipIfNoCassette(t, "metrics")

	rec, err := testutil.NewRecorder("metrics", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()
	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	acclimation, err := client.Metrics.GetHeatAltitudeAcclimation(ctx, date)
	if err != nil {
		t.Fatalf("GetHeatAltitudeAcclimation failed: %v", err)
	}

	if acclimation == nil {
		t.Fatal("expected heat altitude acclimation data, got nil")
	}

	if acclimation.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if acclimation.HeatTrend == "" {
		t.Error("expected HeatTrend to be set")
	}

	if acclimation.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_UserProfile_GetSocialProfile(t *testing.T) {
	skipIfNoCassette(t, "userprofile")

	rec, err := testutil.NewRecorder("userprofile", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	profile, err := client.UserProfile.GetSocialProfile(ctx)
	if err != nil {
		t.Fatalf("GetSocialProfile failed: %v", err)
	}

	if profile == nil {
		t.Fatal("expected social profile data, got nil")
	}

	if profile.ID == 0 {
		t.Error("expected ID to be set")
	}
	if profile.DisplayName == "" {
		t.Error("expected DisplayName to be set")
	}
	if profile.UserName == "" {
		t.Error("expected UserName to be set")
	}
	if profile.ProfileVisibility == "" {
		t.Error("expected ProfileVisibility to be set")
	}

	if profile.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_UserProfile_GetUserSettings(t *testing.T) {
	skipIfNoCassette(t, "userprofile")

	rec, err := testutil.NewRecorder("userprofile", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	settings, err := client.UserProfile.GetUserSettings(ctx)
	if err != nil {
		t.Fatalf("GetUserSettings failed: %v", err)
	}

	if settings == nil {
		t.Fatal("expected user settings data, got nil")
	}

	if settings.ID == 0 {
		t.Error("expected ID to be set")
	}
	if settings.UserData.Gender == "" {
		t.Error("expected Gender to be set")
	}
	if settings.UserData.MeasurementSystem == "" {
		t.Error("expected MeasurementSystem to be set")
	}
	if settings.UserSleep.SleepTime == 0 && settings.UserSleep.WakeTime == 0 {
		t.Error("expected sleep times to be set")
	}

	if settings.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_UserProfile_GetProfileSettings(t *testing.T) {
	skipIfNoCassette(t, "userprofile")

	rec, err := testutil.NewRecorder("userprofile", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	settings, err := client.UserProfile.GetProfileSettings(ctx)
	if err != nil {
		t.Fatalf("GetProfileSettings failed: %v", err)
	}

	if settings == nil {
		t.Fatal("expected profile settings data, got nil")
	}

	if settings.DisplayName == "" {
		t.Error("expected DisplayName to be set")
	}
	if settings.MeasurementSystem == "" {
		t.Error("expected MeasurementSystem to be set")
	}
	if settings.TimeZone == "" {
		t.Error("expected TimeZone to be set")
	}

	if settings.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Device_GetDevices(t *testing.T) {
	skipIfNoCassette(t, "devices")

	rec, err := testutil.NewRecorder("devices", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	devices, err := client.Devices.GetDevices(ctx)
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}

	if len(devices) == 0 {
		t.Fatal("expected devices, got none")
	}

	// Verify first device has expected fields
	first := devices[0]
	if first.DeviceID == 0 {
		t.Error("expected DeviceID to be set")
	}
	if first.ProductDisplayName == "" {
		t.Error("expected ProductDisplayName to be set")
	}
	if first.DeviceStatus == "" {
		t.Error("expected DeviceStatus to be set")
	}
	if len(first.DeviceCategories) == 0 {
		t.Error("expected DeviceCategories to be set")
	}

	// Verify RawJSON is available
	if first.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Device_GetSettings(t *testing.T) {
	skipIfNoCassette(t, "devices")

	rec, err := testutil.NewRecorder("devices", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Device ID from the recorded cassette (anonymized to 12345678)
	deviceID := int64(12345678)

	settings, err := client.Devices.GetSettings(ctx, deviceID)
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}

	if settings == nil {
		t.Fatal("expected device settings, got nil")
	}

	if settings.DeviceID == 0 {
		t.Error("expected DeviceID to be set")
	}
	if settings.TimeFormat == "" {
		t.Error("expected TimeFormat to be set")
	}
	if settings.MeasurementUnits == "" {
		t.Error("expected MeasurementUnits to be set")
	}
	if settings.StartOfWeek == "" {
		t.Error("expected StartOfWeek to be set")
	}
	if len(settings.SupportedLanguages) == 0 {
		t.Error("expected SupportedLanguages to be set")
	}

	// Verify RawJSON is available
	if settings.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Device_GetMessages(t *testing.T) {
	skipIfNoCassette(t, "devices")

	rec, err := testutil.NewRecorder("devices", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	messages, err := client.Devices.GetMessages(ctx)
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}

	if messages == nil {
		t.Fatal("expected device messages, got nil")
	}

	if messages.ServiceHost == "" {
		t.Error("expected ServiceHost to be set")
	}

	// NumOfMessages can be 0, so we just check it's accessible
	_ = messages.NumOfMessages

	// Verify RawJSON is available
	if messages.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Device_GetPrimaryTrainingDevice(t *testing.T) {
	skipIfNoCassette(t, "devices")

	rec, err := testutil.NewRecorder("devices", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	info, err := client.Devices.GetPrimaryTrainingDevice(ctx)
	if err != nil {
		t.Fatalf("GetPrimaryTrainingDevice failed: %v", err)
	}

	if info == nil {
		t.Fatal("expected primary training device info, got nil")
	}

	if info.PrimaryTrainingDevice.DeviceID == 0 {
		t.Error("expected PrimaryTrainingDevice.DeviceID to be set")
	}
	if len(info.WearableDevices.DeviceWeights) == 0 {
		t.Error("expected WearableDevices to have at least one device")
	}
	if len(info.PrimaryTrainingDevices.DeviceWeights) == 0 {
		t.Error("expected PrimaryTrainingDevices to have at least one device")
	}

	// Verify RawJSON is available
	if info.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Wellness_GetDailySpO2(t *testing.T) {
	skipIfNoCassette(t, "wellness_extended")

	rec, err := testutil.NewRecorder("wellness_extended", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	spo2, err := client.Wellness.GetDailySpO2(ctx, date)
	if err != nil {
		t.Fatalf("GetDailySpO2 failed: %v", err)
	}

	if spo2 == nil {
		t.Fatal("expected SpO2 data, got nil")
	}

	if spo2.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if spo2.AverageSpO2 == 0 {
		t.Error("expected AverageSpO2 to be set")
	}
	if spo2.LowestSpO2 == 0 {
		t.Error("expected LowestSpO2 to be set")
	}
	if spo2.LatestSpO2 == 0 {
		t.Error("expected LatestSpO2 to be set")
	}

	if spo2.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Wellness_GetDailyRespiration(t *testing.T) {
	skipIfNoCassette(t, "wellness_extended")

	rec, err := testutil.NewRecorder("wellness_extended", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	resp, err := client.Wellness.GetDailyRespiration(ctx, date)
	if err != nil {
		t.Fatalf("GetDailyRespiration failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected respiration data, got nil")
	}

	if resp.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if resp.LowestRespirationValue == 0 {
		t.Error("expected LowestRespirationValue to be set")
	}
	if resp.HighestRespirationValue == 0 {
		t.Error("expected HighestRespirationValue to be set")
	}
	if len(resp.RespirationValuesArray) == 0 {
		t.Error("expected RespirationValuesArray to have data")
	}

	if resp.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Wellness_GetDailyIntensityMinutes(t *testing.T) {
	skipIfNoCassette(t, "wellness_extended")

	rec, err := testutil.NewRecorder("wellness_extended", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	im, err := client.Wellness.GetDailyIntensityMinutes(ctx, date)
	if err != nil {
		t.Fatalf("GetDailyIntensityMinutes failed: %v", err)
	}

	if im == nil {
		t.Fatal("expected intensity minutes data, got nil")
	}

	if im.CalendarDate == "" {
		t.Error("expected CalendarDate to be set")
	}
	if im.WeekGoal == 0 {
		t.Error("expected WeekGoal to be set")
	}
	// WeeklyTotal and minutes can be 0, so just check they're accessible
	_ = im.WeeklyTotal
	_ = im.ModerateMinutes
	_ = im.VigorousMinutes

	if im.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_GetDetails(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	details, err := client.Activities.GetDetails(ctx, activityID, nil)
	if err != nil {
		t.Fatalf("GetDetails failed: %v", err)
	}

	if details == nil {
		t.Fatal("expected activity details, got nil")
	}

	if details.ActivityID != activityID {
		t.Errorf("ActivityID = %d, want %d", details.ActivityID, activityID)
	}
	if details.MeasurementCount == 0 {
		t.Error("expected MeasurementCount to be set")
	}
	if len(details.MetricDescriptors) == 0 {
		t.Error("expected MetricDescriptors to have data")
	}
	if len(details.ActivityDetailMetrics) == 0 {
		t.Error("expected ActivityDetailMetrics to have data")
	}

	// Verify GetMetricIndex works
	timestampIdx := details.GetMetricIndex("directTimestamp")
	if timestampIdx == -1 {
		t.Error("expected to find directTimestamp metric")
	}

	// Verify RawJSON is available
	if details.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_GetHRTimeInZones(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	hrZones, err := client.Activities.GetHRTimeInZones(ctx, activityID)
	if err != nil {
		t.Fatalf("GetHRTimeInZones failed: %v", err)
	}

	if hrZones == nil {
		t.Fatal("expected HR time in zones, got nil")
	}

	if len(hrZones.Zones) == 0 {
		t.Error("expected Zones to have data")
	}

	// Verify first zone has expected fields
	if len(hrZones.Zones) > 0 {
		firstZone := hrZones.Zones[0]
		if firstZone.ZoneNumber == 0 {
			t.Error("expected ZoneNumber to be set")
		}
		if firstZone.ZoneLowBoundary == 0 {
			t.Error("expected ZoneLowBoundary to be set")
		}

		// Verify conversion method works
		dur := firstZone.DurationInZone()
		if dur < 0 {
			t.Errorf("DurationInZone() = %v, expected non-negative", dur)
		}
	}

	// Verify RawJSON is available
	if hrZones.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_GetPowerTimeInZones(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	powerZones, err := client.Activities.GetPowerTimeInZones(ctx, activityID)
	if err != nil {
		t.Fatalf("GetPowerTimeInZones failed: %v", err)
	}

	if powerZones == nil {
		t.Fatal("expected power time in zones, got nil")
	}

	// Power zones may be empty if activity doesn't have power data
	// Just verify RawJSON is available
	if powerZones.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Activity_GetExerciseSets(t *testing.T) {
	skipIfNoCassette(t, "activities")

	rec, err := testutil.NewRecorder("activities", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Activity ID from the recorded cassette
	activityID := int64(21680374805)

	sets, err := client.Activities.GetExerciseSets(ctx, activityID)
	if err != nil {
		t.Fatalf("GetExerciseSets failed: %v", err)
	}

	if sets == nil {
		t.Fatal("expected exercise sets, got nil")
	}

	if sets.ActivityID != activityID {
		t.Errorf("ActivityID = %d, want %d", sets.ActivityID, activityID)
	}

	// Exercise sets may be null/empty for cardio activities
	// Just verify RawJSON is available
	if sets.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetLatestLactateThreshold(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	lt, err := client.Biometric.GetLatestLactateThreshold(ctx)
	if err != nil {
		t.Fatalf("GetLatestLactateThreshold failed: %v", err)
	}

	if lt == nil {
		t.Fatal("expected lactate threshold data, got nil")
	}

	if len(lt.Entries) == 0 {
		t.Error("expected Entries to have data")
	}

	// Verify helper methods work
	if speed := lt.Speed(); speed == nil {
		t.Error("expected Speed() to return a value")
	}
	if hr := lt.HeartRate(); hr == nil {
		t.Error("expected HeartRate() to return a value")
	}

	if lt.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetCyclingFTP(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	ftp, err := client.Biometric.GetCyclingFTP(ctx)
	if err != nil {
		t.Fatalf("GetCyclingFTP failed: %v", err)
	}

	if ftp == nil {
		t.Fatal("expected FTP data, got nil")
	}

	// FTP value may be nil if user has no cycling FTP data
	// Just verify the struct was populated and RawJSON works
	if ftp.UserProfilePK == 0 {
		t.Error("expected UserProfilePK to be set")
	}

	if ftp.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetPowerToWeight(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	date := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	ptw, err := client.Biometric.GetPowerToWeight(ctx, date)
	if err != nil {
		t.Fatalf("GetPowerToWeight failed: %v", err)
	}

	if ptw == nil {
		t.Fatal("expected power to weight data, got nil")
	}

	if ptw.FunctionalThresholdPower == 0 {
		t.Error("expected FunctionalThresholdPower to be set")
	}
	if ptw.Weight == 0 {
		t.Error("expected Weight to be set")
	}
	if ptw.PowerToWeightRatio == 0 {
		t.Error("expected PowerToWeightRatio to be set")
	}
	if ptw.Sport == "" {
		t.Error("expected Sport to be set")
	}

	if ptw.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetLactateThresholdSpeedRange(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	end := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -30)

	stats, err := client.Biometric.GetLactateThresholdSpeedRange(ctx, start, end)
	if err != nil {
		t.Fatalf("GetLactateThresholdSpeedRange failed: %v", err)
	}

	if stats == nil {
		t.Fatal("expected stats, got nil")
	}

	if len(stats.Stats) == 0 {
		t.Error("expected Stats to have data")
	}

	if len(stats.Stats) > 0 {
		stat := stats.Stats[0]
		if stat.Series == "" {
			t.Error("expected Series to be set")
		}
		if stat.Value == 0 {
			t.Error("expected Value to be set")
		}
	}

	if stats.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetLactateThresholdHRRange(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	end := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -30)

	stats, err := client.Biometric.GetLactateThresholdHRRange(ctx, start, end)
	if err != nil {
		t.Fatalf("GetLactateThresholdHRRange failed: %v", err)
	}

	if stats == nil {
		t.Fatal("expected stats, got nil")
	}

	if len(stats.Stats) == 0 {
		t.Error("expected Stats to have data")
	}

	if stats.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetFTPRange(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	end := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -30)

	stats, err := client.Biometric.GetFTPRange(ctx, start, end)
	if err != nil {
		t.Fatalf("GetFTPRange failed: %v", err)
	}

	if stats == nil {
		t.Fatal("expected stats, got nil")
	}

	if len(stats.Stats) == 0 {
		t.Error("expected Stats to have data")
	}

	if stats.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Biometric_GetHeartRateZones(t *testing.T) {
	skipIfNoCassette(t, "biometric")

	rec, err := testutil.NewRecorder("biometric", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	zones, err := client.Biometric.GetHeartRateZones(ctx)
	if err != nil {
		t.Fatalf("GetHeartRateZones failed: %v", err)
	}

	if zones == nil {
		t.Fatal("expected zones, got nil")
	}

	if len(zones.Zones) == 0 {
		t.Error("expected Zones to have data")
	}

	// Verify first zone has expected fields
	if len(zones.Zones) > 0 {
		zone := zones.Zones[0]
		if zone.TrainingMethod == "" {
			t.Error("expected TrainingMethod to be set")
		}
		if zone.Sport == "" {
			t.Error("expected Sport to be set")
		}
		if zone.MaxHeartRateUsed == 0 {
			t.Error("expected MaxHeartRateUsed to be set")
		}
	}

	if zones.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Workout_List(t *testing.T) {
	skipIfNoCassette(t, "workouts")

	rec, err := testutil.NewRecorder("workouts", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	list, err := client.Workouts.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if list == nil {
		t.Fatal("expected workout list, got nil")
	}

	if len(list.Workouts) == 0 {
		t.Error("expected Workouts to have data")
	}

	// Verify first workout has expected fields
	if len(list.Workouts) > 0 {
		first := list.Workouts[0]
		if first.WorkoutID == 0 {
			t.Error("expected WorkoutID to be set")
		}
		if first.WorkoutName == "" {
			t.Error("expected WorkoutName to be set")
		}
		if first.SportType.SportTypeKey == "" {
			t.Error("expected SportType.SportTypeKey to be set")
		}
	}

	if list.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Workout_Get(t *testing.T) {
	skipIfNoCassette(t, "workouts")

	rec, err := testutil.NewRecorder("workouts", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// First get the list to find a workout ID
	list, err := client.Workouts.List(ctx, 0, 1)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list.Workouts) == 0 {
		t.Skip("no workouts available to test Get")
	}

	workoutID := list.Workouts[0].WorkoutID

	workout, err := client.Workouts.Get(ctx, workoutID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if workout == nil {
		t.Fatal("expected workout, got nil")
	}

	if workout.WorkoutID != workoutID {
		t.Errorf("WorkoutID = %d, want %d", workout.WorkoutID, workoutID)
	}
	if workout.WorkoutName == "" {
		t.Error("expected WorkoutName to be set")
	}
	if workout.SportType.SportTypeKey == "" {
		t.Error("expected SportType.SportTypeKey to be set")
	}
	if len(workout.WorkoutSegments) == 0 {
		t.Error("expected WorkoutSegments to have data")
	}

	if workout.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}

func TestIntegration_Course_DownloadGPX(t *testing.T) {
	skipIfNoCassette(t, "courses_download")

	rec, err := testutil.NewRecorder("courses_download", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Course ID from the recorded cassette (anonymized to 87654321)
	courseID := int64(87654321)

	data, err := client.Courses.DownloadGPX(ctx, courseID)
	if err != nil {
		t.Fatalf("DownloadGPX failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected GPX data, got empty")
	}

	// GPX files are XML and should contain the <gpx tag
	if !strings.Contains(string(data), "<gpx") {
		t.Error("expected GPX data to contain <gpx tag")
	}
}

func TestIntegration_Course_DownloadFIT(t *testing.T) {
	skipIfNoCassette(t, "courses_download")

	rec, err := testutil.NewRecorder("courses_download", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Course ID from the recorded cassette (anonymized to 87654321)
	courseID := int64(87654321)

	data, err := client.Courses.DownloadFIT(ctx, courseID)
	if err != nil {
		t.Fatalf("DownloadFIT failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected FIT data, got empty")
	}

	// FIT files have ".FIT" signature at bytes 8-11
	if len(data) >= 12 && string(data[8:12]) != ".FIT" {
		t.Errorf("expected FIT header signature '.FIT' at bytes 8-11, got %q", string(data[8:12]))
	}
}

func TestIntegration_Calendar_Get(t *testing.T) {
	skipIfNoCassette(t, "calendar")

	rec, err := testutil.NewRecorder("calendar", recorder.ModeReplayOnly)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	defer func() { _ = rec.Stop() }()

	client := newTestClient(t, rec)
	ctx := context.Background()

	// Get calendar for January 2026, week containing day 28
	month := 0
	day := 28
	start := 1
	calendar, err := client.Calendar.Get(ctx, 2026, &CalendarOptions{
		Month: &month,
		Day:   &day,
		Start: &start,
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if calendar == nil {
		t.Fatal("expected calendar, got nil")
	}

	if calendar.StartDate == "" {
		t.Error("expected StartDate to be set")
	}
	if calendar.EndDate == "" {
		t.Error("expected EndDate to be set")
	}
	if calendar.NumOfDaysInMonth == 0 {
		t.Error("expected NumOfDaysInMonth to be set")
	}
	if len(calendar.CalendarItems) == 0 {
		t.Error("expected CalendarItems to have data")
	}

	// Verify we have different item types
	itemTypes := make(map[string]bool)
	for _, item := range calendar.CalendarItems {
		itemTypes[item.ItemType] = true
	}
	if len(itemTypes) < 2 {
		t.Errorf("expected multiple item types, got %v", itemTypes)
	}

	if calendar.RawJSON() == nil {
		t.Error("expected RawJSON to be available")
	}
}
