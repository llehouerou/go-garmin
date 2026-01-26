// Package garmin provides a Go client library for interacting with Garmin services.
package garmin

import (
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
