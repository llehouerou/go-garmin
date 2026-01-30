// endpoint/validator.go
package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidatorConfig configures the endpoint validator.
type ValidatorConfig struct {
	CassetteDir           string
	SkipOrphanedCassettes bool
}

// Validator checks endpoints for completeness.
type Validator struct {
	registry *Registry
	config   ValidatorConfig
}

// NewValidator creates a new endpoint validator.
func NewValidator(registry *Registry, config ValidatorConfig) *Validator {
	return &Validator{registry: registry, config: config}
}

// Validate checks all endpoints and returns any errors found.
func (v *Validator) Validate() []string {
	errors := make([]string, 0, len(v.registry.All()))

	for _, ep := range v.registry.All() {
		errors = append(errors, v.validateEndpoint(ep)...)
	}

	if !v.config.SkipOrphanedCassettes {
		errors = append(errors, v.checkOrphanedCassettes()...)
	}

	return errors
}

func (v *Validator) validateEndpoint(ep *Endpoint) []string {
	var errors []string

	// Must have a handler
	if ep.Handler == nil {
		errors = append(errors, ep.Name+": missing Handler")
	}

	// Must have a cassette (or explicitly "none" for static endpoints)
	if ep.Cassette == "" {
		errors = append(errors, ep.Name+": missing Cassette")
	} else if ep.Cassette != "none" {
		cassettePath := filepath.Join(v.config.CassetteDir, ep.Cassette+".yaml")
		if _, err := os.Stat(cassettePath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("%s: cassette file not found: %s", ep.Name, cassettePath))
		}
	}

	// Must have CLI or MCP (or both)
	if ep.CLICommand == "" && ep.MCPTool == "" {
		errors = append(errors, ep.Name+": must have CLICommand or MCPTool (or both)")
	}

	// Must have descriptions
	if ep.Short == "" {
		errors = append(errors, ep.Name+": missing Short description")
	}
	if ep.Long == "" {
		errors = append(errors, ep.Name+": missing Long description")
	}

	// Path must be set
	if ep.Path == "" {
		errors = append(errors, ep.Name+": missing Path")
	}

	// HTTPMethod must be valid
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}
	if !validMethods[ep.HTTPMethod] {
		errors = append(errors, fmt.Sprintf("%s: invalid HTTPMethod: %s", ep.Name, ep.HTTPMethod))
	}

	// POST/PUT should have Body config
	if (ep.HTTPMethod == "POST" || ep.HTTPMethod == "PUT") && ep.Body == nil {
		errors = append(errors, fmt.Sprintf("%s: %s endpoint should have Body config", ep.Name, ep.HTTPMethod))
	}

	// Params must have descriptions
	for _, p := range ep.Params {
		if p.Description == "" {
			errors = append(errors, fmt.Sprintf("%s: param %s missing description", ep.Name, p.Name))
		}
	}

	// DependsOn validation
	if ep.DependsOn != "" {
		if v.registry.ByName(ep.DependsOn) == nil {
			errors = append(errors, fmt.Sprintf("%s: DependsOn references unknown endpoint: %s", ep.Name, ep.DependsOn))
		}
		if ep.ArgProvider == nil {
			errors = append(errors, ep.Name+": has DependsOn but missing ArgProvider")
		}
	}

	return errors
}

func (v *Validator) checkOrphanedCassettes() []string {
	var errors []string

	usedCassettes := make(map[string]bool)
	for _, ep := range v.registry.All() {
		if ep.Cassette != "" {
			usedCassettes[ep.Cassette] = true
		}
	}

	files, err := filepath.Glob(filepath.Join(v.config.CassetteDir, "*.yaml"))
	if err != nil {
		return errors
	}

	for _, f := range files {
		name := filepath.Base(f)
		name = name[:len(name)-5] // Remove .yaml

		if name == "auth" {
			continue // Skip auth cassette
		}

		if !usedCassettes[name] {
			errors = append(errors, "orphaned cassette (no endpoint references it): "+name)
		}
	}

	return errors
}
