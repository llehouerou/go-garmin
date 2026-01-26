// errors.go
package garmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotAuthenticated   = errors.New("garmin: not authenticated")
	ErrSessionExpired     = errors.New("garmin: session expired, re-login required")
	ErrMFARequired        = errors.New("garmin: MFA required but no handler provided")
	ErrRateLimited        = errors.New("garmin: rate limited, retry later")
	ErrMaxRetriesExceeded = errors.New("garmin: max retries exceeded")
	ErrNotFound           = errors.New("garmin: resource not found")
)

type APIError struct {
	StatusCode int
	Status     string
	Endpoint   string
	Message    string
	Body       json.RawMessage
}

func (e *APIError) Error() string {
	return fmt.Sprintf("garmin: %s %s: %s", e.Status, e.Endpoint, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return errors.Is(err, ErrNotFound)
}

func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return errors.Is(err, ErrRateLimited)
}

func IsAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized || apiErr.StatusCode == http.StatusForbidden
	}
	return errors.Is(err, ErrNotAuthenticated) || errors.Is(err, ErrSessionExpired)
}

func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500
	}
	return false
}

func IsRetryable(err error) bool {
	return IsRateLimited(err) || IsServerError(err)
}
