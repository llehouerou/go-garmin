# Claude Code Instructions

## Verification

Always run `make check` after making changes to verify the code compiles, passes linting, and tests pass.

## Adding New Endpoints

The project uses a **Declarative Endpoint Registry** system. Adding a new endpoint requires updating only ONE file - the endpoint definition. CLI commands and MCP tools are automatically generated.

### Quick Start

1. Add endpoint definition to `endpoint/definitions/<service>.go`
2. Register in `endpoint/definitions/register.go` (if new file)
3. Run `make check` and `make validate-endpoints`

### Step-by-Step Guide

#### 1. Create or Update Endpoint Definition

Add your endpoint to the appropriate file in `endpoint/definitions/`:

```go
// endpoint/definitions/<service>.go
package definitions

import (
    "context"
    "fmt"

    "github.com/llehouerou/go-garmin"
    "github.com/llehouerou/go-garmin/endpoint"
)

var ServiceEndpoints = []endpoint.Endpoint{
    {
        Name:          "GetData",           // Unique identifier
        Service:       "ServiceName",       // Service group
        Cassette:      "cassette_name",     // VCR cassette for testing
        Path:          "/api/path",         // API endpoint path
        HTTPMethod:    "GET",               // HTTP method
        Params: []endpoint.Param{
            {Name: "date", Type: endpoint.ParamTypeDate, Required: false, Description: "Date (YYYY-MM-DD)"},
        },
        CLICommand:    "service",           // CLI command (garmin service ...)
        CLISubcommand: "subcommand",        // CLI subcommand (garmin service subcommand)
        MCPTool:       "get_data",          // MCP tool name
        Short:         "Short description",
        Long:          "Longer description for help text",
        Handler: func(ctx context.Context, c any, args *endpoint.HandlerArgs) (any, error) {
            client, ok := c.(*garmin.Client)
            if !ok {
                return nil, fmt.Errorf("handler received invalid client type: %T", c)
            }
            return client.Service.GetData(ctx, args.Date("date"))
        },
    },
}
```

#### 2. Parameter Types

Available parameter types:
- `endpoint.ParamTypeString` - String parameter
- `endpoint.ParamTypeInt` - Integer parameter
- `endpoint.ParamTypeDate` - Date parameter (YYYY-MM-DD, defaults to today)
- `endpoint.ParamTypeDateRange` - Date range (uses --start and --end flags)
- `endpoint.ParamTypeBool` - Boolean parameter

#### 3. Register Endpoints (if new file)

If you created a new definition file, add it to `endpoint/definitions/register.go`:

```go
func RegisterAll(r *endpoint.Registry) {
    // ... existing registrations ...
    for i := range ServiceEndpoints {
        r.Register(ServiceEndpoints[i])
    }
}
```

#### 4. Validate Completeness

Run the endpoint validator to check for missing fields:

```bash
make validate-endpoints
```

This checks:
- Handler is defined
- Cassette exists (or is "none")
- Short and Long descriptions are set
- Path is specified
- CLI command or MCP tool is defined
- Params have descriptions

#### 5. Create Service Implementation (if new service)

If this is a new Garmin service, create the service file:

```go
// service_<name>.go
type ResponseType struct {
    Field1 string `json:"field1"`
    // ... fields
    raw json.RawMessage
}

func (r *ResponseType) RawJSON() json.RawMessage {
    return r.raw
}

func (s *ServiceName) GetData(ctx context.Context, date time.Time) (*ResponseType, error) {
    path := fmt.Sprintf("/api/path/%s", date.Format("2006-01-02"))
    resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
    // ... handle response
}
```

#### 6. Record Cassette (for tests)

```bash
go run ./cmd/record-fixtures -email=USER -password=PASS
```

#### 7. Verify and Commit

```bash
make check
make validate-endpoints
git add -A
git commit -m "feat: add ServiceName.GetData endpoint"
```

## Endpoint Definition Fields

| Field | Required | Description |
|-------|----------|-------------|
| `Name` | Yes | Unique endpoint identifier |
| `Service` | Yes | Service group name |
| `Cassette` | Yes | VCR cassette name for testing |
| `Path` | Yes | API endpoint path |
| `HTTPMethod` | Yes | HTTP method (GET, POST, PUT, DELETE) |
| `Params` | No | List of parameters |
| `CLICommand` | No* | CLI command name |
| `CLISubcommand` | No | CLI subcommand (for grouped commands) |
| `MCPTool` | No* | MCP tool name |
| `Short` | Yes | Short description |
| `Long` | Yes | Long description |
| `Handler` | Yes | Handler function |
| `DependsOn` | No | Name of endpoint this depends on (for fixture recording) |
| `ArgProvider` | No | Function to extract args from DependsOn result |

*At least one of CLICommand or MCPTool should be set.

## File Structure

```
├── endpoint/
│   ├── endpoint.go          # Core types and registry
│   ├── cli.go               # CLI generator
│   ├── mcp.go               # MCP generator
│   ├── validator.go         # Endpoint validator
│   ├── recorder.go          # Fixture recorder
│   └── definitions/         # Endpoint definitions
│       ├── register.go      # Registration of all endpoints
│       ├── sleep.go
│       ├── wellness.go
│       ├── activities.go
│       └── ...
├── service_<name>.go        # Service implementation
├── cmd/garmin/
│   ├── root.go              # CLI root (uses CLIGenerator)
│   ├── mcp.go               # MCP server (uses MCPGenerator)
│   └── registry.go          # Global endpoint registry
└── testdata/cassettes/      # VCR cassettes
```

## Code Style

- Use pointers for optional fields (`*int`, `*float64`, `*string`)
- Always include `raw json.RawMessage` and `RawJSON()` method
- Use `time.Time` helpers for timestamp conversions
- Follow naming conventions: `GetDaily`, `List`, `Get`, `GetRange`
- For handlers that don't use args, use `_ *endpoint.HandlerArgs`

## Makefile Commands

```bash
make check              # Run imports, lint, and tests
make validate-endpoints # Validate endpoint completeness
make lint               # Run linter only
make test               # Run tests only
```
