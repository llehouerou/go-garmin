// endpoint/endpoint.go
package endpoint

import (
	"context"
	"reflect"
	"time"
)

// ParamType defines the type of a simple parameter.
type ParamType int

const (
	ParamTypeString ParamType = iota
	ParamTypeInt
	ParamTypeDate      // YYYY-MM-DD, defaults to today
	ParamTypeDateRange // provides "start" and "end" params
	ParamTypeBool
)

// Param describes a simple endpoint parameter (query/path).
type Param struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// BodyConfig describes a complex JSON request body.
type BodyConfig struct {
	Type        reflect.Type
	Description string
	Example     string
}

// Endpoint defines a single API endpoint with all its metadata.
type Endpoint struct {
	// Identity
	Name     string
	Service  string
	Cassette string

	// API details
	Path       string
	HTTPMethod string
	Params     []Param
	Body       *BodyConfig

	// Dependencies
	DependsOn   string
	ArgProvider func(dependencyResult any) map[string]any

	// CLI configuration
	CLICommand    string
	CLISubcommand string
	CLIAliases    []string

	// MCP configuration
	MCPTool string

	// Documentation
	Short string
	Long  string

	// Handler
	Handler func(ctx context.Context, client any, args *HandlerArgs) (any, error)
}

// HandlerArgs provides typed access to parsed parameters.
type HandlerArgs struct {
	Params map[string]any
	Body   any
}

// Date returns a time.Time param, or current time if not set.
func (a *HandlerArgs) Date(name string) time.Time {
	if v, ok := a.Params[name].(time.Time); ok {
		return v
	}
	return time.Now()
}

// Int returns an int param, or 0 if not set.
func (a *HandlerArgs) Int(name string) int {
	if v, ok := a.Params[name].(int); ok {
		return v
	}
	return 0
}

// IntOrDefault returns an int param, or the default if not set.
func (a *HandlerArgs) IntOrDefault(name string, def int) int {
	if v, ok := a.Params[name].(int); ok {
		return v
	}
	return def
}

// String returns a string param, or empty if not set.
func (a *HandlerArgs) String(name string) string {
	if v, ok := a.Params[name].(string); ok {
		return v
	}
	return ""
}

// Bool returns a bool param, or false if not set.
func (a *HandlerArgs) Bool(name string) bool {
	if v, ok := a.Params[name].(bool); ok {
		return v
	}
	return false
}
