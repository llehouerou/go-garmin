// service_hrv.go
package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// HRVBaseline represents the baseline HRV values.
type HRVBaseline struct {
	LowUpper      int     `json:"lowUpper"`
	BalancedLow   int     `json:"balancedLow"`
	BalancedUpper int     `json:"balancedUpper"`
	MarkerValue   float64 `json:"markerValue"`
}

// HRVSummary represents the summary of HRV data for a day.
type HRVSummary struct {
	CalendarDate      string      `json:"calendarDate"`
	WeeklyAvg         int         `json:"weeklyAvg"`
	LastNightAvg      int         `json:"lastNightAvg"`
	LastNight5MinHigh int         `json:"lastNight5MinHigh"`
	Baseline          HRVBaseline `json:"baseline"`
	Status            string      `json:"status"`
	FeedbackPhrase    string      `json:"feedbackPhrase"`
	CreateTimeStamp   string      `json:"createTimeStamp"`
}

// HRVReading represents a single HRV reading.
type HRVReading struct {
	HRVValue         int    `json:"hrvValue"`
	ReadingTimeGMT   string `json:"readingTimeGMT"`
	ReadingTimeLocal string `json:"readingTimeLocal"`
}

// DailyHRV represents HRV data for a single day.
type DailyHRV struct {
	UserProfilePK            int64        `json:"userProfilePk"`
	HRVSummary               HRVSummary   `json:"hrvSummary"`
	HRVReadings              []HRVReading `json:"hrvReadings"`
	StartTimestampGMT        string       `json:"startTimestampGMT"`
	EndTimestampGMT          string       `json:"endTimestampGMT"`
	StartTimestampLocal      string       `json:"startTimestampLocal"`
	EndTimestampLocal        string       `json:"endTimestampLocal"`
	SleepStartTimestampGMT   string       `json:"sleepStartTimestampGMT"`
	SleepEndTimestampGMT     string       `json:"sleepEndTimestampGMT"`
	SleepStartTimestampLocal string       `json:"sleepStartTimestampLocal"`
	SleepEndTimestampLocal   string       `json:"sleepEndTimestampLocal"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyHRV) RawJSON() json.RawMessage { return d.raw }

// SetRaw sets the raw JSON response.
func (d *DailyHRV) SetRaw(data json.RawMessage) { d.raw = data }

// HRVRange represents HRV summaries for a date range.
type HRVRange struct {
	HRVSummaries  []HRVSummary `json:"hrvSummaries"`
	UserProfilePK int64        `json:"userProfilePk"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (r *HRVRange) RawJSON() json.RawMessage { return r.raw }

// SetRaw sets the raw JSON response.
func (r *HRVRange) SetRaw(data json.RawMessage) { r.raw = data }

// GetDaily retrieves HRV data for the specified date.
func (s *HRVService) GetDaily(ctx context.Context, date time.Time) (*DailyHRV, error) {
	return fetch[DailyHRV](ctx, s.client, "/hrv-service/hrv/"+date.Format("2006-01-02"))
}

// GetRange retrieves HRV summaries for a date range.
func (s *HRVService) GetRange(ctx context.Context, startDate, endDate time.Time) (*HRVRange, error) {
	return fetch[HRVRange](ctx, s.client, fmt.Sprintf("/hrv-service/hrv/daily/%s/%s",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02")))
}
