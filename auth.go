package garmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
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

const oauthConsumerURL = "https://thegarth.s3.amazonaws.com/oauth_consumer.json"

type oauthConsumer struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

var (
	cachedConsumer *oauthConsumer
	consumerMu     sync.Mutex
)

func fetchOAuthConsumer(ctx context.Context, client *http.Client) (*oauthConsumer, error) {
	consumerMu.Lock()
	defer consumerMu.Unlock()

	if cachedConsumer != nil {
		return cachedConsumer, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, oauthConsumerURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var consumer oauthConsumer
	if err := json.NewDecoder(resp.Body).Decode(&consumer); err != nil {
		return nil, err
	}

	cachedConsumer = &consumer
	return cachedConsumer, nil
}
