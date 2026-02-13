// Package testutil provides testing utilities for the garmin package.
package testutil

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// Patterns for anonymizing response bodies.
var (
	userProfilePKPattern   = regexp.MustCompile(`"userProfilePK"\s*:\s*\d+`)
	userProfilePkPattern   = regexp.MustCompile(`"userProfilePk"\s*:\s*\d+`)
	userProfileIDPattern   = regexp.MustCompile(`"userProfileId"\s*:\s*\d+`)
	userIDPattern          = regexp.MustCompile(`"userId"\s*:\s*\d+`)
	ownerIDPattern         = regexp.MustCompile(`"ownerId"\s*:\s*\d+`)
	profileIDPattern       = regexp.MustCompile(`"profileId"\s*:\s*\d+`)
	bareIDPattern          = regexp.MustCompile(`"id"\s*:\s*\d+`)
	displayNamePattern     = regexp.MustCompile(`"displayName"\s*:\s*"[^"]*"`)
	displaynamePattern     = regexp.MustCompile(`"displayname"\s*:\s*"[^"]*"`)
	ownerDisplayPattern    = regexp.MustCompile(`"ownerDisplayName"\s*:\s*"[^"]*"`)
	fullNamePattern        = regexp.MustCompile(`"fullName"\s*:\s*"[^"]*"`)
	fullnamePattern        = regexp.MustCompile(`"fullname"\s*:\s*"[^"]*"`)
	ownerFullNamePattern   = regexp.MustCompile(`"ownerFullName"\s*:\s*"[^"]*"`)
	userProfileFullPattern = regexp.MustCompile(`"userProfileFullName"\s*:\s*"[^"]*"`)
	firstNamePattern       = regexp.MustCompile(`"firstName"\s*:\s*"[^"]*"`)
	lastNamePattern        = regexp.MustCompile(`"lastName"\s*:\s*"[^"]*"`)
	emailPattern           = regexp.MustCompile(`"email"\s*:\s*"[^"]*"`)
	userNamePattern        = regexp.MustCompile(`"userName"\s*:\s*"[^"]*"`)
	birthDatePattern       = regexp.MustCompile(`"birthDate"\s*:\s*"[^"]*"`)
	locationPattern        = regexp.MustCompile(`"location"\s*:\s*"[^"]*"`)

	// Device-related patterns
	deviceIDPattern     = regexp.MustCompile(`"deviceId"\s*:\s*\d+`)
	unitIDPattern       = regexp.MustCompile(`"unitId"\s*:\s*\d+`)
	serialNumberPattern = regexp.MustCompile(`"serialNumber"\s*:\s*"[^"]*"`)

	// Course-related patterns
	courseIDPattern       = regexp.MustCompile(`"courseId"\s*:\s*\d+`)
	courseNamePattern     = regexp.MustCompile(`"courseName"\s*:\s*"[^"]*"`)
	coursePointIDPattern  = regexp.MustCompile(`"coursePointId"\s*:\s*\d+`)
	coursePKPattern       = regexp.MustCompile(`"coursePk"\s*:\s*\d+`)
	virtualPartnerPattern = regexp.MustCompile(`"virtualPartnerId"\s*:\s*\d+`)
	startLatitudePattern  = regexp.MustCompile(`"startLatitude"\s*:\s*-?[\d.]+`)
	startLongitudePattern = regexp.MustCompile(`"startLongitude"\s*:\s*-?[\d.]+`)

	// Geolocation patterns (covers geoPoints, coursePoints, boundingBox, startPoint)
	latitudePattern         = regexp.MustCompile(`"latitude"\s*:\s*-?[\d.]+`)
	longitudePattern        = regexp.MustCompile(`"longitude"\s*:\s*-?[\d.]+`)
	latPattern              = regexp.MustCompile(`"lat"\s*:\s*-?[\d.]+`)
	lonPattern              = regexp.MustCompile(`"lon"\s*:\s*-?[\d.]+`)
	elevationPattern        = regexp.MustCompile(`"elevation"\s*:\s*-?[\d.]+`)
	derivedElevationPattern = regexp.MustCompile(`"derivedElevation"\s*:\s*-?[\d.]+`)

	// URL path patterns (for anonymizing IDs in request URLs)
	deviceSettingsURLPattern  = regexp.MustCompile(`/device-info/settings/\d+`)
	racePredictionsURLPattern = regexp.MustCompile(`/racepredictions/(latest|daily|monthly)/[a-zA-Z0-9-]+`)
	courseURLPattern          = regexp.MustCompile(`/course-service/course/\d+`)

	// Profile image URLs
	profileImageURLPattern = regexp.MustCompile(`"(ownerProfileImageUrl[^"]*|profileImageUrl[^"]*)"\s*:\s*"https://s3\.amazonaws\.com/garmin-connect-prod/profile_images/[^"]*"`)

	// Profile image filenames (contain UUIDs)
	profileImgNameLargePattern  = regexp.MustCompile(`"profileImgNameLarge"\s*:\s*"[^"]*"`)
	profileImgNameMediumPattern = regexp.MustCompile(`"profileImgNameMedium"\s*:\s*"[^"]*"`)
	profileImgNameSmallPattern  = regexp.MustCompile(`"profileImgNameSmall"\s*:\s*"[^"]*"`)

	// Auth-related patterns
	ticketPattern          = regexp.MustCompile(`ticket=ST-[^&"\\]+`)
	oauth1TokenPattern     = regexp.MustCompile(`oauth_token=[^&\s]+`)
	oauth1SecretPattern    = regexp.MustCompile(`oauth_token_secret=[^&\s]+`)
	accessTokenPattern     = regexp.MustCompile(`"access_token"\s*:\s*"[^"]*"`)
	refreshTokenPattern    = regexp.MustCompile(`"refresh_token"\s*:\s*"[^"]*"`)
	garminGUIDSnakePattern = regexp.MustCompile(`"garmin_guid"\s*:\s*"[^"]*"`)
	garminGUIDCamelPattern = regexp.MustCompile(`"garminGUID"\s*:\s*"[^"]*"`)
)

// CassetteDir is the directory where cassettes are stored.
const CassetteDir = "testdata/cassettes"

// sensitiveHeaders are headers that should be sanitized in recordings.
var sensitiveHeaders = []string{
	"Authorization",
	"Cookie",
	"Set-Cookie",
}

// NewRecorder creates a new VCR recorder for the given cassette name.
// In recording mode, it records HTTP interactions to the cassette file.
// In replay mode, it replays recorded interactions.
func NewRecorder(cassetteName string, mode recorder.Mode) (*recorder.Recorder, error) {
	cassettePath := filepath.Join(CassetteDir, cassetteName)

	// Ensure cassette directory exists
	if mode == recorder.ModeRecordOnly || mode == recorder.ModeRecordOnce {
		if err := os.MkdirAll(CassetteDir, 0o755); err != nil {
			return nil, err
		}
	}

	r, err := recorder.New(
		cassettePath,
		recorder.WithMode(mode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithHook(sanitizeHook, recorder.BeforeSaveHook),
		recorder.WithMatcher(flexibleMatcher),
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewReplayRecorder creates a recorder in replay-only mode for tests.
func NewReplayRecorder(cassetteName string) (*recorder.Recorder, error) {
	return NewRecorder(cassetteName, recorder.ModeReplayOnly)
}

// NewRecordingRecorder creates a recorder that records new interactions.
func NewRecordingRecorder(cassetteName string) (*recorder.Recorder, error) {
	return NewRecorder(cassetteName, recorder.ModeRecordOnly)
}

// HTTPClientWithRecorder returns an http.Client that uses the recorder.
func HTTPClientWithRecorder(r *recorder.Recorder) *http.Client {
	return r.GetDefaultClient()
}

// FakeSessionJSON returns a JSON string representing a fake authenticated session.
// Use this with client.LoadSession() to create pre-authenticated test clients
// that can replay API cassettes without needing to replay the auth flow.
func FakeSessionJSON() string {
	return `{
		"oauth1_token": "fake-oauth1-token",
		"oauth1_secret": "fake-oauth1-secret",
		"oauth2_access_token": "fake-oauth2-access-token",
		"oauth2_refresh_token": "fake-oauth2-refresh-token",
		"oauth2_expiry": "2099-01-01T00:00:00Z",
		"domain": "garmin.com"
	}`
}

// sanitizeHook removes sensitive information from recorded interactions.
func sanitizeHook(i *cassette.Interaction) error {
	// Sanitize request headers
	for _, header := range sensitiveHeaders {
		if i.Request.Headers.Get(header) != "" {
			i.Request.Headers.Set(header, "[REDACTED]")
		}
	}

	// Sanitize response headers
	for _, header := range sensitiveHeaders {
		if i.Response.Headers.Get(header) != "" {
			i.Response.Headers.Set(header, "[REDACTED]")
		}
	}

	// Sanitize OAuth tokens in URL query parameters
	if strings.Contains(i.Request.URL, "ticket=") {
		i.Request.URL = redactQueryParam(i.Request.URL, "ticket")
	}

	// Anonymize device IDs in URL paths
	i.Request.URL = deviceSettingsURLPattern.ReplaceAllString(i.Request.URL, "/device-info/settings/12345678")

	// Anonymize displayName/UUID in race predictions URLs
	i.Request.URL = racePredictionsURLPattern.ReplaceAllString(i.Request.URL, "/racepredictions/$1/anonymous")

	// Anonymize course IDs in URL paths
	i.Request.URL = courseURLPattern.ReplaceAllString(i.Request.URL, "/course-service/course/87654321")

	// Sanitize request body (for login requests)
	if strings.Contains(i.Request.Body, "password") {
		i.Request.Body = "[REDACTED]"
	}

	// Sanitize form data
	sanitizeFormData(i.Request.Form)

	// Anonymize personal information in response body
	i.Response.Body = anonymizeBody(i.Response.Body)

	return nil
}

// sanitizeFormData redacts sensitive form fields.
func sanitizeFormData(form map[string][]string) {
	sensitiveFields := []string{"password", "username", "ticket"}
	for _, field := range sensitiveFields {
		if _, ok := form[field]; ok {
			form[field] = []string{"[REDACTED]"}
		}
	}
}

// anonymizeBody replaces personal information in JSON response bodies.
func anonymizeBody(body string) string {
	// User profile IDs
	body = userProfilePKPattern.ReplaceAllString(body, `"userProfilePK":12345678`)
	body = userProfilePkPattern.ReplaceAllString(body, `"userProfilePk":12345678`)
	body = userProfileIDPattern.ReplaceAllString(body, `"userProfileId":12345678`)
	body = userIDPattern.ReplaceAllString(body, `"userId":12345678`)
	body = ownerIDPattern.ReplaceAllString(body, `"ownerId":12345678`)
	body = profileIDPattern.ReplaceAllString(body, `"profileId":12345678`)
	body = bareIDPattern.ReplaceAllString(body, `"id":12345678`)

	// Display names
	body = displayNamePattern.ReplaceAllString(body, `"displayName":"anonymous"`)
	body = displaynamePattern.ReplaceAllString(body, `"displayname":"anonymous"`)
	body = ownerDisplayPattern.ReplaceAllString(body, `"ownerDisplayName":"anonymous"`)

	// Full names
	body = fullNamePattern.ReplaceAllString(body, `"fullName":"Anonymous User"`)
	body = fullnamePattern.ReplaceAllString(body, `"fullname":"Anonymous User"`)
	body = ownerFullNamePattern.ReplaceAllString(body, `"ownerFullName":"Anonymous User"`)
	body = userProfileFullPattern.ReplaceAllString(body, `"userProfileFullName":"Anonymous User"`)

	body = firstNamePattern.ReplaceAllString(body, `"firstName":"Anonymous"`)
	body = lastNamePattern.ReplaceAllString(body, `"lastName":"User"`)
	body = emailPattern.ReplaceAllString(body, `"email":"anonymous@example.com"`)
	body = userNamePattern.ReplaceAllString(body, `"userName":"anonymous"`)

	// Personal info
	body = birthDatePattern.ReplaceAllString(body, `"birthDate":"1990-01-01"`)
	body = locationPattern.ReplaceAllString(body, `"location":"Anonymous City"`)

	// Device info
	body = deviceIDPattern.ReplaceAllString(body, `"deviceId":12345678`)
	body = unitIDPattern.ReplaceAllString(body, `"unitId":12345678`)
	body = serialNumberPattern.ReplaceAllString(body, `"serialNumber":"ABC123456"`)

	// Course info
	body = courseIDPattern.ReplaceAllString(body, `"courseId":87654321`)
	body = courseNamePattern.ReplaceAllString(body, `"courseName":"Anonymous Course"`)
	body = coursePointIDPattern.ReplaceAllString(body, `"coursePointId":11111111`)
	body = coursePKPattern.ReplaceAllString(body, `"coursePk":87654321`)
	body = virtualPartnerPattern.ReplaceAllString(body, `"virtualPartnerId":87654321`)
	body = startLatitudePattern.ReplaceAllString(body, `"startLatitude":48.8566`)
	body = startLongitudePattern.ReplaceAllString(body, `"startLongitude":2.3522`)

	// Geolocation
	body = latitudePattern.ReplaceAllString(body, `"latitude":48.8566`)
	body = longitudePattern.ReplaceAllString(body, `"longitude":2.3522`)
	body = latPattern.ReplaceAllString(body, `"lat":48.8566`)
	body = lonPattern.ReplaceAllString(body, `"lon":2.3522`)
	body = elevationPattern.ReplaceAllString(body, `"elevation":100.0`)
	body = derivedElevationPattern.ReplaceAllString(body, `"derivedElevation":100.0`)

	// Profile image URLs
	body = profileImageURLPattern.ReplaceAllString(body, `"$1":"https://example.com/profile.png"`)

	// Profile image filenames
	body = profileImgNameLargePattern.ReplaceAllString(body, `"profileImgNameLarge":"anonymous-profile-large.png"`)
	body = profileImgNameMediumPattern.ReplaceAllString(body, `"profileImgNameMedium":"anonymous-profile-medium.png"`)
	body = profileImgNameSmallPattern.ReplaceAllString(body, `"profileImgNameSmall":"anonymous-profile-small.png"`)

	// Auth tokens and tickets
	body = ticketPattern.ReplaceAllString(body, `ticket=[REDACTED]`)
	body = oauth1TokenPattern.ReplaceAllString(body, `oauth_token=[REDACTED]`)
	body = oauth1SecretPattern.ReplaceAllString(body, `oauth_token_secret=[REDACTED]`)
	body = accessTokenPattern.ReplaceAllString(body, `"access_token":"[REDACTED]"`)
	body = refreshTokenPattern.ReplaceAllString(body, `"refresh_token":"[REDACTED]"`)
	body = garminGUIDSnakePattern.ReplaceAllString(body, `"garmin_guid":"00000000-0000-0000-0000-000000000000"`)
	body = garminGUIDCamelPattern.ReplaceAllString(body, `"garminGUID":"00000000-0000-0000-0000-000000000000"`)

	return body
}

// flexibleMatcher matches requests ignoring volatile headers and query params.
func flexibleMatcher(r *http.Request, i cassette.Request) bool {
	// Match on method
	if r.Method != i.Method {
		return false
	}

	// Match on URL path (ignore query string for flexibility)
	reqURL := r.URL.Path
	cassetteURL := extractPath(i.URL)

	return reqURL == cassetteURL
}

// redactQueryParam replaces a query parameter value with [REDACTED].
func redactQueryParam(urlStr, param string) string {
	parts := strings.Split(urlStr, param+"=")
	if len(parts) != 2 {
		return urlStr
	}

	endIdx := strings.IndexAny(parts[1], "&# ")
	if endIdx == -1 {
		return parts[0] + param + "=[REDACTED]"
	}
	return parts[0] + param + "=[REDACTED]" + parts[1][endIdx:]
}

// extractPath extracts just the path from a URL string.
func extractPath(urlStr string) string {
	_, rest, found := strings.Cut(urlStr, "://")
	if !found {
		return urlStr
	}

	_, path, found := strings.Cut(rest, "/")
	if !found {
		return urlStr
	}
	path = "/" + path

	pathOnly, _, _ := strings.Cut(path, "?")
	return pathOnly
}
