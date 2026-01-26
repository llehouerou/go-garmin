// http.go
package garmin

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type retryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func defaultRetryConfig() retryConfig {
	return retryConfig{
		MaxRetries:     3,
		InitialBackoff: time.Second,
		MaxBackoff:     30 * time.Second,
	}
}

type httpTransport struct {
	client      *http.Client
	retry       retryConfig
	rateLimiter *rateLimiter
}

func newHTTPTransport(client *http.Client, retry retryConfig, rl *rateLimiter) *httpTransport {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &httpTransport{
		client:      client,
		retry:       retry,
		rateLimiter: rl,
	}
}

func (t *httpTransport) do(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= t.retry.MaxRetries; attempt++ {
		if t.rateLimiter != nil {
			if err := t.rateLimiter.Wait(req.Context()); err != nil {
				return nil, err
			}
		}

		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = err
			if !t.shouldRetry(0, attempt) {
				return nil, err
			}
			t.backoff(req.Context(), attempt)
			continue
		}

		if !t.shouldRetry(resp.StatusCode, attempt) {
			return resp, nil
		}

		// Drain and close body for retry
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		lastErr = &APIError{StatusCode: resp.StatusCode, Status: resp.Status}

		t.backoff(req.Context(), attempt)
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrMaxRetriesExceeded
}

func (t *httpTransport) shouldRetry(statusCode, attempt int) bool {
	if attempt >= t.retry.MaxRetries {
		return false
	}
	if statusCode == 0 {
		return true // network error
	}
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func (t *httpTransport) backoff(ctx context.Context, attempt int) {
	backoff := t.retry.InitialBackoff * (1 << attempt)
	backoff = min(backoff, t.retry.MaxBackoff)
	// Add jitter (0-25%)
	//nolint:gosec // weak random is acceptable for backoff jitter
	jitter := time.Duration(rand.Int63n(int64(backoff / 4)))
	backoff += jitter

	select {
	case <-time.After(backoff):
	case <-ctx.Done():
	}
}
