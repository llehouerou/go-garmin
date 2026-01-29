// service_fitnessage.go
package garmin

import (
	"context"
	"encoding/json"
	"time"
)

// FitnessAgeValues contains the fitness age metrics for a specific day.
type FitnessAgeValues struct {
	AchievableFitnessAge float64 `json:"achievableFitnessAge"`
	VigorousDaysAvg      float64 `json:"vigorousDaysAvg"`
	FitnessAge           float64 `json:"fitnessAge"`
	RHR                  int     `json:"rhr"`
	BMI                  float64 `json:"bmi"`
}

// FitnessAgeEntry represents fitness age data for a single day.
type FitnessAgeEntry struct {
	CalendarDate string           `json:"calendarDate"`
	Values       FitnessAgeValues `json:"values"`
}

// FitnessAgeStats represents daily fitness age statistics over a date range.
type FitnessAgeStats struct {
	Entries []FitnessAgeEntry
	raw     json.RawMessage
}

// RawJSON returns the original JSON response.
func (f *FitnessAgeStats) RawJSON() json.RawMessage {
	return f.raw
}

// SetRaw sets the raw JSON response.
func (f *FitnessAgeStats) SetRaw(data json.RawMessage) {
	f.raw = data
}

// UnmarshalJSON unmarshals the array response into the Entries field.
func (f *FitnessAgeStats) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Entries)
}

// GetStatsDaily retrieves daily fitness age statistics for a date range.
func (s *FitnessAgeService) GetStatsDaily(ctx context.Context, start, end time.Time) (*FitnessAgeStats, error) {
	path := "/fitnessage-service/stats/daily/" + start.Format("2006-01-02") + "/" + end.Format("2006-01-02")
	return fetch[FitnessAgeStats](ctx, s.client, path)
}
