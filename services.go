package garmin

// WellnessService provides access to wellness-related API endpoints.
type WellnessService struct{ client *Client }

// ActivityService provides access to activity-related API endpoints.
type ActivityService struct{ client *Client }

// MetricsService provides access to metrics-related API endpoints.
type MetricsService struct{ client *Client }

// WeightService provides access to weight-related API endpoints.
type WeightService struct{ client *Client }

// DeviceService provides access to device-related API endpoints.
type DeviceService struct{ client *Client }

// WorkoutService provides access to workout-related API endpoints.
type WorkoutService struct{ client *Client }

// GoalService provides access to goal-related API endpoints.
type GoalService struct{ client *Client }

// BadgeService provides access to badge-related API endpoints.
type BadgeService struct{ client *Client }

// GearService provides access to gear-related API endpoints.
type GearService struct{ client *Client }

// DownloadService provides access to download-related API endpoints.
type DownloadService struct{ client *Client }

// UploadService provides access to upload-related API endpoints.
type UploadService struct{ client *Client }

// HydrationService provides access to hydration-related API endpoints.
type HydrationService struct{ client *Client }

// BloodPressureService provides access to blood pressure-related API endpoints.
type BloodPressureService struct{ client *Client }

// PersonalRecordsService provides access to personal records-related API endpoints.
type PersonalRecordsService struct{ client *Client }

// StepsService provides access to steps-related API endpoints.
type StepsService struct{ client *Client }

// UserProfileService provides access to user profile-related API endpoints.
type UserProfileService struct{ client *Client }

// HRVService provides access to HRV (heart rate variability) related API endpoints.
type HRVService struct{ client *Client }

// BiometricService provides access to biometric-related API endpoints (FTP, lactate threshold, power-to-weight).
type BiometricService struct{ client *Client }

// CalendarService provides access to calendar-related API endpoints.
type CalendarService struct{ client *Client }

// FitnessAgeService provides access to fitness age-related API endpoints.
type FitnessAgeService struct{ client *Client }

// FitnessStatsService provides access to fitness statistics API endpoints.
type FitnessStatsService struct{ client *Client }

// CourseService provides access to course-related API endpoints.
type CourseService struct{ client *Client }
