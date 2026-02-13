package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultCoordinateSystem = "WGS84"

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

// saveGeoPoint is a GeoPoint for the save request, with nullable timestamp.
type saveGeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Distance  float64 `json:"distance"`
	Timestamp *int64  `json:"timestamp"`
}

// saveCourseLine is a CourseLine for the save request, with nullable courseId.
type saveCourseLine struct {
	CourseID         *int64  `json:"courseId"`
	SortOrder        int     `json:"sortOrder"`
	NumberOfPoints   int     `json:"numberOfPoints"`
	DistanceInMeters float64 `json:"distanceInMeters"`
	Bearing          float64 `json:"bearing"`
	Points           any     `json:"points"`
	CoordinateSystem string  `json:"coordinateSystem"`
}

// saveCourseRequest contains only the fields accepted by POST /course-service/course.
type saveCourseRequest struct {
	ActivityTypePK           int              `json:"activityTypePk"`
	HasTurnDetectionDisabled bool             `json:"hasTurnDetectionDisabled"`
	GeoPoints                []saveGeoPoint   `json:"geoPoints"`
	CourseLines              []saveCourseLine `json:"courseLines"`
	BoundingBox              BoundingBox      `json:"boundingBox"`
	CoursePoints             []CoursePoint    `json:"coursePoints"`
	DistanceMeter            float64          `json:"distanceMeter"`
	ElevationGainMeter       float64          `json:"elevationGainMeter"`
	ElevationLossMeter       float64          `json:"elevationLossMeter"`
	StartPoint               StartPoint       `json:"startPoint"`
	ElapsedSeconds           *int             `json:"elapsedSeconds"`
	OpenStreetMap            bool             `json:"openStreetMap"`
	CoordinateSystem         string           `json:"coordinateSystem"`
	RulePK                   int              `json:"rulePK"`
	CourseName               string           `json:"courseName"`
	MatchedToSegments        bool             `json:"matchedToSegments"`
	IncludeLaps              bool             `json:"includeLaps"`
	HasPaceBand              bool             `json:"hasPaceBand"`
	HasPowerGuide            bool             `json:"hasPowerGuide"`
	Favorite                 bool             `json:"favorite"`
	UserProfilePK            int64            `json:"userProfilePk"`
	SpeedMeterPerSecond      *float64         `json:"speedMeterPerSecond"`
	SourceTypeID             int              `json:"sourceTypeId"`
}

func computeBoundingBox(points []saveGeoPoint) BoundingBox {
	first := points[0]
	minLat, maxLat := first.Latitude, first.Latitude
	minLng, maxLng := first.Longitude, first.Longitude
	for _, p := range points[1:] {
		if p.Latitude < minLat {
			minLat = p.Latitude
		}
		if p.Latitude > maxLat {
			maxLat = p.Latitude
		}
		if p.Longitude < minLng {
			minLng = p.Longitude
		}
		if p.Longitude > maxLng {
			maxLng = p.Longitude
		}
	}
	return BoundingBox{
		LowerLeft:           Coordinate{Latitude: minLat, Longitude: minLng},
		UpperRight:          Coordinate{Latitude: maxLat, Longitude: maxLng},
		LowerLeftLatIsSet:   true,
		LowerLeftLongIsSet:  true,
		UpperRightLatIsSet:  true,
		UpperRightLongIsSet: true,
	}
}

func computeElevation(points []saveGeoPoint) (gain, loss float64) {
	for i := 1; i < len(points); i++ {
		diff := points[i].Elevation - points[i-1].Elevation
		if diff > 0 {
			gain += diff
		} else {
			loss -= diff
		}
	}
	return gain, loss
}

func newSaveCourseRequest(d *CourseDetail) saveCourseRequest {
	geoPoints := make([]saveGeoPoint, len(d.GeoPoints))
	for i, p := range d.GeoPoints {
		geoPoints[i] = saveGeoPoint{
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
			Elevation: p.Elevation,
			Distance:  p.Distance,
			Timestamp: nil,
		}
	}

	// Compute metadata from geoPoints when the upload endpoint returns empty values.
	startPoint := d.StartPoint
	bb := d.BoundingBox
	distanceMeter := d.DistanceMeter
	elevationGain := d.ElevationGainMeter
	elevationLoss := d.ElevationLossMeter

	if len(geoPoints) > 0 {
		if startPoint.Latitude == 0 && startPoint.Longitude == 0 {
			first := geoPoints[0]
			startPoint = StartPoint{
				Latitude:  first.Latitude,
				Longitude: first.Longitude,
				Elevation: first.Elevation,
			}
		}
		if !bb.LowerLeftLatIsSet {
			bb = computeBoundingBox(geoPoints)
		}
		if distanceMeter == 0 {
			distanceMeter = geoPoints[len(geoPoints)-1].Distance
		}
		if elevationGain == 0 && elevationLoss == 0 {
			elevationGain, elevationLoss = computeElevation(geoPoints)
		}
	}

	courseLines := make([]saveCourseLine, len(d.CourseLines))
	for i, l := range d.CourseLines {
		cs := defaultCoordinateSystem
		if l.CoordinateSystem != nil {
			cs = *l.CoordinateSystem
		}
		courseLines[i] = saveCourseLine{
			CourseID:         nil,
			SortOrder:        l.SortOrder,
			NumberOfPoints:   l.NumberOfPoints,
			DistanceInMeters: l.DistanceInMeters,
			Bearing:          l.Bearing,
			Points:           l.Points,
			CoordinateSystem: cs,
		}
	}

	return saveCourseRequest{
		ActivityTypePK:           d.ActivityTypePK,
		HasTurnDetectionDisabled: d.HasTurnDetectionDisabled,
		GeoPoints:                geoPoints,
		CourseLines:              courseLines,
		BoundingBox:              bb,
		CoursePoints:             d.CoursePoints,
		DistanceMeter:            distanceMeter,
		ElevationGainMeter:       elevationGain,
		ElevationLossMeter:       elevationLoss,
		StartPoint:               startPoint,
		ElapsedSeconds:           d.ElapsedSeconds,
		OpenStreetMap:            d.OpenStreetMap,
		CoordinateSystem:         d.CoordinateSystem,
		RulePK:                   d.RulePK,
		CourseName:               d.CourseName,
		MatchedToSegments:        d.MatchedToSegments,
		IncludeLaps:              d.IncludeLaps,
		HasPaceBand:              d.HasPaceBand,
		HasPowerGuide:            d.HasPowerGuide,
		Favorite:                 d.Favorite,
		UserProfilePK:            d.UserProfilePK,
		SpeedMeterPerSecond:      d.SpeedMeterPerSecond,
		SourceTypeID:             d.SourceTypeID,
	}
}

// Save creates a course on Garmin Connect.
func (s *CourseService) Save(ctx context.Context, course *CourseDetail) (*CourseDetail, error) {
	req := newSaveCourseRequest(course)
	return send[CourseDetail](ctx, s.client, http.MethodPost, "/course-service/course", req)
}

// Import imports a course from a GPX file and saves it to Garmin Connect.
// If activityType > 0, it overrides the auto-detected activity type.
// If privacy > 0, it sets the privacy rule (1=Public, 2=Private, 4=Group); defaults to 2 (Private).
func (s *CourseService) Import(ctx context.Context, fileName string, content io.Reader, activityType, privacy int) (*CourseDetail, error) {
	detail, err := upload[CourseDetail](ctx, s.client, "/course-service/course/import", "file", fileName, content)
	if err != nil {
		return nil, fmt.Errorf("upload gpx: %w", err)
	}
	if activityType > 0 {
		detail.ActivityTypePK = activityType
	}
	if detail.CoordinateSystem == "" {
		detail.CoordinateSystem = defaultCoordinateSystem
	}
	if detail.SourceTypeID == 0 {
		detail.SourceTypeID = 3 // GPX import
	}
	if privacy > 0 {
		detail.RulePK = privacy
	} else if detail.RulePK == 0 {
		detail.RulePK = 2 // Private
	}
	return s.Save(ctx, detail)
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
