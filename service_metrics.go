// service_metrics.go
package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TrainingReadinessEntry represents a single training readiness measurement.
type TrainingReadinessEntry struct {
	UserProfilePK                     int64   `json:"userProfilePK"`
	CalendarDate                      string  `json:"calendarDate"`
	Timestamp                         string  `json:"timestamp"`
	TimestampLocal                    string  `json:"timestampLocal"`
	DeviceID                          int64   `json:"deviceId"`
	Level                             string  `json:"level"`
	FeedbackLong                      string  `json:"feedbackLong"`
	FeedbackShort                     string  `json:"feedbackShort"`
	Score                             int     `json:"score"`
	SleepScore                        *int    `json:"sleepScore"`
	SleepScoreFactorPercent           int     `json:"sleepScoreFactorPercent"`
	SleepScoreFactorFeedback          string  `json:"sleepScoreFactorFeedback"`
	RecoveryTime                      int     `json:"recoveryTime"`
	RecoveryTimeFactorPercent         int     `json:"recoveryTimeFactorPercent"`
	RecoveryTimeFactorFeedback        string  `json:"recoveryTimeFactorFeedback"`
	AcwrFactorPercent                 int     `json:"acwrFactorPercent"`
	AcwrFactorFeedback                string  `json:"acwrFactorFeedback"`
	AcuteLoad                         int     `json:"acuteLoad"`
	StressHistoryFactorPercent        int     `json:"stressHistoryFactorPercent"`
	StressHistoryFactorFeedback       string  `json:"stressHistoryFactorFeedback"`
	HRVFactorPercent                  int     `json:"hrvFactorPercent"`
	HRVFactorFeedback                 string  `json:"hrvFactorFeedback"`
	HRVWeeklyAverage                  int     `json:"hrvWeeklyAverage"`
	SleepHistoryFactorPercent         int     `json:"sleepHistoryFactorPercent"`
	SleepHistoryFactorFeedback        string  `json:"sleepHistoryFactorFeedback"`
	ValidSleep                        bool    `json:"validSleep"`
	InputContext                      string  `json:"inputContext"`
	PrimaryActivityTracker            bool    `json:"primaryActivityTracker"`
	RecoveryTimeChangePhrase          *string `json:"recoveryTimeChangePhrase"`
	SleepHistoryFactorFeedbackPhrase  *string `json:"sleepHistoryFactorFeedbackPhrase"`
	HRVFactorFeedbackPhrase           *string `json:"hrvFactorFeedbackPhrase"`
	StressHistoryFactorFeedbackPhrase *string `json:"stressHistoryFactorFeedbackPhrase"`
	AcwrFactorFeedbackPhrase          *string `json:"acwrFactorFeedbackPhrase"`
	RecoveryTimeFactorFeedbackPhrase  *string `json:"recoveryTimeFactorFeedbackPhrase"`
	SleepScoreFactorFeedbackPhrase    *string `json:"sleepScoreFactorFeedbackPhrase"`
}

// TrainingReadiness represents the training readiness response.
type TrainingReadiness struct {
	Entries []TrainingReadinessEntry

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (t *TrainingReadiness) RawJSON() json.RawMessage {
	return t.raw
}

// ScoreContributor represents a contributor to an endurance or hill score.
type ScoreContributor struct {
	ActivityTypeID *int    `json:"activityTypeId"`
	Group          *int    `json:"group"`
	Contribution   float64 `json:"contribution"`
}

// EnduranceScore represents the endurance score response.
type EnduranceScore struct {
	UserProfilePK                        int64              `json:"userProfilePK"`
	DeviceID                             int64              `json:"deviceId"`
	CalendarDate                         string             `json:"calendarDate"`
	OverallScore                         int                `json:"overallScore"`
	Classification                       int                `json:"classification"`
	FeedbackPhrase                       int                `json:"feedbackPhrase"`
	PrimaryTrainingDevice                bool               `json:"primaryTrainingDevice"`
	GaugeLowerLimit                      int                `json:"gaugeLowerLimit"`
	ClassificationLowerLimitIntermediate int                `json:"classificationLowerLimitIntermediate"`
	ClassificationLowerLimitTrained      int                `json:"classificationLowerLimitTrained"`
	ClassificationLowerLimitWellTrained  int                `json:"classificationLowerLimitWellTrained"`
	ClassificationLowerLimitExpert       int                `json:"classificationLowerLimitExpert"`
	ClassificationLowerLimitSuperior     int                `json:"classificationLowerLimitSuperior"`
	ClassificationLowerLimitElite        int                `json:"classificationLowerLimitElite"`
	GaugeUpperLimit                      int                `json:"gaugeUpperLimit"`
	Contributors                         []ScoreContributor `json:"contributors"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (e *EnduranceScore) RawJSON() json.RawMessage {
	return e.raw
}

// HillScore represents the hill score response.
type HillScore struct {
	UserProfilePK             int64   `json:"userProfilePK"`
	DeviceID                  int64   `json:"deviceId"`
	CalendarDate              string  `json:"calendarDate"`
	StrengthScore             int     `json:"strengthScore"`
	EnduranceScore            int     `json:"enduranceScore"`
	HillScoreClassificationID int     `json:"hillScoreClassificationId"`
	OverallScore              int     `json:"overallScore"`
	HillScoreFeedbackPhraseID int     `json:"hillScoreFeedbackPhraseId"`
	VO2Max                    float64 `json:"vo2Max"`
	VO2MaxPreciseValue        float64 `json:"vo2MaxPreciseValue"`
	PrimaryTrainingDevice     bool    `json:"primaryTrainingDevice"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (h *HillScore) RawJSON() json.RawMessage {
	return h.raw
}

// HeatAltitudeAcclimation represents heat and altitude acclimation data.
type HeatAltitudeAcclimation struct {
	CalendarDate                      string  `json:"calendarDate"`
	AltitudeAcclimationDate           string  `json:"altitudeAcclimationDate"`
	PreviousAltitudeAcclimationDate   string  `json:"previousAltitudeAcclimationDate"`
	HeatAcclimationDate               string  `json:"heatAcclimationDate"`
	PreviousHeatAcclimationDate       string  `json:"previousHeatAcclimationDate"`
	AltitudeAcclimation               int     `json:"altitudeAcclimation"`
	PreviousAltitudeAcclimation       int     `json:"previousAltitudeAcclimation"`
	HeatAcclimationPercentage         int     `json:"heatAcclimationPercentage"`
	PreviousHeatAcclimationPercentage int     `json:"previousHeatAcclimationPercentage"`
	HeatTrend                         string  `json:"heatTrend"`
	AltitudeTrend                     *string `json:"altitudeTrend"`
	CurrentAltitude                   int     `json:"currentAltitude"`
	PreviousAltitude                  int     `json:"previousAltitude"`
	AcclimationPercentage             int     `json:"acclimationPercentage"`
	PreviousAcclimationPercentage     int     `json:"previousAcclimationPercentage"`
	AltitudeAcclimationLocalTimestamp string  `json:"altitudeAcclimationLocalTimestamp"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (h *HeatAltitudeAcclimation) RawJSON() json.RawMessage {
	return h.raw
}

// VO2MaxGeneric represents generic VO2 max data.
type VO2MaxGeneric struct {
	CalendarDate          string  `json:"calendarDate"`
	VO2MaxPreciseValue    float64 `json:"vo2MaxPreciseValue"`
	VO2MaxValue           float64 `json:"vo2MaxValue"`
	FitnessAge            *int    `json:"fitnessAge"`
	FitnessAgeDescription *string `json:"fitnessAgeDescription"`
	MaxMetCategory        int     `json:"maxMetCategory"`
}

// MaxMetEntry represents a single VO2 max / MET entry.
type MaxMetEntry struct {
	UserID                  int64                    `json:"userId"`
	Generic                 *VO2MaxGeneric           `json:"generic"`
	Cycling                 *VO2MaxGeneric           `json:"cycling"`
	HeatAltitudeAcclimation *HeatAltitudeAcclimation `json:"heatAltitudeAcclimation"`
}

// MaxMetLatest represents the latest VO2 max / MET data.
type MaxMetLatest struct {
	MaxMetEntry

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (m *MaxMetLatest) RawJSON() json.RawMessage {
	return m.raw
}

// MaxMetDaily represents VO2 max / MET data for a date range.
type MaxMetDaily struct {
	Entries []MaxMetEntry

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (m *MaxMetDaily) RawJSON() json.RawMessage {
	return m.raw
}

// AcuteTrainingLoad represents acute training load data.
type AcuteTrainingLoad struct {
	AcwrPercent                    int     `json:"acwrPercent"`
	AcwrStatus                     string  `json:"acwrStatus"`
	AcwrStatusFeedback             string  `json:"acwrStatusFeedback"`
	DailyTrainingLoadAcute         int     `json:"dailyTrainingLoadAcute"`
	MaxTrainingLoadChronic         float64 `json:"maxTrainingLoadChronic"`
	MinTrainingLoadChronic         float64 `json:"minTrainingLoadChronic"`
	DailyTrainingLoadChronic       int     `json:"dailyTrainingLoadChronic"`
	DailyAcuteChronicWorkloadRatio float64 `json:"dailyAcuteChronicWorkloadRatio"`
}

// TrainingStatusData represents training status data for a device.
type TrainingStatusData struct {
	CalendarDate                 string             `json:"calendarDate"`
	SinceDate                    string             `json:"sinceDate"`
	WeeklyTrainingLoad           *int               `json:"weeklyTrainingLoad"`
	TrainingStatus               int                `json:"trainingStatus"`
	Timestamp                    int64              `json:"timestamp"`
	DeviceID                     int64              `json:"deviceId"`
	LoadTunnelMin                *int               `json:"loadTunnelMin"`
	LoadTunnelMax                *int               `json:"loadTunnelMax"`
	LoadLevelTrend               *string            `json:"loadLevelTrend"`
	Sport                        *string            `json:"sport"`
	SubSport                     *string            `json:"subSport"`
	FitnessTrendSport            string             `json:"fitnessTrendSport"`
	FitnessTrend                 int                `json:"fitnessTrend"`
	TrainingStatusFeedbackPhrase string             `json:"trainingStatusFeedbackPhrase"`
	TrainingPaused               bool               `json:"trainingPaused"`
	AcuteTrainingLoadDTO         *AcuteTrainingLoad `json:"acuteTrainingLoadDTO"`
	PrimaryTrainingDevice        bool               `json:"primaryTrainingDevice"`
}

// RecordedDevice represents a recorded device.
type RecordedDevice struct {
	DeviceID   int64  `json:"deviceId"`
	ImageURL   string `json:"imageURL"`
	DeviceName string `json:"deviceName"`
	Category   int    `json:"category"`
}

// TrainingStatusDaily represents daily training status.
type TrainingStatusDaily struct {
	UserID                   int64                          `json:"userId"`
	LatestTrainingStatusData map[string]*TrainingStatusData `json:"latestTrainingStatusData"`
	RecordedDevices          []RecordedDevice               `json:"recordedDevices"`
	ShowSelector             bool                           `json:"showSelector"`
	LastPrimarySyncDate      string                         `json:"lastPrimarySyncDate"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (t *TrainingStatusDaily) RawJSON() json.RawMessage {
	return t.raw
}

// TrainingLoadBalanceData represents training load balance data for a device.
type TrainingLoadBalanceData struct {
	CalendarDate                    string  `json:"calendarDate"`
	DeviceID                        int64   `json:"deviceId"`
	MonthlyLoadAerobicLow           float64 `json:"monthlyLoadAerobicLow"`
	MonthlyLoadAerobicHigh          float64 `json:"monthlyLoadAerobicHigh"`
	MonthlyLoadAnaerobic            float64 `json:"monthlyLoadAnaerobic"`
	MonthlyLoadAerobicLowTargetMin  int     `json:"monthlyLoadAerobicLowTargetMin"`
	MonthlyLoadAerobicLowTargetMax  int     `json:"monthlyLoadAerobicLowTargetMax"`
	MonthlyLoadAerobicHighTargetMin int     `json:"monthlyLoadAerobicHighTargetMin"`
	MonthlyLoadAerobicHighTargetMax int     `json:"monthlyLoadAerobicHighTargetMax"`
	MonthlyLoadAnaerobicTargetMin   int     `json:"monthlyLoadAnaerobicTargetMin"`
	MonthlyLoadAnaerobicTargetMax   int     `json:"monthlyLoadAnaerobicTargetMax"`
	TrainingBalanceFeedbackPhrase   string  `json:"trainingBalanceFeedbackPhrase"`
	PrimaryTrainingDevice           bool    `json:"primaryTrainingDevice"`
}

// TrainingLoadBalance represents training load balance response.
type TrainingLoadBalance struct {
	UserID                           int64                               `json:"userId"`
	MetricsTrainingLoadBalanceDTOMap map[string]*TrainingLoadBalanceData `json:"metricsTrainingLoadBalanceDTOMap"`
	RecordedDevices                  []RecordedDevice                    `json:"recordedDevices"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (t *TrainingLoadBalance) RawJSON() json.RawMessage {
	return t.raw
}

// TrainingStatusAggregated represents aggregated training status.
type TrainingStatusAggregated struct {
	UserID                        int64                    `json:"userId"`
	MostRecentVO2Max              *MaxMetEntry             `json:"mostRecentVO2Max"`
	MostRecentTrainingLoadBalance *TrainingLoadBalance     `json:"mostRecentTrainingLoadBalance"`
	MostRecentTrainingStatus      *TrainingStatusDaily     `json:"mostRecentTrainingStatus"`
	HeatAltitudeAcclimationDTO    *HeatAltitudeAcclimation `json:"heatAltitudeAcclimationDTO"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (t *TrainingStatusAggregated) RawJSON() json.RawMessage {
	return t.raw
}

// GetTrainingReadiness retrieves training readiness data for the specified date.
func (s *MetricsService) GetTrainingReadiness(ctx context.Context, date time.Time) (*TrainingReadiness, error) {
	path := "/metrics-service/metrics/trainingreadiness/" + date.Format("2006-01-02")

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

	var entries []TrainingReadinessEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}

	return &TrainingReadiness{
		Entries: entries,
		raw:     raw,
	}, nil
}

// GetEnduranceScore retrieves endurance score data for the specified date.
func (s *MetricsService) GetEnduranceScore(ctx context.Context, date time.Time) (*EnduranceScore, error) {
	path := "/metrics-service/metrics/endurancescore?calendarDate=" + date.Format("2006-01-02")

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

	var score EnduranceScore
	if err := json.Unmarshal(raw, &score); err != nil {
		return nil, err
	}
	score.raw = raw

	return &score, nil
}

// GetHillScore retrieves hill score data for the specified date.
func (s *MetricsService) GetHillScore(ctx context.Context, date time.Time) (*HillScore, error) {
	path := "/metrics-service/metrics/hillscore?calendarDate=" + date.Format("2006-01-02")

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

	var score HillScore
	if err := json.Unmarshal(raw, &score); err != nil {
		return nil, err
	}
	score.raw = raw

	return &score, nil
}

// GetMaxMetLatest retrieves the latest VO2 max / MET data.
func (s *MetricsService) GetMaxMetLatest(ctx context.Context, date time.Time) (*MaxMetLatest, error) {
	path := "/metrics-service/metrics/maxmet/latest/" + date.Format("2006-01-02")

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

	var maxMet MaxMetLatest
	if err := json.Unmarshal(raw, &maxMet); err != nil {
		return nil, err
	}
	maxMet.raw = raw

	return &maxMet, nil
}

// GetMaxMetDaily retrieves VO2 max / MET data for a date range.
func (s *MetricsService) GetMaxMetDaily(ctx context.Context, startDate, endDate time.Time) (*MaxMetDaily, error) {
	path := fmt.Sprintf("/metrics-service/metrics/maxmet/daily/%s/%s",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

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

	var entries []MaxMetEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}

	return &MaxMetDaily{
		Entries: entries,
		raw:     raw,
	}, nil
}

// GetTrainingStatusAggregated retrieves aggregated training status data.
func (s *MetricsService) GetTrainingStatusAggregated(ctx context.Context, date time.Time) (*TrainingStatusAggregated, error) {
	path := "/metrics-service/metrics/trainingstatus/aggregated/" + date.Format("2006-01-02")

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

	var status TrainingStatusAggregated
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, err
	}
	status.raw = raw

	return &status, nil
}

// GetTrainingStatusDaily retrieves daily training status data.
func (s *MetricsService) GetTrainingStatusDaily(ctx context.Context, date time.Time) (*TrainingStatusDaily, error) {
	path := "/metrics-service/metrics/trainingstatus/daily/" + date.Format("2006-01-02")

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

	var status TrainingStatusDaily
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, err
	}
	status.raw = raw

	return &status, nil
}

// GetTrainingLoadBalance retrieves training load balance data.
func (s *MetricsService) GetTrainingLoadBalance(ctx context.Context, date time.Time) (*TrainingLoadBalance, error) {
	path := "/metrics-service/metrics/trainingloadbalance/latest/" + date.Format("2006-01-02")

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

	var balance TrainingLoadBalance
	if err := json.Unmarshal(raw, &balance); err != nil {
		return nil, err
	}
	balance.raw = raw

	return &balance, nil
}

// GetHeatAltitudeAcclimation retrieves heat and altitude acclimation data.
func (s *MetricsService) GetHeatAltitudeAcclimation(ctx context.Context, date time.Time) (*HeatAltitudeAcclimation, error) {
	path := "/metrics-service/metrics/heataltitudeacclimation/latest/" + date.Format("2006-01-02")

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

	var acclimation HeatAltitudeAcclimation
	if err := json.Unmarshal(raw, &acclimation); err != nil {
		return nil, err
	}
	acclimation.raw = raw

	return &acclimation, nil
}
