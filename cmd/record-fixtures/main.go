// Command record-fixtures records API interactions for testing.
//
// Usage:
//
//	record-fixtures -email=user@example.com -password=secret
//	record-fixtures -email=user@example.com -password=secret -cassette=activities
//
// This command authenticates with Garmin Connect and records API
// responses to cassette files in testdata/cassettes/.
//
// Available cassettes:
//   - sleep_daily
//   - wellness_stress
//   - wellness_body_battery
//   - wellness_heart_rate
//   - wellness_extended
//   - activities
//   - hrv
//   - weight
//   - metrics
//   - userprofile
//   - devices
//   - biometric
//   - workouts
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	garmin "github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/testutil"
)

func main() {
	email := flag.String("email", "", "Garmin Connect email")
	password := flag.String("password", "", "Garmin Connect password")
	date := flag.String("date", "", "Date to record (YYYY-MM-DD, defaults to today)")
	cassette := flag.String("cassette", "", "Record only this cassette (defaults to all)")
	listCassettes := flag.Bool("list", false, "List available cassettes and exit")
	flag.Parse()

	if *listCassettes {
		fmt.Println("Available cassettes:")
		names := getCassetteNames()
		for _, name := range names {
			fmt.Printf("  %s\n", name)
		}
		os.Exit(0)
	}

	if *email == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Usage: record-fixtures -email=EMAIL -password=PASSWORD [-date=YYYY-MM-DD] [-cassette=NAME]")
		fmt.Fprintln(os.Stderr, "       record-fixtures -list")
		os.Exit(1)
	}

	if *cassette != "" {
		if !isValidCassette(*cassette) {
			fmt.Fprintf(os.Stderr, "Unknown cassette: %s\n", *cassette)
			fmt.Fprintln(os.Stderr, "Use -list to see available cassettes")
			os.Exit(1)
		}
	}

	targetDate := time.Now()
	if *date != "" {
		var err error
		targetDate, err = time.Parse("2006-01-02", *date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid date format: %s\n", *date)
			os.Exit(1)
		}
	}

	if err := recordFixtures(*email, *password, targetDate, *cassette); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *cassette != "" {
		fmt.Printf("Done! Cassette '%s' recorded to testdata/cassettes/\n", *cassette)
	} else {
		fmt.Println("Done! All cassettes recorded to testdata/cassettes/")
	}
}

// cassetteRecorder defines a function that records a cassette.
type cassetteRecorder func(ctx context.Context, session []byte, date time.Time) error

// getCassetteRecorders returns a map of cassette names to their recorder functions.
func getCassetteRecorders() map[string]cassetteRecorder {
	return map[string]cassetteRecorder{
		"sleep_daily":           recordSleep,
		"wellness_stress":       recordStress,
		"wellness_body_battery": recordBodyBattery,
		"wellness_heart_rate":   recordHeartRate,
		"wellness_extended":     recordWellnessExtended,
		"activities":            recordActivities,
		"hrv":                   recordHRV,
		"weight":                recordWeight,
		"metrics":               recordMetrics,
		"userprofile":           recordUserProfile,
		"devices":               recordDevices,
		"biometric":             recordBiometric,
		"workouts":              recordWorkouts,
	}
}

// getCassetteNames returns a sorted list of available cassette names.
func getCassetteNames() []string {
	recorders := getCassetteRecorders()
	names := make([]string, 0, len(recorders))
	for name := range recorders {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// isValidCassette checks if the given name is a valid cassette.
func isValidCassette(name string) bool {
	recorders := getCassetteRecorders()
	_, ok := recorders[name]
	return ok
}

func recordFixtures(email, password string, date time.Time, cassette string) error {
	ctx := context.Background()

	// Step 1: Login once and record auth flow
	fmt.Println("Recording authentication...")
	session, err := recordAuth(ctx, email, password)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	// Step 2: Record API calls using the saved session
	recorders := getCassetteRecorders()

	// If a specific cassette is requested, only record that one
	if cassette != "" {
		fmt.Printf("Recording cassette '%s'...\n", cassette)
		recordFn := recorders[cassette]
		if err := recordFn(ctx, session, date); err != nil {
			return fmt.Errorf("%s: %w", cassette, err)
		}
		return nil
	}

	// Record all cassettes
	fmt.Println("Recording all API calls...")
	names := getCassetteNames()
	for _, name := range names {
		fmt.Printf("Recording %s...\n", name)
		recordFn := recorders[name]
		if err := recordFn(ctx, session, date); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}

	return nil
}

// stopRecorder stops the recorder and returns any error.
func stopRecorder(rec *recorder.Recorder) error {
	return rec.Stop()
}

// recordAuth logs in and records the auth flow, returning the session data.
func recordAuth(ctx context.Context, email, password string) ([]byte, error) {
	rec, err := testutil.NewRecordingRecorder("auth")
	if err != nil {
		return nil, err
	}
	defer func() { _ = stopRecorder(rec) }()

	client := garmin.New(garmin.Options{
		HTTPClient: testutil.HTTPClientWithRecorder(rec),
	})

	if err := client.Login(ctx, email, password); err != nil {
		return nil, err
	}

	// Save session for reuse
	var buf bytes.Buffer
	if err := client.SaveSession(&buf); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return buf.Bytes(), nil
}

// loadSession creates a client with the recorded session loaded.
func loadSession(rec *recorder.Recorder, session []byte) (*garmin.Client, error) {
	client := garmin.New(garmin.Options{
		HTTPClient: testutil.HTTPClientWithRecorder(rec),
	})

	if err := client.LoadSession(bytes.NewReader(session)); err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}

	return client, nil
}

func recordSleep(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("sleep_daily")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	client, err := loadSession(rec, session)
	if err != nil {
		return err
	}

	fmt.Printf("  Getting sleep data for %s...\n", date.Format("2006-01-02"))
	_, err = client.Sleep.GetDaily(ctx, date)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordStress(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("wellness_stress")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	client, err := loadSession(rec, session)
	if err != nil {
		return err
	}

	fmt.Printf("  Getting stress data for %s...\n", date.Format("2006-01-02"))
	_, err = client.Wellness.GetDailyStress(ctx, date)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordBodyBattery(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("wellness_body_battery")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	client, err := loadSession(rec, session)
	if err != nil {
		return err
	}

	fmt.Printf("  Getting body battery data for %s...\n", date.Format("2006-01-02"))
	_, err = client.Wellness.GetBodyBatteryEvents(ctx, date)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordHeartRate(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("wellness_heart_rate")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	fmt.Printf("  Getting heart rate data for %s...\n", date.Format("2006-01-02"))
	url := fmt.Sprintf("https://connectapi.%s/wellness-service/wellness/dailyHeartRate/?date=%s",
		authState.Domain, date.Format("2006-01-02"))
	_, err = doAPIRequest(ctx, httpClient, url, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordHRV(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("hrv")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// Record daily HRV
	fmt.Printf("  Getting HRV data for %s...\n", date.Format("2006-01-02"))
	dailyURL := fmt.Sprintf("https://connectapi.%s/hrv-service/hrv/%s",
		authState.Domain, date.Format("2006-01-02"))
	_, err = doAPIRequest(ctx, httpClient, dailyURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	// Record HRV range (last 7 days)
	startDate := date.AddDate(0, 0, -7)
	fmt.Printf("  Getting HRV range from %s to %s...\n", startDate.Format("2006-01-02"), date.Format("2006-01-02"))
	rangeURL := fmt.Sprintf("https://connectapi.%s/hrv-service/hrv/daily/%s/%s",
		authState.Domain, startDate.Format("2006-01-02"), date.Format("2006-01-02"))
	_, err = doAPIRequest(ctx, httpClient, rangeURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordWeight(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("weight")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// Record daily weight
	fmt.Printf("  Getting weight data for %s...\n", date.Format("2006-01-02"))
	dailyURL := fmt.Sprintf("https://connectapi.%s/weight-service/weight/dayview/%s",
		authState.Domain, date.Format("2006-01-02"))
	_, err = doAPIRequest(ctx, httpClient, dailyURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	// Record weight range (last 30 days)
	startDate := date.AddDate(0, 0, -30)
	fmt.Printf("  Getting weight range from %s to %s...\n", startDate.Format("2006-01-02"), date.Format("2006-01-02"))
	rangeURL := fmt.Sprintf("https://connectapi.%s/weight-service/weight/range/%s/%s?includeAll=true",
		authState.Domain, startDate.Format("2006-01-02"), date.Format("2006-01-02"))
	_, err = doAPIRequest(ctx, httpClient, rangeURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: %v\n", err)
	}

	return nil
}

func recordActivities(ctx context.Context, session []byte, _ time.Time) error {
	rec, err := testutil.NewRecordingRecorder("activities")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// Record activities list (get last 5 activities)
	fmt.Println("  Getting activities list...")
	activitiesURL := fmt.Sprintf("https://connectapi.%s/activitylist-service/activities/search/activities?start=0&limit=5", authState.Domain)
	activities, err := doAPIRequest(ctx, httpClient, activitiesURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: activities list: %v\n", err)
		return nil
	}

	if len(activities) == 0 {
		return nil
	}

	activityID, ok := activities[0]["activityId"].(float64)
	if !ok {
		return nil
	}

	// Record details, splits, and weather for the first activity
	recordActivityDetails(ctx, httpClient, authState.Domain, authState.OAuth2AccessToken, int64(activityID))

	return nil
}

func recordActivityDetails(ctx context.Context, client *http.Client, domain, token string, id int64) {
	fmt.Printf("  Getting activity details for %d...\n", id)
	activityURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d", domain, id)
	_, err := doAPIRequest(ctx, client, activityURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity details: %v\n", err)
	}

	fmt.Printf("  Getting activity splits for %d...\n", id)
	splitsURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/splits", domain, id)
	_, err = doAPIRequest(ctx, client, splitsURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity splits: %v\n", err)
	}

	fmt.Printf("  Getting activity weather for %d...\n", id)
	weatherURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/weather", domain, id)
	_, err = doAPIRequest(ctx, client, weatherURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity weather: %v\n", err)
	}

	// Activity extension endpoints
	fmt.Printf("  Getting activity extended details for %d...\n", id)
	detailsURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/details", domain, id)
	_, err = doAPIRequest(ctx, client, detailsURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity extended details: %v\n", err)
	}

	fmt.Printf("  Getting activity HR time in zones for %d...\n", id)
	hrZonesURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/hrTimeInZones", domain, id)
	_, err = doAPIRequest(ctx, client, hrZonesURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity HR zones: %v\n", err)
	}

	fmt.Printf("  Getting activity power time in zones for %d...\n", id)
	powerZonesURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/powerTimeInZones", domain, id)
	_, err = doAPIRequest(ctx, client, powerZonesURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity power zones: %v\n", err)
	}

	fmt.Printf("  Getting activity exercise sets for %d...\n", id)
	exerciseSetsURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d/exerciseSets", domain, id)
	_, err = doAPIRequest(ctx, client, exerciseSetsURL, token)
	if err != nil {
		fmt.Printf("  Warning: activity exercise sets: %v\n", err)
	}
}

func recordMetrics(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("metrics")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)
	dateStr := date.Format("2006-01-02")

	// Training readiness
	fmt.Printf("  Getting training readiness for %s...\n", dateStr)
	trainingReadinessURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/trainingreadiness/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, trainingReadinessURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: training readiness: %v\n", err)
	}

	// Endurance score
	fmt.Printf("  Getting endurance score for %s...\n", dateStr)
	enduranceURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/endurancescore?calendarDate=%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, enduranceURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: endurance score: %v\n", err)
	}

	// Endurance score stats (weekly aggregation, last ~12 weeks)
	statsStartDate := date.AddDate(0, 0, -84)
	fmt.Printf("  Getting endurance score stats from %s to %s...\n", statsStartDate.Format("2006-01-02"), dateStr)
	enduranceStatsURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/endurancescore/stats?startDate=%s&endDate=%s&aggregation=weekly",
		authState.Domain, statsStartDate.Format("2006-01-02"), dateStr)
	_, err = doAPIRequest(ctx, httpClient, enduranceStatsURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: endurance score stats: %v\n", err)
	}

	// Hill score
	fmt.Printf("  Getting hill score for %s...\n", dateStr)
	hillURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/hillscore?calendarDate=%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, hillURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: hill score: %v\n", err)
	}

	// Race predictions - skipped, requires display name from user profile
	// URL: /metrics-service/metrics/racepredictions/latest/{displayName}

	// VO2 max / MET - latest
	fmt.Printf("  Getting latest VO2 max for %s...\n", dateStr)
	maxMetLatestURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/maxmet/latest/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, maxMetLatestURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: maxmet latest: %v\n", err)
	}

	// VO2 max / MET - daily range (last 30 days)
	startDate := date.AddDate(0, 0, -30)
	fmt.Printf("  Getting VO2 max range from %s to %s...\n", startDate.Format("2006-01-02"), dateStr)
	maxMetDailyURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/maxmet/daily/%s/%s",
		authState.Domain, startDate.Format("2006-01-02"), dateStr)
	_, err = doAPIRequest(ctx, httpClient, maxMetDailyURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: maxmet daily: %v\n", err)
	}

	// Training status - aggregated (requires date in path)
	fmt.Printf("  Getting aggregated training status for %s...\n", dateStr)
	trainingStatusAggURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/trainingstatus/aggregated/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, trainingStatusAggURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: training status aggregated: %v\n", err)
	}

	// Training status - daily
	fmt.Printf("  Getting daily training status for %s...\n", dateStr)
	trainingStatusDailyURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/trainingstatus/daily/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, trainingStatusDailyURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: training status daily: %v\n", err)
	}

	// Training load balance
	fmt.Printf("  Getting training load balance for %s...\n", dateStr)
	loadBalanceURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/trainingloadbalance/latest/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, loadBalanceURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: training load balance: %v\n", err)
	}

	// Heat/altitude acclimation
	fmt.Printf("  Getting heat/altitude acclimation for %s...\n", dateStr)
	acclimationURL := fmt.Sprintf("https://connectapi.%s/metrics-service/metrics/heataltitudeacclimation/latest/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, acclimationURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: heat/altitude acclimation: %v\n", err)
	}

	return nil
}

func recordUserProfile(ctx context.Context, session []byte, _ time.Time) error {
	rec, err := testutil.NewRecordingRecorder("userprofile")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// Social profile
	fmt.Println("  Getting social profile...")
	socialProfileURL := fmt.Sprintf("https://connectapi.%s/userprofile-service/socialProfile",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, socialProfileURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: social profile: %v\n", err)
	}

	// User settings
	fmt.Println("  Getting user settings...")
	userSettingsURL := fmt.Sprintf("https://connectapi.%s/userprofile-service/userprofile/user-settings",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, userSettingsURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: user settings: %v\n", err)
	}

	// Profile settings
	fmt.Println("  Getting profile settings...")
	profileSettingsURL := fmt.Sprintf("https://connectapi.%s/userprofile-service/userprofile/settings",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, profileSettingsURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: profile settings: %v\n", err)
	}

	return nil
}

func recordDevices(ctx context.Context, session []byte, _ time.Time) error {
	rec, err := testutil.NewRecordingRecorder("devices")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// List devices
	fmt.Println("  Getting device list...")
	devicesURL := fmt.Sprintf("https://connectapi.%s/device-service/deviceregistration/devices",
		authState.Domain)
	devices, err := doAPIRequest(ctx, httpClient, devicesURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: device list: %v\n", err)
	}

	// Get settings for first device if available
	if len(devices) > 0 {
		if deviceID, ok := devices[0]["deviceId"].(float64); ok {
			fmt.Printf("  Getting device settings for %d...\n", int64(deviceID))
			settingsURL := fmt.Sprintf("https://connectapi.%s/device-service/deviceservice/device-info/settings/%d",
				authState.Domain, int64(deviceID))
			_, err = doAPIRequest(ctx, httpClient, settingsURL, authState.OAuth2AccessToken)
			if err != nil {
				fmt.Printf("  Warning: device settings: %v\n", err)
			}
		}
	}

	// Device messages
	fmt.Println("  Getting device messages...")
	messagesURL := fmt.Sprintf("https://connectapi.%s/device-service/devicemessage/messages",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, messagesURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: device messages: %v\n", err)
	}

	// Primary training device
	fmt.Println("  Getting primary training device...")
	primaryURL := fmt.Sprintf("https://connectapi.%s/web-gateway/device-info/primary-training-device",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, primaryURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: primary training device: %v\n", err)
	}

	return nil
}

func recordWellnessExtended(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("wellness_extended")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)
	dateStr := date.Format("2006-01-02")

	// SpO2 (blood oxygen)
	fmt.Printf("  Getting SpO2 data for %s...\n", dateStr)
	spo2URL := fmt.Sprintf("https://connectapi.%s/wellness-service/wellness/daily/spo2/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, spo2URL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: SpO2: %v\n", err)
	}

	// Respiration
	fmt.Printf("  Getting respiration data for %s...\n", dateStr)
	respirationURL := fmt.Sprintf("https://connectapi.%s/wellness-service/wellness/daily/respiration/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, respirationURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: respiration: %v\n", err)
	}

	// Intensity minutes
	fmt.Printf("  Getting intensity minutes for %s...\n", dateStr)
	imURL := fmt.Sprintf("https://connectapi.%s/wellness-service/wellness/daily/im/%s",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, imURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: intensity minutes: %v\n", err)
	}

	return nil
}

func recordBiometric(ctx context.Context, session []byte, date time.Time) error {
	rec, err := testutil.NewRecordingRecorder("biometric")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)
	dateStr := date.Format("2006-01-02")
	startDate := date.AddDate(0, 0, -30)
	startDateStr := startDate.Format("2006-01-02")

	// Latest Lactate Threshold
	fmt.Println("  Getting latest lactate threshold...")
	lactateURL := fmt.Sprintf("https://connectapi.%s/biometric-service/biometric/latestLactateThreshold",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, lactateURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: lactate threshold: %v\n", err)
	}

	// Latest Cycling FTP
	fmt.Println("  Getting latest cycling FTP...")
	ftpURL := fmt.Sprintf("https://connectapi.%s/biometric-service/biometric/latestFunctionalThresholdPower/CYCLING",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, ftpURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: cycling FTP: %v\n", err)
	}

	// Power to Weight (Running)
	fmt.Printf("  Getting power to weight for %s...\n", dateStr)
	powerToWeightURL := fmt.Sprintf("https://connectapi.%s/biometric-service/biometric/powerToWeight/latest/%s?sport=Running",
		authState.Domain, dateStr)
	_, err = doAPIRequest(ctx, httpClient, powerToWeightURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: power to weight: %v\n", err)
	}

	// Lactate Threshold Speed Range
	fmt.Printf("  Getting lactate threshold speed from %s to %s...\n", startDateStr, dateStr)
	ltSpeedURL := fmt.Sprintf("https://connectapi.%s/biometric-service/stats/lactateThresholdSpeed/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		authState.Domain, startDateStr, dateStr)
	_, err = doAPIRequest(ctx, httpClient, ltSpeedURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: lactate threshold speed: %v\n", err)
	}

	// Lactate Threshold Heart Rate Range
	fmt.Printf("  Getting lactate threshold heart rate from %s to %s...\n", startDateStr, dateStr)
	ltHrURL := fmt.Sprintf("https://connectapi.%s/biometric-service/stats/lactateThresholdHeartRate/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		authState.Domain, startDateStr, dateStr)
	_, err = doAPIRequest(ctx, httpClient, ltHrURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: lactate threshold heart rate: %v\n", err)
	}

	// Functional Threshold Power Range (Running)
	fmt.Printf("  Getting FTP range from %s to %s...\n", startDateStr, dateStr)
	ftpRangeURL := fmt.Sprintf("https://connectapi.%s/biometric-service/stats/functionalThresholdPower/range/%s/%s?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST",
		authState.Domain, startDateStr, dateStr)
	_, err = doAPIRequest(ctx, httpClient, ftpRangeURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: FTP range: %v\n", err)
	}

	// Heart Rate Zones
	fmt.Println("  Getting heart rate zones...")
	hrZonesURL := fmt.Sprintf("https://connectapi.%s/biometric-service/heartRateZones/",
		authState.Domain)
	_, err = doAPIRequest(ctx, httpClient, hrZonesURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: heart rate zones: %v\n", err)
	}

	return nil
}

func recordWorkouts(ctx context.Context, session []byte, _ time.Time) error {
	rec, err := testutil.NewRecordingRecorder("workouts")
	if err != nil {
		return err
	}
	defer func() { _ = stopRecorder(rec) }()

	// Parse session to get OAuth2 token
	var authState struct {
		OAuth2AccessToken string `json:"oauth2_access_token"`
		Domain            string `json:"domain"`
	}
	if err := json.Unmarshal(session, &authState); err != nil {
		return fmt.Errorf("failed to parse session: %w", err)
	}

	httpClient := testutil.HTTPClientWithRecorder(rec)

	// List workouts
	fmt.Println("  Getting workouts list...")
	listURL := fmt.Sprintf("https://connectapi.%s/workout-service/workouts?start=0&limit=10",
		authState.Domain)
	workouts, err := doAPIRequest(ctx, httpClient, listURL, authState.OAuth2AccessToken)
	if err != nil {
		fmt.Printf("  Warning: workouts list: %v\n", err)
		return nil
	}

	// Get first workout details if any exist
	if len(workouts) > 0 {
		if workoutID, ok := workouts[0]["workoutId"].(float64); ok {
			fmt.Printf("  Getting workout %d details...\n", int64(workoutID))
			detailURL := fmt.Sprintf("https://connectapi.%s/workout-service/workout/%d",
				authState.Domain, int64(workoutID))
			_, err = doAPIRequest(ctx, httpClient, detailURL, authState.OAuth2AccessToken)
			if err != nil {
				fmt.Printf("  Warning: workout detail: %v\n", err)
			}
		}
	}

	return nil
}

func doAPIRequest(ctx context.Context, client *http.Client, url, token string) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "GCM-iOS-5.19.1.2")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		// Try single object
		var single map[string]any
		if err := json.Unmarshal(body, &single); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return []map[string]any{single}, nil
	}

	return result, nil
}
