// service_userprofile.go
package garmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// SocialProfile represents the user's social profile information.
type SocialProfile struct {
	ID                            int64    `json:"id"`
	ProfileID                     int64    `json:"profileId"`
	GarminGUID                    string   `json:"garminGUID"`
	DisplayName                   string   `json:"displayName"`
	FullName                      string   `json:"fullName"`
	UserName                      string   `json:"userName"`
	ProfileImageType              string   `json:"profileImageType"`
	ProfileImageURLLarge          string   `json:"profileImageUrlLarge"`
	ProfileImageURLMedium         string   `json:"profileImageUrlMedium"`
	ProfileImageURLSmall          string   `json:"profileImageUrlSmall"`
	HasPremiumSocialIcon          bool     `json:"hasPremiumSocialIcon"`
	Location                      string   `json:"location"`
	FacebookURL                   *string  `json:"facebookUrl"`
	TwitterURL                    *string  `json:"twitterUrl"`
	PersonalWebsite               *string  `json:"personalWebsite"`
	Motivation                    *string  `json:"motivation"`
	Bio                           *string  `json:"bio"`
	PrimaryActivity               *string  `json:"primaryActivity"`
	FavoriteActivityTypes         []string `json:"favoriteActivityTypes"`
	RunningTrainingSpeed          float64  `json:"runningTrainingSpeed"`
	CyclingTrainingSpeed          float64  `json:"cyclingTrainingSpeed"`
	FavoriteCyclingActivityTypes  []string `json:"favoriteCyclingActivityTypes"`
	CyclingClassification         *string  `json:"cyclingClassification"`
	CyclingMaxAvgPower            float64  `json:"cyclingMaxAvgPower"`
	SwimmingTrainingSpeed         float64  `json:"swimmingTrainingSpeed"`
	ProfileVisibility             string   `json:"profileVisibility"`
	ActivityStartVisibility       string   `json:"activityStartVisibility"`
	ActivityMapVisibility         string   `json:"activityMapVisibility"`
	CourseVisibility              string   `json:"courseVisibility"`
	ActivityHeartRateVisibility   string   `json:"activityHeartRateVisibility"`
	ActivityPowerVisibility       string   `json:"activityPowerVisibility"`
	BadgeVisibility               string   `json:"badgeVisibility"`
	ShowAge                       bool     `json:"showAge"`
	ShowWeight                    bool     `json:"showWeight"`
	ShowHeight                    bool     `json:"showHeight"`
	ShowWeightClass               bool     `json:"showWeightClass"`
	ShowAgeRange                  bool     `json:"showAgeRange"`
	ShowGender                    bool     `json:"showGender"`
	ShowActivityClass             bool     `json:"showActivityClass"`
	ShowVO2Max                    bool     `json:"showVO2Max"`
	ShowPersonalRecords           bool     `json:"showPersonalRecords"`
	ShowLast12Months              bool     `json:"showLast12Months"`
	ShowLifetimeTotals            bool     `json:"showLifetimeTotals"`
	ShowUpcomingEvents            bool     `json:"showUpcomingEvents"`
	ShowRecentFavorites           bool     `json:"showRecentFavorites"`
	ShowRecentDevice              bool     `json:"showRecentDevice"`
	ShowRecentGear                bool     `json:"showRecentGear"`
	ShowBadges                    bool     `json:"showBadges"`
	OtherActivity                 *string  `json:"otherActivity"`
	OtherPrimaryActivity          *string  `json:"otherPrimaryActivity"`
	OtherMotivation               *string  `json:"otherMotivation"`
	UserRoles                     []string `json:"userRoles"`
	NameApproved                  bool     `json:"nameApproved"`
	UserProfileFullName           string   `json:"userProfileFullName"`
	MakeGolfScorecardsPrivate     bool     `json:"makeGolfScorecardsPrivate"`
	AllowGolfLiveScoring          bool     `json:"allowGolfLiveScoring"`
	AllowGolfScoringByConnections bool     `json:"allowGolfScoringByConnections"`
	UserLevel                     int      `json:"userLevel"`
	UserPoint                     int      `json:"userPoint"`
	LevelUpdateDate               string   `json:"levelUpdateDate"`
	LevelIsViewed                 bool     `json:"levelIsViewed"`
	LevelPointThreshold           int      `json:"levelPointThreshold"`
	UserPointOffset               int      `json:"userPointOffset"`
	UserPro                       bool     `json:"userPro"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (s *SocialProfile) RawJSON() json.RawMessage {
	return s.raw
}

// FormatSettings represents display format settings.
type FormatSettings struct {
	FormatID      int     `json:"formatId"`
	FormatKey     string  `json:"formatKey"`
	MinFraction   int     `json:"minFraction"`
	MaxFraction   int     `json:"maxFraction"`
	GroupingUsed  bool    `json:"groupingUsed"`
	DisplayFormat *string `json:"displayFormat"`
}

// DayOfWeek represents a day of the week setting.
type DayOfWeek struct {
	DayID              int    `json:"dayId"`
	DayName            string `json:"dayName"`
	SortOrder          int    `json:"sortOrder"`
	IsPossibleFirstDay bool   `json:"isPossibleFirstDay"`
}

// HydrationContainer represents a hydration container setting.
type HydrationContainer struct {
	Name   *string `json:"name"`
	Volume int     `json:"volume"`
	Unit   string  `json:"unit"`
}

// WeatherLocation represents weather location settings.
type WeatherLocation struct {
	UseFixedLocation *bool    `json:"useFixedLocation"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	LocationName     *string  `json:"locationName"`
	ISOCountryCode   *string  `json:"isoCountryCode"`
	PostalCode       *string  `json:"postalCode"`
}

// UserData represents the user's personal data within settings.
type UserData struct {
	Gender                         string               `json:"gender"`
	Weight                         float64              `json:"weight"`
	Height                         float64              `json:"height"`
	TimeFormat                     string               `json:"timeFormat"`
	BirthDate                      string               `json:"birthDate"`
	MeasurementSystem              string               `json:"measurementSystem"`
	ActivityLevel                  int                  `json:"activityLevel"`
	Handedness                     string               `json:"handedness"`
	PowerFormat                    FormatSettings       `json:"powerFormat"`
	HeartRateFormat                FormatSettings       `json:"heartRateFormat"`
	FirstDayOfWeek                 DayOfWeek            `json:"firstDayOfWeek"`
	VO2MaxRunning                  *float64             `json:"vo2MaxRunning"`
	VO2MaxCycling                  *float64             `json:"vo2MaxCycling"`
	LactateThresholdSpeed          float64              `json:"lactateThresholdSpeed"`
	LactateThresholdHeartRate      int                  `json:"lactateThresholdHeartRate"`
	DiveNumber                     *int                 `json:"diveNumber"`
	IntensityMinutesCalcMethod     string               `json:"intensityMinutesCalcMethod"`
	ModerateIntensityMinutesHrZone int                  `json:"moderateIntensityMinutesHrZone"`
	VigorousIntensityMinutesHrZone int                  `json:"vigorousIntensityMinutesHrZone"`
	HydrationMeasurementUnit       string               `json:"hydrationMeasurementUnit"`
	HydrationContainers            []HydrationContainer `json:"hydrationContainers"`
	HydrationAutoGoalEnabled       bool                 `json:"hydrationAutoGoalEnabled"`
	FirstbeatMaxStressScore        *int                 `json:"firstbeatMaxStressScore"`
	FirstbeatCyclingLtTimestamp    *int64               `json:"firstbeatCyclingLtTimestamp"`
	FirstbeatRunningLtTimestamp    *int64               `json:"firstbeatRunningLtTimestamp"`
	ThresholdHeartRateAutoDetected bool                 `json:"thresholdHeartRateAutoDetected"`
	FTPAutoDetected                bool                 `json:"ftpAutoDetected"`
	TrainingStatusPausedDate       *string              `json:"trainingStatusPausedDate"`
	WeatherLocation                WeatherLocation      `json:"weatherLocation"`
	GolfDistanceUnit               string               `json:"golfDistanceUnit"`
	GolfElevationUnit              *string              `json:"golfElevationUnit"`
	GolfSpeedUnit                  *string              `json:"golfSpeedUnit"`
	ExternalBottomTime             *int                 `json:"externalBottomTime"`
	AvailableTrainingDays          []string             `json:"availableTrainingDays"`
	PreferredLongTrainingDays      []string             `json:"preferredLongTrainingDays"`
	VirtualCaddieDataSource        *string              `json:"virtualCaddieDataSource"`
	NumberDivesAutomatically       *bool                `json:"numberDivesAutomatically"`
}

// UserSleep represents the user's sleep settings.
type UserSleep struct {
	SleepTime        int  `json:"sleepTime"`
	DefaultSleepTime bool `json:"defaultSleepTime"`
	WakeTime         int  `json:"wakeTime"`
	DefaultWakeTime  bool `json:"defaultWakeTime"`
}

// SleepWindow represents a sleep window for a specific day.
type SleepWindow struct {
	SleepWindowFrequency              string `json:"sleepWindowFrequency"`
	StartSleepTimeSecondsFromMidnight int    `json:"startSleepTimeSecondsFromMidnight"`
	EndSleepTimeSecondsFromMidnight   int    `json:"endSleepTimeSecondsFromMidnight"`
}

// UserSettings represents the user settings response.
type UserSettings struct {
	ID               int64         `json:"id"`
	UserData         UserData      `json:"userData"`
	UserSleep        UserSleep     `json:"userSleep"`
	ConnectDate      *string       `json:"connectDate"`
	SourceType       *string       `json:"sourceType"`
	UserSleepWindows []SleepWindow `json:"userSleepWindows"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (u *UserSettings) RawJSON() json.RawMessage {
	return u.raw
}

// ProfileSettings represents the profile settings response.
type ProfileSettings struct {
	DisplayName               string               `json:"displayName"`
	PreferredLocale           string               `json:"preferredLocale"`
	MeasurementSystem         string               `json:"measurementSystem"`
	FirstDayOfWeek            DayOfWeek            `json:"firstDayOfWeek"`
	NumberFormat              string               `json:"numberFormat"`
	TimeFormat                FormatSettings       `json:"timeFormat"`
	DateFormat                FormatSettings       `json:"dateFormat"`
	PowerFormat               FormatSettings       `json:"powerFormat"`
	HeartRateFormat           FormatSettings       `json:"heartRateFormat"`
	TimeZone                  string               `json:"timeZone"`
	HydrationMeasurementUnit  string               `json:"hydrationMeasurementUnit"`
	HydrationContainers       []HydrationContainer `json:"hydrationContainers"`
	GolfDistanceUnit          string               `json:"golfDistanceUnit"`
	GolfElevationUnit         *string              `json:"golfElevationUnit"`
	GolfSpeedUnit             *string              `json:"golfSpeedUnit"`
	AvailableTrainingDays     []string             `json:"availableTrainingDays"`
	PreferredLongTrainingDays []string             `json:"preferredLongTrainingDays"`

	raw json.RawMessage
}

// RawJSON returns the original JSON response.
func (p *ProfileSettings) RawJSON() json.RawMessage {
	return p.raw
}

// GetSocialProfile retrieves the user's social profile.
func (s *UserProfileService) GetSocialProfile(ctx context.Context) (*SocialProfile, error) {
	resp, err := s.client.doAPI(ctx, http.MethodGet, "/userprofile-service/socialProfile", http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile SocialProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}
	profile.raw = body

	return &profile, nil
}

// GetUserSettings retrieves the user's settings.
func (s *UserProfileService) GetUserSettings(ctx context.Context) (*UserSettings, error) {
	resp, err := s.client.doAPI(ctx, http.MethodGet, "/userprofile-service/userprofile/user-settings", http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var settings UserSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, err
	}
	settings.raw = body

	return &settings, nil
}

// GetProfileSettings retrieves the user's profile settings.
func (s *UserProfileService) GetProfileSettings(ctx context.Context) (*ProfileSettings, error) {
	resp, err := s.client.doAPI(ctx, http.MethodGet, "/userprofile-service/userprofile/settings", http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var settings ProfileSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, err
	}
	settings.raw = body

	return &settings, nil
}
