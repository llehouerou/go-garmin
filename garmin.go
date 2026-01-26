// Package garmin provides a Go client library for interacting with Garmin services.
package garmin

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultDomain = "garmin.com"
)

// Options configures the Garmin client.
type Options struct {
	HTTPClient *http.Client
	MFAHandler func() (string, error)
	RateLimit  *RateLimitConfig
	Domain     string // "garmin.com" or "garmin.cn"
}

// Client is the main entry point for interacting with Garmin services.
type Client struct {
	Wellness        *WellnessService
	Activities      *ActivityService
	Metrics         *MetricsService
	Weight          *WeightService
	Devices         *DeviceService
	Workouts        *WorkoutService
	Goals           *GoalService
	Badges          *BadgeService
	Gear            *GearService
	Download        *DownloadService
	Upload          *UploadService
	Hydration       *HydrationService
	BloodPressure   *BloodPressureService
	PersonalRecords *PersonalRecordsService
	Steps           *StepsService
	UserProfile     *UserProfileService

	opts      Options
	transport *httpTransport
	auth      *authState
}

// New creates a new Garmin client with the provided options.
func New(opts Options) *Client {
	if opts.Domain == "" {
		opts.Domain = defaultDomain
	}

	rlConfig := DefaultRateLimitConfig()
	if opts.RateLimit != nil {
		rlConfig = *opts.RateLimit
	}

	c := &Client{
		opts:      opts,
		transport: newHTTPTransport(opts.HTTPClient, defaultRetryConfig(), newRateLimiter(rlConfig)),
		auth:      &authState{Domain: opts.Domain},
	}

	// Initialize services
	c.Wellness = &WellnessService{client: c}
	c.Activities = &ActivityService{client: c}
	c.Metrics = &MetricsService{client: c}
	c.Weight = &WeightService{client: c}
	c.Devices = &DeviceService{client: c}
	c.Workouts = &WorkoutService{client: c}
	c.Goals = &GoalService{client: c}
	c.Badges = &BadgeService{client: c}
	c.Gear = &GearService{client: c}
	c.Download = &DownloadService{client: c}
	c.Upload = &UploadService{client: c}
	c.Hydration = &HydrationService{client: c}
	c.BloodPressure = &BloodPressureService{client: c}
	c.PersonalRecords = &PersonalRecordsService{client: c}
	c.Steps = &StepsService{client: c}
	c.UserProfile = &UserProfileService{client: c}

	return c
}

// SaveSession persists the authentication state to the provided writer.
func (c *Client) SaveSession(w io.Writer) error {
	return c.auth.save(w)
}

// LoadSession restores the authentication state from the provided reader.
func (c *Client) LoadSession(r io.Reader) error {
	return c.auth.load(r)
}

// doAPI performs an authenticated API request to Garmin Connect.
//
//nolint:unused // Will be used by service implementations
func (c *Client) doAPI(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if !c.auth.isAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	if c.auth.isExpired() {
		if err := c.refreshOAuth2(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("https://connectapi.%s%s", c.auth.Domain, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.auth.OAuth2AccessToken)
	req.Header.Set("User-Agent", "GCM-iOS-5.19.1.2")

	return c.transport.do(req)
}

// refreshOAuth2 re-exchanges OAuth1 for a fresh OAuth2 token.
//
//nolint:unused // Will be used by service implementations
func (c *Client) refreshOAuth2(ctx context.Context) error {
	sso, err := newSSOClient(c.auth.Domain, c.transport.client.Timeout)
	if err != nil {
		return err
	}

	consumer, err := fetchOAuthConsumer(ctx, c.transport.client)
	if err != nil {
		return err
	}

	oauth1 := &OAuth1Token{
		Token:    c.auth.OAuth1Token,
		Secret:   c.auth.OAuth1Secret,
		MFAToken: c.auth.MFAToken,
	}

	oauth2, err := sso.exchangeOAuth1ForOAuth2(ctx, oauth1, consumer)
	if err != nil {
		return err
	}

	c.auth.OAuth2AccessToken = oauth2.AccessToken
	c.auth.OAuth2RefreshToken = oauth2.RefreshToken
	c.auth.OAuth2Expiry = oauth2.Expiry
	c.auth.OAuth2Scope = oauth2.Scope

	return nil
}
