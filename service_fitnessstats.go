// service_fitnessstats.go
package garmin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FitnessMetric represents a metric type for fitness stats.
type FitnessMetric string

const (
	MetricCalories                FitnessMetric = "calories"
	MetricDistance                FitnessMetric = "distance"
	MetricDuration                FitnessMetric = "duration"
	MetricAvgSpeed                FitnessMetric = "avgSpeed"
	MetricMaxHR                   FitnessMetric = "maxHr"
	MetricAvgHR                   FitnessMetric = "avgHr"
	MetricElevationGain           FitnessMetric = "elevationGain"
	MetricAvgRunCadence           FitnessMetric = "avgRunCadence"
	MetricAvgGroundContactBalance FitnessMetric = "avgGroundContactBalance"
	MetricAvgStrideLength         FitnessMetric = "avgStrideLength"
	MetricAvgVerticalOscillation  FitnessMetric = "avgVerticalOscillation"
	MetricAvgVerticalRatio        FitnessMetric = "avgVerticalRatio"
	MetricAvgGroundContactTime    FitnessMetric = "avgGroundContactTime"
	MetricStartLocal              FitnessMetric = "startLocal"
	MetricActivityType            FitnessMetric = "activityType"
	MetricActivitySubType         FitnessMetric = "activitySubType"
	MetricName                    FitnessMetric = "name"
	MetricAerobicTrainingEffect   FitnessMetric = "aerobicTrainingEffect"
	MetricAnaerobicTrainingEffect FitnessMetric = "anaerobicTrainingEffect"
)

// WeekStartDay represents the first day of the week for aggregation.
type WeekStartDay string

const (
	WeekStartMonday   WeekStartDay = "monday"
	WeekStartSunday   WeekStartDay = "sunday"
	WeekStartSaturday WeekStartDay = "saturday"
)

// FitnessStatsOptions configures the fitness stats query.
type FitnessStatsOptions struct {
	// Aggregation type (daily, weekly, monthly, yearly). Required for aggregated stats.
	Aggregation Aggregation
	// StartDate for the stats range. Required.
	StartDate time.Time
	// EndDate for the stats range. Required.
	EndDate time.Time
	// Metrics to include in the response. At least one required.
	Metrics []FitnessMetric
	// UserFirstDay sets the first day of the week for weekly aggregation.
	UserFirstDay WeekStartDay
	// GroupByActivityType groups stats by activity type (e.g., running, hiking).
	GroupByActivityType bool
	// GroupByParentActivityType groups stats by parent activity type.
	GroupByParentActivityType bool
	// StandardizedUnits uses standardized units in the response.
	StandardizedUnits bool
	// MinimumDistance filters activities by minimum distance (in meters).
	MinimumDistance *float64
	// ActivityType filters by activity type (e.g., "running", "hiking"). Used with /activity/all endpoint.
	ActivityType string
}

// MetricValue represents statistics for a single metric.
type MetricValue struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
}

// FitnessStatEntry represents fitness statistics for a single time period.
type FitnessStatEntry struct {
	// Date is the date or start date of the aggregation period.
	Date string `json:"date"`
	// CountOfActivities is the number of activities in this period.
	CountOfActivities int `json:"countOfActivities"`
	// Stats contains metrics grouped by activity type (or "all" if not grouped).
	// The outer map key is the activity type (e.g., "all", "running", "hiking").
	// The inner map key is the metric name (e.g., "duration", "distance", "calories").
	Stats map[string]map[string]MetricValue `json:"stats"`
}

// FitnessStats represents fitness statistics over a date range.
type FitnessStats struct {
	Entries []FitnessStatEntry
	raw     json.RawMessage
}

// RawJSON returns the original JSON response.
func (f *FitnessStats) RawJSON() json.RawMessage {
	return f.raw
}

// SetRaw sets the raw JSON response.
func (f *FitnessStats) SetRaw(data json.RawMessage) {
	f.raw = data
}

// UnmarshalJSON unmarshals the array response into the Entries field.
func (f *FitnessStats) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Entries)
}

// FitnessStatsActivity represents a single activity from the /activity/all endpoint.
type FitnessStatsActivity struct {
	ActivityID              int64    `json:"activityId"`
	Name                    string   `json:"name,omitempty"`
	StartLocal              string   `json:"startLocal,omitempty"`
	ActivityType            string   `json:"activityType,omitempty"`
	ActivitySubType         string   `json:"activitySubType,omitempty"`
	AerobicTrainingEffect   *float64 `json:"aerobicTrainingEffect,omitempty"`
	AnaerobicTrainingEffect *float64 `json:"anaerobicTrainingEffect,omitempty"`
	Duration                *float64 `json:"duration,omitempty"`
	Distance                *float64 `json:"distance,omitempty"`
	Calories                *float64 `json:"calories,omitempty"`
	AvgSpeed                *float64 `json:"avgSpeed,omitempty"`
	MaxHR                   *float64 `json:"maxHr,omitempty"`
	AvgHR                   *float64 `json:"avgHr,omitempty"`
	ElevationGain           *float64 `json:"elevationGain,omitempty"`
	AvgRunCadence           *float64 `json:"avgRunCadence,omitempty"`
	AvgGroundContactBalance *float64 `json:"avgGroundContactBalance,omitempty"`
	AvgStrideLength         *float64 `json:"avgStrideLength,omitempty"`
	AvgVerticalOscillation  *float64 `json:"avgVerticalOscillation,omitempty"`
	AvgVerticalRatio        *float64 `json:"avgVerticalRatio,omitempty"`
	AvgGroundContactTime    *float64 `json:"avgGroundContactTime,omitempty"`
}

// FitnessStatsActivities represents a list of activities from the /activity/all endpoint.
type FitnessStatsActivities struct {
	Activities []FitnessStatsActivity
	raw        json.RawMessage
}

// RawJSON returns the original JSON response.
func (f *FitnessStatsActivities) RawJSON() json.RawMessage {
	return f.raw
}

// SetRaw sets the raw JSON response.
func (f *FitnessStatsActivities) SetRaw(data json.RawMessage) {
	f.raw = data
}

// UnmarshalJSON unmarshals the array response into the Activities field.
func (f *FitnessStatsActivities) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Activities)
}

// buildQuery constructs the query string from options.
func (opts *FitnessStatsOptions) buildQuery() string {
	params := url.Values{}

	params.Set("aggregation", string(opts.Aggregation))
	params.Set("startDate", opts.StartDate.Format("2006-01-02"))
	params.Set("endDate", opts.EndDate.Format("2006-01-02"))

	if opts.UserFirstDay != "" {
		params.Set("userFirstDay", string(opts.UserFirstDay))
	} else {
		params.Set("userFirstDay", string(WeekStartMonday))
	}

	params.Set("groupByActivityType", strconv.FormatBool(opts.GroupByActivityType))
	params.Set("groupByParentActivityType", strconv.FormatBool(opts.GroupByParentActivityType))
	params.Set("standardizedUnits", strconv.FormatBool(opts.StandardizedUnits))

	if opts.MinimumDistance != nil {
		params.Set("minimumDistance", fmt.Sprintf("%f", *opts.MinimumDistance))
	}

	// Build the base query string
	query := params.Encode()

	// Metrics need special handling because they appear multiple times
	// and url.Values.Encode() doesn't handle multiple values in the order we want
	metricParts := make([]string, 0, len(opts.Metrics))
	for _, m := range opts.Metrics {
		metricParts = append(metricParts, "metric="+string(m))
	}
	if len(metricParts) > 0 {
		query += "&" + strings.Join(metricParts, "&")
	}

	return query
}

// buildQueryAll constructs the query string for the /activity/all endpoint.
func (opts *FitnessStatsOptions) buildQueryAll() string {
	params := url.Values{}

	params.Set("startDate", opts.StartDate.Format("2006-01-02"))
	params.Set("endDate", opts.EndDate.Format("2006-01-02"))

	if opts.ActivityType != "" {
		params.Set("activityType", opts.ActivityType)
	}

	// Build the base query string
	query := params.Encode()

	// Metrics need special handling because they appear multiple times
	metricParts := make([]string, 0, len(opts.Metrics))
	for _, m := range opts.Metrics {
		metricParts = append(metricParts, "metric="+string(m))
	}
	if len(metricParts) > 0 {
		query += "&" + strings.Join(metricParts, "&")
	}

	return query
}

// GetActivityStats retrieves aggregated fitness statistics for activities.
func (s *FitnessStatsService) GetActivityStats(ctx context.Context, opts *FitnessStatsOptions) (*FitnessStats, error) {
	if opts == nil {
		return nil, errors.New("options are required")
	}
	if len(opts.Metrics) == 0 {
		return nil, errors.New("at least one metric is required")
	}

	path := "/fitnessstats-service/activity?" + opts.buildQuery()
	return fetch[FitnessStats](ctx, s.client, path)
}

// GetAllActivities retrieves individual activity data without aggregation.
func (s *FitnessStatsService) GetAllActivities(ctx context.Context, opts *FitnessStatsOptions) (*FitnessStatsActivities, error) {
	if opts == nil {
		return nil, errors.New("options are required")
	}
	if len(opts.Metrics) == 0 {
		return nil, errors.New("at least one metric is required")
	}

	path := "/fitnessstats-service/activity/all?" + opts.buildQueryAll()
	return fetch[FitnessStatsActivities](ctx, s.client, path)
}
