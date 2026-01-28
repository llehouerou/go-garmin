// endpoint/recorder.go
package endpoint

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// RecorderConfig configures the fixture recorder.
type RecorderConfig struct {
	Session      []byte
	Date         time.Time
	NewRecorder  func(name string) (*recorder.Recorder, error)
	HTTPClient   func(rec *recorder.Recorder) *http.Client
	ClientLoader func(httpClient *http.Client, session []byte) (any, error)
}

// FixtureRecorder records API interactions for all registered endpoints.
type FixtureRecorder struct {
	registry     *Registry
	session      []byte
	date         time.Time
	newRecorder  func(name string) (*recorder.Recorder, error)
	httpClient   func(rec *recorder.Recorder) *http.Client
	clientLoader func(httpClient *http.Client, session []byte) (any, error)
}

// NewFixtureRecorder creates a new fixture recorder.
func NewFixtureRecorder(registry *Registry, cfg RecorderConfig) *FixtureRecorder {
	return &FixtureRecorder{
		registry:     registry,
		session:      cfg.Session,
		date:         cfg.Date,
		newRecorder:  cfg.NewRecorder,
		httpClient:   cfg.HTTPClient,
		clientLoader: cfg.ClientLoader,
	}
}

// RecordAll records fixtures for all endpoints in the registry.
func (f *FixtureRecorder) RecordAll(ctx context.Context) error {
	byCassette := make(map[string][]*Endpoint)
	for _, ep := range f.registry.endpoints {
		if ep.Cassette == "" {
			continue
		}
		byCassette[ep.Cassette] = append(byCassette[ep.Cassette], ep)
	}

	for cassette, endpoints := range byCassette {
		if err := f.recordCassette(ctx, cassette, endpoints); err != nil {
			return fmt.Errorf("%s: %w", cassette, err)
		}
	}

	return nil
}

// RecordCassette records a specific cassette.
func (f *FixtureRecorder) RecordCassette(ctx context.Context, cassetteName string) error {
	endpoints := f.endpointsForCassette(cassetteName)

	if len(endpoints) == 0 {
		return fmt.Errorf("no endpoints found for cassette: %s", cassetteName)
	}

	return f.recordCassette(ctx, cassetteName, endpoints)
}

func (f *FixtureRecorder) recordCassette(ctx context.Context, cassette string, endpoints []*Endpoint) error {
	fmt.Printf("Recording cassette: %s\n", cassette)

	rec, err := f.newRecorder(cassette)
	if err != nil {
		return err
	}
	defer func() { _ = rec.Stop() }()

	httpClient := f.httpClient(rec)
	client, err := f.clientLoader(httpClient, f.session)
	if err != nil {
		return fmt.Errorf("failed to load client: %w", err)
	}

	results := make(map[string]any)
	sorted := f.sortByDependencies(endpoints)

	for _, ep := range sorted {
		args := f.buildDefaultArgs(ep)

		if ep.DependsOn != "" {
			depResult, ok := results[ep.DependsOn]
			if ok && ep.ArgProvider != nil {
				extraArgs := ep.ArgProvider(depResult)
				if extraArgs == nil {
					fmt.Printf("  Skipping %s: no data from %s\n", ep.Name, ep.DependsOn)
					continue
				}
				maps.Copy(args.Params, extraArgs)
			}
		}

		fmt.Printf("  Recording %s...\n", ep.Name)
		result, err := ep.Handler(ctx, client, args)
		if err != nil {
			fmt.Printf("  Warning: %s: %v\n", ep.Name, err)
			continue
		}
		results[ep.Name] = result
	}

	return nil
}

func (f *FixtureRecorder) buildDefaultArgs(ep *Endpoint) *HandlerArgs {
	args := &HandlerArgs{Params: make(map[string]any)}

	for _, p := range ep.Params {
		switch p.Type {
		case ParamTypeDate:
			args.Params[p.Name] = f.date
		case ParamTypeDateRange:
			args.Params["start"] = f.date.AddDate(0, 0, -7)
			args.Params["end"] = f.date
		case ParamTypeInt:
			switch p.Name {
			case "limit":
				args.Params[p.Name] = 10
			default:
				args.Params[p.Name] = 0
			}
		case ParamTypeString:
			args.Params[p.Name] = ""
		case ParamTypeBool:
			args.Params[p.Name] = false
		}
	}

	return args
}

// ListCassettes returns all unique cassette names in the registry.
func (f *FixtureRecorder) ListCassettes() []string {
	seen := make(map[string]bool)
	var cassettes []string
	for _, ep := range f.registry.endpoints {
		if ep.Cassette != "" && !seen[ep.Cassette] {
			seen[ep.Cassette] = true
			cassettes = append(cassettes, ep.Cassette)
		}
	}
	return cassettes
}

func (f *FixtureRecorder) endpointsForCassette(cassetteName string) []*Endpoint {
	var endpoints []*Endpoint
	for _, ep := range f.registry.endpoints {
		if ep.Cassette == cassetteName {
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints
}

func (f *FixtureRecorder) sortByDependencies(endpoints []*Endpoint) []*Endpoint {
	// Build dependency graph
	byName := make(map[string]*Endpoint)
	for _, ep := range endpoints {
		byName[ep.Name] = ep
	}

	// Topological sort
	var result []*Endpoint
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(ep *Endpoint)
	visit = func(ep *Endpoint) {
		if visited[ep.Name] {
			return
		}
		if temp[ep.Name] {
			return // cycle, skip
		}
		temp[ep.Name] = true

		if ep.DependsOn != "" {
			if dep, ok := byName[ep.DependsOn]; ok {
				visit(dep)
			}
		}

		temp[ep.Name] = false
		visited[ep.Name] = true
		result = append(result, ep)
	}

	for _, ep := range endpoints {
		visit(ep)
	}

	return result
}
