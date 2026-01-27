// endpoint/validator_test.go
package endpoint

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidator_ValidEndpoint(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Service:    "Sleep",
		Cassette:   "sleep_daily",
		Path:       "/sleep-service/sleep",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		MCPTool:    "get_sleep",
		Short:      "Get sleep",
		Long:       "Get sleep data",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	// Create temp cassette dir with file
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidator_MissingHandler(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		// Handler missing
	})

	v := NewValidator(r, ValidatorConfig{CassetteDir: t.TempDir()})
	errs := v.Validate()

	if !containsError(errs, "missing Handler") {
		t.Errorf("expected 'missing Handler' error, got: %v", errs)
	}
}

func TestValidator_MissingCassette(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "nonexistent",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	v := NewValidator(r, ValidatorConfig{CassetteDir: t.TempDir()})
	errs := v.Validate()

	if !containsError(errs, "cassette file not found") {
		t.Errorf("expected 'cassette file not found' error, got: %v", errs)
	}
}

func TestValidator_NoCLIOrMCP(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
		// No CLICommand or MCPTool
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "must have CLICommand or MCPTool") {
		t.Errorf("expected 'must have CLICommand or MCPTool' error, got: %v", errs)
	}
}

func TestValidator_InvalidHTTPMethod(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "INVALID",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "invalid HTTPMethod") {
		t.Errorf("expected 'invalid HTTPMethod' error, got: %v", errs)
	}
}

func TestValidator_EmptyCassette(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "", // Empty cassette
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	v := NewValidator(r, ValidatorConfig{CassetteDir: t.TempDir()})
	errs := v.Validate()

	if !containsError(errs, "missing Cassette") {
		t.Errorf("expected 'missing Cassette' error, got: %v", errs)
	}
}

func TestValidator_MissingShortDescription(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "", // Missing short
		Long:       "Long description",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "missing Short description") {
		t.Errorf("expected 'missing Short description' error, got: %v", errs)
	}
}

func TestValidator_MissingLongDescription(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "", // Missing long
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "missing Long description") {
		t.Errorf("expected 'missing Long description' error, got: %v", errs)
	}
}

func TestValidator_MissingPath(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "", // Missing path
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "missing Path") {
		t.Errorf("expected 'missing Path' error, got: %v", errs)
	}
}

func TestValidator_POSTWithoutBody(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "CreateWorkout",
		Cassette:   "create_workout",
		Path:       "/workout",
		HTTPMethod: "POST",
		CLICommand: "workout",
		Short:      "Create workout",
		Long:       "Create a new workout",
		Body:       nil, // Missing body for POST
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "create_workout.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "POST endpoint should have Body config") {
		t.Errorf("expected 'POST endpoint should have Body config' error, got: %v", errs)
	}
}

func TestValidator_PUTWithoutBody(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "UpdateWorkout",
		Cassette:   "update_workout",
		Path:       "/workout",
		HTTPMethod: "PUT",
		CLICommand: "workout",
		Short:      "Update workout",
		Long:       "Update an existing workout",
		Body:       nil, // Missing body for PUT
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "update_workout.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "PUT endpoint should have Body config") {
		t.Errorf("expected 'PUT endpoint should have Body config' error, got: %v", errs)
	}
}

func TestValidator_ParamMissingDescription(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Params: []Param{
			{Name: "date", Type: ParamTypeDate, Description: ""}, // Missing description
		},
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "param date missing description") {
		t.Errorf("expected 'param date missing description' error, got: %v", errs)
	}
}

func TestValidator_DependsOnUnknownEndpoint(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:        "GetActivity",
		Cassette:    "activity",
		Path:        "/activity",
		HTTPMethod:  "GET",
		CLICommand:  "activity",
		Short:       "Get activity",
		Long:        "Get activity details",
		DependsOn:   "NonExistent", // Unknown endpoint
		ArgProvider: func(_ any) map[string]any { return nil },
		Handler:     func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "activity.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "DependsOn references unknown endpoint: NonExistent") {
		t.Errorf("expected 'DependsOn references unknown endpoint' error, got: %v", errs)
	}
}

func TestValidator_DependsOnWithoutArgProvider(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "ListActivities",
		Cassette:   "activities_list",
		Path:       "/activities",
		HTTPMethod: "GET",
		CLICommand: "activity",
		Short:      "List activities",
		Long:       "List all activities",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})
	r.Register(Endpoint{
		Name:        "GetActivity",
		Cassette:    "activity",
		Path:        "/activity",
		HTTPMethod:  "GET",
		CLICommand:  "activity",
		Short:       "Get activity",
		Long:        "Get activity details",
		DependsOn:   "ListActivities",
		ArgProvider: nil, // Missing ArgProvider
		Handler:     func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "activities_list.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "activity.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "has DependsOn but missing ArgProvider") {
		t.Errorf("expected 'has DependsOn but missing ArgProvider' error, got: %v", errs)
	}
}

func TestValidator_OrphanedCassette(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	// Create the registered cassette
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}
	// Create an orphaned cassette
	if err := os.WriteFile(filepath.Join(tmpDir, "orphaned_cassette.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write orphaned cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	if !containsError(errs, "orphaned cassette (no endpoint references it): orphaned_cassette") {
		t.Errorf("expected 'orphaned cassette' error, got: %v", errs)
	}
}

func TestValidator_AuthCassetteSkipped(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		Cassette:   "sleep_daily",
		Path:       "/test",
		HTTPMethod: "GET",
		CLICommand: "sleep",
		Short:      "Short",
		Long:       "Long",
		Handler:    func(_ context.Context, _ any, _ *HandlerArgs) (any, error) { return struct{}{}, nil },
	})

	tmpDir := t.TempDir()
	// Create the registered cassette
	if err := os.WriteFile(filepath.Join(tmpDir, "sleep_daily.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write cassette file: %v", err)
	}
	// Create auth cassette - should be skipped, not reported as orphaned
	if err := os.WriteFile(filepath.Join(tmpDir, "auth.yaml"), []byte(""), 0o600); err != nil {
		t.Fatalf("failed to write auth cassette file: %v", err)
	}

	v := NewValidator(r, ValidatorConfig{CassetteDir: tmpDir})
	errs := v.Validate()

	// Should have no errors - auth cassette should be skipped
	if containsError(errs, "orphaned cassette") {
		t.Errorf("auth cassette should be skipped, got errors: %v", errs)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func containsError(errs []string, substr string) bool {
	for _, e := range errs {
		if strings.Contains(e, substr) {
			return true
		}
	}
	return false
}
