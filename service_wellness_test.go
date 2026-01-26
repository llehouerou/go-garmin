// service_wellness_test.go
package garmin

import (
	"encoding/json"
	"testing"
	"time"
)

const testDate = "2026-01-26"

func TestDailySleepConversions(t *testing.T) {
	sleep := &DailySleep{
		DailySleepDTO: DailySleepDTO{
			CalendarDate:        testDate,
			SleepStartTimestamp: 1737853200000, // 2026-01-26 01:00:00 UTC
			SleepEndTimestamp:   1737882000000, // 2026-01-26 09:00:00 UTC
			SleepSeconds:        28800,         // 8 hours
		},
	}

	if sleep.Duration() != 8*time.Hour {
		t.Errorf("Duration() = %v, want 8h", sleep.Duration())
	}

	start := sleep.SleepStart().UTC()
	if start.Hour() != 1 {
		t.Errorf("SleepStart hour = %d, want 1", start.Hour())
	}

	end := sleep.SleepEnd().UTC()
	if end.Hour() != 9 {
		t.Errorf("SleepEnd hour = %d, want 9", end.Hour())
	}
}

func TestDailySleepHasData(t *testing.T) {
	id := int64(123)
	sleepWithData := &DailySleep{
		DailySleepDTO: DailySleepDTO{ID: &id},
	}
	if !sleepWithData.HasData() {
		t.Error("HasData() should return true when ID is set")
	}

	sleepWithoutData := &DailySleep{}
	if sleepWithoutData.HasData() {
		t.Error("HasData() should return false when ID is nil")
	}
}

func TestDailySleepRawJSON(t *testing.T) {
	rawJSON := `{"dailySleepDTO":{"calendarDate":"2026-01-26","sleepTimeSeconds":28800}}`

	var sleep DailySleep
	if err := json.Unmarshal([]byte(rawJSON), &sleep); err != nil {
		t.Fatal(err)
	}
	sleep.raw = json.RawMessage(rawJSON)

	if string(sleep.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestDailySleepJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"dailySleepDTO": {
			"id": 123456789,
			"calendarDate": "2026-01-26",
			"sleepStartTimestampGMT": 1737853200000,
			"sleepEndTimestampGMT": 1737882000000,
			"sleepTimeSeconds": 28800,
			"deepSleepSeconds": 7200,
			"lightSleepSeconds": 14400,
			"remSleepSeconds": 5400,
			"awakeSleepSeconds": 1800,
			"averageSpO2Value": 96.5
		},
		"remSleepData": true,
		"bodyBatteryChange": 45
	}`

	var sleep DailySleep
	if err := json.Unmarshal([]byte(rawJSON), &sleep); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sleep.DailySleepDTO.CalendarDate != testDate {
		t.Errorf("CalendarDate = %s, want %s", sleep.DailySleepDTO.CalendarDate, testDate)
	}
	if sleep.DailySleepDTO.DeepSleepSeconds == nil || *sleep.DailySleepDTO.DeepSleepSeconds != 7200 {
		t.Errorf("DeepSleepSeconds = %v, want 7200", sleep.DailySleepDTO.DeepSleepSeconds)
	}
	if sleep.DailySleepDTO.LightSleepSeconds == nil || *sleep.DailySleepDTO.LightSleepSeconds != 14400 {
		t.Errorf("LightSleepSeconds = %v, want 14400", sleep.DailySleepDTO.LightSleepSeconds)
	}
	if sleep.DailySleepDTO.REMSleepSeconds == nil || *sleep.DailySleepDTO.REMSleepSeconds != 5400 {
		t.Errorf("REMSleepSeconds = %v, want 5400", sleep.DailySleepDTO.REMSleepSeconds)
	}
	if sleep.DailySleepDTO.AverageSpO2 == nil || *sleep.DailySleepDTO.AverageSpO2 != 96.5 {
		t.Errorf("AverageSpO2 = %v, want 96.5", sleep.DailySleepDTO.AverageSpO2)
	}
	if !sleep.REMSleepData {
		t.Error("REMSleepData should be true")
	}
	if sleep.BodyBatteryChange == nil || *sleep.BodyBatteryChange != 45 {
		t.Errorf("BodyBatteryChange = %v, want 45", sleep.BodyBatteryChange)
	}
}

func TestDailySleepOptionalSpO2(t *testing.T) {
	// Test when SpO2 is null/missing
	rawJSON := `{"dailySleepDTO":{"calendarDate":"2026-01-26","sleepTimeSeconds":28800}}`

	var sleep DailySleep
	if err := json.Unmarshal([]byte(rawJSON), &sleep); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sleep.DailySleepDTO.AverageSpO2 != nil {
		t.Error("AverageSpO2 should be nil when not present")
	}
}

func TestDailyStressJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"calendarDate": "2026-01-26",
		"maxStressLevel": 85,
		"avgStressLevel": 42,
		"stressChartValueOffset": 0,
		"stressChartYAxisOrigin": 0,
		"stressValuesArray": [[1737853200000, 12], [1737856800000, 25]],
		"bodyBatteryValuesArray": [[1737853200000, "charging", 45, 1.0]]
	}`

	var stress DailyStress
	if err := json.Unmarshal([]byte(rawJSON), &stress); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if stress.CalendarDate != testDate {
		t.Errorf("CalendarDate = %s, want %s", stress.CalendarDate, testDate)
	}
	if stress.MaxStressLevel != 85 {
		t.Errorf("MaxStressLevel = %d, want 85", stress.MaxStressLevel)
	}
	if stress.AvgStressLevel != 42 {
		t.Errorf("AvgStressLevel = %d, want 42", stress.AvgStressLevel)
	}
	if len(stress.StressValuesArray) != 2 {
		t.Errorf("StressValuesArray length = %d, want 2", len(stress.StressValuesArray))
	}
	if len(stress.BodyBatteryValuesArray) != 1 {
		t.Errorf("BodyBatteryValuesArray length = %d, want 1", len(stress.BodyBatteryValuesArray))
	}
}

func TestDailyStressRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-26","maxStressLevel":85,"avgStressLevel":42}`

	var stress DailyStress
	if err := json.Unmarshal([]byte(rawJSON), &stress); err != nil {
		t.Fatal(err)
	}
	stress.raw = json.RawMessage(rawJSON)

	if string(stress.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestBodyBatteryEventJSONUnmarshal(t *testing.T) {
	rawJSON := `[{
		"event": {
			"eventType": "sleep",
			"eventStartTimeGmt": "2026-01-25T23:00:00.000",
			"timezoneOffset": -18000000,
			"durationInMilliseconds": 28800000,
			"bodyBatteryImpact": 45,
			"feedbackType": "good_sleep",
			"shortFeedback": "Good sleep restored your Body Battery"
		},
		"activityName": null,
		"activityType": null,
		"activityId": null,
		"averageStress": 15.5,
		"stressValuesArray": [[1737853200000, 12]],
		"bodyBatteryValuesArray": [[1737853200000, "charging", 45, 1.0]]
	}]`

	var events []BodyBatteryEvent
	if err := json.Unmarshal([]byte(rawJSON), &events); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Event == nil {
		t.Fatal("Event should not be nil")
	}
	if event.Event.EventType != "sleep" {
		t.Errorf("EventType = %s, want sleep", event.Event.EventType)
	}
	if event.Event.BodyBatteryImpact != 45 {
		t.Errorf("BodyBatteryImpact = %d, want 45", event.Event.BodyBatteryImpact)
	}
	if event.AverageStress == nil || *event.AverageStress != 15.5 {
		t.Errorf("AverageStress = %v, want 15.5", event.AverageStress)
	}
}

func TestBodyBatteryEventsRawJSON(t *testing.T) {
	rawJSON := `[{"event":{"eventType":"sleep"}}]`

	var events []BodyBatteryEvent
	if err := json.Unmarshal([]byte(rawJSON), &events); err != nil {
		t.Fatal(err)
	}

	bb := &BodyBatteryEvents{
		Events: events,
		raw:    json.RawMessage(rawJSON),
	}

	if string(bb.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
