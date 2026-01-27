package garmin

import (
	"encoding/json"
	"testing"
)

func TestSportTypeJSONUnmarshal(t *testing.T) {
	const sportTypeRunning = "running"
	rawJSON := `{"sportTypeId":1,"sportTypeKey":"running","displayOrder":1}`

	var st SportType
	if err := json.Unmarshal([]byte(rawJSON), &st); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if st.SportTypeID != 1 {
		t.Errorf("SportTypeID = %d, want 1", st.SportTypeID)
	}
	if st.SportTypeKey != sportTypeRunning {
		t.Errorf("SportTypeKey = %s, want running", st.SportTypeKey)
	}
	if st.DisplayOrder != 1 {
		t.Errorf("DisplayOrder = %d, want 1", st.DisplayOrder)
	}
}

func TestWorkoutStepJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"type": "ExecutableStepDTO",
		"stepId": 12345,
		"stepOrder": 1,
		"stepType": {"stepTypeId": 1, "stepTypeKey": "warmup"},
		"endCondition": {"conditionTypeId": 2, "conditionTypeKey": "time"},
		"endConditionValue": 300.0,
		"targetType": {"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"}
	}`

	var step WorkoutStep
	if err := json.Unmarshal([]byte(rawJSON), &step); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if step.Type != "ExecutableStepDTO" {
		t.Errorf("Type = %s, want ExecutableStepDTO", step.Type)
	}
	if step.StepID != 12345 {
		t.Errorf("StepID = %d, want 12345", step.StepID)
	}
	if step.StepOrder != 1 {
		t.Errorf("StepOrder = %d, want 1", step.StepOrder)
	}
	if step.StepType == nil || step.StepType.StepTypeKey != "warmup" {
		t.Errorf("StepType.StepTypeKey = %v, want warmup", step.StepType)
	}
	if step.EndCondition == nil || step.EndCondition.ConditionTypeKey != "time" {
		t.Errorf("EndCondition.ConditionTypeKey = %v, want time", step.EndCondition)
	}
	if step.EndConditionValue == nil || *step.EndConditionValue != 300.0 {
		t.Errorf("EndConditionValue = %v, want 300.0", step.EndConditionValue)
	}
}

func TestWorkoutStepRepeatGroupJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"type": "RepeatGroupDTO",
		"stepOrder": 2,
		"numberOfIterations": 5,
		"workoutSteps": [
			{"type": "ExecutableStepDTO", "stepOrder": 1},
			{"type": "ExecutableStepDTO", "stepOrder": 2}
		],
		"smartRepeat": false
	}`

	var step WorkoutStep
	if err := json.Unmarshal([]byte(rawJSON), &step); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if step.Type != "RepeatGroupDTO" {
		t.Errorf("Type = %s, want RepeatGroupDTO", step.Type)
	}
	if step.NumberOfIterations == nil || *step.NumberOfIterations != 5 {
		t.Errorf("NumberOfIterations = %v, want 5", step.NumberOfIterations)
	}
	if len(step.WorkoutSteps) != 2 {
		t.Errorf("len(WorkoutSteps) = %d, want 2", len(step.WorkoutSteps))
	}
}

func TestWorkoutSegmentJSONUnmarshal(t *testing.T) {
	const sportTypeRunning = "running"
	rawJSON := `{
		"segmentOrder": 1,
		"sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
		"workoutSteps": [
			{"type": "ExecutableStepDTO", "stepOrder": 1}
		]
	}`

	var seg WorkoutSegment
	if err := json.Unmarshal([]byte(rawJSON), &seg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if seg.SegmentOrder != 1 {
		t.Errorf("SegmentOrder = %d, want 1", seg.SegmentOrder)
	}
	if seg.SportType.SportTypeKey != sportTypeRunning {
		t.Errorf("SportType.SportTypeKey = %s, want running", seg.SportType.SportTypeKey)
	}
	if len(seg.WorkoutSteps) != 1 {
		t.Errorf("len(WorkoutSteps) = %d, want 1", len(seg.WorkoutSteps))
	}
}

func TestWorkoutJSONUnmarshal(t *testing.T) {
	const sportTypeRunning = "running"
	rawJSON := `{
		"workoutId": 987654321,
		"ownerId": 12345678,
		"workoutName": "Easy Run",
		"description": "30 minute easy run",
		"sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
		"estimatedDurationInSecs": 1800,
		"workoutSegments": [
			{
				"segmentOrder": 1,
				"sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
				"workoutSteps": []
			}
		],
		"shared": false
	}`

	var workout Workout
	if err := json.Unmarshal([]byte(rawJSON), &workout); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if workout.WorkoutID != 987654321 {
		t.Errorf("WorkoutID = %d, want 987654321", workout.WorkoutID)
	}
	if workout.OwnerID != 12345678 {
		t.Errorf("OwnerID = %d, want 12345678", workout.OwnerID)
	}
	if workout.WorkoutName != "Easy Run" {
		t.Errorf("WorkoutName = %s, want Easy Run", workout.WorkoutName)
	}
	if workout.Description != "30 minute easy run" {
		t.Errorf("Description = %s, want 30 minute easy run", workout.Description)
	}
	if workout.EstimatedDurationInSecs != 1800 {
		t.Errorf("EstimatedDurationInSecs = %d, want 1800", workout.EstimatedDurationInSecs)
	}
	if workout.SportType.SportTypeKey != sportTypeRunning {
		t.Errorf("SportType.SportTypeKey = %s, want running", workout.SportType.SportTypeKey)
	}
	if len(workout.WorkoutSegments) != 1 {
		t.Errorf("len(WorkoutSegments) = %d, want 1", len(workout.WorkoutSegments))
	}
	if workout.Shared {
		t.Error("Shared = true, want false")
	}
}

func TestWorkoutRawJSON(t *testing.T) {
	rawJSON := `{"workoutId":123}`
	workout := &Workout{raw: json.RawMessage(rawJSON)}

	if string(workout.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestWorkoutSummaryJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"workoutId": 123456,
		"ownerId": 789,
		"workoutName": "Interval Training",
		"sportType": {"sportTypeId": 1, "sportTypeKey": "running"},
		"estimatedDurationInSecs": 3600,
		"shared": true
	}`

	var summary WorkoutSummary
	if err := json.Unmarshal([]byte(rawJSON), &summary); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if summary.WorkoutID != 123456 {
		t.Errorf("WorkoutID = %d, want 123456", summary.WorkoutID)
	}
	if summary.WorkoutName != "Interval Training" {
		t.Errorf("WorkoutName = %s, want Interval Training", summary.WorkoutName)
	}
	if !summary.Shared {
		t.Error("Shared = false, want true")
	}
}

func TestWorkoutSummaryRawJSON(t *testing.T) {
	rawJSON := `{"workoutId":456}`
	summary := &WorkoutSummary{raw: json.RawMessage(rawJSON)}

	if string(summary.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestWorkoutListRawJSON(t *testing.T) {
	rawJSON := `[{"workoutId":1},{"workoutId":2}]`
	list := &WorkoutList{raw: json.RawMessage(rawJSON)}

	if string(list.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestScheduledWorkoutJSONUnmarshal(t *testing.T) {
	rawJSON := `{
		"workoutScheduleId": 111222,
		"workoutId": 333444,
		"workoutName": "Morning Run",
		"date": "2026-01-27",
		"calendarDate": "2026-01-27T00:00:00.000"
	}`

	var scheduled ScheduledWorkout
	if err := json.Unmarshal([]byte(rawJSON), &scheduled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if scheduled.WorkoutScheduleID != 111222 {
		t.Errorf("WorkoutScheduleID = %d, want 111222", scheduled.WorkoutScheduleID)
	}
	if scheduled.WorkoutID != 333444 {
		t.Errorf("WorkoutID = %d, want 333444", scheduled.WorkoutID)
	}
	if scheduled.WorkoutName != "Morning Run" {
		t.Errorf("WorkoutName = %s, want Morning Run", scheduled.WorkoutName)
	}
	if scheduled.Date != "2026-01-27" {
		t.Errorf("Date = %s, want 2026-01-27", scheduled.Date)
	}
}

func TestScheduledWorkoutRawJSON(t *testing.T) {
	rawJSON := `{"workoutScheduleId":789}`
	scheduled := &ScheduledWorkout{raw: json.RawMessage(rawJSON)}

	if string(scheduled.RawJSON()) != rawJSON {
		t.Error("RawJSON should return original JSON")
	}
}

func TestWorkoutJSONMarshal(t *testing.T) {
	const sportTypeRunning = "running"
	workout := &Workout{
		WorkoutName: "Test Workout",
		SportType: SportType{
			SportTypeID:  1,
			SportTypeKey: sportTypeRunning,
		},
		EstimatedDurationInSecs: 1800,
		WorkoutSegments: []WorkoutSegment{
			{
				SegmentOrder: 1,
				SportType: SportType{
					SportTypeID:  1,
					SportTypeKey: sportTypeRunning,
				},
				WorkoutSteps: []WorkoutStep{
					{
						Type:      "ExecutableStepDTO",
						StepOrder: 1,
					},
				},
			},
		},
	}

	data, err := json.Marshal(workout)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify it can be unmarshaled back
	var unmarshaled Workout
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.WorkoutName != workout.WorkoutName {
		t.Errorf("WorkoutName = %s, want %s", unmarshaled.WorkoutName, workout.WorkoutName)
	}
	if unmarshaled.EstimatedDurationInSecs != workout.EstimatedDurationInSecs {
		t.Errorf("EstimatedDurationInSecs = %d, want %d", unmarshaled.EstimatedDurationInSecs, workout.EstimatedDurationInSecs)
	}
}
