package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GeoPoint represents a GPS track point with coordinates, elevation, distance, and timestamp.
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Distance  float64 `json:"distance"`
	Timestamp int64   `json:"timestamp"`
}

// CoursePoint represents a named waypoint on a course.
type CoursePoint struct {
	CoursePointID    int64   `json:"coursePointId"`
	CourseID         int64   `json:"courseId"`
	Name             string  `json:"name"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Elevation        float64 `json:"elevation"`
	Distance         float64 `json:"distance"`
	PointType        string  `json:"pointType"`
	SortOrder        int     `json:"sortOrder"`
	DerivedElevation float64 `json:"derivedElevation"`
}

// Coordinate represents a latitude/longitude pair.
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// BoundingBox represents the geographic bounds of a course.
type BoundingBox struct {
	Center              *Coordinate `json:"center"`
	LowerLeft           Coordinate  `json:"lowerLeft"`
	UpperRight          Coordinate  `json:"upperRight"`
	LowerLeftLatIsSet   bool        `json:"lowerLeftLatIsSet"`
	LowerLeftLongIsSet  bool        `json:"lowerLeftLongIsSet"`
	UpperRightLatIsSet  bool        `json:"upperRightLatIsSet"`
	UpperRightLongIsSet bool        `json:"upperRightLongIsSet"`
}

// CourseLine represents a segment of a course.
type CourseLine struct {
	CourseID                 int64   `json:"courseId"`
	SortOrder                int     `json:"sortOrder"`
	NumberOfPoints           int     `json:"numberOfPoints"`
	DistanceInMeters         float64 `json:"distanceInMeters"`
	Bearing                  float64 `json:"bearing"`
	Points                   any     `json:"points"`
	CoordinateSystem         *string `json:"coordinateSystem"`
	OriginalCoordinateSystem *string `json:"originalCoordinateSystem"`
}

// StartPoint represents the starting point of a course.
type StartPoint struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Elevation float64  `json:"elevation"`
	Distance  *float64 `json:"distance"`
	Timestamp *int64   `json:"timestamp"`
}

// CourseDetail represents the detailed response for a specific course.
type CourseDetail struct {
	CourseID                 int64         `json:"courseId"`
	CourseName               string        `json:"courseName"`
	Description              *string       `json:"description"`
	OpenStreetMap            bool          `json:"openStreetMap"`
	MatchedToSegments        bool          `json:"matchedToSegments"`
	UserProfilePK            int64         `json:"userProfilePk"`
	UserGroupPK              *int64        `json:"userGroupPk"`
	RulePK                   int           `json:"rulePK"`
	FirstName                string        `json:"firstName"`
	LastName                 string        `json:"lastName"`
	DisplayName              string        `json:"displayName"`
	GeoRoutePK               *int64        `json:"geoRoutePk"`
	SourceTypeID             int           `json:"sourceTypeId"`
	SourcePK                 *int64        `json:"sourcePk"`
	DistanceMeter            float64       `json:"distanceMeter"`
	ElevationGainMeter       float64       `json:"elevationGainMeter"`
	ElevationLossMeter       float64       `json:"elevationLossMeter"`
	StartPoint               StartPoint    `json:"startPoint"`
	CoursePoints             []CoursePoint `json:"coursePoints"`
	BoundingBox              BoundingBox   `json:"boundingBox"`
	HasShareableEvent        bool          `json:"hasShareableEvent"`
	HasTurnDetectionDisabled bool          `json:"hasTurnDetectionDisabled"`
	ActivityTypePK           int           `json:"activityTypePk"`
	VirtualPartnerID         int64         `json:"virtualPartnerId"`
	IncludeLaps              bool          `json:"includeLaps"`
	ElapsedSeconds           *int          `json:"elapsedSeconds"`
	SpeedMeterPerSecond      *float64      `json:"speedMeterPerSecond"`
	CreateDate               string        `json:"createDate"`
	UpdateDate               string        `json:"updateDate"`
	CourseLines              []CourseLine  `json:"courseLines"`
	CoordinateSystem         string        `json:"coordinateSystem"`
	TargetCoordinateSystem   string        `json:"targetCoordinateSystem"`
	OriginalCoordinateSystem string        `json:"originalCoordinateSystem"`
	Consumer                 *string       `json:"consumer"`
	ElevationSource          int           `json:"elevationSource"`
	HasPaceBand              bool          `json:"hasPaceBand"`
	HasPowerGuide            bool          `json:"hasPowerGuide"`
	Favorite                 bool          `json:"favorite"`
	StartNote                *string       `json:"startNote"`
	FinishNote               *string       `json:"finishNote"`
	CutoffDuration           *int          `json:"cutoffDuration"`
	GeoPoints                []GeoPoint    `json:"geoPoints"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (d *CourseDetail) RawJSON() json.RawMessage {
	return d.raw
}

// SetRaw sets the raw JSON data.
func (d *CourseDetail) SetRaw(data json.RawMessage) {
	d.raw = data
}

// Get retrieves detailed information about a specific course.
func (s *CourseService) Get(ctx context.Context, courseID int64) (*CourseDetail, error) {
	path := fmt.Sprintf("/course-service/course/%d", courseID)
	return fetch[CourseDetail](ctx, s.client, path)
}

// DownloadGPX downloads the course as a GPX file.
func (s *CourseService) DownloadGPX(ctx context.Context, courseID int64) ([]byte, error) {
	path := fmt.Sprintf("/course-service/course/gpx/%d", courseID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return io.ReadAll(resp.Body)
}

// DownloadFIT downloads the course as a FIT file.
func (s *CourseService) DownloadFIT(ctx context.Context, courseID int64) ([]byte, error) {
	path := fmt.Sprintf("/course-service/course/fit/%d/0?elevation=true", courseID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return io.ReadAll(resp.Body)
}
