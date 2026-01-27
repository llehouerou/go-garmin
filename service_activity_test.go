package garmin

import (
	"encoding/json"
	"testing"
	"time"
)

func TestActivityConversions(t *testing.T) {
	activity := &Activity{
		StartTimeGMT: "2026-01-25 13:59:36",
		Duration:     2356.998,
		Distance:     5692.08,
	}

	// Test StartTime
	start := activity.StartTime()
	if start.Year() != 2026 || start.Month() != 1 || start.Day() != 25 {
		t.Errorf("StartTime() = %v, want 2026-01-25", start)
	}
	if start.Hour() != 13 || start.Minute() != 59 {
		t.Errorf("StartTime() hour:min = %d:%d, want 13:59", start.Hour(), start.Minute())
	}

	// Test DurationTime
	dur := activity.DurationTime()
	if dur < 39*time.Minute || dur > 40*time.Minute {
		t.Errorf("DurationTime() = %v, want ~39 minutes", dur)
	}

	// Test DistanceKm
	distKm := activity.DistanceKm()
	if distKm < 5.69 || distKm > 5.70 {
		t.Errorf("DistanceKm() = %v, want ~5.69", distKm)
	}

	// Test DistanceMiles
	distMi := activity.DistanceMiles()
	if distMi < 3.53 || distMi > 3.54 {
		t.Errorf("DistanceMiles() = %v, want ~3.53", distMi)
	}
}

func TestActivityAveragePacePerKm(t *testing.T) {
	activity := &Activity{
		Duration: 2400, // 40 minutes
		Distance: 8000, // 8 km
	}

	pace := activity.AveragePacePerKm()
	expected := 5 * time.Minute // 5 min/km
	if pace != expected {
		t.Errorf("AveragePacePerKm() = %v, want %v", pace, expected)
	}

	// Test zero distance
	zeroActivity := &Activity{Duration: 1000, Distance: 0}
	if zeroActivity.AveragePacePerKm() != 0 {
		t.Error("AveragePacePerKm() should return 0 for zero distance")
	}
}

func TestActivityRawJSON(t *testing.T) {
	rawJSON := `{"activityId":123,"activityName":"Test Run"}`

	var activity Activity
	if err := json.Unmarshal([]byte(rawJSON), &activity); err != nil {
		t.Fatal(err)
	}
	activity.raw = json.RawMessage(rawJSON)

	if string(activity.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestActivityDetailConversions(t *testing.T) {
	detail := &ActivityDetail{
		SummaryDTO: ActivitySummary{
			StartTimeGMT: "2026-01-25T13:59:36.0",
			Duration:     2356.998,
			Distance:     5692.08,
		},
	}

	// Test StartTime
	start := detail.StartTime()
	if start.Year() != 2026 || start.Month() != 1 || start.Day() != 25 {
		t.Errorf("StartTime() = %v, want 2026-01-25", start)
	}

	// Test DurationTime
	dur := detail.DurationTime()
	if dur < 39*time.Minute || dur > 40*time.Minute {
		t.Errorf("DurationTime() = %v, want ~39 minutes", dur)
	}

	// Test DistanceKm
	distKm := detail.DistanceKm()
	if distKm < 5.69 || distKm > 5.70 {
		t.Errorf("DistanceKm() = %v, want ~5.69", distKm)
	}
}

func TestActivityDetailRawJSON(t *testing.T) {
	rawJSON := `{"activityId":123,"activityName":"Test Run"}`

	var detail ActivityDetail
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		t.Fatal(err)
	}
	detail.raw = json.RawMessage(rawJSON)

	if string(detail.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestActivityJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"activityId": 21661023200,
		"activityName": "Morning Run",
		"startTimeLocal": "2026-01-25 17:59:36",
		"startTimeGMT": "2026-01-25 13:59:36",
		"activityType": {
			"typeId": 1,
			"typeKey": "running",
			"parentTypeId": 17,
			"isHidden": false,
			"restricted": false,
			"trimmable": true
		},
		"eventType": {
			"typeId": 9,
			"typeKey": "uncategorized",
			"sortOrder": 10
		},
		"distance": 5692.08,
		"duration": 2356.998,
		"elevationGain": 66.0,
		"elevationLoss": 65.0,
		"averageSpeed": 2.415,
		"maxSpeed": 3.191,
		"calories": 437.0,
		"averageHR": 143.0,
		"maxHR": 160.0,
		"steps": 5912,
		"privacy": {
			"typeId": 2,
			"typeKey": "private"
		},
		"splitSummaries": [
			{
				"noOfSplits": 1,
				"totalAscent": 66.0,
				"duration": 2356.998,
				"splitType": "INTERVAL_ACTIVE",
				"distance": 5692.08
			}
		]
	}`

	var activity Activity
	if err := json.Unmarshal([]byte(rawJSON), &activity); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if activity.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", activity.ActivityID)
	}
	if activity.ActivityName != "Morning Run" {
		t.Errorf("ActivityName = %s, want Morning Run", activity.ActivityName)
	}
	if activity.ActivityType.TypeKey != "running" {
		t.Errorf("ActivityType.TypeKey = %s, want running", activity.ActivityType.TypeKey)
	}
	if activity.EventType.TypeKey != "uncategorized" {
		t.Errorf("EventType.TypeKey = %s, want uncategorized", activity.EventType.TypeKey)
	}
	if activity.Distance != 5692.08 {
		t.Errorf("Distance = %f, want 5692.08", activity.Distance)
	}
	if activity.Steps != 5912 {
		t.Errorf("Steps = %d, want 5912", activity.Steps)
	}
	if activity.Privacy.TypeKey != "private" {
		t.Errorf("Privacy.TypeKey = %s, want private", activity.Privacy.TypeKey)
	}
	if len(activity.SplitSummaries) != 1 {
		t.Errorf("SplitSummaries length = %d, want 1", len(activity.SplitSummaries))
	}
	if activity.SplitSummaries[0].SplitType != "INTERVAL_ACTIVE" {
		t.Errorf("SplitSummaries[0].SplitType = %s, want INTERVAL_ACTIVE", activity.SplitSummaries[0].SplitType)
	}
}

func TestActivityDetailJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"activityId": 21661023200,
		"activityUUID": {"uuid": "e4ed3e69-d34e-477f-80d0-dddda2da652f"},
		"activityName": "Morning Run",
		"userProfileId": 12345678,
		"isMultiSportParent": false,
		"activityTypeDTO": {
			"typeId": 1,
			"typeKey": "running",
			"parentTypeId": 17
		},
		"eventTypeDTO": {
			"typeId": 9,
			"typeKey": "uncategorized",
			"sortOrder": 10
		},
		"accessControlRuleDTO": {
			"typeId": 2,
			"typeKey": "private"
		},
		"timeZoneUnitDTO": {
			"unitId": 125,
			"unitKey": "Asia/Dubai",
			"factor": 0.0,
			"timeZone": "Asia/Dubai"
		},
		"summaryDTO": {
			"startTimeLocal": "2026-01-25T17:59:36.0",
			"startTimeGMT": "2026-01-25T13:59:36.0",
			"distance": 5692.08,
			"duration": 2356.998,
			"calories": 437.0,
			"averageHR": 143.0,
			"maxHR": 160.0,
			"steps": 5912
		},
		"locationName": "Saint-Pierre"
	}`

	var detail ActivityDetail
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if detail.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", detail.ActivityID)
	}
	if detail.ActivityUUID.UUID != "e4ed3e69-d34e-477f-80d0-dddda2da652f" {
		t.Errorf("ActivityUUID.UUID = %s, want e4ed3e69-d34e-477f-80d0-dddda2da652f", detail.ActivityUUID.UUID)
	}
	if detail.UserProfileID != 12345678 {
		t.Errorf("UserProfileID = %d, want 12345678", detail.UserProfileID)
	}
	if detail.ActivityTypeDTO.TypeKey != "running" {
		t.Errorf("ActivityTypeDTO.TypeKey = %s, want running", detail.ActivityTypeDTO.TypeKey)
	}
	if detail.TimeZoneUnitDTO.TimeZone != "Asia/Dubai" {
		t.Errorf("TimeZoneUnitDTO.TimeZone = %s, want Asia/Dubai", detail.TimeZoneUnitDTO.TimeZone)
	}
	if detail.SummaryDTO.Distance != 5692.08 {
		t.Errorf("SummaryDTO.Distance = %f, want 5692.08", detail.SummaryDTO.Distance)
	}
	if detail.LocationName != "Saint-Pierre" {
		t.Errorf("LocationName = %s, want Saint-Pierre", detail.LocationName)
	}
}

func TestListOptions(t *testing.T) {
	// Test default values
	opts := &ListOptions{}
	if opts.Start != 0 {
		t.Errorf("Default Start = %d, want 0", opts.Start)
	}
	if opts.Limit != 0 {
		t.Errorf("Default Limit = %d, want 0", opts.Limit)
	}

	// Test custom values
	opts = &ListOptions{Start: 10, Limit: 50}
	if opts.Start != 10 {
		t.Errorf("Start = %d, want 10", opts.Start)
	}
	if opts.Limit != 50 {
		t.Errorf("Limit = %d, want 50", opts.Limit)
	}
}
