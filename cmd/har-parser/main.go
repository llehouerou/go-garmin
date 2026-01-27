// har-parser extracts Garmin API endpoints from HAR files.
//
// Usage:
//
//	go run ./cmd/har-parser <har-file> [--compare endpoints.md] [--output schemas/]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// HAR file structure (partial)
type HAR struct {
	Log struct {
		Entries []Entry `json:"entries"`
	} `json:"log"`
}

type Entry struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Headers     []Header  `json:"headers"`
	PostData    *PostData `json:"postData,omitempty"`
	QueryString []NVP     `json:"queryString"`
}

type Response struct {
	Status     int      `json:"status"`
	StatusText string   `json:"statusText"`
	Headers    []Header `json:"headers"`
	Content    Content  `json:"content"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PostData struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type Content struct {
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type NVP struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Endpoint represents a discovered API endpoint
type Endpoint struct {
	Method           string
	Path             string
	NormalizedPath   string
	QueryParams      []string
	RequestBodyType  string
	ResponseBodyType string
	StatusCode       int
	Examples         []Example
}

type Example struct {
	RequestBody  string
	ResponseBody string
	QueryString  string
}

// Patterns to normalize path parameters
var pathParamPatterns = []*regexp.Regexp{
	// UUIDs
	regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`),
	// Large numeric IDs (activity IDs, workout IDs, etc.)
	regexp.MustCompile(`/\d{8,}/`),
	// Dates in path: 2024-01-27
	regexp.MustCompile(`/\d{4}-\d{2}-\d{2}(/|$)`),
	// Device IDs (8+ digits)
	regexp.MustCompile(`/\d{7,}(/|$)`),
}

var pathParamReplacements = []string{
	"{uuid}",
	"/{id}/",
	"/{date}/",
	"/{deviceId}/",
}

func main() {
	compareFile := flag.String("compare", "", "Compare with ENDPOINTS.md file")
	outputDir := flag.String("output", "", "Output directory for schemas")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <har-file> [--compare endpoints.md] [--output schemas/] [--json]\n", os.Args[0])
		os.Exit(1)
	}

	harFile := flag.Arg(0)

	// Parse HAR file
	data, err := os.ReadFile(harFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading HAR file: %v\n", err)
		os.Exit(1)
	}

	var har HAR
	if err := json.Unmarshal(data, &har); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HAR file: %v\n", err)
		os.Exit(1)
	}

	// Extract Garmin endpoints
	endpoints := extractGarminEndpoints(har)

	// Group by normalized path
	grouped := groupEndpoints(endpoints)

	if *jsonOutput {
		outputJSON(grouped)
	} else {
		outputMarkdown(grouped)
	}

	// Compare with ENDPOINTS.md if provided
	if *compareFile != "" {
		compareWithEndpointsMD(*compareFile, grouped)
	}

	// Output schemas if directory provided
	if *outputDir != "" {
		outputSchemas(*outputDir, grouped)
	}
}

func extractGarminEndpoints(har HAR) []Endpoint {
	var endpoints []Endpoint

	for i := range har.Log.Entries {
		entry := &har.Log.Entries[i]
		// Filter for Garmin API requests
		if !isGarminAPI(entry.Request.URL) {
			continue
		}

		// Skip non-API content types
		contentType := getHeader(entry.Response.Headers, "content-type")
		if contentType != "" && !strings.Contains(contentType, "json") && !strings.Contains(contentType, "octet-stream") {
			continue
		}

		endpoint := parseEndpoint(*entry)
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

func isGarminAPI(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Match Garmin API domains
	host := u.Host
	if strings.Contains(host, "connectapi.garmin.com") ||
		strings.Contains(host, "connect.garmin.com/gc-api") ||
		strings.Contains(host, "connect.garmin.com/proxy") ||
		(strings.Contains(host, "garmin.com") && strings.Contains(u.Path, "-service")) {
		return true
	}

	return false
}

func parseEndpoint(entry Entry) Endpoint {
	u, _ := url.Parse(entry.Request.URL)

	// Extract path, removing gc-api or proxy prefix
	path := u.Path
	path = strings.TrimPrefix(path, "/gc-api")
	path = strings.TrimPrefix(path, "/proxy")

	// Normalize path parameters
	normalizedPath := normalizePath(path)

	// Extract query params
	queryParams := make([]string, 0, len(entry.Request.QueryString))
	for _, q := range entry.Request.QueryString {
		queryParams = append(queryParams, q.Name)
	}
	sort.Strings(queryParams)

	// Determine request body type
	var requestBodyType string
	if entry.Request.PostData != nil && entry.Request.PostData.Text != "" {
		requestBodyType = detectJSONType(entry.Request.PostData.Text)
	}

	// Determine response body type
	var responseBodyType string
	if entry.Response.Content.Text != "" {
		responseBodyType = detectJSONType(entry.Response.Content.Text)
	}

	// Build query string for example
	var queryString string
	if len(entry.Request.QueryString) > 0 {
		parts := make([]string, 0, len(entry.Request.QueryString))
		for _, q := range entry.Request.QueryString {
			parts = append(parts, fmt.Sprintf("%s=%s", q.Name, q.Value))
		}
		queryString = strings.Join(parts, "&")
	}

	return Endpoint{
		Method:           entry.Request.Method,
		Path:             path,
		NormalizedPath:   normalizedPath,
		QueryParams:      queryParams,
		RequestBodyType:  requestBodyType,
		ResponseBodyType: responseBodyType,
		StatusCode:       entry.Response.Status,
		Examples: []Example{{
			RequestBody:  entry.Request.PostData.GetText(),
			ResponseBody: entry.Response.Content.Text,
			QueryString:  queryString,
		}},
	}
}

func (p *PostData) GetText() string {
	if p == nil {
		return ""
	}
	return p.Text
}

func normalizePath(path string) string {
	normalized := path

	// Apply each pattern
	for i, pattern := range pathParamPatterns {
		if i < len(pathParamReplacements) {
			normalized = pattern.ReplaceAllString(normalized, pathParamReplacements[i])
		}
	}

	// Clean up double slashes
	normalized = strings.ReplaceAll(normalized, "//", "/")

	return normalized
}

func detectJSONType(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	if strings.HasPrefix(text, "[") {
		return "array"
	}
	if strings.HasPrefix(text, "{") {
		return "object"
	}
	return "unknown"
}

func getHeader(headers []Header, name string) string {
	for _, h := range headers {
		if strings.EqualFold(h.Name, name) {
			return h.Value
		}
	}
	return ""
}

// GroupedEndpoint represents multiple calls to the same endpoint
type GroupedEndpoint struct {
	Method         string
	NormalizedPath string
	QueryParams    []string
	Examples       []Example
	StatusCodes    []int
}

func groupEndpoints(endpoints []Endpoint) map[string]*GroupedEndpoint {
	grouped := make(map[string]*GroupedEndpoint)

	for i := range endpoints {
		ep := &endpoints[i]
		key := fmt.Sprintf("%s %s", ep.Method, ep.NormalizedPath)

		if existing, ok := grouped[key]; ok {
			// Add example if different
			existing.Examples = append(existing.Examples, ep.Examples...)
			// Add status code if not seen
			if !slices.Contains(existing.StatusCodes, ep.StatusCode) {
				existing.StatusCodes = append(existing.StatusCodes, ep.StatusCode)
			}
			// Merge query params
			for _, qp := range ep.QueryParams {
				if !slices.Contains(existing.QueryParams, qp) {
					existing.QueryParams = append(existing.QueryParams, qp)
				}
			}
		} else {
			grouped[key] = &GroupedEndpoint{
				Method:         ep.Method,
				NormalizedPath: ep.NormalizedPath,
				QueryParams:    ep.QueryParams,
				Examples:       ep.Examples,
				StatusCodes:    []int{ep.StatusCode},
			}
		}
	}

	return grouped
}

func outputMarkdown(grouped map[string]*GroupedEndpoint) {
	// Group by service
	services := make(map[string][]*GroupedEndpoint)

	for _, ep := range grouped {
		service := extractService(ep.NormalizedPath)
		services[service] = append(services[service], ep)
	}

	// Sort services
	serviceNames := make([]string, 0, len(services))
	for name := range services {
		serviceNames = append(serviceNames, name)
	}
	sort.Strings(serviceNames)

	fmt.Println("# Discovered Garmin API Endpoints")
	fmt.Println()

	for _, serviceName := range serviceNames {
		eps := services[serviceName]

		// Sort endpoints within service
		sort.Slice(eps, func(i, j int) bool {
			return eps[i].NormalizedPath < eps[j].NormalizedPath
		})

		fmt.Printf("## %s\n\n", serviceName)
		fmt.Println("| Method | Endpoint | Query Params | Status |")
		fmt.Println("|--------|----------|--------------|--------|")

		for _, ep := range eps {
			queryParams := ""
			if len(ep.QueryParams) > 0 {
				queryParams = strings.Join(ep.QueryParams, ", ")
			}

			statusCodeStrs := make([]string, 0, len(ep.StatusCodes))
			for _, sc := range ep.StatusCodes {
				statusCodeStrs = append(statusCodeStrs, strconv.Itoa(sc))
			}
			statusCodes := strings.Join(statusCodeStrs, ", ")

			fmt.Printf("| %s | `%s` | %s | %s |\n",
				ep.Method, ep.NormalizedPath, queryParams, statusCodes)
		}
		fmt.Println()
	}
}

func outputJSON(grouped map[string]*GroupedEndpoint) {
	// Convert to sorted slice
	endpoints := make([]*GroupedEndpoint, 0, len(grouped))
	for _, ep := range grouped {
		endpoints = append(endpoints, ep)
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].NormalizedPath < endpoints[j].NormalizedPath
	})

	output := struct {
		TotalEndpoints int                `json:"total_endpoints"`
		Endpoints      []*GroupedEndpoint `json:"endpoints"`
	}{
		TotalEndpoints: len(endpoints),
		Endpoints:      endpoints,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output)
}

func extractService(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 0 {
		service := parts[0]
		if strings.HasSuffix(service, "-service") {
			return service
		}
		return service
	}
	return "unknown"
}

func compareWithEndpointsMD(filename string, discovered map[string]*GroupedEndpoint) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read %s: %v\n", filename, err)
		return
	}

	content := string(data)

	fmt.Println("\n---")
	fmt.Println("# Comparison with ENDPOINTS.md")
	fmt.Println()

	var missing []string
	var found []string

	for key, ep := range discovered {
		// Check if endpoint exists in ENDPOINTS.md
		// Look for the path pattern
		pathPattern := ep.NormalizedPath

		if strings.Contains(content, pathPattern) || strings.Contains(content, strings.ReplaceAll(pathPattern, "{id}", "{")) {
			found = append(found, key)
		} else {
			missing = append(missing, key)
		}
	}

	sort.Strings(missing)
	sort.Strings(found)

	fmt.Printf("## Found in ENDPOINTS.md (%d)\n\n", len(found))
	for _, ep := range found {
		fmt.Printf("- %s\n", ep)
	}

	fmt.Printf("\n## Missing from ENDPOINTS.md (%d)\n\n", len(missing))
	for _, key := range missing {
		ep := discovered[key]
		queryParams := ""
		if len(ep.QueryParams) > 0 {
			queryParams = "?" + strings.Join(ep.QueryParams, "&")
		}
		fmt.Printf("- `%s %s%s`\n", ep.Method, ep.NormalizedPath, queryParams)
	}
}

func outputSchemas(dir string, grouped map[string]*GroupedEndpoint) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		return
	}

	for _, ep := range grouped {
		if len(ep.Examples) == 0 {
			continue
		}

		// Use first example with content
		var responseBody string
		for _, ex := range ep.Examples {
			if ex.ResponseBody != "" {
				responseBody = ex.ResponseBody
				break
			}
		}

		if responseBody == "" {
			continue
		}

		// Try to parse and pretty-print JSON
		var parsed any
		if err := json.Unmarshal([]byte(responseBody), &parsed); err != nil {
			continue
		}

		// Generate filename from path
		safePath := strings.ReplaceAll(ep.NormalizedPath, "/", "_")
		safePath = strings.ReplaceAll(safePath, "{", "")
		safePath = strings.ReplaceAll(safePath, "}", "")
		safePath = strings.Trim(safePath, "_")

		filename := filepath.Join(dir, fmt.Sprintf("%s_%s.json", ep.Method, safePath))

		pretty, _ := json.MarshalIndent(parsed, "", "  ")
		if err := os.WriteFile(filename, pretty, 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not write %s: %v\n", filename, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Schemas written to %s/\n", dir)
}
