package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llehouerou/go-garmin"
)

func sessionPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
	}
	return filepath.Join(configDir, "garmin", "session.json")
}

func loadClient() (*garmin.Client, error) {
	path := sessionPath()
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("not logged in, run: garmin login")
	}
	defer f.Close()

	client := garmin.New(garmin.Options{})
	if err := client.LoadSession(f); err != nil {
		return nil, fmt.Errorf("session corrupted: %w", err)
	}
	return client, nil
}

func saveClient(client *garmin.Client) error {
	path := sessionPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	return client.SaveSession(f)
}

func removeSession() error {
	return os.Remove(sessionPath())
}
