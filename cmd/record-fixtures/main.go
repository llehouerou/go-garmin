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
	"flag"
	"fmt"
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
