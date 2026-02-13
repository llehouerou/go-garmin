// endpoint/mcp.go
package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPGenerator creates MCP tools from the registry.
type MCPGenerator struct {
	registry *Registry
	client   any
}

// NewMCPGenerator creates a new MCP generator.
func NewMCPGenerator(registry *Registry, client any) *MCPGenerator {
	return &MCPGenerator{registry: registry, client: client}
}

// RegisterTools adds all endpoint tools to the MCP server.
func (g *MCPGenerator) RegisterTools(s *server.MCPServer) {
	for _, ep := range g.registry.endpoints {
		if ep.MCPTool == "" || ep.RawOutput {
			continue
		}
		g.registerTool(s, ep)
	}
}

func (g *MCPGenerator) registerTool(s *server.MCPServer, ep *Endpoint) {
	opts := []mcp.ToolOption{
		mcp.WithDescription(ep.Long),
	}

	for _, p := range ep.Params {
		opts = append(opts, g.paramToMCPOptions(p)...)
	}

	if ep.Body != nil {
		opts = append(opts, g.bodyToMCPOption(ep.Body))
	}

	tool := mcp.NewTool(ep.MCPTool, opts...)
	s.AddTool(tool, g.createHandler(ep))
}

func (g *MCPGenerator) paramToMCPOptions(p Param) []mcp.ToolOption {
	var opts []mcp.ToolOption

	switch p.Type {
	case ParamTypeString:
		if p.Required {
			opts = append(opts, mcp.WithString(p.Name, mcp.Required(), mcp.Description(p.Description)))
		} else {
			opts = append(opts, mcp.WithString(p.Name, mcp.Description(p.Description)))
		}

	case ParamTypeInt:
		if p.Required {
			opts = append(opts, mcp.WithNumber(p.Name, mcp.Required(), mcp.Description(p.Description)))
		} else {
			opts = append(opts, mcp.WithNumber(p.Name, mcp.Description(p.Description)))
		}

	case ParamTypeDate:
		desc := p.Description
		if !strings.Contains(strings.ToLower(desc), "yyyy-mm-dd") {
			desc += " (YYYY-MM-DD)"
		}
		if p.Required {
			opts = append(opts, mcp.WithString(p.Name, mcp.Required(), mcp.Description(desc)))
		} else {
			opts = append(opts, mcp.WithString(p.Name, mcp.Description(desc)))
		}

	case ParamTypeDateRange:
		opts = append(opts,
			mcp.WithString("start", mcp.Description("Start date (YYYY-MM-DD)")),
			mcp.WithString("end", mcp.Description("End date (YYYY-MM-DD)")),
		)

	case ParamTypeBool:
		opts = append(opts, mcp.WithBoolean(p.Name, mcp.Description(p.Description)))
	}

	return opts
}

func (g *MCPGenerator) bodyToMCPOption(body *BodyConfig) mcp.ToolOption {
	paramName := g.typeToParamName(body.Type)
	desc := body.Description
	if body.Example != "" {
		desc += "\n\nExample:\n" + body.Example
	}
	return mcp.WithString(paramName, mcp.Required(), mcp.Description(desc))
}

func (g *MCPGenerator) typeToParamName(t reflect.Type) string {
	if t == nil {
		return "body"
	}
	name := t.Name()
	if name == "" {
		return "body"
	}
	// Convert to camelCase (e.g., "Workout" -> "workout")
	return strings.ToLower(name[:1]) + name[1:]
}

func (g *MCPGenerator) createHandler(ep *Endpoint) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := g.parseRequest(request, ep)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if ep.Body != nil {
			body, err := g.parseBodyFromRequest(request, ep.Body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			args.Body = body
		}

		result, err := ep.Handler(ctx, g.client, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return g.jsonResult(result), nil
	}
}

func (g *MCPGenerator) parseRequest(request mcp.CallToolRequest, ep *Endpoint) (*HandlerArgs, error) {
	args := &HandlerArgs{Params: make(map[string]any)}

	for _, p := range ep.Params {
		switch p.Type {
		case ParamTypeString:
			if v, err := request.RequireString(p.Name); err == nil {
				args.Params[p.Name] = v
			} else if p.Required {
				return nil, fmt.Errorf("missing required parameter: %s", p.Name)
			}

		case ParamTypeInt:
			if v, err := request.RequireFloat(p.Name); err == nil {
				args.Params[p.Name] = int(v)
			} else if p.Required {
				return nil, fmt.Errorf("missing required parameter: %s", p.Name)
			}

		case ParamTypeDate:
			if v, err := request.RequireString(p.Name); err == nil && v != "" {
				t, err := time.Parse("2006-01-02", v)
				if err != nil {
					return nil, fmt.Errorf("invalid date format for %s: %w", p.Name, err)
				}
				args.Params[p.Name] = t
			} else {
				args.Params[p.Name] = time.Now()
			}

		case ParamTypeDateRange:
			if start, err := request.RequireString("start"); err == nil && start != "" {
				t, err := time.Parse("2006-01-02", start)
				if err != nil {
					return nil, fmt.Errorf("invalid start date: %w", err)
				}
				args.Params["start"] = t
			}
			if end, err := request.RequireString("end"); err == nil && end != "" {
				t, err := time.Parse("2006-01-02", end)
				if err != nil {
					return nil, fmt.Errorf("invalid end date: %w", err)
				}
				args.Params["end"] = t
			}

		case ParamTypeBool:
			if v, err := request.RequireBool(p.Name); err == nil {
				args.Params[p.Name] = v
			}
		}
	}

	return args, nil
}

func (g *MCPGenerator) parseBodyFromRequest(request mcp.CallToolRequest, body *BodyConfig) (any, error) {
	paramName := g.typeToParamName(body.Type)
	jsonStr, err := request.RequireString(paramName)
	if err != nil {
		return nil, fmt.Errorf("missing required parameter: %s", paramName)
	}

	bodyPtr := reflect.New(body.Type).Interface()
	if err := json.Unmarshal([]byte(jsonStr), bodyPtr); err != nil {
		return nil, fmt.Errorf("invalid JSON for %s: %w", paramName, err)
	}

	return bodyPtr, nil
}

func (g *MCPGenerator) jsonResult(v any) *mcp.CallToolResult {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error())
	}
	return mcp.NewToolResultText(string(data))
}
