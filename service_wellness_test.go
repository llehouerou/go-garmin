// service_wellness_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const testDate = "2026-01-26"

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

func TestDailyHeartRateJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"calendarDate": "2026-01-27",
		"startTimestampGMT": "2026-01-26T20:00:00.0",
		"endTimestampGMT": "2026-01-27T06:04:00.0",
		"startTimestampLocal": "2026-01-27T00:00:00.0",
		"endTimestampLocal": "2026-01-28T00:00:00.0",
		"maxHeartRate": 119,
		"minHeartRate": 50,
		"restingHeartRate": 51,
		"lastSevenDaysAvgRestingHeartRate": 54,
		"heartRateValueDescriptors": [
			{"key": "timestamp", "index": 0},
			{"key": "heartrate", "index": 1}
		],
		"heartRateValues": [[1769457600000, 51], [1769457720000, 52]]
	}`

	var hr DailyHeartRate
	if err := json.Unmarshal([]byte(rawJSON), &hr); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if hr.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", hr.UserProfilePK)
	}
	if hr.CalendarDate != "2026-01-27" {
		t.Errorf("CalendarDate = %s, want 2026-01-27", hr.CalendarDate)
	}
	if hr.MaxHeartRate != 119 {
		t.Errorf("MaxHeartRate = %d, want 119", hr.MaxHeartRate)
	}
	if hr.MinHeartRate != 50 {
		t.Errorf("MinHeartRate = %d, want 50", hr.MinHeartRate)
	}
	if hr.RestingHeartRate != 51 {
		t.Errorf("RestingHeartRate = %d, want 51", hr.RestingHeartRate)
	}
	if hr.LastSevenDaysAvgRestingHeartRate != 54 {
		t.Errorf("LastSevenDaysAvgRestingHeartRate = %d, want 54", hr.LastSevenDaysAvgRestingHeartRate)
	}
	if len(hr.HeartRateValueDescriptors) != 2 {
		t.Errorf("HeartRateValueDescriptors length = %d, want 2", len(hr.HeartRateValueDescriptors))
	}
	if len(hr.HeartRateValues) != 2 {
		t.Errorf("HeartRateValues length = %d, want 2", len(hr.HeartRateValues))
	}
	if hr.HeartRateValues[0][1] != 51 {
		t.Errorf("HeartRateValues[0][1] = %d, want 51", hr.HeartRateValues[0][1])
	}
}

func TestDailyHeartRateRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-27","maxHeartRate":119}`

	var hr DailyHeartRate
	if err := json.Unmarshal([]byte(rawJSON), &hr); err != nil {
		t.Fatal(err)
	}
	hr.raw = json.RawMessage(rawJSON)

	if string(hr.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
