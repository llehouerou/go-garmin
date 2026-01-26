// errors_test.go
package garmin

import (
	"errors"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Status:     "404 Not Found",
		Endpoint:   "/wellness-service/wellness/dailySleep",
		Message:    "No data found",
	}

	if err.Error() != "garmin: 404 Not Found /wellness-service/wellness/dailySleep: No data found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
	if IsRateLimited(err) {
		t.Error("expected IsRateLimited to return false")
	}
}

func TestSentinelErrors(t *testing.T) {
	if !errors.Is(ErrNotAuthenticated, ErrNotAuthenticated) {
		t.Error("sentinel error identity failed")
	}
}
