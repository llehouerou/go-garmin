// service_biometric.go
package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LactateThresholdEntry represents a single lactate threshold measurement.
// The API returns separate entries for speed and heart rate.
type LactateThresholdEntry struct {
	UserProfilePK    int64    `json:"userProfilePK"`
	Version          int64    `json:"version"`
	CalendarDate     string   `json:"calendarDate"`
	Sequence         int64    `json:"sequence"`
	Speed            *float64 `json:"speed"`    // m/s, null if this is HR entry
	HeartRate        *int     `json:"hearRate"` // Note: Garmin API has typo "hearRate"
	HeartRateCycling *int     `json:"heartRateCycling"`
}

// LactateThreshold represents the latest lactate threshold data.
type LactateThreshold struct {
	Entries []LactateThresholdEntry
	raw     json.RawMessage
}

// RawJSON returns the original JSON response.
func (lt *LactateThreshold) RawJSON() json.RawMessage {
	return lt.raw
}

// SetRaw sets the raw JSON response.
func (lt *LactateThreshold) SetRaw(data json.RawMessage) {
	lt.raw = data
}

// UnmarshalJSON unmarshals the array response into the Entries field.
func (lt *LactateThreshold) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &lt.Entries)
}

// Speed returns the lactate threshold speed in m/s, or nil if not available.
func (lt *LactateThreshold) Speed() *float64 {
	for _, e := range lt.Entries {
		if e.Speed != nil {
			return e.Speed
		}
	}
	return nil
}

// HeartRate returns the lactate threshold heart rate in bpm, or nil if not available.
func (lt *LactateThreshold) HeartRate() *int {
	for _, e := range lt.Entries {
		if e.HeartRate != nil {
			return e.HeartRate
		}
	}
	return nil
}

// FunctionalThresholdPower represents cycling FTP data.
type FunctionalThresholdPower struct {
	UserProfilePK            int64   `json:"userProfilePK"`
	Version                  *int64  `json:"version"`
	CalendarDate             *string `json:"calendarDate"`
	IsStale                  *bool   `json:"isStale"`
	Sequence                 *int64  `json:"sequence"`
	Sport                    *string `json:"sport"`
	FunctionalThresholdPower *int    `json:"functionalThresholdPower"` // Watts
	BiometricSourceType      *string `json:"biometricSourceType"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (ftp *FunctionalThresholdPower) RawJSON() json.RawMessage {
	return ftp.raw
}

// SetRaw sets the raw JSON response.
func (ftp *FunctionalThresholdPower) SetRaw(data json.RawMessage) {
	ftp.raw = data
}

// PowerToWeight represents power-to-weight ratio data.
type PowerToWeight struct {
	UserProfilePK            int64   `json:"userProfilePk"`
	CalendarDate             string  `json:"calendarDate"`
	Origin                   string  `json:"origin"`
	Sport                    string  `json:"sport"`
	FunctionalThresholdPower int     `json:"functionalThresholdPower"` // Watts
	Weight                   float64 `json:"weight"`                   // kg
	PowerToWeightRatio       float64 `json:"powerToWeight"`            // W/kg
	FTPCreateTime            string  `json:"ftpCreateTime"`
	WeightCreateTime         string  `json:"weightCreateTime"`
	IsStale                  bool    `json:"isStale"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (ptw *PowerToWeight) RawJSON() json.RawMessage {
	return ptw.raw
}

// SetRaw sets the raw JSON response.
func (ptw *PowerToWeight) SetRaw(data json.RawMessage) {
	ptw.raw = data
}

// HeartRateZone represents heart rate zone configuration for a sport.
type HeartRateZone struct {
	TrainingMethod                string `json:"trainingMethod"` // HR_RESERVE, HR_MAX, etc.
	RestingHeartRateUsed          int    `json:"restingHeartRateUsed"`
	LactateThresholdHeartRateUsed int    `json:"lactateThresholdHeartRateUsed"`
	Zone1Floor                    int    `json:"zone1Floor"`
	Zone2Floor                    int    `json:"zone2Floor"`
	Zone3Floor                    int    `json:"zone3Floor"`
	Zone4Floor                    int    `json:"zone4Floor"`
	Zone5Floor                    int    `json:"zone5Floor"`
	MaxHeartRateUsed              int    `json:"maxHeartRateUsed"`
	RestingHrAutoUpdateUsed       bool   `json:"restingHrAutoUpdateUsed"`
	Sport                         string `json:"sport"` // DEFAULT, RUNNING, CYCLING, etc.
	ChangeState                   string `json:"changeState"`
}

// HeartRateZones represents a collection of heart rate zone configurations.
type HeartRateZones struct {
	Zones []HeartRateZone
	raw   json.RawMessage
}

// RawJSON returns the original JSON response.
func (hrz *HeartRateZones) RawJSON() json.RawMessage {
	return hrz.raw
}

// SetRaw sets the raw JSON response.
func (hrz *HeartRateZones) SetRaw(data json.RawMessage) {
	hrz.raw = data
}

// UnmarshalJSON unmarshals the array response into the Zones field.
func (hrz *HeartRateZones) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &hrz.Zones)
}

// BiometricStat represents a single biometric statistic entry.
type BiometricStat struct {
	From        string  `json:"from"`
	Until       string  `json:"until"`
	Series      string  `json:"series"`
	Value       float64 `json:"value"`
	UpdatedDate string  `json:"updatedDate"`
}

// BiometricStats represents a collection of biometric statistics.
type BiometricStats struct {
	Stats []BiometricStat
	raw   json.RawMessage
}

// RawJSON returns the original JSON response.
func (bs *BiometricStats) RawJSON() json.RawMessage {
	return bs.raw
}

// SetRaw sets the raw JSON response.
func (bs *BiometricStats) SetRaw(data json.RawMessage) {
	bs.raw = data
}

// UnmarshalJSON unmarshals the array response into the Stats field.
func (bs *BiometricStats) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &bs.Stats)
}

// GetLatestLactateThreshold retrieves the latest lactate threshold data.
func (s *BiometricService) GetLatestLactateThreshold(ctx context.Context) (*LactateThreshold, error) {
	return fetch[LactateThreshold](ctx, s.client, "/biometric-service/biometric/latestLactateThreshold")
}

// GetCyclingFTP retrieves the latest cycling Functional Threshold Power.
func (s *BiometricService) GetCyclingFTP(ctx context.Context) (*FunctionalThresholdPower, error) {
	return fetch[FunctionalThresholdPower](ctx, s.client, "/biometric-service/biometric/latestFunctionalThresholdPower/CYCLING")
}

// GetPowerToWeight retrieves the power-to-weight ratio for running on the given date.
func (s *BiometricService) GetPowerToWeight(ctx context.Context, date time.Time) (*PowerToWeight, error) {
	path := fmt.Sprintf("/biometric-service/biometric/powerToWeight/latest/%s?sport=Running", date.Format("2006-01-02"))
	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// API returns an array with one element
	var results []PowerToWeight
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	results[0].raw = body
	return &results[0], nil
}

// GetHeartRateZones retrieves heart rate zone configurations for all sports.
func (s *BiometricService) GetHeartRateZones(ctx context.Context) (*HeartRateZones, error) {
	return fetch[HeartRateZones](ctx, s.client, "/biometric-service/heartRateZones/")
}

// GetLactateThresholdSpeedRange retrieves lactate threshold speed stats for a date range.
func (s *BiometricService) GetLactateThresholdSpeedRange(ctx context.Context, start, end time.Time) (*BiometricStats, error) {
	path := fmt.Sprintf("/biometric-service/stats/lactateThresholdSpeed/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		start.Format("2006-01-02"), end.Format("2006-01-02"))
	return s.getBiometricStats(ctx, path)
}

// GetLactateThresholdHRRange retrieves lactate threshold heart rate stats for a date range.
func (s *BiometricService) GetLactateThresholdHRRange(ctx context.Context, start, end time.Time) (*BiometricStats, error) {
	path := fmt.Sprintf("/biometric-service/stats/lactateThresholdHeartRate/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		start.Format("2006-01-02"), end.Format("2006-01-02"))
	return s.getBiometricStats(ctx, path)
}

// GetFTPRange retrieves Functional Threshold Power stats for a date range.
func (s *BiometricService) GetFTPRange(ctx context.Context, start, end time.Time) (*BiometricStats, error) {
	path := fmt.Sprintf("/biometric-service/stats/functionalThresholdPower/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		start.Format("2006-01-02"), end.Format("2006-01-02"))
	return s.getBiometricStats(ctx, path)
}

func (s *BiometricService) getBiometricStats(ctx context.Context, path string) (*BiometricStats, error) {
	return fetch[BiometricStats](ctx, s.client, path)
}
