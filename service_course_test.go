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

func TestCourseDetailJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"courseId": 87654321,
		"courseName": "Test Course",
		"description": "A test course description",
		"openStreetMap": false,
		"matchedToSegments": false,
		"userProfilePk": 12345678,
		"userGroupPk": null,
		"rulePK": 2,
		"firstName": "Anonymous",
		"lastName": "User",
		"displayName": "anonymous",
		"geoRoutePk": null,
		"sourceTypeId": 3,
		"sourcePk": null,
		"distanceMeter": 7217.69,
		"elevationGainMeter": 277.86,
		"elevationLossMeter": 280.95,
		"startPoint": {
			"latitude": 48.8566,
			"longitude": 2.3522,
			"elevation": 100.0,
			"distance": null,
			"timestamp": null
		},
		"coursePoints": [
			{
				"coursePointId": 11111111,
				"courseId": 87654321,
				"name": "Start",
				"latitude": 48.8566,
				"longitude": 2.3522,
				"elevation": 100.0,
				"distance": 0.0,
				"pointType": "START",
				"sortOrder": 1,
				"derivedElevation": 100.0
			}
		],
		"boundingBox": {
			"center": null,
			"lowerLeft": {"latitude": 48.0, "longitude": 2.0},
			"upperRight": {"latitude": 49.0, "longitude": 3.0},
			"lowerLeftLatIsSet": true,
			"lowerLeftLongIsSet": true,
			"upperRightLatIsSet": true,
			"upperRightLongIsSet": true
		},
		"hasShareableEvent": false,
		"hasTurnDetectionDisabled": false,
		"activityTypePk": 3,
		"virtualPartnerId": 87654321,
		"includeLaps": false,
		"elapsedSeconds": null,
		"speedMeterPerSecond": null,
		"createDate": "2026-02-13T07:25:43.0",
		"updateDate": "2026-02-13T07:25:43.0",
		"courseLines": [
			{
				"courseId": 87654321,
				"sortOrder": 1,
				"numberOfPoints": 52,
				"distanceInMeters": 1187.61,
				"bearing": 0.0,
				"points": null,
				"coordinateSystem": null,
				"originalCoordinateSystem": null
			}
		],
		"coordinateSystem": "WGS84",
		"targetCoordinateSystem": "WGS84",
		"originalCoordinateSystem": "WGS84",
		"consumer": null,
		"elevationSource": 3,
		"hasPaceBand": false,
		"hasPowerGuide": false,
		"favorite": false,
		"startNote": null,
		"finishNote": null,
		"cutoffDuration": null,
		"geoPoints": [
			{"latitude": 48.8566, "longitude": 2.3522, "elevation": 100.0, "distance": 0.0, "timestamp": 1770967543000},
			{"latitude": 48.8567, "longitude": 2.3523, "elevation": 101.0, "distance": 25.5, "timestamp": 1770967573596}
		]
	}`

	var detail CourseDetail
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Basic fields
	if detail.CourseID != 87654321 {
		t.Errorf("CourseID = %d, want 87654321", detail.CourseID)
	}
	if detail.CourseName != "Test Course" {
		t.Errorf("CourseName = %s, want Test Course", detail.CourseName)
	}
	if detail.Description == nil || *detail.Description != "A test course description" {
		t.Errorf("Description = %v, want 'A test course description'", detail.Description)
	}
	if detail.OpenStreetMap {
		t.Error("OpenStreetMap = true, want false")
	}
	if detail.MatchedToSegments {
		t.Error("MatchedToSegments = true, want false")
	}
	if detail.UserProfilePK != 12345678 {
		t.Errorf("UserProfilePK = %d, want 12345678", detail.UserProfilePK)
	}
	if detail.RulePK != 2 {
		t.Errorf("RulePK = %d, want 2", detail.RulePK)
	}
	if detail.FirstName != "Anonymous" {
		t.Errorf("FirstName = %s, want Anonymous", detail.FirstName)
	}
	if detail.LastName != "User" {
		t.Errorf("LastName = %s, want User", detail.LastName)
	}
	if detail.DisplayName != testAnonymousName {
		t.Errorf("DisplayName = %s, want %s", detail.DisplayName, testAnonymousName)
	}

	// Distance/elevation (different field names from list response)
	if detail.DistanceMeter != 7217.69 {
		t.Errorf("DistanceMeter = %f, want 7217.69", detail.DistanceMeter)
	}
	if detail.ElevationGainMeter != 277.86 {
		t.Errorf("ElevationGainMeter = %f, want 277.86", detail.ElevationGainMeter)
	}
	if detail.ElevationLossMeter != 280.95 {
		t.Errorf("ElevationLossMeter = %f, want 280.95", detail.ElevationLossMeter)
	}

	// StartPoint
	if detail.StartPoint.Latitude != 48.8566 {
		t.Errorf("StartPoint.Latitude = %f, want 48.8566", detail.StartPoint.Latitude)
	}
	if detail.StartPoint.Longitude != 2.3522 {
		t.Errorf("StartPoint.Longitude = %f, want 2.3522", detail.StartPoint.Longitude)
	}
	if detail.StartPoint.Elevation != 100.0 {
		t.Errorf("StartPoint.Elevation = %f, want 100.0", detail.StartPoint.Elevation)
	}

	// CoursePoints
	if len(detail.CoursePoints) != 1 {
		t.Fatalf("len(CoursePoints) = %d, want 1", len(detail.CoursePoints))
	}
	if detail.CoursePoints[0].CoursePointID != 11111111 {
		t.Errorf("CoursePoints[0].CoursePointID = %d, want 11111111", detail.CoursePoints[0].CoursePointID)
	}
	if detail.CoursePoints[0].Name != "Start" {
		t.Errorf("CoursePoints[0].Name = %s, want Start", detail.CoursePoints[0].Name)
	}
	if detail.CoursePoints[0].PointType != "START" {
		t.Errorf("CoursePoints[0].PointType = %s, want START", detail.CoursePoints[0].PointType)
	}

	// BoundingBox
	if detail.BoundingBox.LowerLeft.Latitude != 48.0 {
		t.Errorf("BoundingBox.LowerLeft.Latitude = %f, want 48.0", detail.BoundingBox.LowerLeft.Latitude)
	}
	if detail.BoundingBox.UpperRight.Longitude != 3.0 {
		t.Errorf("BoundingBox.UpperRight.Longitude = %f, want 3.0", detail.BoundingBox.UpperRight.Longitude)
	}

	// Activity/partner
	if detail.ActivityTypePK != 3 {
		t.Errorf("ActivityTypePK = %d, want 3", detail.ActivityTypePK)
	}
	if detail.VirtualPartnerID != 87654321 {
		t.Errorf("VirtualPartnerID = %d, want 87654321", detail.VirtualPartnerID)
	}

	// Dates (string format, not timestamp)
	if detail.CreateDate != "2026-02-13T07:25:43.0" {
		t.Errorf("CreateDate = %s, want 2026-02-13T07:25:43.0", detail.CreateDate)
	}
	if detail.UpdateDate != "2026-02-13T07:25:43.0" {
		t.Errorf("UpdateDate = %s, want 2026-02-13T07:25:43.0", detail.UpdateDate)
	}

	// CourseLines
	if len(detail.CourseLines) != 1 {
		t.Fatalf("len(CourseLines) = %d, want 1", len(detail.CourseLines))
	}
	if detail.CourseLines[0].SortOrder != 1 {
		t.Errorf("CourseLines[0].SortOrder = %d, want 1", detail.CourseLines[0].SortOrder)
	}
	if detail.CourseLines[0].NumberOfPoints != 52 {
		t.Errorf("CourseLines[0].NumberOfPoints = %d, want 52", detail.CourseLines[0].NumberOfPoints)
	}
	if detail.CourseLines[0].DistanceInMeters != 1187.61 {
		t.Errorf("CourseLines[0].DistanceInMeters = %f, want 1187.61", detail.CourseLines[0].DistanceInMeters)
	}

	// GeoPoints
	if len(detail.GeoPoints) != 2 {
		t.Fatalf("len(GeoPoints) = %d, want 2", len(detail.GeoPoints))
	}
	if detail.GeoPoints[0].Distance != 0.0 {
		t.Errorf("GeoPoints[0].Distance = %f, want 0.0", detail.GeoPoints[0].Distance)
	}
	if detail.GeoPoints[1].Distance != 25.5 {
		t.Errorf("GeoPoints[1].Distance = %f, want 25.5", detail.GeoPoints[1].Distance)
	}
	if detail.GeoPoints[0].Timestamp != 1770967543000 {
		t.Errorf("GeoPoints[0].Timestamp = %d, want 1770967543000", detail.GeoPoints[0].Timestamp)
	}

	// Boolean flags
	if detail.IncludeLaps {
		t.Error("IncludeLaps = true, want false")
	}
	if detail.Favorite {
		t.Error("Favorite = true, want false")
	}
	if detail.CoordinateSystem != "WGS84" {
		t.Errorf("CoordinateSystem = %s, want WGS84", detail.CoordinateSystem)
	}
}

func TestCourseDetailNullCoursePoints(t *testing.T) {
	rawJSON := `{
		"courseId": 123,
		"courseName": "Test",
		"coursePoints": null,
		"geoPoints": null,
		"courseLines": null
	}`

	var detail CourseDetail
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if detail.CoursePoints != nil {
		t.Errorf("CoursePoints = %v, want nil", detail.CoursePoints)
	}
	if detail.GeoPoints != nil {
		t.Errorf("GeoPoints = %v, want nil", detail.GeoPoints)
	}
	if detail.CourseLines != nil {
		t.Errorf("CourseLines = %v, want nil", detail.CourseLines)
	}
}

func TestCourseDetailRawJSON(t *testing.T) {
	rawJSON := `{"courseId":123,"courseName":"Test"}`

	var detail CourseDetail
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		t.Fatal(err)
	}
	detail.raw = json.RawMessage(rawJSON)

	if string(detail.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}
