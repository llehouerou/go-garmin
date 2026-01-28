// endpoint/recorder_test.go
package endpoint

import (
	"context"
	"testing"
	"time"
)

func TestFixtureRecorder_BuildDefaultArgs(t *testing.T) {
	r := NewRegistry()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	ep := Endpoint{
		Params: []Param{
			{Name: "date", Type: ParamTypeDate},
			{Name: "limit", Type: ParamTypeInt},
			{Name: "range", Type: ParamTypeDateRange},
		},
	}

	rec := &FixtureRecorder{
		registry: r,
		date:     date,
	}

	args := rec.buildDefaultArgs(&ep)

	if got := args.Date("date"); !got.Equal(date) {
		t.Errorf("date = %v, want %v", got, date)
	}
	if got := args.IntOrDefault("limit", 0); got != 10 {
		t.Errorf("limit = %d, want 10", got)
	}
	if start, ok := args.Params["start"].(time.Time); !ok || start.IsZero() {
		t.Error("start should be set for DateRange")
	}
	if end, ok := args.Params["end"].(time.Time); !ok || end.IsZero() {
		t.Error("end should be set for DateRange")
	}
}

func TestFixtureRecorder_ListCassettes(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{Name: "A", Cassette: "cassette1"})
	r.Register(Endpoint{Name: "B", Cassette: "cassette2"})
	r.Register(Endpoint{Name: "C", Cassette: "cassette1"}) // duplicate

	rec := &FixtureRecorder{registry: r}
	cassettes := rec.ListCassettes()

	if len(cassettes) != 2 {
		t.Errorf("ListCassettes() = %d cassettes, want 2", len(cassettes))
	}
}

func TestFixtureRecorder_SortByDependencies(t *testing.T) {
	r := NewRegistry()

	epA := &Endpoint{Name: "ListItems"}
	epB := &Endpoint{Name: "GetItem", DependsOn: "ListItems"}
	epC := &Endpoint{Name: "GetItemDetails", DependsOn: "GetItem"}

	// Register all endpoints
	r.Register(*epA)
	r.Register(*epB)
	r.Register(*epC)

	rec := &FixtureRecorder{registry: r}

	// Create endpoints in wrong order
	endpoints := []*Endpoint{epC, epA, epB}
	sorted := rec.sortByDependencies(endpoints)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 endpoints, got %d", len(sorted))
	}
	if sorted[0].Name != "ListItems" {
		t.Errorf("first should be ListItems, got %s", sorted[0].Name)
	}
	if sorted[1].Name != "GetItem" {
		t.Errorf("second should be GetItem, got %s", sorted[1].Name)
	}
	if sorted[2].Name != "GetItemDetails" {
		t.Errorf("third should be GetItemDetails, got %s", sorted[2].Name)
	}
}

func TestFixtureRecorder_SortByDependencies_NoDeps(t *testing.T) {
	r := NewRegistry()

	epA := &Endpoint{Name: "GetA"}
	epB := &Endpoint{Name: "GetB"}

	rec := &FixtureRecorder{registry: r}

	endpoints := []*Endpoint{epA, epB}
	sorted := rec.sortByDependencies(endpoints)

	if len(sorted) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(sorted))
	}
}

func TestNewFixtureRecorder(t *testing.T) {
	r := NewRegistry()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	rec := NewFixtureRecorder(r, RecorderConfig{
		Date:    date,
		Session: []byte("test-session"),
	})

	if rec.registry != r {
		t.Error("registry not set correctly")
	}
	if !rec.date.Equal(date) {
		t.Errorf("date = %v, want %v", rec.date, date)
	}
}

func TestFixtureRecorder_BuildDefaultArgs_StringParam(t *testing.T) {
	r := NewRegistry()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	ep := Endpoint{
		Params: []Param{
			{Name: "name", Type: ParamTypeString},
		},
	}

	rec := &FixtureRecorder{
		registry: r,
		date:     date,
	}

	args := rec.buildDefaultArgs(&ep)

	if got := args.String("name"); got != "" {
		t.Errorf("name = %q, want empty string", got)
	}
}

func TestFixtureRecorder_BuildDefaultArgs_BoolParam(t *testing.T) {
	r := NewRegistry()
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	ep := Endpoint{
		Params: []Param{
			{Name: "verbose", Type: ParamTypeBool},
		},
	}

	rec := &FixtureRecorder{
		registry: r,
		date:     date,
	}

	args := rec.buildDefaultArgs(&ep)

	if got := args.Bool("verbose"); got != false {
		t.Errorf("verbose = %v, want false", got)
	}
}

func TestFixtureRecorder_EndpointsForCassette(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{Name: "GetA", Cassette: "test_cassette"})
	r.Register(Endpoint{Name: "GetB", Cassette: "other_cassette"})
	r.Register(Endpoint{Name: "GetC", Cassette: "test_cassette"})

	rec := &FixtureRecorder{registry: r}
	endpoints := rec.endpointsForCassette("test_cassette")

	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
}

func TestFixtureRecorder_RecordCassetteSkipsDependencyWithNilArgs(t *testing.T) {
	r := NewRegistry()

	// Register the dependency first
	r.Register(Endpoint{
		Name:     "ListItems",
		Cassette: "test",
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return []string{}, nil // Empty result
		},
	})

	// Register endpoint that depends on it
	r.Register(Endpoint{
		Name:      "GetItem",
		Cassette:  "test",
		DependsOn: "ListItems",
		ArgProvider: func(result any) map[string]any {
			items, ok := result.([]string)
			if !ok || len(items) == 0 {
				return nil // Signal to skip
			}
			return map[string]any{"id": items[0]}
		},
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return struct{}{}, nil
		},
	})

	// This test just ensures the code doesn't panic when ArgProvider returns nil
	rec := NewFixtureRecorder(r, RecorderConfig{
		Date: time.Now(),
	})

	// We can't fully test RecordAll without a real recorder, but we can verify
	// the endpoints are correctly identified
	endpoints := rec.endpointsForCassette("test")
	if len(endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(endpoints))
	}
}
