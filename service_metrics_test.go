// service_metrics_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const testDateMetrics = "2026-01-27"

func TestTrainingReadinessJSONUnmarshal(t *testing.T) {
	rawJSON := `[{
		"userProfilePK": 12345678,
		"calendarDate": "2026-01-27",
		"timestamp": "2026-01-27T02:40:15.0",
		"level": "HIGH",
		"feedbackShort": "RESTED_AND_READY",
		"score": 91,
		"sleepScore": 96,
		"recoveryTime": 0,
		"acuteLoad": 249,
		"hrvWeeklyAverage": 53,
		"validSleep": true,
		"primaryActivityTracker": true
	}]`

	var entries []TrainingReadinessEntry
	if err := json.Unmarshal([]byte(rawJSON), &entries); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Entries length = %d, want 1", len(entries))
	}
	if entries[0].CalendarDate != testDateMetrics {
		t.Errorf("CalendarDate = %s, want %s", entries[0].CalendarDate, testDateMetrics)
	}
	if entries[0].Level != "HIGH" {
		t.Errorf("Level = %s, want HIGH", entries[0].Level)
	}
	if entries[0].Score != 91 {
		t.Errorf("Score = %d, want 91", entries[0].Score)
	}
	if entries[0].SleepScore == nil || *entries[0].SleepScore != 96 {
		t.Errorf("SleepScore = %v, want 96", entries[0].SleepScore)
	}
	if entries[0].HRVWeeklyAverage != 53 {
		t.Errorf("HRVWeeklyAverage = %d, want 53", entries[0].HRVWeeklyAverage)
	}
}

func TestEnduranceScoreJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"deviceId": 3490854130,
		"calendarDate": "2026-01-27",
		"overallScore": 5549,
		"classification": 2,
		"feedbackPhrase": 30,
		"primaryTrainingDevice": true,
		"gaugeLowerLimit": 3570,
		"gaugeUpperLimit": 10320,
		"contributors": [
			{"activityTypeId": null, "group": 0, "contribution": 94.29},
			{"activityTypeId": 3, "group": null, "contribution": 5.71}
		]
	}`

	var score EnduranceScore
	if err := json.Unmarshal([]byte(rawJSON), &score); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if score.CalendarDate != testDateMetrics {
		t.Errorf("CalendarDate = %s, want %s", score.CalendarDate, testDateMetrics)
	}
	if score.OverallScore != 5549 {
		t.Errorf("OverallScore = %d, want 5549", score.OverallScore)
	}
	if score.Classification != 2 {
		t.Errorf("Classification = %d, want 2", score.Classification)
	}
	if len(score.Contributors) != 2 {
		t.Fatalf("Contributors length = %d, want 2", len(score.Contributors))
	}
	if score.Contributors[0].Contribution != 94.29 {
		t.Errorf("Contributors[0].Contribution = %f, want 94.29", score.Contributors[0].Contribution)
	}
}

func TestHillScoreJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userProfilePK": 12345678,
		"deviceId": 3490854130,
		"calendarDate": "2026-01-27",
		"strengthScore": 7,
		"enduranceScore": 11,
		"hillScoreClassificationId": 2,
		"overallScore": 32,
		"hillScoreFeedbackPhraseId": 4,
		"vo2Max": 50.0,
		"vo2MaxPreciseValue": 50.1,
		"primaryTrainingDevice": true
	}`

	var score HillScore
	if err := json.Unmarshal([]byte(rawJSON), &score); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if score.CalendarDate != testDateMetrics {
		t.Errorf("CalendarDate = %s, want %s", score.CalendarDate, testDateMetrics)
	}
	if score.OverallScore != 32 {
		t.Errorf("OverallScore = %d, want 32", score.OverallScore)
	}
	if score.VO2Max != 50.0 {
		t.Errorf("VO2Max = %f, want 50.0", score.VO2Max)
	}
	if score.VO2MaxPreciseValue != 50.1 {
		t.Errorf("VO2MaxPreciseValue = %f, want 50.1", score.VO2MaxPreciseValue)
	}
}

func TestMaxMetEntryJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userId": 12345678,
		"generic": {
			"calendarDate": "2026-01-25",
			"vo2MaxPreciseValue": 50.1,
			"vo2MaxValue": 50.0,
			"fitnessAge": null,
			"maxMetCategory": 0
		},
		"cycling": null,
		"heatAltitudeAcclimation": {
			"calendarDate": "2026-01-27",
			"heatAcclimationPercentage": 100,
			"heatTrend": "ACCLIMATIZED",
			"altitudeAcclimation": 0
		}
	}`

	var entry MaxMetEntry
	if err := json.Unmarshal([]byte(rawJSON), &entry); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if entry.UserID != 12345678 {
		t.Errorf("UserID = %d, want 12345678", entry.UserID)
	}
	if entry.Generic == nil {
		t.Fatal("Generic is nil")
	}
	if entry.Generic.VO2MaxValue != 50.0 {
		t.Errorf("Generic.VO2MaxValue = %f, want 50.0", entry.Generic.VO2MaxValue)
	}
	if entry.HeatAltitudeAcclimation == nil {
		t.Fatal("HeatAltitudeAcclimation is nil")
	}
	if entry.HeatAltitudeAcclimation.HeatAcclimationPercentage != 100 {
		t.Errorf("HeatAcclimationPercentage = %d, want 100", entry.HeatAltitudeAcclimation.HeatAcclimationPercentage)
	}
}

func TestTrainingStatusDailyJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userId": 12345678,
		"latestTrainingStatusData": {
			"3490854130": {
				"calendarDate": "2026-01-27",
				"trainingStatus": 7,
				"trainingStatusFeedbackPhrase": "PRODUCTIVE_1",
				"trainingPaused": false,
				"acuteTrainingLoadDTO": {
					"acwrPercent": 47,
					"acwrStatus": "OPTIMAL",
					"dailyTrainingLoadAcute": 249,
					"dailyAcuteChronicWorkloadRatio": 1.1
				},
				"primaryTrainingDevice": true
			}
		},
		"recordedDevices": [
			{"deviceId": 3490854130, "deviceName": "Forerunner 965", "category": 0}
		],
		"showSelector": false
	}`

	var status TrainingStatusDaily
	if err := json.Unmarshal([]byte(rawJSON), &status); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if status.UserID != 12345678 {
		t.Errorf("UserID = %d, want 12345678", status.UserID)
	}
	if len(status.LatestTrainingStatusData) != 1 {
		t.Fatalf("LatestTrainingStatusData length = %d, want 1", len(status.LatestTrainingStatusData))
	}

	data := status.LatestTrainingStatusData["3490854130"]
	if data == nil {
		t.Fatal("Device data is nil")
	}
	if data.TrainingStatus != 7 {
		t.Errorf("TrainingStatus = %d, want 7", data.TrainingStatus)
	}
	if data.AcuteTrainingLoadDTO == nil {
		t.Fatal("AcuteTrainingLoadDTO is nil")
	}
	if data.AcuteTrainingLoadDTO.AcwrStatus != "OPTIMAL" {
		t.Errorf("AcwrStatus = %s, want OPTIMAL", data.AcuteTrainingLoadDTO.AcwrStatus)
	}
}

func TestTrainingLoadBalanceJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"userId": 12345678,
		"metricsTrainingLoadBalanceDTOMap": {
			"3490854130": {
				"calendarDate": "2026-01-27",
				"deviceId": 3490854130,
				"monthlyLoadAerobicLow": 187.12,
				"monthlyLoadAerobicHigh": 209.70,
				"monthlyLoadAnaerobic": 0.0,
				"trainingBalanceFeedbackPhrase": "AEROBIC_HIGH_SHORTAGE",
				"primaryTrainingDevice": true
			}
		},
		"recordedDevices": [
			{"deviceId": 3490854130, "deviceName": "Forerunner 965", "category": 0}
		]
	}`

	var balance TrainingLoadBalance
	if err := json.Unmarshal([]byte(rawJSON), &balance); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if balance.UserID != 12345678 {
		t.Errorf("UserID = %d, want 12345678", balance.UserID)
	}

	data := balance.MetricsTrainingLoadBalanceDTOMap["3490854130"]
	if data == nil {
		t.Fatal("Device data is nil")
	}
	if data.TrainingBalanceFeedbackPhrase != "AEROBIC_HIGH_SHORTAGE" {
		t.Errorf("TrainingBalanceFeedbackPhrase = %s, want AEROBIC_HIGH_SHORTAGE", data.TrainingBalanceFeedbackPhrase)
	}
}

func TestHeatAltitudeAcclimationJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"calendarDate": "2026-01-27",
		"altitudeAcclimationDate": "2026-01-26",
		"heatAcclimationDate": "2026-01-26",
		"altitudeAcclimation": 0,
		"heatAcclimationPercentage": 100,
		"heatTrend": "ACCLIMATIZED",
		"altitudeTrend": null,
		"currentAltitude": 0
	}`

	var acclimation HeatAltitudeAcclimation
	if err := json.Unmarshal([]byte(rawJSON), &acclimation); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if acclimation.CalendarDate != testDateMetrics {
		t.Errorf("CalendarDate = %s, want %s", acclimation.CalendarDate, testDateMetrics)
	}
	if acclimation.HeatAcclimationPercentage != 100 {
		t.Errorf("HeatAcclimationPercentage = %d, want 100", acclimation.HeatAcclimationPercentage)
	}
	if acclimation.HeatTrend != "ACCLIMATIZED" {
		t.Errorf("HeatTrend = %s, want ACCLIMATIZED", acclimation.HeatTrend)
	}
}

func TestTrainingReadinessRawJSON(t *testing.T) {
	rawJSON := `[{"score":91}]`

	var entries []TrainingReadinessEntry
	if err := json.Unmarshal([]byte(rawJSON), &entries); err != nil {
		t.Fatal(err)
	}

	tr := &TrainingReadiness{
		Entries: entries,
		raw:     json.RawMessage(rawJSON),
	}

	if string(tr.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
