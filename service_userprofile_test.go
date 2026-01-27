// service_userprofile_test.go
package garmin

import (
	"encoding/json"
	"testing"
)

func TestSocialProfileJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"id": 12345678,
		"profileId": 12345678,
		"garminGUID": "00000000-0000-0000-0000-000000000000",
		"displayName": "anonymous",
		"fullName": "Anonymous User",
		"userName": "anonymous",
		"profileImageType": "UPLOADED_PHOTO",
		"profileImageUrlLarge": "https://example.com/profile.png",
		"profileImageUrlMedium": "https://example.com/profile.png",
		"profileImageUrlSmall": "https://example.com/profile.png",
		"hasPremiumSocialIcon": false,
		"location": "Anonymous City",
		"facebookUrl": null,
		"twitterUrl": null,
		"personalWebsite": null,
		"motivation": null,
		"bio": null,
		"primaryActivity": null,
		"favoriteActivityTypes": [],
		"runningTrainingSpeed": 0.0,
		"cyclingTrainingSpeed": 0.0,
		"favoriteCyclingActivityTypes": [],
		"cyclingClassification": null,
		"cyclingMaxAvgPower": 0.0,
		"swimmingTrainingSpeed": 0.0,
		"profileVisibility": "private",
		"activityStartVisibility": "public",
		"activityMapVisibility": "public",
		"courseVisibility": "public",
		"activityHeartRateVisibility": "public",
		"activityPowerVisibility": "public",
		"badgeVisibility": "following",
		"showAge": true,
		"showWeight": true,
		"showHeight": true,
		"showWeightClass": false,
		"showAgeRange": false,
		"showGender": true,
		"showActivityClass": false,
		"showVO2Max": true,
		"showPersonalRecords": true,
		"showLast12Months": true,
		"showLifetimeTotals": true,
		"showUpcomingEvents": true,
		"showRecentFavorites": true,
		"showRecentDevice": true,
		"showRecentGear": false,
		"showBadges": true,
		"otherActivity": null,
		"otherPrimaryActivity": null,
		"otherMotivation": null,
		"userRoles": ["ROLE_CONNECTUSER", "ROLE_FITNESS_USER"],
		"nameApproved": true,
		"userProfileFullName": "Anonymous User",
		"makeGolfScorecardsPrivate": true,
		"allowGolfLiveScoring": false,
		"allowGolfScoringByConnections": true,
		"userLevel": 4,
		"userPoint": 203,
		"levelUpdateDate": "2025-08-24T04:43:23.0",
		"levelIsViewed": false,
		"levelPointThreshold": 300,
		"userPointOffset": 0,
		"userPro": false
	}`

	var profile SocialProfile
	if err := json.Unmarshal([]byte(rawJSON), &profile); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if profile.ID != 12345678 {
		t.Errorf("ID = %d, want 12345678", profile.ID)
	}
	if profile.ProfileID != 12345678 {
		t.Errorf("ProfileID = %d, want 12345678", profile.ProfileID)
	}
	if profile.GarminGUID != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("GarminGUID = %s, want 00000000-0000-0000-0000-000000000000", profile.GarminGUID)
	}
	if profile.DisplayName != "anonymous" {
		t.Errorf("DisplayName = %s, want anonymous", profile.DisplayName)
	}
	if profile.FullName != "Anonymous User" {
		t.Errorf("FullName = %s, want Anonymous User", profile.FullName)
	}
	if profile.ProfileVisibility != "private" {
		t.Errorf("ProfileVisibility = %s, want private", profile.ProfileVisibility)
	}
	if profile.UserLevel != 4 {
		t.Errorf("UserLevel = %d, want 4", profile.UserLevel)
	}
	if len(profile.UserRoles) != 2 {
		t.Errorf("UserRoles length = %d, want 2", len(profile.UserRoles))
	}
}

func TestUserSettingsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"id": 12345678,
		"userData": {
			"gender": "MALE",
			"weight": 74200.0,
			"height": 190.0,
			"timeFormat": "time_twelve_hr",
			"birthDate": "1990-01-01",
			"measurementSystem": "metric",
			"activityLevel": 6,
			"handedness": "RIGHT",
			"powerFormat": {
				"formatId": 30,
				"formatKey": "watt",
				"minFraction": 0,
				"maxFraction": 0,
				"groupingUsed": true,
				"displayFormat": null
			},
			"heartRateFormat": {
				"formatId": 21,
				"formatKey": "bpm",
				"minFraction": 0,
				"maxFraction": 0,
				"groupingUsed": false,
				"displayFormat": null
			},
			"firstDayOfWeek": {
				"dayId": 3,
				"dayName": "monday",
				"sortOrder": 3,
				"isPossibleFirstDay": true
			},
			"vo2MaxRunning": 50.0,
			"vo2MaxCycling": null,
			"lactateThresholdSpeed": 0.33611017,
			"lactateThresholdHeartRate": 166,
			"diveNumber": null,
			"intensityMinutesCalcMethod": "AUTO",
			"moderateIntensityMinutesHrZone": 3,
			"vigorousIntensityMinutesHrZone": 4,
			"hydrationMeasurementUnit": "cup",
			"hydrationContainers": [
				{"name": null, "volume": 1, "unit": "cup"}
			],
			"hydrationAutoGoalEnabled": true,
			"firstbeatMaxStressScore": null,
			"firstbeatCyclingLtTimestamp": null,
			"firstbeatRunningLtTimestamp": 1132123569,
			"thresholdHeartRateAutoDetected": true,
			"ftpAutoDetected": true,
			"trainingStatusPausedDate": null,
			"weatherLocation": {
				"useFixedLocation": null,
				"latitude": null,
				"longitude": null,
				"locationName": null,
				"isoCountryCode": null,
				"postalCode": null
			},
			"golfDistanceUnit": "statute_us",
			"golfElevationUnit": null,
			"golfSpeedUnit": null,
			"externalBottomTime": null,
			"availableTrainingDays": ["FRIDAY", "MONDAY"],
			"preferredLongTrainingDays": ["SUNDAY"],
			"virtualCaddieDataSource": null,
			"numberDivesAutomatically": null
		},
		"userSleep": {
			"sleepTime": 82800,
			"defaultSleepTime": false,
			"wakeTime": 25200,
			"defaultWakeTime": false
		},
		"connectDate": null,
		"sourceType": null,
		"userSleepWindows": [
			{
				"sleepWindowFrequency": "SUNDAY",
				"startSleepTimeSecondsFromMidnight": 82800,
				"endSleepTimeSecondsFromMidnight": 25200
			}
		]
	}`

	var settings UserSettings
	if err := json.Unmarshal([]byte(rawJSON), &settings); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if settings.ID != 12345678 {
		t.Errorf("ID = %d, want 12345678", settings.ID)
	}
	if settings.UserData.Gender != "MALE" {
		t.Errorf("Gender = %s, want MALE", settings.UserData.Gender)
	}
	if settings.UserData.Weight != 74200.0 {
		t.Errorf("Weight = %f, want 74200.0", settings.UserData.Weight)
	}
	if settings.UserData.Height != 190.0 {
		t.Errorf("Height = %f, want 190.0", settings.UserData.Height)
	}
	if settings.UserData.MeasurementSystem != "metric" {
		t.Errorf("MeasurementSystem = %s, want metric", settings.UserData.MeasurementSystem)
	}
	if settings.UserData.FirstDayOfWeek.DayName != "monday" {
		t.Errorf("FirstDayOfWeek.DayName = %s, want monday", settings.UserData.FirstDayOfWeek.DayName)
	}
	if settings.UserData.VO2MaxRunning == nil || *settings.UserData.VO2MaxRunning != 50.0 {
		t.Errorf("VO2MaxRunning = %v, want 50.0", settings.UserData.VO2MaxRunning)
	}
	if settings.UserSleep.SleepTime != 82800 {
		t.Errorf("SleepTime = %d, want 82800", settings.UserSleep.SleepTime)
	}
	if len(settings.UserSleepWindows) != 1 {
		t.Errorf("UserSleepWindows length = %d, want 1", len(settings.UserSleepWindows))
	}
}

func TestProfileSettingsJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"displayName": "anonymous",
		"preferredLocale": "fr",
		"measurementSystem": "metric",
		"firstDayOfWeek": {
			"dayId": 3,
			"dayName": "monday",
			"sortOrder": 3,
			"isPossibleFirstDay": true
		},
		"numberFormat": "decimal_period",
		"timeFormat": {
			"formatId": 32,
			"formatKey": "time_twelve_hr",
			"minFraction": 0,
			"maxFraction": 0,
			"groupingUsed": false,
			"displayFormat": "h:mm a"
		},
		"dateFormat": {
			"formatId": 23,
			"formatKey": "mmddyyyy",
			"minFraction": 0,
			"maxFraction": 0,
			"groupingUsed": false,
			"displayFormat": "EEE, MMM d, yyyy"
		},
		"powerFormat": {
			"formatId": 30,
			"formatKey": "watt",
			"minFraction": 0,
			"maxFraction": 0,
			"groupingUsed": true,
			"displayFormat": null
		},
		"heartRateFormat": {
			"formatId": 21,
			"formatKey": "bpm",
			"minFraction": 0,
			"maxFraction": 0,
			"groupingUsed": false,
			"displayFormat": null
		},
		"timeZone": "Asia/Dubai",
		"hydrationMeasurementUnit": "cup",
		"hydrationContainers": [
			{"name": null, "volume": 1, "unit": "cup"}
		],
		"golfDistanceUnit": "statute_us",
		"golfElevationUnit": null,
		"golfSpeedUnit": null,
		"availableTrainingDays": null,
		"preferredLongTrainingDays": null
	}`

	var settings ProfileSettings
	if err := json.Unmarshal([]byte(rawJSON), &settings); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if settings.DisplayName != "anonymous" {
		t.Errorf("DisplayName = %s, want anonymous", settings.DisplayName)
	}
	if settings.PreferredLocale != "fr" {
		t.Errorf("PreferredLocale = %s, want fr", settings.PreferredLocale)
	}
	if settings.MeasurementSystem != "metric" {
		t.Errorf("MeasurementSystem = %s, want metric", settings.MeasurementSystem)
	}
	if settings.TimeZone != "Asia/Dubai" {
		t.Errorf("TimeZone = %s, want Asia/Dubai", settings.TimeZone)
	}
	if settings.FirstDayOfWeek.DayName != "monday" {
		t.Errorf("FirstDayOfWeek.DayName = %s, want monday", settings.FirstDayOfWeek.DayName)
	}
	if settings.TimeFormat.FormatKey != "time_twelve_hr" {
		t.Errorf("TimeFormat.FormatKey = %s, want time_twelve_hr", settings.TimeFormat.FormatKey)
	}
	if settings.DateFormat.FormatKey != "mmddyyyy" {
		t.Errorf("DateFormat.FormatKey = %s, want mmddyyyy", settings.DateFormat.FormatKey)
	}
}

func TestSocialProfileRawJSON(t *testing.T) {
	rawJSON := `{"id":12345678,"displayName":"anonymous"}`

	var profile SocialProfile
	if err := json.Unmarshal([]byte(rawJSON), &profile); err != nil {
		t.Fatal(err)
	}

	profile.raw = json.RawMessage(rawJSON)

	if string(profile.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
