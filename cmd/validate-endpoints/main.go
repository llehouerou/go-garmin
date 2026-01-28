// Command validate-endpoints validates all endpoint definitions in the registry.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/llehouerou/go-garmin/endpoint"
	"github.com/llehouerou/go-garmin/endpoint/definitions"
	"github.com/llehouerou/go-garmin/testutil"
)

func main() {
	skipOrphans := flag.Bool("skip-orphaned-cassettes", false, "Skip checking for orphaned cassettes")
	flag.Parse()

	registry := endpoint.NewRegistry()
	definitions.RegisterAll(registry)

	validator := endpoint.NewValidator(registry, endpoint.ValidatorConfig{
		CassetteDir:           testutil.CassetteDir,
		SkipOrphanedCassettes: *skipOrphans,
	})

	errors := validator.Validate()

	if len(errors) > 0 {
		fmt.Println("Endpoint validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(1)
	}

	printSummary(registry)
}

func printSummary(registry *endpoint.Registry) {
	var total, withCLI, withMCP, withBoth int
	total = len(registry.All())

	for _, ep := range registry.All() {
		hasCLI := ep.CLICommand != ""
		hasMCP := ep.MCPTool != ""
		if hasCLI {
			withCLI++
		}
		if hasMCP {
			withMCP++
		}
		if hasCLI && hasMCP {
			withBoth++
		}
	}

	fmt.Printf("All %d endpoints validated successfully\n\n", total)
	fmt.Println("Endpoint Summary:")
	fmt.Printf("  Total endpoints:    %d\n", total)
	fmt.Printf("  With CLI command:   %d\n", withCLI)
	fmt.Printf("  With MCP tool:      %d\n", withMCP)
	fmt.Printf("  With both:          %d\n", withBoth)
}
