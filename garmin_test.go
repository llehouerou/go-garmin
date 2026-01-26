package garmin

import (
	"bytes"
	"testing"
)

func TestClientCreation(t *testing.T) {
	client := New(Options{})
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Wellness == nil {
		t.Error("expected Wellness service to be initialized")
	}
	if client.Activities == nil {
		t.Error("expected Activities service to be initialized")
	}
}

func TestClientSessionPersistence(t *testing.T) {
	client := New(Options{})
	client.auth = &authState{
		OAuth1Token:       "test-token",
		OAuth1Secret:      "test-secret",
		OAuth2AccessToken: "test-access",
		Domain:            "garmin.com",
	}

	var buf bytes.Buffer
	if err := client.SaveSession(&buf); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	client2 := New(Options{})
	if err := client2.LoadSession(&buf); err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	if client2.auth.OAuth1Token != "test-token" {
		t.Errorf("token mismatch: got %s", client2.auth.OAuth1Token)
	}
}

func TestClientDefaultDomain(t *testing.T) {
	client := New(Options{})
	if client.opts.Domain != "garmin.com" {
		t.Errorf("expected default domain 'garmin.com', got %s", client.opts.Domain)
	}
}

func TestClientCustomDomain(t *testing.T) {
	client := New(Options{Domain: "garmin.cn"})
	if client.opts.Domain != "garmin.cn" {
		t.Errorf("expected domain 'garmin.cn', got %s", client.opts.Domain)
	}
}

func TestClientCustomRateLimit(t *testing.T) {
	customConfig := &RateLimitConfig{
		RequestsPerMinute: 30,
		BurstSize:         10,
	}
	client := New(Options{RateLimit: customConfig})
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	// Verify the client was created with custom rate limit config
	if client.transport == nil {
		t.Error("expected transport to be initialized")
	}
}

func TestAllServicesInitialized(t *testing.T) {
	client := New(Options{})

	services := []struct {
		name    string
		service any
	}{
		{"Wellness", client.Wellness},
		{"Activities", client.Activities},
		{"Metrics", client.Metrics},
		{"Weight", client.Weight},
		{"Devices", client.Devices},
		{"Workouts", client.Workouts},
		{"Goals", client.Goals},
		{"Badges", client.Badges},
		{"Gear", client.Gear},
		{"Download", client.Download},
		{"Upload", client.Upload},
		{"Hydration", client.Hydration},
		{"BloodPressure", client.BloodPressure},
		{"PersonalRecords", client.PersonalRecords},
		{"Steps", client.Steps},
		{"UserProfile", client.UserProfile},
	}

	for _, s := range services {
		if s.service == nil {
			t.Errorf("expected %s service to be initialized", s.name)
		}
	}
}
