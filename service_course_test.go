package garmin

import (
	"encoding/json"
	"testing"
)

func TestCourseJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"courseId": 432134584,
		"userProfileId": 12345678,
		"displayName": "anonymous",
		"userGroupId": null,
		"geoRoutePk": null,
		"activityType": {
			"typeId": 3,
			"typeKey": "hiking",
			"parentTypeId": 17,
			"isHidden": false,
			"restricted": false,
			"trimmable": false
		},
		"courseName": "Test Course",
		"courseDescription": null,
		"createdDate": 1770967543000,
		"updatedDate": 1770967543000,
		"privacyRule": {
			"typeId": 2,
			"typeKey": "private"
		},
		"distanceInMeters": 7217.69,
		"elevationGainInMeters": 277.86,
		"elevationLossInMeters": 280.95,
		"startLatitude": -21.3136395,
		"startLongitude": 55.5420436,
		"speedInMetersPerSecond": 0.0,
		"sourceTypeId": 3,
		"sourcePk": null,
		"elapsedSeconds": null,
		"coordinateSystem": "WGS84",
		"originalCoordinateSystem": "WGS84",
		"consumer": null,
		"elevationSource": 3,
		"hasShareableEvent": false,
		"hasPaceBand": false,
		"hasPowerGuide": false,
		"favorite": false,
		"hasTurnDetectionDisabled": false,
		"curatedCourseId": null,
		"startNote": null,
		"finishNote": null,
		"cutoffDuration": null,
		"createdDateFormatted": "2026-02-13 07:25:43.0 GMT",
		"updatedDateFormatted": "2026-02-13 07:25:43.0 GMT",
		"activityTypeId": {
			"typeId": 3,
			"typeKey": "hiking",
			"parentTypeId": 17,
			"isHidden": false,
			"restricted": false,
			"trimmable": false
		},
		"public": false
	}`

	var course Course
	if err := json.Unmarshal([]byte(rawJSON), &course); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if course.CourseID != 432134584 {
		t.Errorf("CourseID = %d, want 432134584", course.CourseID)
	}
	if course.UserProfileID != 12345678 {
		t.Errorf("UserProfileID = %d, want 12345678", course.UserProfileID)
	}
	if course.DisplayName != testAnonymousName {
		t.Errorf("DisplayName = %s, want %s", course.DisplayName, testAnonymousName)
	}
	if course.CourseName != "Test Course" {
		t.Errorf("CourseName = %s, want Test Course", course.CourseName)
	}
	if course.CourseDescription != nil {
		t.Errorf("CourseDescription = %v, want nil", course.CourseDescription)
	}
	if course.ActivityType.TypeKey != "hiking" {
		t.Errorf("ActivityType.TypeKey = %s, want hiking", course.ActivityType.TypeKey)
	}
	if course.ActivityType.TypeID != 3 {
		t.Errorf("ActivityType.TypeID = %d, want 3", course.ActivityType.TypeID)
	}
	if course.CreatedDate != 1770967543000 {
		t.Errorf("CreatedDate = %d, want 1770967543000", course.CreatedDate)
	}
	if course.PrivacyRule.TypeKey != testPrivateTypeKey {
		t.Errorf("PrivacyRule.TypeKey = %s, want %s", course.PrivacyRule.TypeKey, testPrivateTypeKey)
	}
	if course.DistanceInMeters != 7217.69 {
		t.Errorf("DistanceInMeters = %f, want 7217.69", course.DistanceInMeters)
	}
	if course.ElevationGainInMeters != 277.86 {
		t.Errorf("ElevationGainInMeters = %f, want 277.86", course.ElevationGainInMeters)
	}
	if course.ElevationLossInMeters != 280.95 {
		t.Errorf("ElevationLossInMeters = %f, want 280.95", course.ElevationLossInMeters)
	}
	if course.StartLatitude != -21.3136395 {
		t.Errorf("StartLatitude = %f, want -21.3136395", course.StartLatitude)
	}
	if course.StartLongitude != 55.5420436 {
		t.Errorf("StartLongitude = %f, want 55.5420436", course.StartLongitude)
	}
	if course.CoordinateSystem != "WGS84" {
		t.Errorf("CoordinateSystem = %s, want WGS84", course.CoordinateSystem)
	}
	if course.Favorite {
		t.Error("Favorite = true, want false")
	}
	if course.Public {
		t.Error("Public = true, want false")
	}
	if course.HasShareableEvent {
		t.Error("HasShareableEvent = true, want false")
	}
}

func TestCourseRawJSON(t *testing.T) {
	rawJSON := `{"courseId":123,"courseName":"Test"}`

	var course Course
	if err := json.Unmarshal([]byte(rawJSON), &course); err != nil {
		t.Fatal(err)
	}
	course.raw = json.RawMessage(rawJSON)

	if string(course.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestCoursesForUserResponseJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"coursesForUser": [
			{
				"courseId": 111,
				"courseName": "Course 1",
				"distanceInMeters": 5000.0,
				"public": false
			},
			{
				"courseId": 222,
				"courseName": "Course 2",
				"distanceInMeters": 10000.0,
				"public": true
			}
		]
	}`

	var resp CoursesForUserResponse
	if err := json.Unmarshal([]byte(rawJSON), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(resp.CoursesForUser) != 2 {
		t.Fatalf("len(CoursesForUser) = %d, want 2", len(resp.CoursesForUser))
	}
	if resp.CoursesForUser[0].CourseID != 111 {
		t.Errorf("CoursesForUser[0].CourseID = %d, want 111", resp.CoursesForUser[0].CourseID)
	}
	if resp.CoursesForUser[0].CourseName != "Course 1" {
		t.Errorf("CoursesForUser[0].CourseName = %s, want Course 1", resp.CoursesForUser[0].CourseName)
	}
	if resp.CoursesForUser[1].CourseID != 222 {
		t.Errorf("CoursesForUser[1].CourseID = %d, want 222", resp.CoursesForUser[1].CourseID)
	}
	if !resp.CoursesForUser[1].Public {
		t.Error("CoursesForUser[1].Public = false, want true")
	}
}

func TestCoursesForUserResponseRawJSON(t *testing.T) {
	rawJSON := `{"coursesForUser":[]}`
	resp := &CoursesForUserResponse{raw: json.RawMessage(rawJSON)}

	if string(resp.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
