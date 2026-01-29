// service_weight.go
package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// WeightEntry represents a single weight measurement.
type WeightEntry struct {
	SamplePK       *int64   `json:"samplePk"`
	Date           *int64   `json:"date"`
	CalendarDate   string   `json:"calendarDate"`
	Weight         *float64 `json:"weight"`
	BMI            *float64 `json:"bmi"`
	BodyFat        *float64 `json:"bodyFat"`
	BodyWater      *float64 `json:"bodyWater"`
	BoneMass       *float64 `json:"boneMass"`
	MuscleMass     *float64 `json:"muscleMass"`
	PhysiqueRating *float64 `json:"physiqueRating"`
	VisceralFat    *float64 `json:"visceralFat"`
	MetabolicAge   *int     `json:"metabolicAge"`
	SourceType     *string  `json:"sourceType"`
	TimestampGMT   *int64   `json:"timestampGMT"`
	WeightDelta    *float64 `json:"weightDelta"`
}

// WeightKg returns the weight in kilograms.
func (w *WeightEntry) WeightKg() float64 {
	if w.Weight == nil {
		return 0
	}
	return *w.Weight / 1000
}

// WeightLbs returns the weight in pounds.
func (w *WeightEntry) WeightLbs() float64 {
	if w.Weight == nil {
		return 0
	}
	return *w.Weight / 1000 * 2.20462
}

// WeightAverage represents average weight data over a period.
type WeightAverage struct {
	From           int64    `json:"from"`
	Until          int64    `json:"until"`
	Weight         *float64 `json:"weight"`
	BMI            *float64 `json:"bmi"`
	BodyFat        *float64 `json:"bodyFat"`
	BodyWater      *float64 `json:"bodyWater"`
	BoneMass       *float64 `json:"boneMass"`
	MuscleMass     *float64 `json:"muscleMass"`
	PhysiqueRating *float64 `json:"physiqueRating"`
	VisceralFat    *float64 `json:"visceralFat"`
	MetabolicAge   *int     `json:"metabolicAge"`
}

// DailyWeight represents weight data for a single day.
type DailyWeight struct {
	StartDate      string        `json:"startDate"`
	EndDate        string        `json:"endDate"`
	DateWeightList []WeightEntry `json:"dateWeightList"`
	TotalAverage   WeightAverage `json:"totalAverage"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyWeight) RawJSON() json.RawMessage { return d.raw }

// SetRaw sets the raw JSON response.
func (d *DailyWeight) SetRaw(data json.RawMessage) { d.raw = data }

// DailyWeightSummary represents a summary of weight entries for a single day.
type DailyWeightSummary struct {
	SummaryDate        string        `json:"summaryDate"`
	NumOfWeightEntries int           `json:"numOfWeightEntries"`
	MinWeight          float64       `json:"minWeight"`
	MaxWeight          float64       `json:"maxWeight"`
	LatestWeight       WeightEntry   `json:"latestWeight"`
	AllWeightMetrics   []WeightEntry `json:"allWeightMetrics"`
}

// WeightRange represents weight data for a date range.
type WeightRange struct {
	DailyWeightSummaries []DailyWeightSummary `json:"dailyWeightSummaries"`
	TotalAverage         WeightAverage        `json:"totalAverage"`
	PreviousDateWeight   WeightEntry          `json:"previousDateWeight"`
	NextDateWeight       WeightEntry          `json:"nextDateWeight"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (r *WeightRange) RawJSON() json.RawMessage { return r.raw }

// SetRaw sets the raw JSON response.
func (r *WeightRange) SetRaw(data json.RawMessage) { r.raw = data }

// GetDaily retrieves weight data for the specified date.
func (s *WeightService) GetDaily(ctx context.Context, date time.Time) (*DailyWeight, error) {
	return fetch[DailyWeight](ctx, s.client, "/weight-service/weight/dayview/"+date.Format("2006-01-02"))
}

// GetRange retrieves weight data for a date range.
func (s *WeightService) GetRange(ctx context.Context, startDate, endDate time.Time) (*WeightRange, error) {
	return fetch[WeightRange](ctx, s.client, fmt.Sprintf("/weight-service/weight/range/%s/%s?includeAll=true",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02")))
}
