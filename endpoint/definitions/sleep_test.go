// endpoint/definitions/sleep_test.go
package definitions

import (
	"testing"

	"github.com/llehouerou/go-garmin/endpoint"
)

func TestSleepEndpoints_Registered(t *testing.T) {
	if len(SleepEndpoints) == 0 {
		t.Fatal("SleepEndpoints should not be empty")
	}

	ep := SleepEndpoints[0]
	if ep.Name != "GetDailySleep" {
		t.Errorf("Name = %q, want GetDailySleep", ep.Name)
	}
	if ep.CLICommand != "sleep" {
		t.Errorf("CLICommand = %q, want sleep", ep.CLICommand)
	}
	if ep.MCPTool != "get_sleep" {
		t.Errorf("MCPTool = %q, want get_sleep", ep.MCPTool)
	}
	if ep.Handler == nil {
		t.Error("Handler should not be nil")
	}
}

func TestSleepEndpoints_HasDateParam(t *testing.T) {
	ep := SleepEndpoints[0]

	if len(ep.Params) != 1 {
		t.Fatalf("Params len = %d, want 1", len(ep.Params))
	}

	p := ep.Params[0]
	if p.Name != "date" {
		t.Errorf("Param name = %q, want date", p.Name)
	}
	if p.Type != endpoint.ParamTypeDate {
		t.Errorf("Param type = %v, want ParamTypeDate", p.Type)
	}
}
