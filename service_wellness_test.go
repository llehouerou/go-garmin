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
