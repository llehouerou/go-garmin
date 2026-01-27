// service_wellness_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const (
	testDate    = "2026-01-26"
	testDateNew = "2026-01-27"
)

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
	if hr.CalendarDate != testDateNew {
		t.Errorf("CalendarDate = %s, want %s", hr.CalendarDate, testDateNew)
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

func TestDailySpO2JSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"calendarDate": "2026-01-27",
		"startTimestampGMT": "2026-01-26T20:00:00.0",
		"endTimestampGMT": "2026-01-27T06:04:00.0",
		"startTimestampLocal": "2026-01-27T00:00:00.0",
		"endTimestampLocal": "2026-01-27T10:04:00.0",
		"sleepStartTimestampGMT": "2026-01-26T19:27:28.0",
		"sleepEndTimestampGMT": "2026-01-27T02:40:15.0",
		"sleepStartTimestampLocal": "2026-01-26T23:27:28.0",
		"sleepEndTimestampLocal": "2026-01-27T06:40:15.0",
		"averageSpO2": 96.0,
		"lowestSpO2": 88,
		"lastSevenDaysAvgSpO2": 96.85,
		"latestSpO2": 97,
		"latestSpO2TimestampGMT": "2026-01-27T02:40:00.0",
		"latestSpO2TimestampLocal": "2026-01-27T06:40:00.0",
		"avgSleepSpO2": 97.0,
		"spO2ValueDescriptorsDTOList": [
			{"spo2ValueDescriptorIndex": 0, "spo2ValueDescriptorKey": "timestamp"},
			{"spo2ValueDescriptorIndex": 1, "spo2ValueDescriptorKey": "spo2Reading"}
		],
		"spO2HourlyAverages": [[1769457600000, 96], [1769461200000, 95]]
	}`

	var spo2 DailySpO2
	if err := json.Unmarshal([]byte(rawJSON), &spo2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if spo2.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", spo2.UserProfilePK)
	}
	if spo2.CalendarDate != testDateNew {
		t.Errorf("CalendarDate = %s, want %s", spo2.CalendarDate, testDateNew)
	}
	if spo2.AverageSpO2 != 96.0 {
		t.Errorf("AverageSpO2 = %f, want 96.0", spo2.AverageSpO2)
	}
	if spo2.LowestSpO2 != 88 {
		t.Errorf("LowestSpO2 = %d, want 88", spo2.LowestSpO2)
	}
	if spo2.LatestSpO2 != 97 {
		t.Errorf("LatestSpO2 = %d, want 97", spo2.LatestSpO2)
	}
	if spo2.AvgSleepSpO2 != 97.0 {
		t.Errorf("AvgSleepSpO2 = %f, want 97.0", spo2.AvgSleepSpO2)
	}
	if len(spo2.SpO2ValueDescriptors) != 2 {
		t.Errorf("SpO2ValueDescriptors length = %d, want 2", len(spo2.SpO2ValueDescriptors))
	}
	if len(spo2.SpO2HourlyAverages) != 2 {
		t.Errorf("SpO2HourlyAverages length = %d, want 2", len(spo2.SpO2HourlyAverages))
	}
}

func TestDailySpO2RawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-27","averageSpO2":96.0}`

	var spo2 DailySpO2
	if err := json.Unmarshal([]byte(rawJSON), &spo2); err != nil {
		t.Fatal(err)
	}
	spo2.raw = json.RawMessage(rawJSON)

	if string(spo2.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestDailyRespirationJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"calendarDate": "2026-01-27",
		"startTimestampGMT": "2026-01-26T20:00:00.0",
		"endTimestampGMT": "2026-01-27T06:04:00.0",
		"startTimestampLocal": "2026-01-27T00:00:00.0",
		"endTimestampLocal": "2026-01-27T10:04:00.0",
		"sleepStartTimestampGMT": "2026-01-26T19:27:28.0",
		"sleepEndTimestampGMT": "2026-01-27T02:40:15.0",
		"sleepStartTimestampLocal": "2026-01-26T23:27:28.0",
		"sleepEndTimestampLocal": "2026-01-27T06:40:15.0",
		"lowestRespirationValue": 8.0,
		"highestRespirationValue": 21.0,
		"avgWakingRespirationValue": 15.0,
		"avgSleepRespirationValue": 12.0,
		"respirationValueDescriptorsDTOList": [
			{"key": "timestamp", "index": 0},
			{"key": "respiration", "index": 1}
		],
		"respirationValuesArray": [[1769457720000, 11.0], [1769457840000, 12.0]],
		"respirationAveragesValueDescriptorDTOList": [
			{"respirationAveragesValueDescriptorIndex": 0, "respirationAveragesValueDescriptionKey": "timestamp"}
		],
		"respirationAveragesValuesArray": [[1769461200000, 11.38, 14.0, 8.0]],
		"respirationVersion": 200
	}`

	var resp DailyRespiration
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if resp.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", resp.UserProfilePK)
	}
	if resp.CalendarDate != testDateNew {
		t.Errorf("CalendarDate = %s, want %s", resp.CalendarDate, testDateNew)
	}
	if resp.LowestRespirationValue != 8.0 {
		t.Errorf("LowestRespirationValue = %f, want 8.0", resp.LowestRespirationValue)
	}
	if resp.HighestRespirationValue != 21.0 {
		t.Errorf("HighestRespirationValue = %f, want 21.0", resp.HighestRespirationValue)
	}
	if resp.AvgWakingRespirationValue != 15.0 {
		t.Errorf("AvgWakingRespirationValue = %f, want 15.0", resp.AvgWakingRespirationValue)
	}
	if resp.AvgSleepRespirationValue != 12.0 {
		t.Errorf("AvgSleepRespirationValue = %f, want 12.0", resp.AvgSleepRespirationValue)
	}
	if len(resp.RespirationValueDescriptors) != 2 {
		t.Errorf("RespirationValueDescriptors length = %d, want 2", len(resp.RespirationValueDescriptors))
	}
	if len(resp.RespirationValuesArray) != 2 {
		t.Errorf("RespirationValuesArray length = %d, want 2", len(resp.RespirationValuesArray))
	}
	if resp.RespirationVersion != 200 {
		t.Errorf("RespirationVersion = %d, want 200", resp.RespirationVersion)
	}
}

func TestDailyRespirationRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-27","lowestRespirationValue":8.0}`

	var resp DailyRespiration
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		t.Fatal(err)
	}
	resp.raw = json.RawMessage(rawJSON)

	if string(resp.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestDailyIntensityMinutesJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"calendarDate": "2026-01-27",
		"startTimestampGMT": "2026-01-26T20:00:00.0",
		"endTimestampGMT": "2026-01-27T06:04:00.0",
		"startTimestampLocal": "2026-01-27T00:00:00.0",
		"endTimestampLocal": "2026-01-27T10:04:00.0",
		"weeklyModerate": 14,
		"weeklyVigorous": 0,
		"weeklyTotal": 14,
		"weekGoal": 240,
		"dayOfGoalMet": null,
		"startDayMinutes": 7,
		"endDayMinutes": 14,
		"moderateMinutes": 7,
		"vigorousMinutes": 0,
		"imValueDescriptorsDTOList": [
			{"index": 0, "key": "timestamp"},
			{"index": 1, "key": "value"}
		],
		"imValuesArray": [[1769486400000, 7]]
	}`

	var im DailyIntensityMinutes
	if err := json.Unmarshal([]byte(rawJSON), &im); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if im.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", im.UserProfilePK)
	}
	if im.CalendarDate != testDateNew {
		t.Errorf("CalendarDate = %s, want %s", im.CalendarDate, testDateNew)
	}
	if im.WeeklyModerate != 14 {
		t.Errorf("WeeklyModerate = %d, want 14", im.WeeklyModerate)
	}
	if im.WeeklyVigorous != 0 {
		t.Errorf("WeeklyVigorous = %d, want 0", im.WeeklyVigorous)
	}
	if im.WeeklyTotal != 14 {
		t.Errorf("WeeklyTotal = %d, want 14", im.WeeklyTotal)
	}
	if im.WeekGoal != 240 {
		t.Errorf("WeekGoal = %d, want 240", im.WeekGoal)
	}
	if im.ModerateMinutes != 7 {
		t.Errorf("ModerateMinutes = %d, want 7", im.ModerateMinutes)
	}
	if im.VigorousMinutes != 0 {
		t.Errorf("VigorousMinutes = %d, want 0", im.VigorousMinutes)
	}
	if len(im.IMValueDescriptors) != 2 {
		t.Errorf("IMValueDescriptors length = %d, want 2", len(im.IMValueDescriptors))
	}
	if len(im.IMValuesArray) != 1 {
		t.Errorf("IMValuesArray length = %d, want 1", len(im.IMValuesArray))
	}
}

func TestDailyIntensityMinutesRawJSON(t *testing.T) {
	rawJSON := `{"calendarDate":"2026-01-27","weeklyTotal":14}`

	var im DailyIntensityMinutes
	if err := json.Unmarshal([]byte(rawJSON), &im); err != nil {
		t.Fatal(err)
	}
	im.raw = json.RawMessage(rawJSON)

	if string(im.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
