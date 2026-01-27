// endpoint/endpoint_test.go
package endpoint

import (
	"testing"
	"time"
)

func TestHandlerArgs_Date(t *testing.T) {
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	args := &HandlerArgs{
		Params: map[string]any{"date": date},
	}

	got := args.Date("date")
	if !got.Equal(date) {
		t.Errorf("Date() = %v, want %v", got, date)
	}
}

func TestHandlerArgs_DateDefault(t *testing.T) {
	args := &HandlerArgs{Params: map[string]any{}}

	got := args.Date("missing")
	if got.IsZero() {
		t.Error("Date() should return current time for missing key, got zero")
	}
}

func TestHandlerArgs_Int(t *testing.T) {
	args := &HandlerArgs{
		Params: map[string]any{"limit": 10},
	}

	if got := args.Int("limit"); got != 10 {
		t.Errorf("Int() = %d, want 10", got)
	}
	if got := args.Int("missing"); got != 0 {
		t.Errorf("Int() for missing = %d, want 0", got)
	}
}

func TestHandlerArgs_IntOrDefault(t *testing.T) {
	args := &HandlerArgs{
		Params: map[string]any{"limit": 10},
	}

	if got := args.IntOrDefault("limit", 20); got != 10 {
		t.Errorf("IntOrDefault() = %d, want 10", got)
	}
	if got := args.IntOrDefault("missing", 20); got != 20 {
		t.Errorf("IntOrDefault() for missing = %d, want 20", got)
	}
}

func TestHandlerArgs_String(t *testing.T) {
	args := &HandlerArgs{
		Params: map[string]any{"name": "test"},
	}

	if got := args.String("name"); got != "test" {
		t.Errorf("String() = %q, want %q", got, "test")
	}
	if got := args.String("missing"); got != "" {
		t.Errorf("String() for missing = %q, want empty", got)
	}
}

func TestHandlerArgs_Bool(t *testing.T) {
	args := &HandlerArgs{
		Params: map[string]any{"enabled": true},
	}

	if got := args.Bool("enabled"); got != true {
		t.Errorf("Bool() = %v, want true", got)
	}
	if got := args.Bool("missing"); got != false {
		t.Errorf("Bool() for missing = %v, want false", got)
	}
}
