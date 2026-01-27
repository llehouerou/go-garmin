// endpoint/endpoint_test.go
package endpoint

import (
	"fmt"
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

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	ep := Endpoint{
		Name:       "GetSleep",
		CLICommand: "sleep",
		MCPTool:    "get_sleep",
	}
	r.Register(ep)

	if len(r.All()) != 1 {
		t.Errorf("All() len = %d, want 1", len(r.All()))
	}
}

func TestRegistry_ByCLI(t *testing.T) {
	r := NewRegistry()

	r.Register(Endpoint{Name: "GetSleep", CLICommand: "sleep", MCPTool: "get_sleep"})
	r.Register(Endpoint{Name: "ListWorkouts", CLICommand: "workouts", CLISubcommand: "list"})

	byCLI := r.ByCLI()
	if _, ok := byCLI["sleep"]; !ok {
		t.Error("expected 'sleep' in ByCLI")
	}
	if _, ok := byCLI["workouts:list"]; !ok {
		t.Error("expected 'workouts:list' in ByCLI")
	}
}

func TestRegistry_ByMCP(t *testing.T) {
	r := NewRegistry()

	r.Register(Endpoint{Name: "GetSleep", MCPTool: "get_sleep"})

	byMCP := r.ByMCP()
	if _, ok := byMCP["get_sleep"]; !ok {
		t.Error("expected 'get_sleep' in ByMCP")
	}
}

func TestRegistry_ByName(t *testing.T) {
	r := NewRegistry()

	r.Register(Endpoint{Name: "GetSleep"})

	if ep := r.ByName("GetSleep"); ep == nil {
		t.Error("expected to find 'GetSleep' by name")
	}
	if ep := r.ByName("NotFound"); ep != nil {
		t.Error("expected nil for unknown name")
	}
}

func TestRegistry_PointerStability(t *testing.T) {
	r := NewRegistry()

	r.Register(Endpoint{Name: "First"})
	first := r.ByName("First")

	// Register many more to force reallocation
	for i := range 100 {
		r.Register(Endpoint{Name: fmt.Sprintf("Endpoint%d", i)})
	}

	// Verify first pointer is still valid
	if got := r.ByName("First"); got != first || got.Name != "First" {
		t.Error("pointer to First became invalid after reallocation")
	}
}
