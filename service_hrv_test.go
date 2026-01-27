// service_hrv_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const testDateHRV = "2026-01-27"

func TestDailyHRVJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePk": 12345678,
		"hrvSummary": {
			"calendarDate": "2026-01-27",
			"weeklyAvg": 53,
			"lastNightAvg": 56,
			"lastNight5MinHigh": 72,
			"baseline": {
				"lowUpper": 47,
				"balancedLow": 49,
				"balancedUpper": 56,
				"markerValue": 0.53570557
			},
			"status": "BALANCED",
			"feedbackPhrase": "HRV_BALANCED_2",
			"createTimeStamp": "2026-01-27T02:40:24.187"
		},
		"hrvReadings": [
			{"hrvValue": 47, "readingTimeGMT": "2026-01-26T19:30:37.0", "readingTimeLocal": "2026-01-26T23:30:37.0"},
			{"hrvValue": 53, "readingTimeGMT": "2026-01-26T19:35:37.0", "readingTimeLocal": "2026-01-26T23:35:37.0"}
		],
		"startTimestampGMT": "2026-01-26T19:27:28.0",
		"endTimestampGMT": "2026-01-27T02:35:37.0",
		"startTimestampLocal": "2026-01-26T23:27:28.0",
		"endTimestampLocal": "2026-01-27T06:35:37.0",
		"sleepStartTimestampGMT": "2026-01-26T19:27:28.0",
		"sleepEndTimestampGMT": "2026-01-27T02:40:15.0",
		"sleepStartTimestampLocal": "2026-01-26T23:27:28.0",
		"sleepEndTimestampLocal": "2026-01-27T06:40:15.0"
	}`

	var hrv DailyHRV
	if err := json.Unmarshal([]byte(rawJSON), &hrv); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if hrv.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", hrv.UserProfilePK)
	}
	if hrv.HRVSummary.CalendarDate != testDateHRV {
		t.Errorf("HRVSummary.CalendarDate = %s, want %s", hrv.HRVSummary.CalendarDate, testDateHRV)
	}
	if hrv.HRVSummary.WeeklyAvg != 53 {
		t.Errorf("HRVSummary.WeeklyAvg = %d, want 53", hrv.HRVSummary.WeeklyAvg)
	}
	if hrv.HRVSummary.LastNightAvg != 56 {
		t.Errorf("HRVSummary.LastNightAvg = %d, want 56", hrv.HRVSummary.LastNightAvg)
	}
	if hrv.HRVSummary.Status != "BALANCED" {
		t.Errorf("HRVSummary.Status = %s, want BALANCED", hrv.HRVSummary.Status)
	}
	if hrv.HRVSummary.Baseline.BalancedLow != 49 {
		t.Errorf("HRVSummary.Baseline.BalancedLow = %d, want 49", hrv.HRVSummary.Baseline.BalancedLow)
	}
	if len(hrv.HRVReadings) != 2 {
		t.Errorf("HRVReadings length = %d, want 2", len(hrv.HRVReadings))
	}
	if hrv.HRVReadings[0].HRVValue != 47 {
		t.Errorf("HRVReadings[0].HRVValue = %d, want 47", hrv.HRVReadings[0].HRVValue)
	}
}

func TestDailyHRVRawJSON(t *testing.T) {
	rawJSON := `{"userProfilePk":12345678,"hrvSummary":{"status":"BALANCED"}}`

	var hrv DailyHRV
	if err := json.Unmarshal([]byte(rawJSON), &hrv); err != nil {
		t.Fatal(err)
	}
	hrv.raw = json.RawMessage(rawJSON)

	if string(hrv.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestHRVRangeJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"hrvSummaries": [
			{
				"calendarDate": "2026-01-20",
				"weeklyAvg": 48,
				"lastNightAvg": 45,
				"lastNight5MinHigh": 62,
				"baseline": {"lowUpper": 47, "balancedLow": 49, "balancedUpper": 57, "markerValue": 0.14430237},
				"status": "UNBALANCED",
				"feedbackPhrase": "HRV_UNBALANCED_12",
				"createTimeStamp": "2026-01-20T04:46:11.947"
			},
			{
				"calendarDate": "2026-01-21",
				"weeklyAvg": 48,
				"lastNightAvg": 54,
				"lastNight5MinHigh": 74,
				"baseline": {"lowUpper": 47, "balancedLow": 49, "balancedUpper": 57, "markerValue": 0.1401825},
				"status": "UNBALANCED",
				"feedbackPhrase": "HRV_UNBALANCED_11",
				"createTimeStamp": "2026-01-21T04:16:29.20"
			}
		],
		"userProfilePk": 12345678
	}`

	var hrvRange HRVRange
	if err := json.Unmarshal([]byte(rawJSON), &hrvRange); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if hrvRange.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", hrvRange.UserProfilePK)
	}
	if len(hrvRange.HRVSummaries) != 2 {
		t.Fatalf("HRVSummaries length = %d, want 2", len(hrvRange.HRVSummaries))
	}
	if hrvRange.HRVSummaries[0].CalendarDate != "2026-01-20" {
		t.Errorf("HRVSummaries[0].CalendarDate = %s, want 2026-01-20", hrvRange.HRVSummaries[0].CalendarDate)
	}
	if hrvRange.HRVSummaries[0].Status != "UNBALANCED" {
		t.Errorf("HRVSummaries[0].Status = %s, want UNBALANCED", hrvRange.HRVSummaries[0].Status)
	}
	if hrvRange.HRVSummaries[1].CalendarDate != "2026-01-21" {
		t.Errorf("HRVSummaries[1].CalendarDate = %s, want 2026-01-21", hrvRange.HRVSummaries[1].CalendarDate)
	}
}

func TestHRVRangeRawJSON(t *testing.T) {
	rawJSON := `{"hrvSummaries":[],"userProfilePk":12345678}`

	var hrvRange HRVRange
	if err := json.Unmarshal([]byte(rawJSON), &hrvRange); err != nil {
		t.Fatal(err)
	}
	hrvRange.raw = json.RawMessage(rawJSON)

	if string(hrvRange.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
