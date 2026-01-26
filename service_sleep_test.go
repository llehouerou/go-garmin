// service_sleep_test.go
package garmin

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDailySleepConversions(t *testing.T) {
	sleep := &DailySleep{
		DailySleepDTO: DailySleepDTO{
			CalendarDate:        "2026-01-26",
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

	if sleep.DailySleepDTO.CalendarDate != "2026-01-26" {
		t.Errorf("CalendarDate = %s, want 2026-01-26", sleep.DailySleepDTO.CalendarDate)
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
