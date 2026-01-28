package definitions

import (
	"context"
	"testing"
	"time"
)

func TestUtilityEndpoints_Registered(t *testing.T) {
	if len(UtilityEndpoints) == 0 {
		t.Fatal("UtilityEndpoints should not be empty")
	}

	ep := UtilityEndpoints[0]
	if ep.Name != "GetCurrentDate" {
		t.Errorf("Name = %q, want GetCurrentDate", ep.Name)
	}
	if ep.MCPTool != "get_current_date" {
		t.Errorf("MCPTool = %q, want get_current_date", ep.MCPTool)
	}
	if ep.Handler == nil {
		t.Error("Handler should not be nil")
	}
}

func TestUtilityEndpoints_GetCurrentDate_Handler(t *testing.T) {
	ep := UtilityEndpoints[0]

	result, err := ep.Handler(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want map[string]any", result)
	}

	// Check required fields exist
	requiredFields := []string{"date", "year", "month", "month_name", "day", "weekday", "iso8601"}
	for _, field := range requiredFields {
		if _, ok := m[field]; !ok {
			t.Errorf("missing field %q in result", field)
		}
	}

	// Verify date format
	dateStr, ok := m["date"].(string)
	if !ok {
		t.Fatalf("date type = %T, want string", m["date"])
	}
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		t.Errorf("date %q is not in YYYY-MM-DD format: %v", dateStr, err)
	}

	// Verify year is reasonable
	year, ok := m["year"].(int)
	if !ok {
		t.Fatalf("year type = %T, want int", m["year"])
	}
	if year < 2020 || year > 2100 {
		t.Errorf("year = %d, want reasonable year", year)
	}

	// Verify month is valid
	month, ok := m["month"].(int)
	if !ok {
		t.Fatalf("month type = %T, want int", m["month"])
	}
	if month < 1 || month > 12 {
		t.Errorf("month = %d, want 1-12", month)
	}

	// Verify day is valid
	day, ok := m["day"].(int)
	if !ok {
		t.Fatalf("day type = %T, want int", m["day"])
	}
	if day < 1 || day > 31 {
		t.Errorf("day = %d, want 1-31", day)
	}

	// Verify weekday is a valid day name
	weekday, ok := m["weekday"].(string)
	if !ok {
		t.Fatalf("weekday type = %T, want string", m["weekday"])
	}
	validWeekdays := map[string]bool{
		"Sunday": true, "Monday": true, "Tuesday": true, "Wednesday": true,
		"Thursday": true, "Friday": true, "Saturday": true,
	}
	if !validWeekdays[weekday] {
		t.Errorf("weekday = %q, want valid weekday name", weekday)
	}

	// Verify iso8601 format
	iso8601, ok := m["iso8601"].(string)
	if !ok {
		t.Fatalf("iso8601 type = %T, want string", m["iso8601"])
	}
	if _, err := time.Parse(time.RFC3339, iso8601); err != nil {
		t.Errorf("iso8601 %q is not in RFC3339 format: %v", iso8601, err)
	}
}
