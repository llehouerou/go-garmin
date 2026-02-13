// endpoint/cli.go
package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// CLIGenerator creates cobra commands from the registry.
type CLIGenerator struct {
	registry *Registry
	client   any
	output   io.Writer
}

// NewCLIGenerator creates a new CLI generator.
func NewCLIGenerator(registry *Registry) *CLIGenerator {
	return &CLIGenerator{
		registry: registry,
		output:   os.Stdout,
	}
}

// SetClient sets the client for handlers.
func (g *CLIGenerator) SetClient(client any) {
	g.client = client
}

// SetOutput sets the output writer (for testing).
func (g *CLIGenerator) SetOutput(w io.Writer) {
	g.output = w
}

// GenerateCommands creates all cobra commands from the registry.
func (g *CLIGenerator) GenerateCommands() []*cobra.Command {
	commandGroups := make(map[string][]*Endpoint)
	var simpleCommands []*Endpoint

	for _, ep := range g.registry.endpoints {
		if ep.CLICommand == "" {
			continue
		}
		if ep.CLISubcommand == "" {
			simpleCommands = append(simpleCommands, ep)
		} else {
			commandGroups[ep.CLICommand] = append(commandGroups[ep.CLICommand], ep)
		}
	}

	commands := make([]*cobra.Command, 0, len(simpleCommands)+len(commandGroups))

	for _, ep := range simpleCommands {
		commands = append(commands, g.createCommand(ep))
	}

	for parentName, endpoints := range commandGroups {
		parent := &cobra.Command{
			Use:   parentName,
			Short: parentName + " commands",
		}
		for _, ep := range endpoints {
			parent.AddCommand(g.createSubcommand(ep))
		}
		commands = append(commands, parent)
	}

	return commands
}

func (g *CLIGenerator) createCommand(ep *Endpoint) *cobra.Command {
	cmd := &cobra.Command{
		Use:     g.buildUse(ep.CLICommand, ep.Params),
		Short:   ep.Short,
		Long:    ep.Long,
		Aliases: ep.CLIAliases,
		RunE:    g.createRunFunc(ep),
	}
	g.addFlags(cmd, ep)
	return cmd
}

func (g *CLIGenerator) createSubcommand(ep *Endpoint) *cobra.Command {
	cmd := &cobra.Command{
		Use:   g.buildUse(ep.CLISubcommand, ep.Params),
		Short: ep.Short,
		Long:  ep.Long,
		RunE:  g.createRunFunc(ep),
	}
	g.addFlags(cmd, ep)
	return cmd
}

func (g *CLIGenerator) buildUse(base string, params []Param) string {
	var sb strings.Builder
	sb.WriteString(base)
	for _, p := range params {
		if p.Type == ParamTypeDateRange {
			continue // DateRange uses flags, not positional
		}
		sb.WriteString(" ")
		if p.Required {
			sb.WriteString("<")
			sb.WriteString(p.Name)
			sb.WriteString(">")
		} else {
			sb.WriteString("[")
			sb.WriteString(p.Name)
			sb.WriteString("]")
		}
	}
	return sb.String()
}

func (g *CLIGenerator) addFlags(cmd *cobra.Command, ep *Endpoint) {
	for _, p := range ep.Params {
		if p.Required {
			continue // Required params are positional
		}
		switch p.Type {
		case ParamTypeString:
			cmd.Flags().String(p.Name, "", p.Description)
		case ParamTypeInt:
			cmd.Flags().Int(p.Name, 0, p.Description)
		case ParamTypeDate:
			// Date is positional optional, no flag needed
		case ParamTypeDateRange:
			cmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
			cmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
		case ParamTypeBool:
			cmd.Flags().Bool(p.Name, false, p.Description)
		}
	}

	if ep.Body != nil {
		cmd.Flags().StringP("file", "f", "", "Read JSON body from file")
		cmd.Flags().String("json", "", "JSON body as string")
	}

	if ep.RawOutput {
		cmd.Flags().StringP("output", "o", "", "Output file path")
	}
}

func (g *CLIGenerator) createRunFunc(ep *Endpoint) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		handlerArgs, err := g.parseArgs(cmd, args, ep)
		if err != nil {
			return err
		}

		if ep.Body != nil {
			body, err := g.parseBody(cmd, ep)
			if err != nil {
				return err
			}
			handlerArgs.Body = body
		}

		result, err := ep.Handler(cmd.Context(), g.client, handlerArgs)
		if err != nil {
			return err
		}

		if ep.RawOutput {
			data, ok := result.([]byte)
			if !ok {
				return fmt.Errorf("raw output handler must return []byte, got %T", result)
			}
			if outputPath, _ := cmd.Flags().GetString("output"); outputPath != "" {
				return os.WriteFile(outputPath, data, 0o600)
			}
			_, err := g.output.Write(data)
			return err
		}

		return g.printJSON(result)
	}
}

func (g *CLIGenerator) parseArgs(cmd *cobra.Command, args []string, ep *Endpoint) (*HandlerArgs, error) {
	handlerArgs := &HandlerArgs{Params: make(map[string]any)}
	argIndex := 0

	for _, p := range ep.Params {
		switch p.Type {
		case ParamTypeDateRange:
			// Always from flags
			if start, _ := cmd.Flags().GetString("start"); start != "" {
				t, err := time.Parse("2006-01-02", start)
				if err != nil {
					return nil, fmt.Errorf("invalid start date: %w", err)
				}
				handlerArgs.Params["start"] = t
			}
			if end, _ := cmd.Flags().GetString("end"); end != "" {
				t, err := time.Parse("2006-01-02", end)
				if err != nil {
					return nil, fmt.Errorf("invalid end date: %w", err)
				}
				handlerArgs.Params["end"] = t
			}

		case ParamTypeDate:
			// Positional arg or default to today
			if argIndex < len(args) && args[argIndex] != "" {
				t, err := time.Parse("2006-01-02", args[argIndex])
				if err != nil {
					return nil, fmt.Errorf("invalid date %s: %w", p.Name, err)
				}
				handlerArgs.Params[p.Name] = t
				argIndex++
			} else {
				handlerArgs.Params[p.Name] = time.Now()
			}

		case ParamTypeInt:
			if p.Required {
				if argIndex >= len(args) {
					return nil, fmt.Errorf("missing required argument: %s", p.Name)
				}
				var v int
				if _, err := fmt.Sscanf(args[argIndex], "%d", &v); err != nil {
					return nil, fmt.Errorf("invalid integer for %s: %w", p.Name, err)
				}
				handlerArgs.Params[p.Name] = v
				argIndex++
			} else {
				if v, _ := cmd.Flags().GetInt(p.Name); v != 0 {
					handlerArgs.Params[p.Name] = v
				}
			}

		case ParamTypeString:
			if p.Required {
				if argIndex >= len(args) {
					return nil, fmt.Errorf("missing required argument: %s", p.Name)
				}
				handlerArgs.Params[p.Name] = args[argIndex]
				argIndex++
			} else {
				if v, _ := cmd.Flags().GetString(p.Name); v != "" {
					handlerArgs.Params[p.Name] = v
				}
			}

		case ParamTypeBool:
			if v, _ := cmd.Flags().GetBool(p.Name); v {
				handlerArgs.Params[p.Name] = v
			}
		}
	}

	return handlerArgs, nil
}

var errNoJSONBody = errors.New("no JSON body provided (use --json, --file, or pipe to stdin)")

func (g *CLIGenerator) parseBody(cmd *cobra.Command, ep *Endpoint) (any, error) {
	jsonData, err := g.readJSONData(cmd)
	if err != nil {
		return nil, err
	}

	if len(jsonData) == 0 {
		return nil, errNoJSONBody
	}

	bodyPtr := reflect.New(ep.Body.Type).Interface()
	if err := json.Unmarshal(jsonData, bodyPtr); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return bodyPtr, nil
}

func (g *CLIGenerator) readJSONData(cmd *cobra.Command) ([]byte, error) {
	if jsonStr, _ := cmd.Flags().GetString("json"); jsonStr != "" {
		return []byte(jsonStr), nil
	}

	if file, _ := cmd.Flags().GetString("file"); file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return data, nil
	}

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read stdin: %w", err)
		}
		return data, nil
	}

	return nil, nil
}

func (g *CLIGenerator) printJSON(v any) error {
	enc := json.NewEncoder(g.output)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
