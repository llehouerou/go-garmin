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
		CalendarDate:        testDate,
		SleepStartTimestamp: 1737853200000, // 2026-01-26 01:00:00 UTC
		SleepEndTimestamp:   1737882000000, // 2026-01-26 09:00:00 UTC
		SleepSeconds:        28800,         // 8 hours
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

func TestDailySleepRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-26","sleepTimeSeconds":28800}`

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
		"calendarDate": "2026-01-26",
		"sleepStartTimestampGMT": 1737853200000,
		"sleepEndTimestampGMT": 1737882000000,
		"sleepTimeSeconds": 28800,
		"deepSleepSeconds": 7200,
		"lightSleepSeconds": 14400,
		"remSleepSeconds": 5400,
		"awakeSleepSeconds": 1800,
		"averageSpO2Value": 96.5
	}`

	var sleep DailySleep
	if err := json.Unmarshal([]byte(rawJSON), &sleep); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sleep.CalendarDate != testDate {
		t.Errorf("CalendarDate = %s, want %s", sleep.CalendarDate, testDate)
	}
	if sleep.DeepSleepSeconds != 7200 {
		t.Errorf("DeepSleepSeconds = %d, want 7200", sleep.DeepSleepSeconds)
	}
	if sleep.LightSleepSeconds != 14400 {
		t.Errorf("LightSleepSeconds = %d, want 14400", sleep.LightSleepSeconds)
	}
	if sleep.REMSleepSeconds != 5400 {
		t.Errorf("REMSleepSeconds = %d, want 5400", sleep.REMSleepSeconds)
	}
	if sleep.AwakeSeconds != 1800 {
		t.Errorf("AwakeSeconds = %d, want 1800", sleep.AwakeSeconds)
	}
	if sleep.AverageSpO2 == nil || *sleep.AverageSpO2 != 96.5 {
		t.Errorf("AverageSpO2 = %v, want 96.5", sleep.AverageSpO2)
	}
}

func TestDailySleepOptionalSpO2(t *testing.T) {
	// Test when SpO2 is null/missing
	rawJSON := `{"calendarDate":"2026-01-26","sleepTimeSeconds":28800}`

	var sleep DailySleep
	if err := json.Unmarshal([]byte(rawJSON), &sleep); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sleep.AverageSpO2 != nil {
		t.Error("AverageSpO2 should be nil when not present")
	}
}

func TestDailyStressJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"calendarDate": "2026-01-26",
		"overallStressLevel": 42,
		"highStressDuration": 3600,
		"mediumStressDuration": 7200,
		"lowStressDuration": 14400,
		"restStressDuration": 18000
	}`

	var stress DailyStress
	if err := json.Unmarshal([]byte(rawJSON), &stress); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if stress.CalendarDate != testDate {
		t.Errorf("CalendarDate = %s, want %s", stress.CalendarDate, testDate)
	}
	if stress.OverallStressLevel != 42 {
		t.Errorf("OverallStressLevel = %d, want 42", stress.OverallStressLevel)
	}
	if stress.HighStressDuration != 3600 {
		t.Errorf("HighStressDuration = %d, want 3600", stress.HighStressDuration)
	}
	if stress.MedStressDuration != 7200 {
		t.Errorf("MedStressDuration = %d, want 7200", stress.MedStressDuration)
	}
}

func TestDailyStressRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-26","overallStressLevel":42}`

	var stress DailyStress
	if err := json.Unmarshal([]byte(rawJSON), &stress); err != nil {
		t.Fatal(err)
	}
	stress.raw = json.RawMessage(rawJSON)

	if string(stress.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestBodyBatteryReportJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"date": "2026-01-26",
		"charged": 60,
		"drained": 45,
		"startOfDayBodyBattery": 85,
		"endOfDayBodyBattery": 40,
		"maxBodyBattery": 100,
		"minBodyBattery": 25
	}`

	var battery BodyBatteryReport
	if err := json.Unmarshal([]byte(rawJSON), &battery); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if battery.Date != testDate {
		t.Errorf("Date = %s, want %s", battery.Date, testDate)
	}
	if battery.Charged != 60 {
		t.Errorf("Charged = %d, want 60", battery.Charged)
	}
	if battery.Drained != 45 {
		t.Errorf("Drained = %d, want 45", battery.Drained)
	}
	if battery.StartLevel != 85 {
		t.Errorf("StartLevel = %d, want 85", battery.StartLevel)
	}
	if battery.EndLevel != 40 {
		t.Errorf("EndLevel = %d, want 40", battery.EndLevel)
	}
	if battery.HighestLevel != 100 {
		t.Errorf("HighestLevel = %d, want 100", battery.HighestLevel)
	}
	if battery.LowestLevel != 25 {
		t.Errorf("LowestLevel = %d, want 25", battery.LowestLevel)
	}
}

func TestBodyBatteryReportRawJSON(t *testing.T) {
	rawJSON := `{"date":"2026-01-26","charged":60}`

	var battery BodyBatteryReport
	if err := json.Unmarshal([]byte(rawJSON), &battery); err != nil {
		t.Fatal(err)
	}
	battery.raw = json.RawMessage(rawJSON)

	if string(battery.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
