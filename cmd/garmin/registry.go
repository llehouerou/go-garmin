package main

import (
	"github.com/llehouerou/go-garmin/endpoint"
	"github.com/llehouerou/go-garmin/endpoint/definitions"
)

// endpointRegistry is the global registry for all endpoint definitions.
var endpointRegistry *endpoint.Registry

func init() {
	endpointRegistry = endpoint.NewRegistry()
	definitions.RegisterAll(endpointRegistry)
}
