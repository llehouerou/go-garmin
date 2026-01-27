package garmin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ActivityType represents the type of activity.
type ActivityType struct {
	TypeID       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeID int    `json:"parentTypeId"`
	IsHidden     bool   `json:"isHidden"`
	Restricted   bool   `json:"restricted"`
	Trimmable    bool   `json:"trimmable"`
}

// EventType represents the event type of an activity.
type EventType struct {
	TypeID    int    `json:"typeId"`
	TypeKey   string `json:"typeKey"`
	SortOrder int    `json:"sortOrder"`
}

// Privacy represents the privacy settings for an activity.
type Privacy struct {
	TypeID  int    `json:"typeId"`
	TypeKey string `json:"typeKey"`
}

// SplitSummary represents a split segment within an activity.
type SplitSummary struct {
	NoOfSplits           int     `json:"noOfSplits"`
	TotalAscent          float64 `json:"totalAscent"`
	Duration             float64 `json:"duration"`
	SplitType            string  `json:"splitType"`
	NumClimbSends        int     `json:"numClimbSends"`
	MaxElevationGain     float64 `json:"maxElevationGain"`
	AverageElevationGain float64 `json:"averageElevationGain"`
	MaxDistance          int     `json:"maxDistance"`
	Distance             float64 `json:"distance"`
	AverageSpeed         float64 `json:"averageSpeed"`
	MaxSpeed             float64 `json:"maxSpeed"`
	NumFalls             int     `json:"numFalls"`
	ElevationLoss        float64 `json:"elevationLoss"`

	// Additional fields from detailed response
	MovingDuration           *float64 `json:"movingDuration,omitempty"`
	ElevationGain            *float64 `json:"elevationGain,omitempty"`
	AverageMovingSpeed       *float64 `json:"averageMovingSpeed,omitempty"`
	Calories                 *float64 `json:"calories,omitempty"`
	BMRCalories              *float64 `json:"bmrCalories,omitempty"`
	AverageHR                *float64 `json:"averageHR,omitempty"`
	MaxHR                    *float64 `json:"maxHR,omitempty"`
	AverageRunCadence        *float64 `json:"averageRunCadence,omitempty"`
	MaxRunCadence            *float64 `json:"maxRunCadence,omitempty"`
	AveragePower             *float64 `json:"averagePower,omitempty"`
	MaxPower                 *float64 `json:"maxPower,omitempty"`
	NormalizedPower          *float64 `json:"normalizedPower,omitempty"`
	GroundContactTime        *float64 `json:"groundContactTime,omitempty"`
	StrideLength             *float64 `json:"strideLength,omitempty"`
	VerticalOscillation      *float64 `json:"verticalOscillation,omitempty"`
	VerticalRatio            *float64 `json:"verticalRatio,omitempty"`
	TotalExerciseReps        *int     `json:"totalExerciseReps,omitempty"`
	AvgVerticalSpeed         *float64 `json:"avgVerticalSpeed,omitempty"`
	AvgGradeAdjustedSpeed    *float64 `json:"avgGradeAdjustedSpeed,omitempty"`
	MaxDistanceWithPrecision *float64 `json:"maxDistanceWithPrecision,omitempty"`
	AvgStepFrequency         *float64 `json:"avgStepFrequency,omitempty"`
	AvgStepLength            *float64 `json:"avgStepLength,omitempty"`
}

// SummarizedDiveInfo represents diving information for an activity.
type SummarizedDiveInfo struct {
	SummarizedDiveGases []any `json:"summarizedDiveGases"`
}

// Activity represents a Garmin activity from the activity list.
type Activity struct {
	// Core identification
	ActivityID   int64  `json:"activityId"`
	ActivityName string `json:"activityName"`
	ActivityUUID string `json:"activityUUID"`

	// Timestamps
	StartTimeLocal string `json:"startTimeLocal"`
	StartTimeGMT   string `json:"startTimeGMT"`
	EndTimeGMT     string `json:"endTimeGMT,omitempty"`
	BeginTimestamp int64  `json:"beginTimestamp"`

	// Activity classification
	ActivityType ActivityType `json:"activityType"`
	EventType    EventType    `json:"eventType"`
	SportTypeID  int          `json:"sportTypeId"`

	// Distance and duration
	Distance        float64 `json:"distance"`
	Duration        float64 `json:"duration"`
	ElapsedDuration float64 `json:"elapsedDuration"`
	MovingDuration  float64 `json:"movingDuration"`

	// Elevation
	ElevationGain float64 `json:"elevationGain"`
	ElevationLoss float64 `json:"elevationLoss"`
	MinElevation  float64 `json:"minElevation"`
	MaxElevation  float64 `json:"maxElevation"`

	// Speed
	AverageSpeed float64 `json:"averageSpeed"`
	MaxSpeed     float64 `json:"maxSpeed"`

	// Location
	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`
	EndLatitude    float64 `json:"endLatitude,omitempty"`
	EndLongitude   float64 `json:"endLongitude,omitempty"`
	LocationName   string  `json:"locationName"`
	TimeZoneID     int     `json:"timeZoneId"`

	// Content flags
	HasPolyline bool `json:"hasPolyline"`
	HasImages   bool `json:"hasImages"`
	HasVideo    bool `json:"hasVideo"`
	HasSplits   bool `json:"hasSplits"`
	HasHeatMap  bool `json:"hasHeatMap"`

	// Owner information
	OwnerID                    int64    `json:"ownerId"`
	OwnerDisplayName           string   `json:"ownerDisplayName"`
	OwnerFullName              string   `json:"ownerFullName"`
	OwnerProfileImageURLSmall  string   `json:"ownerProfileImageUrlSmall"`
	OwnerProfileImageURLMedium string   `json:"ownerProfileImageUrlMedium"`
	OwnerProfileImageURLLarge  string   `json:"ownerProfileImageUrlLarge"`
	UserRoles                  []string `json:"userRoles"`
	UserPro                    bool     `json:"userPro"`

	// Calories
	Calories    float64 `json:"calories"`
	BMRCalories float64 `json:"bmrCalories"`

	// Heart rate
	AverageHR float64 `json:"averageHR"`
	MaxHR     float64 `json:"maxHR"`

	// Cadence
	AverageRunningCadenceInStepsPerMinute float64 `json:"averageRunningCadenceInStepsPerMinute"`
	MaxRunningCadenceInStepsPerMinute     float64 `json:"maxRunningCadenceInStepsPerMinute"`
	MaxDoubleCadence                      float64 `json:"maxDoubleCadence"`

	// Steps
	Steps int `json:"steps"`

	// Power
	AvgPower  float64 `json:"avgPower"`
	MaxPower  float64 `json:"maxPower"`
	NormPower float64 `json:"normPower"`

	// Training effect
	AerobicTrainingEffect          float64 `json:"aerobicTrainingEffect"`
	AnaerobicTrainingEffect        float64 `json:"anaerobicTrainingEffect"`
	TrainingEffectLabel            string  `json:"trainingEffectLabel"`
	AerobicTrainingEffectMessage   string  `json:"aerobicTrainingEffectMessage"`
	AnaerobicTrainingEffectMessage string  `json:"anaerobicTrainingEffectMessage"`
	ActivityTrainingLoad           float64 `json:"activityTrainingLoad"`

	// Running dynamics
	AvgVerticalOscillation float64 `json:"avgVerticalOscillation"`
	AvgGroundContactTime   float64 `json:"avgGroundContactTime"`
	AvgStrideLength        float64 `json:"avgStrideLength"`
	AvgVerticalRatio       float64 `json:"avgVerticalRatio"`
	AvgGradeAdjustedSpeed  float64 `json:"avgGradeAdjustedSpeed"`

	// Physiological
	VO2MaxValue      float64 `json:"vO2MaxValue"`
	WaterEstimated   float64 `json:"waterEstimated"`
	MaxVerticalSpeed float64 `json:"maxVerticalSpeed"`

	// Respiration
	MinRespirationRate float64 `json:"minRespirationRate"`
	MaxRespirationRate float64 `json:"maxRespirationRate"`
	AvgRespirationRate float64 `json:"avgRespirationRate"`

	// Intensity minutes
	ModerateIntensityMinutes int `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes int `json:"vigorousIntensityMinutes"`

	// Body battery
	DifferenceBodyBattery int `json:"differenceBodyBattery"`

	// Device info
	DeviceID     int64  `json:"deviceId"`
	Manufacturer string `json:"manufacturer"`
	LapCount     int    `json:"lapCount"`

	// Privacy
	Privacy Privacy `json:"privacy"`

	// Splits and pace
	SplitSummaries         []SplitSummary `json:"splitSummaries"`
	MinActivityLapDuration float64        `json:"minActivityLapDuration"`

	// Fastest splits (in seconds)
	FastestSplit1000 float64 `json:"fastestSplit_1000,omitempty"`
	FastestSplit1609 float64 `json:"fastestSplit_1609,omitempty"`
	FastestSplit5000 float64 `json:"fastestSplit_5000,omitempty"`

	// HR time in zones (in seconds)
	HRTimeInZone1 float64 `json:"hrTimeInZone_1"`
	HRTimeInZone2 float64 `json:"hrTimeInZone_2"`
	HRTimeInZone3 float64 `json:"hrTimeInZone_3"`
	HRTimeInZone4 float64 `json:"hrTimeInZone_4"`
	HRTimeInZone5 float64 `json:"hrTimeInZone_5"`

	// Power time in zones (in seconds)
	PowerTimeInZone1 float64 `json:"powerTimeInZone_1"`
	PowerTimeInZone2 float64 `json:"powerTimeInZone_2"`
	PowerTimeInZone3 float64 `json:"powerTimeInZone_3"`
	PowerTimeInZone4 float64 `json:"powerTimeInZone_4"`
	PowerTimeInZone5 float64 `json:"powerTimeInZone_5"`

	// Dive info
	SummarizedDiveInfo SummarizedDiveInfo `json:"summarizedDiveInfo"`
	QualifyingDive     bool               `json:"qualifyingDive"`
	DecoDive           bool               `json:"decoDive"`

	// Flags
	PR                 bool `json:"pr"`
	AutoCalcCalories   bool `json:"autoCalcCalories"`
	Favorite           bool `json:"favorite"`
	ElevationCorrected bool `json:"elevationCorrected"`
	AtpActivity        bool `json:"atpActivity"`
	ManualActivity     bool `json:"manualActivity"`
	Purposeful         bool `json:"purposeful"`
	Parent             bool `json:"parent"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (a *Activity) RawJSON() json.RawMessage {
	return a.raw
}

// StartTime returns the activity start time parsed from StartTimeGMT.
func (a *Activity) StartTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", a.StartTimeGMT)
	return t
}

// DurationTime returns the activity duration as a time.Duration.
func (a *Activity) DurationTime() time.Duration {
	return time.Duration(a.Duration * float64(time.Second))
}

// DistanceKm returns the distance in kilometers.
func (a *Activity) DistanceKm() float64 {
	return a.Distance / 1000
}

// DistanceMiles returns the distance in miles.
func (a *Activity) DistanceMiles() float64 {
	return a.Distance / 1609.344
}

// AveragePacePerKm returns the average pace per kilometer.
func (a *Activity) AveragePacePerKm() time.Duration {
	if a.Distance == 0 {
		return 0
	}
	return time.Duration(a.Duration / (a.Distance / 1000) * float64(time.Second))
}

// ActivityUUIDObject represents the UUID structure in detailed activity response.
type ActivityUUIDObject struct {
	UUID string `json:"uuid"`
}

// AccessControlRule represents the access control settings for an activity.
type AccessControlRule struct {
	TypeID  int    `json:"typeId"`
	TypeKey string `json:"typeKey"`
}

// TimeZoneUnit represents timezone information.
type TimeZoneUnit struct {
	UnitID   int     `json:"unitId"`
	UnitKey  string  `json:"unitKey"`
	Factor   float64 `json:"factor"`
	TimeZone string  `json:"timeZone"`
}

// FileFormat represents the file format of an activity.
type FileFormat struct {
	FormatID  int    `json:"formatId"`
	FormatKey string `json:"formatKey"`
}

// UserInfo represents user information in activity metadata.
type UserInfo struct {
	UserProfilePK         int64  `json:"userProfilePk"`
	DisplayName           string `json:"displayname"`
	FullName              string `json:"fullname"`
	ProfileImageURLLarge  string `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string `json:"profileImageUrlSmall"`
	UserPro               bool   `json:"userPro"`
}

// Sensor represents a sensor used during an activity.
type Sensor struct {
	Manufacturer      string  `json:"manufacturer"`
	SerialNumber      int64   `json:"serialNumber"`
	SourceType        string  `json:"sourceType"`
	AntplusDeviceType string  `json:"antplusDeviceType"`
	SoftwareVersion   float64 `json:"softwareVersion"`
	BatteryStatus     string  `json:"batteryStatus"`
	BatteryLevel      int     `json:"batteryLevel"`
}

// DeviceMetaData represents device metadata.
type DeviceMetaData struct {
	DeviceID        string `json:"deviceId"`
	DeviceTypePK    int    `json:"deviceTypePk"`
	DeviceVersionPK int    `json:"deviceVersionPk"`
}

// ActivityMetadata represents metadata for an activity.
type ActivityMetadata struct {
	IsOriginal                      bool           `json:"isOriginal"`
	DeviceApplicationInstallationID int            `json:"deviceApplicationInstallationId"`
	AgentApplicationInstallationID  *int           `json:"agentApplicationInstallationId"`
	AgentString                     *string        `json:"agentString"`
	FileFormat                      FileFormat     `json:"fileFormat"`
	AssociatedCourseID              *int64         `json:"associatedCourseId"`
	LastUpdateDate                  string         `json:"lastUpdateDate"`
	UploadedDate                    string         `json:"uploadedDate"`
	VideoURL                        *string        `json:"videoUrl"`
	HasPolyline                     bool           `json:"hasPolyline"`
	HasChartData                    bool           `json:"hasChartData"`
	HasHRTimeInZones                bool           `json:"hasHrTimeInZones"`
	HasPowerTimeInZones             bool           `json:"hasPowerTimeInZones"`
	UserInfoDTO                     UserInfo       `json:"userInfoDto"`
	ChildIDs                        []int64        `json:"childIds"`
	ChildActivityTypes              []ActivityType `json:"childActivityTypes"`
	Sensors                         []Sensor       `json:"sensors"`
	ActivityImages                  []any          `json:"activityImages"`
	Manufacturer                    string         `json:"manufacturer"`
	DiveNumber                      *int           `json:"diveNumber"`
	LapCount                        int            `json:"lapCount"`
	AssociatedWorkoutID             *int64         `json:"associatedWorkoutId"`
	IsAtpActivity                   *bool          `json:"isAtpActivity"`
	DeviceMetaDataDTO               DeviceMetaData `json:"deviceMetaDataDTO"`
	HasIntensityIntervals           bool           `json:"hasIntensityIntervals"`
	HasSplits                       bool           `json:"hasSplits"`
	EBikeMaxAssistModes             *int           `json:"eBikeMaxAssistModes"`
	EBikeBatteryUsage               *float64       `json:"eBikeBatteryUsage"`
	EBikeBatteryRemaining           *float64       `json:"eBikeBatteryRemaining"`
	EBikeAssistModeInfoDTOList      []any          `json:"eBikeAssistModeInfoDTOList"`
	HasRunPowerWindData             bool           `json:"hasRunPowerWindData"`
	CalendarEventInfo               *any           `json:"calendarEventInfo"`
	GroupRideUUID                   *string        `json:"groupRideUUID"`
	HasHeatMap                      bool           `json:"hasHeatMap"`
	SpecializedWorkoutCategories    []string       `json:"specializedWorkoutCategories"`
	TrainingPlanID                  *int64         `json:"trainingPlanId"`
	PersonalRecord                  bool           `json:"personalRecord"`
	GCJ02                           bool           `json:"gcj02"`
	RunPowerWindDataEnabled         bool           `json:"runPowerWindDataEnabled"`
	ElevationCorrected              bool           `json:"elevationCorrected"`
	Trimmed                         bool           `json:"trimmed"`
	ManualActivity                  bool           `json:"manualActivity"`
	AutoCalcCalories                bool           `json:"autoCalcCalories"`
	Favorite                        bool           `json:"favorite"`
}

// ActivitySummary represents the summary data for a detailed activity.
type ActivitySummary struct {
	StartTimeLocal string  `json:"startTimeLocal"`
	StartTimeGMT   string  `json:"startTimeGMT"`
	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`
	EndLatitude    float64 `json:"endLatitude"`
	EndLongitude   float64 `json:"endLongitude"`

	// Distance and duration
	Distance        float64 `json:"distance"`
	Duration        float64 `json:"duration"`
	MovingDuration  float64 `json:"movingDuration"`
	ElapsedDuration float64 `json:"elapsedDuration"`

	// Elevation
	ElevationGain float64 `json:"elevationGain"`
	ElevationLoss float64 `json:"elevationLoss"`
	MaxElevation  float64 `json:"maxElevation"`
	MinElevation  float64 `json:"minElevation"`

	// Speed
	AverageSpeed       float64 `json:"averageSpeed"`
	AverageMovingSpeed float64 `json:"averageMovingSpeed"`
	MaxSpeed           float64 `json:"maxSpeed"`
	MaxVerticalSpeed   float64 `json:"maxVerticalSpeed"`

	// Calories
	Calories    float64 `json:"calories"`
	BMRCalories float64 `json:"bmrCalories"`

	// Heart rate
	AverageHR float64 `json:"averageHR"`
	MaxHR     float64 `json:"maxHR"`
	MinHR     float64 `json:"minHR"`

	// Cadence
	AverageRunCadence float64 `json:"averageRunCadence"`
	MaxRunCadence     float64 `json:"maxRunCadence"`

	// Power
	AveragePower    float64 `json:"averagePower"`
	MaxPower        float64 `json:"maxPower"`
	MinPower        float64 `json:"minPower"`
	NormalizedPower float64 `json:"normalizedPower"`
	TotalWork       float64 `json:"totalWork"`

	// Running dynamics
	GroundContactTime   float64 `json:"groundContactTime"`
	StrideLength        float64 `json:"strideLength"`
	VerticalOscillation float64 `json:"verticalOscillation"`
	VerticalRatio       float64 `json:"verticalRatio"`

	// Training effect
	TrainingEffect                 float64 `json:"trainingEffect"`
	AnaerobicTrainingEffect        float64 `json:"anaerobicTrainingEffect"`
	AerobicTrainingEffectMessage   string  `json:"aerobicTrainingEffectMessage"`
	AnaerobicTrainingEffectMessage string  `json:"anaerobicTrainingEffectMessage"`
	TrainingEffectLabel            string  `json:"trainingEffectLabel"`
	ActivityTrainingLoad           float64 `json:"activityTrainingLoad"`

	// Respiration
	WaterEstimated     float64 `json:"waterEstimated"`
	MinRespirationRate float64 `json:"minRespirationRate"`
	MaxRespirationRate float64 `json:"maxRespirationRate"`
	AvgRespirationRate float64 `json:"avgRespirationRate"`

	// Workout feedback
	MinActivityLapDuration       float64 `json:"minActivityLapDuration"`
	DirectWorkoutFeel            int     `json:"directWorkoutFeel"`
	DirectWorkoutRPE             int     `json:"directWorkoutRpe"`
	DirectWorkoutComplianceScore int     `json:"directWorkoutComplianceScore"`

	// Intensity
	ModerateIntensityMinutes int `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes int `json:"vigorousIntensityMinutes"`

	// Steps
	Steps int `json:"steps"`

	// Stamina
	BeginPotentialStamina float64 `json:"beginPotentialStamina"`
	EndPotentialStamina   float64 `json:"endPotentialStamina"`
	MinAvailableStamina   float64 `json:"minAvailableStamina"`

	// Grade adjusted
	AvgGradeAdjustedSpeed float64 `json:"avgGradeAdjustedSpeed"`

	// Body battery
	DifferenceBodyBattery int `json:"differenceBodyBattery"`
}

// ActivityDetail represents detailed information about a single activity.
type ActivityDetail struct {
	ActivityID           int64              `json:"activityId"`
	ActivityUUID         ActivityUUIDObject `json:"activityUUID"`
	ActivityName         string             `json:"activityName"`
	UserProfileID        int64              `json:"userProfileId"`
	IsMultiSportParent   bool               `json:"isMultiSportParent"`
	ActivityTypeDTO      ActivityType       `json:"activityTypeDTO"`
	EventTypeDTO         EventType          `json:"eventTypeDTO"`
	AccessControlRuleDTO AccessControlRule  `json:"accessControlRuleDTO"`
	TimeZoneUnitDTO      TimeZoneUnit       `json:"timeZoneUnitDTO"`
	MetadataDTO          ActivityMetadata   `json:"metadataDTO"`
	SummaryDTO           ActivitySummary    `json:"summaryDTO"`
	LocationName         string             `json:"locationName"`
	SplitSummaries       []SplitSummary     `json:"splitSummaries"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (a *ActivityDetail) RawJSON() json.RawMessage {
	return a.raw
}

// StartTime returns the activity start time parsed from SummaryDTO.
func (a *ActivityDetail) StartTime() time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05.0", a.SummaryDTO.StartTimeGMT)
	return t
}

// DurationTime returns the activity duration as a time.Duration.
func (a *ActivityDetail) DurationTime() time.Duration {
	return time.Duration(a.SummaryDTO.Duration * float64(time.Second))
}

// DistanceKm returns the distance in kilometers.
func (a *ActivityDetail) DistanceKm() float64 {
	return a.SummaryDTO.Distance / 1000
}

// ListOptions specifies options for listing activities.
type ListOptions struct {
	Start int // Starting index (0-based)
	Limit int // Maximum number of activities to return
}

// List retrieves a list of activities.
func (s *ActivityService) List(ctx context.Context, opts *ListOptions) ([]Activity, error) {
	start := 0
	limit := 20
	if opts != nil {
		if opts.Start > 0 {
			start = opts.Start
		}
		if opts.Limit > 0 {
			limit = opts.Limit
		}
	}

	path := fmt.Sprintf("/activitylist-service/activities/search/activities?start=%d&limit=%d", start, limit)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var activities []Activity
	if err := json.Unmarshal(raw, &activities); err != nil {
		return nil, err
	}

	// Store raw JSON in each activity
	var rawActivities []json.RawMessage
	if err := json.Unmarshal(raw, &rawActivities); err == nil {
		for i := range activities {
			if i < len(rawActivities) {
				activities[i].raw = rawActivities[i]
			}
		}
	}

	return activities, nil
}

// Get retrieves detailed information about a specific activity.
func (s *ActivityService) Get(ctx context.Context, activityID int64) (*ActivityDetail, error) {
	path := fmt.Sprintf("/activity-service/activity/%d", activityID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var activity ActivityDetail
	if err := json.Unmarshal(raw, &activity); err != nil {
		return nil, err
	}
	activity.raw = raw

	return &activity, nil
}

// WeatherStation represents the weather station that provided the data.
type WeatherStation struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Timezone *string `json:"timezone"`
}

// WeatherType represents the type of weather conditions.
type WeatherType struct {
	WeatherTypePK *int    `json:"weatherTypePk"`
	Desc          string  `json:"desc"`
	Image         *string `json:"image"`
}

// ActivityWeather represents weather data for an activity.
type ActivityWeather struct {
	IssueDate                 string         `json:"issueDate"`
	Temp                      int            `json:"temp"`
	ApparentTemp              int            `json:"apparentTemp"`
	DewPoint                  int            `json:"dewPoint"`
	RelativeHumidity          int            `json:"relativeHumidity"`
	WindDirection             int            `json:"windDirection"`
	WindDirectionCompassPoint string         `json:"windDirectionCompassPoint"`
	WindSpeed                 int            `json:"windSpeed"`
	WindGust                  *int           `json:"windGust"`
	Latitude                  float64        `json:"latitude"`
	Longitude                 float64        `json:"longitude"`
	WeatherStationDTO         WeatherStation `json:"weatherStationDTO"`
	WeatherTypeDTO            WeatherType    `json:"weatherTypeDTO"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (w *ActivityWeather) RawJSON() json.RawMessage {
	return w.raw
}

// TempCelsius returns the temperature in Celsius.
func (w *ActivityWeather) TempCelsius() float64 {
	return float64(w.Temp-32) * 5 / 9
}

// ApparentTempCelsius returns the apparent temperature in Celsius.
func (w *ActivityWeather) ApparentTempCelsius() float64 {
	return float64(w.ApparentTemp-32) * 5 / 9
}

// GetWeather retrieves weather data for a specific activity.
func (s *ActivityService) GetWeather(ctx context.Context, activityID int64) (*ActivityWeather, error) {
	path := fmt.Sprintf("/activity-service/activity/%d/weather", activityID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather ActivityWeather
	if err := json.Unmarshal(raw, &weather); err != nil {
		return nil, err
	}
	weather.raw = raw

	return &weather, nil
}

// SectionType represents the type of event section.
type SectionType struct {
	ID             int    `json:"id"`
	Key            string `json:"key"`
	SectionTypeKey string `json:"sectionTypeKey"`
}

// ActivityEvent represents an event that occurred during an activity.
type ActivityEvent struct {
	StartTimeGMT            string      `json:"startTimeGMT"`
	StartTimeGMTDoubleValue float64     `json:"startTimeGMTDoubleValue"`
	SectionTypeDTO          SectionType `json:"sectionTypeDTO"`
}

// Lap represents a single lap within an activity.
type Lap struct {
	StartTimeGMT                 string  `json:"startTimeGMT"`
	StartLatitude                float64 `json:"startLatitude"`
	StartLongitude               float64 `json:"startLongitude"`
	Distance                     float64 `json:"distance"`
	Duration                     float64 `json:"duration"`
	MovingDuration               float64 `json:"movingDuration"`
	ElapsedDuration              float64 `json:"elapsedDuration"`
	ElevationGain                float64 `json:"elevationGain"`
	ElevationLoss                float64 `json:"elevationLoss"`
	MaxElevation                 float64 `json:"maxElevation"`
	MinElevation                 float64 `json:"minElevation"`
	AverageSpeed                 float64 `json:"averageSpeed"`
	AverageMovingSpeed           float64 `json:"averageMovingSpeed"`
	MaxSpeed                     float64 `json:"maxSpeed"`
	Calories                     float64 `json:"calories"`
	BMRCalories                  float64 `json:"bmrCalories"`
	AverageHR                    float64 `json:"averageHR"`
	MaxHR                        float64 `json:"maxHR"`
	AverageRunCadence            float64 `json:"averageRunCadence"`
	MaxRunCadence                float64 `json:"maxRunCadence"`
	AveragePower                 float64 `json:"averagePower"`
	MaxPower                     float64 `json:"maxPower"`
	MinPower                     float64 `json:"minPower"`
	NormalizedPower              float64 `json:"normalizedPower"`
	TotalWork                    float64 `json:"totalWork"`
	GroundContactTime            float64 `json:"groundContactTime"`
	StrideLength                 float64 `json:"strideLength"`
	VerticalOscillation          float64 `json:"verticalOscillation"`
	VerticalRatio                float64 `json:"verticalRatio"`
	EndLatitude                  float64 `json:"endLatitude"`
	EndLongitude                 float64 `json:"endLongitude"`
	MaxVerticalSpeed             float64 `json:"maxVerticalSpeed"`
	MaxRespirationRate           float64 `json:"maxRespirationRate"`
	AvgRespirationRate           float64 `json:"avgRespirationRate"`
	DirectWorkoutComplianceScore *int    `json:"directWorkoutComplianceScore,omitempty"`
	AvgGradeAdjustedSpeed        float64 `json:"avgGradeAdjustedSpeed"`
	LapIndex                     int     `json:"lapIndex"`
	WktStepIndex                 *int    `json:"wktStepIndex,omitempty"`
	WktIndex                     *int    `json:"wktIndex,omitempty"`
	IntensityType                string  `json:"intensityType"`
	MessageIndex                 int     `json:"messageIndex"`
}

// DurationTime returns the lap duration as a time.Duration.
func (l *Lap) DurationTime() time.Duration {
	return time.Duration(l.Duration * float64(time.Second))
}

// DistanceKm returns the lap distance in kilometers.
func (l *Lap) DistanceKm() float64 {
	return l.Distance / 1000
}

// AveragePacePerKm returns the average pace per kilometer for this lap.
func (l *Lap) AveragePacePerKm() time.Duration {
	if l.Distance == 0 {
		return 0
	}
	return time.Duration(l.Duration / (l.Distance / 1000) * float64(time.Second))
}

// ActivitySplits represents the splits/laps data for an activity.
type ActivitySplits struct {
	ActivityID int64           `json:"activityId"`
	LapDTOs    []Lap           `json:"lapDTOs"`
	EventDTOs  []ActivityEvent `json:"eventDTOs"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (s *ActivitySplits) RawJSON() json.RawMessage {
	return s.raw
}

// GetSplits retrieves splits/laps data for a specific activity.
func (s *ActivityService) GetSplits(ctx context.Context, activityID int64) (*ActivitySplits, error) {
	path := fmt.Sprintf("/activity-service/activity/%d/splits", activityID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var splits ActivitySplits
	if err := json.Unmarshal(raw, &splits); err != nil {
		return nil, err
	}
	splits.raw = raw

	return &splits, nil
}

// DownloadFIT downloads the original FIT file for an activity.
func (s *ActivityService) DownloadFIT(ctx context.Context, activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/files/activity/%d", activityID)

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

// DownloadTCX exports and downloads the activity as a TCX file.
func (s *ActivityService) DownloadTCX(ctx context.Context, activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/export/tcx/activity/%d", activityID)

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

// DownloadGPX exports and downloads the activity as a GPX file.
func (s *ActivityService) DownloadGPX(ctx context.Context, activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/export/gpx/activity/%d", activityID)

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

// DownloadKML exports and downloads the activity as a KML file.
func (s *ActivityService) DownloadKML(ctx context.Context, activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/export/kml/activity/%d", activityID)

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

// DownloadCSV exports and downloads the activity as a CSV file.
func (s *ActivityService) DownloadCSV(ctx context.Context, activityID int64) ([]byte, error) {
	path := fmt.Sprintf("/download-service/export/csv/activity/%d", activityID)

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
