package garmin

import (
	"encoding/json"
	"io"
	"time"
)

type authState struct {
	OAuth1Token        string    `json:"oauth1_token"`
	OAuth1Secret       string    `json:"oauth1_secret"`
	MFAToken           string    `json:"mfa_token,omitempty"`
	OAuth2AccessToken  string    `json:"oauth2_access_token"`
	OAuth2RefreshToken string    `json:"oauth2_refresh_token"`
	OAuth2Expiry       time.Time `json:"oauth2_expiry"`
	OAuth2Scope        string    `json:"oauth2_scope,omitempty"`
	Domain             string    `json:"domain"`
}

func (a *authState) save(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(a)
}

func (a *authState) load(r io.Reader) error {
	return json.NewDecoder(r).Decode(a)
}

func (a *authState) isExpired() bool {
	// Consider expired if within 5 minutes of expiry
	return time.Now().Add(5 * time.Minute).After(a.OAuth2Expiry)
}

func (a *authState) isAuthenticated() bool {
	return a.OAuth1Token != "" && a.OAuth2AccessToken != ""
}
