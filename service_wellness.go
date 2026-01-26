// service_wellness.go
package garmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

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
