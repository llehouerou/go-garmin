// endpoint/definitions/register.go
package definitions

import "github.com/llehouerou/go-garmin/endpoint"

// RegisterAll registers all endpoint definitions with the registry.
func RegisterAll(r *endpoint.Registry) {
	for i := range SleepEndpoints {
		r.Register(SleepEndpoints[i])
	}
	for i := range WellnessEndpoints {
		r.Register(WellnessEndpoints[i])
	}
	for i := range HRVEndpoints {
		r.Register(HRVEndpoints[i])
	}
	for i := range WeightEndpoints {
		r.Register(WeightEndpoints[i])
	}
	for i := range DeviceEndpoints {
		r.Register(DeviceEndpoints[i])
	}
	for i := range UserProfileEndpoints {
		r.Register(UserProfileEndpoints[i])
	}
	for i := range ActivityEndpoints {
		r.Register(ActivityEndpoints[i])
	}
	for i := range BiometricEndpoints {
		r.Register(BiometricEndpoints[i])
	}
	for i := range MetricsEndpoints {
		r.Register(MetricsEndpoints[i])
	}
	for i := range WorkoutEndpoints {
		r.Register(WorkoutEndpoints[i])
	}
	for i := range UtilityEndpoints {
		r.Register(UtilityEndpoints[i])
	}
	for i := range CalendarEndpoints {
		r.Register(CalendarEndpoints[i])
	}
	for i := range FitnessAgeEndpoints {
		r.Register(FitnessAgeEndpoints[i])
	}
}
