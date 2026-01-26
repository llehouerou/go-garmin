// service_wellness.go
package garmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// DailySleep represents sleep data for a single day.
type DailySleep struct {
	CalendarDate        string   `json:"calendarDate"`
	SleepStartTimestamp int64    `json:"sleepStartTimestampGMT"`
	SleepEndTimestamp   int64    `json:"sleepEndTimestampGMT"`
	SleepSeconds        int      `json:"sleepTimeSeconds"`
	DeepSleepSeconds    int      `json:"deepSleepSeconds"`
	LightSleepSeconds   int      `json:"lightSleepSeconds"`
	REMSleepSeconds     int      `json:"remSleepSeconds"`
	AwakeSeconds        int      `json:"awakeSleepSeconds"`
	AverageSpO2         *float64 `json:"averageSpO2Value"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailySleep) RawJSON() json.RawMessage {
	return d.raw
}

// SleepStart returns the sleep start time.
func (d *DailySleep) SleepStart() time.Time {
	return time.UnixMilli(d.SleepStartTimestamp)
}

// SleepEnd returns the sleep end time.
func (d *DailySleep) SleepEnd() time.Time {
	return time.UnixMilli(d.SleepEndTimestamp)
}

// Duration returns the total sleep duration.
func (d *DailySleep) Duration() time.Duration {
	return time.Duration(d.SleepSeconds) * time.Second
}

// DailyStress represents stress data for a single day.
type DailyStress struct {
	CalendarDate       string `json:"calendarDate"`
	OverallStressLevel int    `json:"overallStressLevel"`
	HighStressDuration int    `json:"highStressDuration"`
	MedStressDuration  int    `json:"mediumStressDuration"`
	LowStressDuration  int    `json:"lowStressDuration"`
	RestStressDuration int    `json:"restStressDuration"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyStress) RawJSON() json.RawMessage {
	return d.raw
}

// BodyBatteryReport represents body battery data for a day.
type BodyBatteryReport struct {
	Date         string `json:"date"`
	Charged      int    `json:"charged"`
	Drained      int    `json:"drained"`
	StartLevel   int    `json:"startOfDayBodyBattery"`
	EndLevel     int    `json:"endOfDayBodyBattery"`
	HighestLevel int    `json:"maxBodyBattery"`
	LowestLevel  int    `json:"minBodyBattery"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (b *BodyBatteryReport) RawJSON() json.RawMessage {
	return b.raw
}

// GetDailySleep retrieves sleep data for the specified date.
func (s *WellnessService) GetDailySleep(ctx context.Context, date time.Time) (*DailySleep, error) {
	path := "/wellness-service/wellness/dailySleepData/" + date.Format("2006-01-02")

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

// GetBodyBattery retrieves body battery data for the specified date.
func (s *WellnessService) GetBodyBattery(ctx context.Context, date time.Time) (*BodyBatteryReport, error) {
	path := "/wellness-service/wellness/bodyBattery/reports/daily/" + date.Format("2006-01-02")

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

	var battery BodyBatteryReport
	if err := json.Unmarshal(raw, &battery); err != nil {
		return nil, err
	}
	battery.raw = raw

	return &battery, nil
}
