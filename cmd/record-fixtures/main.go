// Command record-fixtures records API interactions for testing.
//
// Usage:
//
//	record-fixtures -email=user@example.com -password=secret
//
// This command authenticates with Garmin Connect and records API
// responses to cassette files in testdata/cassettes/.
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
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	garmin "github.com/llehouerou/go-garmin"
	"github.com/llehouerou/go-garmin/testutil"
)

func main() {
	email := flag.String("email", "", "Garmin Connect email")
	password := flag.String("password", "", "Garmin Connect password")
	date := flag.String("date", "", "Date to record (YYYY-MM-DD, defaults to today)")
	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Usage: record-fixtures -email=EMAIL -password=PASSWORD [-date=YYYY-MM-DD]")
		os.Exit(1)
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

	if err := recordFixtures(*email, *password, targetDate); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done! Cassettes recorded to testdata/cassettes/")
}

func recordFixtures(email, password string, date time.Time) error {
	ctx := context.Background()

	// Step 1: Login once and record auth flow
	fmt.Println("Recording authentication...")
	session, err := recordAuth(ctx, email, password)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	// Step 2: Record API calls using the saved session
	fmt.Println("Recording API calls...")

	if err := recordSleep(ctx, session, date); err != nil {
		return fmt.Errorf("sleep: %w", err)
	}

	if err := recordStress(ctx, session, date); err != nil {
		return fmt.Errorf("stress: %w", err)
	}

	if err := recordBodyBattery(ctx, session, date); err != nil {
		return fmt.Errorf("body_battery: %w", err)
	}

	if err := recordActivities(ctx, session); err != nil {
		return fmt.Errorf("activities: %w", err)
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

func recordActivities(ctx context.Context, session []byte) error {
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
	}

	// If we got activities, record details for the first one
	if len(activities) > 0 {
		if activityID, ok := activities[0]["activityId"].(float64); ok {
			fmt.Printf("  Getting activity details for %d...\n", int64(activityID))
			activityURL := fmt.Sprintf("https://connectapi.%s/activity-service/activity/%d", authState.Domain, int64(activityID))
			_, err := doAPIRequest(ctx, httpClient, activityURL, authState.OAuth2AccessToken)
			if err != nil {
				fmt.Printf("  Warning: activity details: %v\n", err)
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
