package garmin

import (
	"context"
	"encoding/json"
)

// CourseActivityType represents the activity type of a course.
type CourseActivityType struct {
	TypeID       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeID int    `json:"parentTypeId"`
	IsHidden     bool   `json:"isHidden"`
	Restricted   bool   `json:"restricted"`
	Trimmable    bool   `json:"trimmable"`
}

// CoursePrivacyRule represents the privacy setting of a course.
type CoursePrivacyRule struct {
	TypeID  int    `json:"typeId"`
	TypeKey string `json:"typeKey"`
}

// Course represents a Garmin course/route.
type Course struct {
	CourseID                 int64              `json:"courseId"`
	UserProfileID            int64              `json:"userProfileId"`
	DisplayName              string             `json:"displayName"`
	UserGroupID              *int64             `json:"userGroupId"`
	GeoRoutePK               *int64             `json:"geoRoutePk"`
	ActivityType             CourseActivityType `json:"activityType"`
	CourseName               string             `json:"courseName"`
	CourseDescription        *string            `json:"courseDescription"`
	CreatedDate              int64              `json:"createdDate"`
	UpdatedDate              int64              `json:"updatedDate"`
	PrivacyRule              CoursePrivacyRule  `json:"privacyRule"`
	DistanceInMeters         float64            `json:"distanceInMeters"`
	ElevationGainInMeters    float64            `json:"elevationGainInMeters"`
	ElevationLossInMeters    float64            `json:"elevationLossInMeters"`
	StartLatitude            float64            `json:"startLatitude"`
	StartLongitude           float64            `json:"startLongitude"`
	SpeedInMetersPerSecond   float64            `json:"speedInMetersPerSecond"`
	SourceTypeID             int                `json:"sourceTypeId"`
	SourcePK                 *int64             `json:"sourcePk"`
	ElapsedSeconds           *int               `json:"elapsedSeconds"`
	CoordinateSystem         string             `json:"coordinateSystem"`
	OriginalCoordinateSystem string             `json:"originalCoordinateSystem"`
	Consumer                 *string            `json:"consumer"`
	ElevationSource          int                `json:"elevationSource"`
	HasShareableEvent        bool               `json:"hasShareableEvent"`
	HasPaceBand              bool               `json:"hasPaceBand"`
	HasPowerGuide            bool               `json:"hasPowerGuide"`
	Favorite                 bool               `json:"favorite"`
	HasTurnDetectionDisabled bool               `json:"hasTurnDetectionDisabled"`
	CuratedCourseID          *int64             `json:"curatedCourseId"`
	StartNote                *string            `json:"startNote"`
	FinishNote               *string            `json:"finishNote"`
	CutoffDuration           *int               `json:"cutoffDuration"`
	CreatedDateFormatted     string             `json:"createdDateFormatted"`
	UpdatedDateFormatted     string             `json:"updatedDateFormatted"`
	Public                   bool               `json:"public"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (c *Course) RawJSON() json.RawMessage {
	return c.raw
}

// SetRaw sets the raw JSON data.
func (c *Course) SetRaw(data json.RawMessage) {
	c.raw = data
}

// CoursesForUserResponse represents the API response for owner courses.
type CoursesForUserResponse struct {
	CoursesForUser []Course `json:"coursesForUser"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (r *CoursesForUserResponse) RawJSON() json.RawMessage {
	return r.raw
}

// SetRaw sets the raw JSON data.
func (r *CoursesForUserResponse) SetRaw(data json.RawMessage) {
	r.raw = data
}

// ListOwner retrieves all courses owned by the authenticated user.
func (s *CourseService) ListOwner(ctx context.Context) (*CoursesForUserResponse, error) {
	path := "/web-gateway/course/owner"
	return fetch[CoursesForUserResponse](ctx, s.client, path)
}
