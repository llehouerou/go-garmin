// endpoint/cli_test.go
package endpoint

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestCLIGenerator_GenerateCommands_SimpleCommand(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		CLICommand: "sleep",
		Short:      "Get sleep data",
		Params: []Param{
			{Name: "date", Type: ParamTypeDate, Description: "The date"},
		},
		Handler: func(_ context.Context, _ any, args *HandlerArgs) (any, error) {
			return map[string]string{"date": args.Date("date").Format("2006-01-02")}, nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]
	if cmd.Use != "sleep [date]" {
		t.Errorf("Use = %q, want 'sleep [date]'", cmd.Use)
	}
	if cmd.Short != "Get sleep data" {
		t.Errorf("Short = %q, want 'Get sleep data'", cmd.Short)
	}
}

func TestCLIGenerator_GenerateCommands_GroupedCommands(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:          "ListWorkouts",
		CLICommand:    "workouts",
		CLISubcommand: "list",
		Short:         "List workouts",
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return []string{}, nil
		},
	})
	r.Register(Endpoint{
		Name:          "GetWorkout",
		CLICommand:    "workouts",
		CLISubcommand: "get",
		Short:         "Get workout",
		Params: []Param{
			{Name: "id", Type: ParamTypeInt, Required: true, Description: "Workout ID"},
		},
		Handler: func(_ context.Context, _ any, args *HandlerArgs) (any, error) {
			return map[string]int{"id": args.Int("id")}, nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	if len(commands) != 1 {
		t.Fatalf("expected 1 parent command, got %d", len(commands))
	}

	parent := commands[0]
	if parent.Use != "workouts" {
		t.Errorf("parent Use = %q, want 'workouts'", parent.Use)
	}

	subcommands := parent.Commands()
	if len(subcommands) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(subcommands))
	}
}

func TestCLIGenerator_ExecuteCommand(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetSleep",
		CLICommand: "sleep",
		Short:      "Get sleep data",
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return map[string]string{"result": "ok"}, nil
		},
	})

	var output bytes.Buffer
	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	gen.SetOutput(&output)
	commands := gen.GenerateCommands()

	root := &cobra.Command{Use: "test"}
	for _, cmd := range commands {
		root.AddCommand(cmd)
	}

	root.SetArgs([]string{"sleep"})
	err := root.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if output.Len() == 0 {
		t.Error("expected output, got empty")
	}
}

func TestCLIGenerator_SkipsEndpointsWithoutCLICommand(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:    "GetSleep",
		MCPTool: "get_sleep",
		// CLICommand not set
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return struct{}{}, nil
		},
	})

	gen := NewCLIGenerator(r)
	commands := gen.GenerateCommands()

	if len(commands) != 0 {
		t.Fatalf("expected 0 commands, got %d", len(commands))
	}
}

func TestCLIGenerator_RequiredParam(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetActivity",
		CLICommand: "activity",
		Short:      "Get activity",
		Params: []Param{
			{Name: "id", Type: ParamTypeInt, Required: true, Description: "Activity ID"},
		},
		Handler: func(_ context.Context, _ any, args *HandlerArgs) (any, error) {
			return map[string]int{"id": args.Int("id")}, nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]
	if cmd.Use != "activity <id>" {
		t.Errorf("Use = %q, want 'activity <id>'", cmd.Use)
	}
}

func TestCLIGenerator_DateRangeFlags(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "GetWeight",
		CLICommand: "weight",
		Short:      "Get weight",
		Params: []Param{
			{Name: "range", Type: ParamTypeDateRange, Description: "Date range"},
		},
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return struct{}{}, nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]
	startFlag := cmd.Flags().Lookup("start")
	if startFlag == nil {
		t.Error("expected 'start' flag to exist")
	}
	endFlag := cmd.Flags().Lookup("end")
	if endFlag == nil {
		t.Error("expected 'end' flag to exist")
	}
}

func TestCLIGenerator_RawOutput_WritesBytes(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "DownloadGPX",
		CLICommand: "download",
		Short:      "Download GPX",
		RawOutput:  true,
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return []byte("<gpx>test data</gpx>"), nil
		},
	})

	var output bytes.Buffer
	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	gen.SetOutput(&output)
	commands := gen.GenerateCommands()

	root := &cobra.Command{Use: "test"}
	for _, cmd := range commands {
		root.AddCommand(cmd)
	}

	root.SetArgs([]string{"download"})
	err := root.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Should write raw bytes, not JSON
	if output.String() != "<gpx>test data</gpx>" {
		t.Errorf("output = %q, want %q", output.String(), "<gpx>test data</gpx>")
	}
}

func TestCLIGenerator_RawOutput_OutputFlag(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "DownloadGPX",
		CLICommand: "download",
		Short:      "Download GPX",
		RawOutput:  true,
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return []byte("file content"), nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	// Verify --output/-o flag exists
	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}
	cmd := commands[0]
	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("expected 'output' flag to exist")
	}
	if outputFlag.Shorthand != "o" {
		t.Errorf("output flag shorthand = %q, want %q", outputFlag.Shorthand, "o")
	}

	// Test writing to file
	tmpFile := t.TempDir() + "/test.gpx"

	root := &cobra.Command{Use: "test"}
	root.AddCommand(cmd)
	root.SetArgs([]string{"download", "--output", tmpFile})
	err := root.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if string(data) != "file content" {
		t.Errorf("file content = %q, want %q", string(data), "file content")
	}
}

func TestCLIGenerator_Aliases(t *testing.T) {
	r := NewRegistry()
	r.Register(Endpoint{
		Name:       "ListDevices",
		CLICommand: "devices",
		CLIAliases: []string{"device", "dev"},
		Short:      "List devices",
		Handler: func(_ context.Context, _ any, _ *HandlerArgs) (any, error) {
			return struct{}{}, nil
		},
	})

	gen := NewCLIGenerator(r)
	gen.SetClient(nil)
	commands := gen.GenerateCommands()

	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]
	if len(cmd.Aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(cmd.Aliases))
	}
	if cmd.Aliases[0] != "device" || cmd.Aliases[1] != "dev" {
		t.Errorf("Aliases = %v, want [device dev]", cmd.Aliases)
	}
}
