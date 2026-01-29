// service_sleep.go
package garmin

import (
	"context"
	"encoding/json"
	"time"
)

// SleepService provides access to sleep data from the Garmin sleep service.
type SleepService struct {
	client *Client
}

// DailySleepDTO represents the inner sleep data from the API.
type DailySleepDTO struct {
	ID                  *int64   `json:"id"`
	CalendarDate        string   `json:"calendarDate"`
	SleepStartTimestamp int64    `json:"sleepStartTimestampGMT"`
	SleepEndTimestamp   int64    `json:"sleepEndTimestampGMT"`
	SleepSeconds        int      `json:"sleepTimeSeconds"`
	DeepSleepSeconds    *int     `json:"deepSleepSeconds"`
	LightSleepSeconds   *int     `json:"lightSleepSeconds"`
	REMSleepSeconds     *int     `json:"remSleepSeconds"`
	AwakeSeconds        *int     `json:"awakeSleepSeconds"`
	AverageSpO2         *float64 `json:"averageSpO2Value"`
	AwakeCount          *int     `json:"awakeCount"`
	AvgSleepStress      *float64 `json:"avgSleepStress"`
}

// DailySleep represents sleep data for a single day.
type DailySleep struct {
	DailySleepDTO     DailySleepDTO `json:"dailySleepDTO"`
	REMSleepData      bool          `json:"remSleepData"`
	BodyBatteryChange *int          `json:"bodyBatteryChange"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailySleep) RawJSON() json.RawMessage { return d.raw }

// SetRaw sets the raw JSON response.
func (d *DailySleep) SetRaw(data json.RawMessage) { d.raw = data }

// SleepStart returns the sleep start time.
func (d *DailySleep) SleepStart() time.Time {
	return time.UnixMilli(d.DailySleepDTO.SleepStartTimestamp)
}

// SleepEnd returns the sleep end time.
func (d *DailySleep) SleepEnd() time.Time {
	return time.UnixMilli(d.DailySleepDTO.SleepEndTimestamp)
}

// Duration returns the total sleep duration.
func (d *DailySleep) Duration() time.Duration {
	return time.Duration(d.DailySleepDTO.SleepSeconds) * time.Second
}

// HasData returns true if actual sleep data was recorded.
func (d *DailySleep) HasData() bool {
	return d.DailySleepDTO.ID != nil
}

// GetDaily retrieves sleep data for the specified date.
func (s *SleepService) GetDaily(ctx context.Context, date time.Time) (*DailySleep, error) {
	return fetch[DailySleep](ctx, s.client, "/sleep-service/sleep/dailySleepData?date="+date.Format("2006-01-02"))
}
