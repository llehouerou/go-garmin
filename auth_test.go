package garmin

import (
	"bytes"
	"testing"
	"time"
)

func TestAuthStateSerialization(t *testing.T) {
	original := &authState{
		OAuth1Token:        "token1",
		OAuth1Secret:       "secret1",
		OAuth2AccessToken:  "access",
		OAuth2RefreshToken: "refresh",
		OAuth2Expiry:       time.Date(2026, 1, 26, 12, 0, 0, 0, time.UTC),
		Domain:             "garmin.com",
	}

	var buf bytes.Buffer
	if err := original.save(&buf); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded := &authState{}
	if err := loaded.load(&buf); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.OAuth1Token != original.OAuth1Token {
		t.Errorf("OAuth1Token mismatch: got %s, want %s", loaded.OAuth1Token, original.OAuth1Token)
	}
	if loaded.Domain != original.Domain {
		t.Errorf("Domain mismatch: got %s, want %s", loaded.Domain, original.Domain)
	}
}

func TestAuthStateExpiry(t *testing.T) {
	auth := &authState{
		OAuth2Expiry: time.Now().Add(-time.Hour),
	}
	if !auth.isExpired() {
		t.Error("expected token to be expired")
	}

	auth.OAuth2Expiry = time.Now().Add(time.Hour)
	if auth.isExpired() {
		t.Error("expected token to not be expired")
	}
}

func TestAuthStateIsAuthenticated(t *testing.T) {
	auth := &authState{}
	if auth.isAuthenticated() {
		t.Error("expected empty auth state to not be authenticated")
	}

	auth.OAuth1Token = "token1"
	if auth.isAuthenticated() {
		t.Error("expected auth state with only OAuth1Token to not be authenticated")
	}

	auth.OAuth2AccessToken = "access"
	if !auth.isAuthenticated() {
		t.Error("expected auth state with both tokens to be authenticated")
	}
}
