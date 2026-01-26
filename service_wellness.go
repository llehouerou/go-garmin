// service_wellness.go
package garmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

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
func (d *DailySleep) RawJSON() json.RawMessage {
	return d.raw
}

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

// DailyStress represents stress and body battery data for a single day.
type DailyStress struct {
	CalendarDate           string  `json:"calendarDate"`
	MaxStressLevel         int     `json:"maxStressLevel"`
	AvgStressLevel         int     `json:"avgStressLevel"`
	StressChartValueOffset int     `json:"stressChartValueOffset"`
	StressChartYAxisOrigin int     `json:"stressChartYAxisOrigin"`
	StressValuesArray      [][]int `json:"stressValuesArray"`
	BodyBatteryValuesArray [][]any `json:"bodyBatteryValuesArray"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyStress) RawJSON() json.RawMessage {
	return d.raw
}

// BodyBatteryEvent represents a single body battery event (sleep, activity, etc).
type BodyBatteryEvent struct {
	Event *struct {
		EventType         string `json:"eventType"`
		EventStartTimeGMT string `json:"eventStartTimeGmt"`
		TimezoneOffset    int64  `json:"timezoneOffset"`
		DurationMs        int64  `json:"durationInMilliseconds"`
		BodyBatteryImpact int    `json:"bodyBatteryImpact"`
		FeedbackType      string `json:"feedbackType"`
		ShortFeedback     string `json:"shortFeedback"`
	} `json:"event"`
	ActivityName           *string  `json:"activityName"`
	ActivityType           *string  `json:"activityType"`
	ActivityID             any      `json:"activityId"`
	AverageStress          *float64 `json:"averageStress"`
	StressValuesArray      [][]int  `json:"stressValuesArray"`
	BodyBatteryValuesArray [][]any  `json:"bodyBatteryValuesArray"`
}

// BodyBatteryEvents represents all body battery events for a day.
type BodyBatteryEvents struct {
	Events []BodyBatteryEvent
	raw    json.RawMessage
}

// RawJSON returns the original JSON response.
func (b *BodyBatteryEvents) RawJSON() json.RawMessage {
	return b.raw
}

// GetDailySleep retrieves sleep data for the specified date.
func (s *WellnessService) GetDailySleep(ctx context.Context, date time.Time) (*DailySleep, error) {
	path := "/sleep-service/sleep/dailySleepData?date=" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sleep DailySleep
	if err := json.Unmarshal(raw, &sleep); err != nil {
		return nil, err
	}
	sleep.raw = raw

	return &sleep, nil
}

// GetDailyStress retrieves stress data for the specified date.
func (s *WellnessService) GetDailyStress(ctx context.Context, date time.Time) (*DailyStress, error) {
	path := "/wellness-service/wellness/dailyStress/" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stress DailyStress
	if err := json.Unmarshal(raw, &stress); err != nil {
		return nil, err
	}
	stress.raw = raw

	return &stress, nil
}

// GetBodyBatteryEvents retrieves body battery events for the specified date.
func (s *WellnessService) GetBodyBatteryEvents(ctx context.Context, date time.Time) (*BodyBatteryEvents, error) {
	path := "/wellness-service/wellness/bodyBattery/events/" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []BodyBatteryEvent
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, err
	}

	return &BodyBatteryEvents{Events: events, raw: raw}, nil
}
