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

func TestActivityWeatherJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"issueDate": "2026-01-25T14:00:00",
		"temp": 82,
		"apparentTemp": 86,
		"dewPoint": 71,
		"relativeHumidity": 70,
		"windDirection": 180,
		"windDirectionCompassPoint": "S",
		"windSpeed": 5,
		"windGust": null,
		"latitude": 25.123,
		"longitude": 55.456,
		"weatherStationDTO": {
			"id": "STATION123",
			"name": "Test Weather Station",
			"timezone": null
		},
		"weatherTypeDTO": {
			"weatherTypePk": 2,
			"desc": "Partly Cloudy",
			"image": null
		}
	}`

	var weather ActivityWeather
	if err := json.Unmarshal([]byte(rawJSON), &weather); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if weather.IssueDate != "2026-01-25T14:00:00" {
		t.Errorf("IssueDate = %s, want 2026-01-25T14:00:00", weather.IssueDate)
	}
	if weather.Temp != 82 {
		t.Errorf("Temp = %d, want 82", weather.Temp)
	}
	if weather.ApparentTemp != 86 {
		t.Errorf("ApparentTemp = %d, want 86", weather.ApparentTemp)
	}
	if weather.RelativeHumidity != 70 {
		t.Errorf("RelativeHumidity = %d, want 70", weather.RelativeHumidity)
	}
	if weather.WindDirectionCompassPoint != "S" {
		t.Errorf("WindDirectionCompassPoint = %s, want S", weather.WindDirectionCompassPoint)
	}
	if weather.WindSpeed != 5 {
		t.Errorf("WindSpeed = %d, want 5", weather.WindSpeed)
	}
	if weather.WindGust != nil {
		t.Errorf("WindGust = %v, want nil", weather.WindGust)
	}
	if weather.WeatherStationDTO.ID != "STATION123" {
		t.Errorf("WeatherStationDTO.ID = %s, want STATION123", weather.WeatherStationDTO.ID)
	}
	if weather.WeatherTypeDTO.Desc != "Partly Cloudy" {
		t.Errorf("WeatherTypeDTO.Desc = %s, want Partly Cloudy", weather.WeatherTypeDTO.Desc)
	}
}

func TestActivityWeatherConversions(t *testing.T) {
	weather := &ActivityWeather{
		Temp:         82, // 82°F = ~27.78°C
		ApparentTemp: 86, // 86°F = 30°C
	}

	// Test TempCelsius
	tempC := weather.TempCelsius()
	if tempC < 27.7 || tempC > 27.8 {
		t.Errorf("TempCelsius() = %v, want ~27.78", tempC)
	}

	// Test ApparentTempCelsius
	apparentC := weather.ApparentTempCelsius()
	if apparentC != 30 {
		t.Errorf("ApparentTempCelsius() = %v, want 30", apparentC)
	}

	// Test freezing point conversion (32°F = 0°C)
	freezing := &ActivityWeather{Temp: 32, ApparentTemp: 32}
	if freezing.TempCelsius() != 0 {
		t.Errorf("TempCelsius(32°F) = %v, want 0", freezing.TempCelsius())
	}
}

func TestActivityWeatherRawJSON(t *testing.T) {
	rawJSON := `{"temp":82,"apparentTemp":86}`

	var weather ActivityWeather
	if err := json.Unmarshal([]byte(rawJSON), &weather); err != nil {
		t.Fatal(err)
	}
	weather.raw = json.RawMessage(rawJSON)

	if string(weather.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestLapConversions(t *testing.T) {
	lap := &Lap{
		Duration: 300,  // 5 minutes
		Distance: 1000, // 1 km
	}

	// Test DurationTime
	dur := lap.DurationTime()
	if dur != 5*time.Minute {
		t.Errorf("DurationTime() = %v, want 5m", dur)
	}

	// Test DistanceKm
	distKm := lap.DistanceKm()
	if distKm != 1.0 {
		t.Errorf("DistanceKm() = %v, want 1.0", distKm)
	}

	// Test AveragePacePerKm
	pace := lap.AveragePacePerKm()
	if pace != 5*time.Minute {
		t.Errorf("AveragePacePerKm() = %v, want 5m", pace)
	}

	// Test zero distance
	zeroLap := &Lap{Duration: 1000, Distance: 0}
	if zeroLap.AveragePacePerKm() != 0 {
		t.Error("AveragePacePerKm() should return 0 for zero distance")
	}
}

func TestActivitySplitsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"activityId": 21661023200,
		"lapDTOs": [
			{
				"startTimeGMT": "2026-01-25 13:59:36",
				"startLatitude": 25.123,
				"startLongitude": 55.456,
				"distance": 1000.0,
				"duration": 300.0,
				"averageSpeed": 3.33,
				"maxSpeed": 4.0,
				"averageHR": 145.0,
				"maxHR": 160.0,
				"lapIndex": 0,
				"intensityType": "ACTIVE",
				"messageIndex": 0
			},
			{
				"startTimeGMT": "2026-01-25 14:04:36",
				"startLatitude": 25.124,
				"startLongitude": 55.457,
				"distance": 1000.0,
				"duration": 295.0,
				"averageSpeed": 3.39,
				"maxSpeed": 4.1,
				"averageHR": 150.0,
				"maxHR": 165.0,
				"lapIndex": 1,
				"intensityType": "ACTIVE",
				"messageIndex": 1
			}
		],
		"eventDTOs": [
			{
				"startTimeGMT": "2026-01-25 13:59:36",
				"startTimeGMTDoubleValue": 1737813576000.0,
				"sectionTypeDTO": {
					"id": 1,
					"key": "START",
					"sectionTypeKey": "START"
				}
			}
		]
	}`

	var splits ActivitySplits
	if err := json.Unmarshal([]byte(rawJSON), &splits); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if splits.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", splits.ActivityID)
	}
	if len(splits.LapDTOs) != 2 {
		t.Fatalf("LapDTOs length = %d, want 2", len(splits.LapDTOs))
	}

	firstLap := splits.LapDTOs[0]
	if firstLap.Distance != 1000.0 {
		t.Errorf("LapDTOs[0].Distance = %f, want 1000.0", firstLap.Distance)
	}
	if firstLap.Duration != 300.0 {
		t.Errorf("LapDTOs[0].Duration = %f, want 300.0", firstLap.Duration)
	}
	if firstLap.LapIndex != 0 {
		t.Errorf("LapDTOs[0].LapIndex = %d, want 0", firstLap.LapIndex)
	}
	if firstLap.IntensityType != "ACTIVE" {
		t.Errorf("LapDTOs[0].IntensityType = %s, want ACTIVE", firstLap.IntensityType)
	}

	if len(splits.EventDTOs) != 1 {
		t.Fatalf("EventDTOs length = %d, want 1", len(splits.EventDTOs))
	}
	if splits.EventDTOs[0].SectionTypeDTO.Key != "START" {
		t.Errorf("EventDTOs[0].SectionTypeDTO.Key = %s, want START", splits.EventDTOs[0].SectionTypeDTO.Key)
	}
}

func TestActivitySplitsRawJSON(t *testing.T) {
	rawJSON := `{"activityId":123,"lapDTOs":[],"eventDTOs":[]}`

	var splits ActivitySplits
	if err := json.Unmarshal([]byte(rawJSON), &splits); err != nil {
		t.Fatal(err)
	}
	splits.raw = json.RawMessage(rawJSON)

	if string(splits.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestTimeInZoneDurationInZone(t *testing.T) {
	zone := &TimeInZone{
		ZoneNumber:      1,
		SecsInZone:      300.5,
		ZoneLowBoundary: 94,
	}

	dur := zone.DurationInZone()
	expected := 300*time.Second + 500*time.Millisecond
	if dur != expected {
		t.Errorf("DurationInZone() = %v, want %v", dur, expected)
	}
}

func TestHRTimeInZonesJSONUnmarshal(t *testing.T) {
	rawJSON := `[
		{"zoneNumber":1,"secsInZone":2.999,"zoneLowBoundary":94},
		{"zoneNumber":2,"secsInZone":48.997,"zoneLowBoundary":113},
		{"zoneNumber":3,"secsInZone":1161.79,"zoneLowBoundary":132},
		{"zoneNumber":4,"secsInZone":1022.05,"zoneLowBoundary":151},
		{"zoneNumber":5,"secsInZone":102.166,"zoneLowBoundary":170}
	]`

	var zones []TimeInZone
	if err := json.Unmarshal([]byte(rawJSON), &zones); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	hrZones := &HRTimeInZones{Zones: zones, raw: json.RawMessage(rawJSON)}

	if len(hrZones.Zones) != 5 {
		t.Errorf("Zones length = %d, want 5", len(hrZones.Zones))
	}

	// Check first zone
	if hrZones.Zones[0].ZoneNumber != 1 {
		t.Errorf("Zones[0].ZoneNumber = %d, want 1", hrZones.Zones[0].ZoneNumber)
	}
	if hrZones.Zones[0].ZoneLowBoundary != 94 {
		t.Errorf("Zones[0].ZoneLowBoundary = %d, want 94", hrZones.Zones[0].ZoneLowBoundary)
	}

	// Check zone 3 with significant time
	if hrZones.Zones[2].SecsInZone < 1161 || hrZones.Zones[2].SecsInZone > 1162 {
		t.Errorf("Zones[2].SecsInZone = %f, want ~1161.79", hrZones.Zones[2].SecsInZone)
	}
}

func TestHRTimeInZonesRawJSON(t *testing.T) {
	rawJSON := `[{"zoneNumber":1,"secsInZone":100,"zoneLowBoundary":94}]`

	var zones []TimeInZone
	if err := json.Unmarshal([]byte(rawJSON), &zones); err != nil {
		t.Fatal(err)
	}

	hrZones := &HRTimeInZones{Zones: zones, raw: json.RawMessage(rawJSON)}

	if string(hrZones.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestPowerTimeInZonesJSONUnmarshal(t *testing.T) {
	rawJSON := `[
		{"zoneNumber":1,"secsInZone":500.0,"zoneLowBoundary":0},
		{"zoneNumber":2,"secsInZone":300.0,"zoneLowBoundary":150},
		{"zoneNumber":3,"secsInZone":200.0,"zoneLowBoundary":200}
	]`

	var zones []TimeInZone
	if err := json.Unmarshal([]byte(rawJSON), &zones); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	powerZones := &PowerTimeInZones{Zones: zones, raw: json.RawMessage(rawJSON)}

	if len(powerZones.Zones) != 3 {
		t.Errorf("Zones length = %d, want 3", len(powerZones.Zones))
	}

	if powerZones.Zones[1].ZoneLowBoundary != 150 {
		t.Errorf("Zones[1].ZoneLowBoundary = %d, want 150", powerZones.Zones[1].ZoneLowBoundary)
	}
}

func TestPowerTimeInZonesRawJSON(t *testing.T) {
	rawJSON := `[{"zoneNumber":1,"secsInZone":100,"zoneLowBoundary":0}]`

	var zones []TimeInZone
	if err := json.Unmarshal([]byte(rawJSON), &zones); err != nil {
		t.Fatal(err)
	}

	powerZones := &PowerTimeInZones{Zones: zones, raw: json.RawMessage(rawJSON)}

	if string(powerZones.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestActivityDetailsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"activityId": 21661023200,
		"measurementCount": 100,
		"metricsCount": 5,
		"totalMetricsCount": 5,
		"metricDescriptors": [
			{"metricsIndex": 0, "key": "directTimestamp", "unit": {"id": 1, "key": "ms", "factor": 1.0}},
			{"metricsIndex": 1, "key": "directHeartRate", "unit": {"id": 2, "key": "bpm", "factor": 1.0}},
			{"metricsIndex": 2, "key": "directSpeed", "unit": {"id": 3, "key": "mps", "factor": 1.0}}
		],
		"activityDetailMetrics": [
			{"metrics": [0, 120, 2.5]},
			{"metrics": [1000, 125, 2.6]}
		]
	}`

	var details ActivityDetails
	if err := json.Unmarshal([]byte(rawJSON), &details); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	details.raw = json.RawMessage(rawJSON)

	if details.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", details.ActivityID)
	}
	if details.MeasurementCount != 100 {
		t.Errorf("MeasurementCount = %d, want 100", details.MeasurementCount)
	}
	if details.MetricsCount != 5 {
		t.Errorf("MetricsCount = %d, want 5", details.MetricsCount)
	}
	if len(details.MetricDescriptors) != 3 {
		t.Errorf("MetricDescriptors length = %d, want 3", len(details.MetricDescriptors))
	}
	if details.MetricDescriptors[1].Key != "directHeartRate" {
		t.Errorf("MetricDescriptors[1].Key = %s, want directHeartRate", details.MetricDescriptors[1].Key)
	}
	if details.MetricDescriptors[1].Unit.Key != "bpm" {
		t.Errorf("MetricDescriptors[1].Unit.Key = %s, want bpm", details.MetricDescriptors[1].Unit.Key)
	}
	if len(details.ActivityDetailMetrics) != 2 {
		t.Errorf("ActivityDetailMetrics length = %d, want 2", len(details.ActivityDetailMetrics))
	}
}

func TestActivityDetailsGetMetricIndex(t *testing.T) {
	details := &ActivityDetails{
		MetricDescriptors: []MetricDescriptor{
			{MetricsIndex: 0, Key: "directTimestamp"},
			{MetricsIndex: 1, Key: "directHeartRate"},
			{MetricsIndex: 2, Key: "directSpeed"},
		},
	}

	// Test finding existing keys
	if idx := details.GetMetricIndex("directHeartRate"); idx != 1 {
		t.Errorf("GetMetricIndex(directHeartRate) = %d, want 1", idx)
	}
	if idx := details.GetMetricIndex("directSpeed"); idx != 2 {
		t.Errorf("GetMetricIndex(directSpeed) = %d, want 2", idx)
	}

	// Test missing key
	if idx := details.GetMetricIndex("nonexistent"); idx != -1 {
		t.Errorf("GetMetricIndex(nonexistent) = %d, want -1", idx)
	}
}

func TestActivityDetailsRawJSON(t *testing.T) {
	rawJSON := `{"activityId":123,"measurementCount":0,"metricsCount":0,"totalMetricsCount":0,"metricDescriptors":[],"activityDetailMetrics":[]}`

	var details ActivityDetails
	if err := json.Unmarshal([]byte(rawJSON), &details); err != nil {
		t.Fatal(err)
	}
	details.raw = json.RawMessage(rawJSON)

	if string(details.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestExerciseSetsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"activityId": 21661023200,
		"exerciseSets": [
			{
				"setType": "ACTIVE",
				"category": "STRENGTH",
				"exerciseName": "BENCH_PRESS",
				"weight": 60.0,
				"repetitionCount": 10,
				"duration": 30.5,
				"startTime": "2026-01-25T14:00:00",
				"messageIndex": 0
			},
			{
				"setType": "REST",
				"category": "STRENGTH",
				"exerciseName": "BENCH_PRESS",
				"weight": null,
				"repetitionCount": null,
				"duration": 60.0,
				"startTime": "2026-01-25T14:00:30",
				"messageIndex": 1
			}
		]
	}`

	var sets ExerciseSets
	if err := json.Unmarshal([]byte(rawJSON), &sets); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	sets.raw = json.RawMessage(rawJSON)

	if sets.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", sets.ActivityID)
	}
	if len(sets.ExerciseSets) != 2 {
		t.Fatalf("ExerciseSets length = %d, want 2", len(sets.ExerciseSets))
	}

	// Check active set
	activeSet := sets.ExerciseSets[0]
	if activeSet.SetType != "ACTIVE" {
		t.Errorf("ExerciseSets[0].SetType = %s, want ACTIVE", activeSet.SetType)
	}
	if activeSet.Category != "STRENGTH" {
		t.Errorf("ExerciseSets[0].Category = %s, want STRENGTH", activeSet.Category)
	}
	if activeSet.ExerciseName != "BENCH_PRESS" {
		t.Errorf("ExerciseSets[0].ExerciseName = %s, want BENCH_PRESS", activeSet.ExerciseName)
	}
	if activeSet.Weight == nil || *activeSet.Weight != 60.0 {
		t.Errorf("ExerciseSets[0].Weight = %v, want 60.0", activeSet.Weight)
	}
	if activeSet.RepetitionCount == nil || *activeSet.RepetitionCount != 10 {
		t.Errorf("ExerciseSets[0].RepetitionCount = %v, want 10", activeSet.RepetitionCount)
	}

	// Check rest set with null values
	restSet := sets.ExerciseSets[1]
	if restSet.SetType != "REST" {
		t.Errorf("ExerciseSets[1].SetType = %s, want REST", restSet.SetType)
	}
	if restSet.Weight != nil {
		t.Errorf("ExerciseSets[1].Weight = %v, want nil", restSet.Weight)
	}
	if restSet.RepetitionCount != nil {
		t.Errorf("ExerciseSets[1].RepetitionCount = %v, want nil", restSet.RepetitionCount)
	}
}

func TestExerciseSetsEmptyArray(t *testing.T) {
	rawJSON := `{"activityId":21661023200,"exerciseSets":null}`

	var sets ExerciseSets
	if err := json.Unmarshal([]byte(rawJSON), &sets); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if sets.ActivityID != 21661023200 {
		t.Errorf("ActivityID = %d, want 21661023200", sets.ActivityID)
	}
	if sets.ExerciseSets != nil {
		t.Errorf("ExerciseSets = %v, want nil", sets.ExerciseSets)
	}
}

func TestExerciseSetsRawJSON(t *testing.T) {
	rawJSON := `{"activityId":123,"exerciseSets":[]}`

	var sets ExerciseSets
	if err := json.Unmarshal([]byte(rawJSON), &sets); err != nil {
		t.Fatal(err)
	}
	sets.raw = json.RawMessage(rawJSON)

	if string(sets.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
