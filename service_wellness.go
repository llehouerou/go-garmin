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

// HeartRateValueDescriptor describes the format of heart rate values.
type HeartRateValueDescriptor struct {
	Key   string `json:"key"`
	Index int    `json:"index"`
}

// DailyHeartRate represents heart rate data for a single day.
type DailyHeartRate struct {
	UserProfilePK                    int64                      `json:"userProfilePK"`
	CalendarDate                     string                     `json:"calendarDate"`
	StartTimestampGMT                string                     `json:"startTimestampGMT"`
	EndTimestampGMT                  string                     `json:"endTimestampGMT"`
	StartTimestampLocal              string                     `json:"startTimestampLocal"`
	EndTimestampLocal                string                     `json:"endTimestampLocal"`
	MaxHeartRate                     int                        `json:"maxHeartRate"`
	MinHeartRate                     int                        `json:"minHeartRate"`
	RestingHeartRate                 int                        `json:"restingHeartRate"`
	LastSevenDaysAvgRestingHeartRate int                        `json:"lastSevenDaysAvgRestingHeartRate"`
	HeartRateValueDescriptors        []HeartRateValueDescriptor `json:"heartRateValueDescriptors"`
	HeartRateValues                  [][]int64                  `json:"heartRateValues"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyHeartRate) RawJSON() json.RawMessage {
	return d.raw
}

// SpO2ValueDescriptor describes the format of SpO2 values.
type SpO2ValueDescriptor struct {
	Index int    `json:"spo2ValueDescriptorIndex"`
	Key   string `json:"spo2ValueDescriptorKey"`
}

// DailySpO2 represents blood oxygen (SpO2) data for a single day.
type DailySpO2 struct {
	UserProfilePK            int64                 `json:"userProfilePK"`
	CalendarDate             string                `json:"calendarDate"`
	StartTimestampGMT        string                `json:"startTimestampGMT"`
	EndTimestampGMT          string                `json:"endTimestampGMT"`
	StartTimestampLocal      string                `json:"startTimestampLocal"`
	EndTimestampLocal        string                `json:"endTimestampLocal"`
	SleepStartTimestampGMT   string                `json:"sleepStartTimestampGMT"`
	SleepEndTimestampGMT     string                `json:"sleepEndTimestampGMT"`
	SleepStartTimestampLocal string                `json:"sleepStartTimestampLocal"`
	SleepEndTimestampLocal   string                `json:"sleepEndTimestampLocal"`
	AverageSpO2              float64               `json:"averageSpO2"`
	LowestSpO2               int                   `json:"lowestSpO2"`
	LastSevenDaysAvgSpO2     float64               `json:"lastSevenDaysAvgSpO2"`
	LatestSpO2               int                   `json:"latestSpO2"`
	LatestSpO2TimestampGMT   string                `json:"latestSpO2TimestampGMT"`
	LatestSpO2TimestampLocal string                `json:"latestSpO2TimestampLocal"`
	AvgSleepSpO2             float64               `json:"avgSleepSpO2"`
	SpO2ValueDescriptors     []SpO2ValueDescriptor `json:"spO2ValueDescriptorsDTOList"`
	SpO2HourlyAverages       [][]any               `json:"spO2HourlyAverages"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailySpO2) RawJSON() json.RawMessage {
	return d.raw
}

// RespirationValueDescriptor describes the format of respiration values.
type RespirationValueDescriptor struct {
	Key   string `json:"key"`
	Index int    `json:"index"`
}

// RespirationAveragesDescriptor describes the format of respiration averages.
type RespirationAveragesDescriptor struct {
	Index int    `json:"respirationAveragesValueDescriptorIndex"`
	Key   string `json:"respirationAveragesValueDescriptionKey"`
}

// DailyRespiration represents respiration data for a single day.
type DailyRespiration struct {
	UserProfilePK                  int64                           `json:"userProfilePK"`
	CalendarDate                   string                          `json:"calendarDate"`
	StartTimestampGMT              string                          `json:"startTimestampGMT"`
	EndTimestampGMT                string                          `json:"endTimestampGMT"`
	StartTimestampLocal            string                          `json:"startTimestampLocal"`
	EndTimestampLocal              string                          `json:"endTimestampLocal"`
	SleepStartTimestampGMT         string                          `json:"sleepStartTimestampGMT"`
	SleepEndTimestampGMT           string                          `json:"sleepEndTimestampGMT"`
	SleepStartTimestampLocal       string                          `json:"sleepStartTimestampLocal"`
	SleepEndTimestampLocal         string                          `json:"sleepEndTimestampLocal"`
	LowestRespirationValue         float64                         `json:"lowestRespirationValue"`
	HighestRespirationValue        float64                         `json:"highestRespirationValue"`
	AvgWakingRespirationValue      float64                         `json:"avgWakingRespirationValue"`
	AvgSleepRespirationValue       float64                         `json:"avgSleepRespirationValue"`
	RespirationValueDescriptors    []RespirationValueDescriptor    `json:"respirationValueDescriptorsDTOList"`
	RespirationValuesArray         [][]float64                     `json:"respirationValuesArray"`
	RespirationAveragesDescriptors []RespirationAveragesDescriptor `json:"respirationAveragesValueDescriptorDTOList"`
	RespirationAveragesValuesArray [][]any                         `json:"respirationAveragesValuesArray"`
	RespirationVersion             int                             `json:"respirationVersion"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyRespiration) RawJSON() json.RawMessage {
	return d.raw
}

// IntensityMinutesValueDescriptor describes the format of intensity minutes values.
type IntensityMinutesValueDescriptor struct {
	Index int    `json:"index"`
	Key   string `json:"key"`
}

// DailyIntensityMinutes represents intensity minutes data for a single day.
type DailyIntensityMinutes struct {
	UserProfilePK       int64                             `json:"userProfilePK"`
	CalendarDate        string                            `json:"calendarDate"`
	StartTimestampGMT   string                            `json:"startTimestampGMT"`
	EndTimestampGMT     string                            `json:"endTimestampGMT"`
	StartTimestampLocal string                            `json:"startTimestampLocal"`
	EndTimestampLocal   string                            `json:"endTimestampLocal"`
	WeeklyModerate      int                               `json:"weeklyModerate"`
	WeeklyVigorous      int                               `json:"weeklyVigorous"`
	WeeklyTotal         int                               `json:"weeklyTotal"`
	WeekGoal            int                               `json:"weekGoal"`
	DayOfGoalMet        *string                           `json:"dayOfGoalMet"`
	StartDayMinutes     int                               `json:"startDayMinutes"`
	EndDayMinutes       int                               `json:"endDayMinutes"`
	ModerateMinutes     int                               `json:"moderateMinutes"`
	VigorousMinutes     int                               `json:"vigorousMinutes"`
	IMValueDescriptors  []IntensityMinutesValueDescriptor `json:"imValueDescriptorsDTOList"`
	IMValuesArray       [][]int64                         `json:"imValuesArray"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (d *DailyIntensityMinutes) RawJSON() json.RawMessage {
	return d.raw
}

// GetDailyHeartRate retrieves heart rate data for the specified date.
func (s *WellnessService) GetDailyHeartRate(ctx context.Context, date time.Time) (*DailyHeartRate, error) {
	path := "/wellness-service/wellness/dailyHeartRate/?date=" + date.Format("2006-01-02")

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

	var hr DailyHeartRate
	if err := json.Unmarshal(raw, &hr); err != nil {
		return nil, err
	}
	hr.raw = raw

	return &hr, nil
}

// GetDailySpO2 retrieves blood oxygen (SpO2) data for the specified date.
func (s *WellnessService) GetDailySpO2(ctx context.Context, date time.Time) (*DailySpO2, error) {
	path := "/wellness-service/wellness/daily/spo2/" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var spo2 DailySpO2
	if err := json.Unmarshal(raw, &spo2); err != nil {
		return nil, err
	}
	spo2.raw = raw

	return &spo2, nil
}

// GetDailyRespiration retrieves respiration data for the specified date.
func (s *WellnessService) GetDailyRespiration(ctx context.Context, date time.Time) (*DailyRespiration, error) {
	path := "/wellness-service/wellness/daily/respiration/" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp2 DailyRespiration
	if err := json.Unmarshal(raw, &resp2); err != nil {
		return nil, err
	}
	resp2.raw = raw

	return &resp2, nil
}

// GetDailyIntensityMinutes retrieves intensity minutes data for the specified date.
func (s *WellnessService) GetDailyIntensityMinutes(ctx context.Context, date time.Time) (*DailyIntensityMinutes, error) {
	path := "/wellness-service/wellness/daily/im/" + date.Format("2006-01-02")

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var im DailyIntensityMinutes
	if err := json.Unmarshal(raw, &im); err != nil {
		return nil, err
	}
	im.raw = raw

	return &im, nil
}
