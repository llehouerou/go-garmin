// fetch.go
package garmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RawSetter is implemented by response types that store raw JSON.
type RawSetter interface {
	SetRaw(json.RawMessage)
}

// fetch performs a GET request and unmarshals the response into T.
// Returns ErrNotFound if the response status is 204 No Content or 404 Not Found.
// Usage: fetch[DailySleep, *DailySleep](ctx, client, path)
func fetch[T any, PT interface {
	*T
	RawSetter
}](ctx context.Context, c *Client, path string) (*T, error) {
	resp, err := c.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(T)
	if err := json.Unmarshal(raw, result); err != nil {
		return nil, err
	}
	PT(result).SetRaw(raw)

	return result, nil
}

// send performs a POST/PUT/PATCH request with a JSON body and unmarshals the response into T.
// Returns APIError if the response status is not in the 2xx range.
// Usage: send[Workout, *Workout](ctx, client, http.MethodPost, path, body)
func send[T any, PT interface {
	*T
	RawSetter
}, R any](ctx context.Context, c *Client, method, path string, body R) (*T, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doAPIWithBody(ctx, method, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: raw}
	}

	result := new(T)
	if err := json.Unmarshal(raw, result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	PT(result).SetRaw(raw)

	return result, nil
}

// upload performs a multipart file upload and unmarshals the response into T.
// Returns APIError if the response status is not in the 2xx range.
func upload[T any, PT interface {
	*T
	RawSetter
}](ctx context.Context, c *Client, path, fieldName, fileName string, content io.Reader) (*T, error) {
	resp, err := c.doAPIMultipart(ctx, path, fieldName, fileName, content)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: raw}
	}

	result := new(T)
	if err := json.Unmarshal(raw, result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	PT(result).SetRaw(raw)

	return result, nil
}

// sendEmpty performs a request expecting no response body (e.g., DELETE).
// Returns APIError if the response status is not in the 2xx range.
func sendEmpty(ctx context.Context, c *Client, method, path string) error {
	resp, err := c.doAPI(ctx, method, path, http.NoBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: raw}
	}

	return nil
}
