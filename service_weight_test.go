// service_weight_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

const testDateWeight = "2026-01-27"

func TestWeightEntryConversions(t *testing.T) {
	weight := float64(74200) // 74.2 kg in grams
	entry := &WeightEntry{
		Weight: &weight,
	}

	// Test WeightKg
	kg := entry.WeightKg()
	if kg != 74.2 {
		t.Errorf("WeightKg() = %v, want 74.2", kg)
	}

	// Test WeightLbs
	lbs := entry.WeightLbs()
	if lbs < 163.5 || lbs > 163.6 {
		t.Errorf("WeightLbs() = %v, want ~163.58", lbs)
	}

	// Test nil weight
	nilEntry := &WeightEntry{}
	if nilEntry.WeightKg() != 0 {
		t.Error("WeightKg() should return 0 for nil weight")
	}
	if nilEntry.WeightLbs() != 0 {
		t.Error("WeightLbs() should return 0 for nil weight")
	}
}

func TestDailyWeightJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"startDate": "2026-01-27",
		"endDate": "2026-01-27",
		"dateWeightList": [
			{
				"samplePk": 1769494024116,
				"date": 1769508416496,
				"calendarDate": "2026-01-27",
				"weight": 74200.0,
				"bmi": null,
				"bodyFat": null,
				"sourceType": "MANUAL",
				"timestampGMT": 1769494016496
			}
		],
		"totalAverage": {
			"from": 1769472000000,
			"until": 1769558399999,
			"weight": 74200.0
		}
	}`

	var daily DailyWeight
	if err := json.Unmarshal([]byte(rawJSON), &daily); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if daily.StartDate != testDateWeight {
		t.Errorf("StartDate = %s, want %s", daily.StartDate, testDateWeight)
	}
	if daily.EndDate != testDateWeight {
		t.Errorf("EndDate = %s, want %s", daily.EndDate, testDateWeight)
	}
	if len(daily.DateWeightList) != 1 {
		t.Fatalf("DateWeightList length = %d, want 1", len(daily.DateWeightList))
	}
	if daily.DateWeightList[0].CalendarDate != testDateWeight {
		t.Errorf("DateWeightList[0].CalendarDate = %s, want %s", daily.DateWeightList[0].CalendarDate, testDateWeight)
	}
	if daily.DateWeightList[0].Weight == nil || *daily.DateWeightList[0].Weight != 74200.0 {
		t.Errorf("DateWeightList[0].Weight = %v, want 74200.0", daily.DateWeightList[0].Weight)
	}
	if daily.TotalAverage.Weight == nil || *daily.TotalAverage.Weight != 74200.0 {
		t.Errorf("TotalAverage.Weight = %v, want 74200.0", daily.TotalAverage.Weight)
	}
}

func TestDailyWeightRawJSON(t *testing.T) {
	rawJSON := `{"startDate":"2026-01-27","endDate":"2026-01-27"}`

	var daily DailyWeight
	if err := json.Unmarshal([]byte(rawJSON), &daily); err != nil {
		t.Fatal(err)
	}
	daily.raw = json.RawMessage(rawJSON)

	if string(daily.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestWeightRangeJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"dailyWeightSummaries": [
			{
				"summaryDate": "2026-01-27",
				"numOfWeightEntries": 1,
				"minWeight": 74200.0,
				"maxWeight": 74200.0,
				"latestWeight": {
					"samplePk": 1769494024116,
					"calendarDate": "2026-01-27",
					"weight": 74200.0,
					"sourceType": "MANUAL"
				},
				"allWeightMetrics": [
					{
						"samplePk": 1769494024116,
						"calendarDate": "2026-01-27",
						"weight": 74200.0
					}
				]
			}
		],
		"totalAverage": {
			"from": 1766880000000,
			"until": 1769558399999,
			"weight": 74200.0
		},
		"previousDateWeight": {
			"samplePk": 1764645477134,
			"calendarDate": "2025-12-02",
			"weight": 75200.0
		},
		"nextDateWeight": {
			"samplePk": null,
			"weight": null
		}
	}`

	var weightRange WeightRange
	if err := json.Unmarshal([]byte(rawJSON), &weightRange); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(weightRange.DailyWeightSummaries) != 1 {
		t.Fatalf("DailyWeightSummaries length = %d, want 1", len(weightRange.DailyWeightSummaries))
	}

	summary := weightRange.DailyWeightSummaries[0]
	if summary.SummaryDate != testDateWeight {
		t.Errorf("SummaryDate = %s, want %s", summary.SummaryDate, testDateWeight)
	}
	if summary.NumOfWeightEntries != 1 {
		t.Errorf("NumOfWeightEntries = %d, want 1", summary.NumOfWeightEntries)
	}
	if summary.MinWeight != 74200.0 {
		t.Errorf("MinWeight = %f, want 74200.0", summary.MinWeight)
	}

	if weightRange.PreviousDateWeight.CalendarDate != "2025-12-02" {
		t.Errorf("PreviousDateWeight.CalendarDate = %s, want 2025-12-02", weightRange.PreviousDateWeight.CalendarDate)
	}
	if weightRange.PreviousDateWeight.Weight == nil || *weightRange.PreviousDateWeight.Weight != 75200.0 {
		t.Errorf("PreviousDateWeight.Weight = %v, want 75200.0", weightRange.PreviousDateWeight.Weight)
	}
}

func TestWeightRangeRawJSON(t *testing.T) {
	rawJSON := `{"dailyWeightSummaries":[]}`

	var weightRange WeightRange
	if err := json.Unmarshal([]byte(rawJSON), &weightRange); err != nil {
		t.Fatal(err)
	}
	weightRange.raw = json.RawMessage(rawJSON)

	if string(weightRange.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
