package main

import (
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// parseOptionalDate parses a date string from request arguments.
// Returns today if the argument is missing or empty.
//
//nolint:unused,unparam // Helper for MCP tool handlers to be implemented
func parseOptionalDate(request mcp.CallToolRequest, key string) (time.Time, error) {
	dateStr, _ := request.RequireString(key) // Ignore error - missing key means use default
	if dateStr == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", dateStr)
}

// parseRequiredDate parses a required date string from request arguments.
//
//nolint:unused // Helper for MCP tool handlers to be implemented
func parseRequiredDate(request mcp.CallToolRequest, key string) (time.Time, error) {
	dateStr, err := request.RequireString(key)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("2006-01-02", dateStr)
}

// jsonResult converts any value to a JSON MCP tool result.
// If marshaling fails, returns an error result instead.
//
//nolint:unused // Helper for MCP tool handlers to be implemented
func jsonResult(v any) *mcp.CallToolResult {
	data, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error())
	}
	return mcp.NewToolResultText(string(data))
}

// errorResult creates an MCP error result.
//
//nolint:unused // Helper for MCP tool handlers to be implemented
func errorResult(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(err.Error())
}
