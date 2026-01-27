// endpoint/definitions/register.go
package definitions

import "github.com/llehouerou/go-garmin/endpoint"

// RegisterAll registers all endpoint definitions with the registry.
func RegisterAll(r *endpoint.Registry) {
	for i := range SleepEndpoints {
		r.Register(SleepEndpoints[i])
	}
	// Additional services will be added here as we migrate them
}
