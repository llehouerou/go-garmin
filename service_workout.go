package garmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SportType represents a sport type for workouts.
type SportType struct {
	SportTypeID  int    `json:"sportTypeId"`
	SportTypeKey string `json:"sportTypeKey"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// StepTypeInfo represents a workout step type.
type StepTypeInfo struct {
	StepTypeID   int    `json:"stepTypeId"`
	StepTypeKey  string `json:"stepTypeKey"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// EndCondition represents the end condition of a workout step.
type EndCondition struct {
	ConditionTypeID  int    `json:"conditionTypeId"`
	ConditionTypeKey string `json:"conditionTypeKey"`
	DisplayOrder     int    `json:"displayOrder,omitempty"`
	Displayable      bool   `json:"displayable,omitempty"`
}

// TargetType represents the target type of a workout step.
type TargetType struct {
	WorkoutTargetTypeID  int    `json:"workoutTargetTypeId"`
	WorkoutTargetTypeKey string `json:"workoutTargetTypeKey"`
	DisplayOrder         int    `json:"displayOrder,omitempty"`
}

// WorkoutStep represents a single step in a workout.
type WorkoutStep struct {
	Type              string        `json:"type"` // ExecutableStepDTO or RepeatGroupDTO
	StepID            int64         `json:"stepId,omitempty"`
	StepOrder         int           `json:"stepOrder"`
	StepType          *StepTypeInfo `json:"stepType,omitempty"`
	EndCondition      *EndCondition `json:"endCondition,omitempty"`
	EndConditionValue *float64      `json:"endConditionValue,omitempty"`
	TargetType        *TargetType   `json:"targetType,omitempty"`
	TargetValueOne    *float64      `json:"targetValueOne,omitempty"`
	TargetValueTwo    *float64      `json:"targetValueTwo,omitempty"`
	ChildStepID       *int64        `json:"childStepId,omitempty"`

	// For repeat groups
	NumberOfIterations *int          `json:"numberOfIterations,omitempty"`
	WorkoutSteps       []WorkoutStep `json:"workoutSteps,omitempty"`
	SmartRepeat        bool          `json:"smartRepeat,omitempty"`
}

// WorkoutSegment represents a segment of a workout.
type WorkoutSegment struct {
	SegmentOrder int           `json:"segmentOrder"`
	SportType    SportType     `json:"sportType"`
	WorkoutSteps []WorkoutStep `json:"workoutSteps"`
}

// WorkoutAuthor represents the author of a workout.
type WorkoutAuthor struct {
	UserProfilePK *int64 `json:"userProfilePk,omitempty"`
}

// Workout represents a Garmin workout.
type Workout struct {
	WorkoutID               int64            `json:"workoutId,omitempty"`
	OwnerID                 int64            `json:"ownerId,omitempty"`
	WorkoutName             string           `json:"workoutName"`
	Description             string           `json:"description,omitempty"`
	UpdateDate              string           `json:"updateDate,omitempty"`
	CreatedDate             string           `json:"createdDate,omitempty"`
	SportType               SportType        `json:"sportType"`
	TrainingPlanID          *int64           `json:"trainingPlanId,omitempty"`
	Author                  *WorkoutAuthor   `json:"author,omitempty"`
	EstimatedDurationInSecs int              `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMtrs float64          `json:"estimatedDistanceInMeters,omitempty"`
	AvgTrainingSpeed        float64          `json:"avgTrainingSpeed,omitempty"`
	WorkoutSegments         []WorkoutSegment `json:"workoutSegments"`
	Locale                  string           `json:"locale,omitempty"`
	PoolLength              float64          `json:"poolLength,omitempty"`
	PoolLengthUnit          string           `json:"poolLengthUnit,omitempty"`
	WorkoutProvider         string           `json:"workoutProvider,omitempty"`
	WorkoutSourceID         string           `json:"workoutSourceId,omitempty"`
	Consumer                string           `json:"consumer,omitempty"`
	AtpPlanID               *int64           `json:"atpPlanId,omitempty"`
	Shared                  bool             `json:"shared,omitempty"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (w *Workout) RawJSON() json.RawMessage {
	return w.raw
}

// WorkoutSummary represents a workout in the list response.
type WorkoutSummary struct {
	WorkoutID               int64     `json:"workoutId"`
	OwnerID                 int64     `json:"ownerId"`
	WorkoutName             string    `json:"workoutName"`
	Description             string    `json:"description,omitempty"`
	UpdateDate              string    `json:"updateDate,omitempty"`
	CreatedDate             string    `json:"createdDate,omitempty"`
	SportType               SportType `json:"sportType"`
	EstimatedDurationInSecs int       `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMtrs float64   `json:"estimatedDistanceInMeters,omitempty"`
	AvgTrainingSpeed        float64   `json:"avgTrainingSpeed,omitempty"`
	Shared                  bool      `json:"shared,omitempty"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (w *WorkoutSummary) RawJSON() json.RawMessage {
	return w.raw
}

// WorkoutList represents a list of workouts.
type WorkoutList struct {
	Workouts []WorkoutSummary
	raw      json.RawMessage
}

// RawJSON returns the raw JSON response.
func (w *WorkoutList) RawJSON() json.RawMessage {
	return w.raw
}

// ScheduledWorkout represents a scheduled workout.
type ScheduledWorkout struct {
	WorkoutScheduleID int64  `json:"workoutScheduleId"`
	WorkoutID         int64  `json:"workoutId"`
	WorkoutName       string `json:"workoutName,omitempty"`
	Date              string `json:"date"` // YYYY-MM-DD
	CalendarDate      string `json:"calendarDate,omitempty"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (s *ScheduledWorkout) RawJSON() json.RawMessage {
	return s.raw
}

// List returns a list of workouts with pagination.
func (s *WorkoutService) List(ctx context.Context, start, limit int) (*WorkoutList, error) {
	path := fmt.Sprintf("/workout-service/workouts?start=%d&limit=%d", start, limit)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list workouts: status %d: %s", resp.StatusCode, string(body))
	}

	var summaries []WorkoutSummary
	if err := json.Unmarshal(body, &summaries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workouts: %w", err)
	}

	// Store raw for each
	for i := range summaries {
		summaries[i].raw = body
	}

	return &WorkoutList{
		Workouts: summaries,
		raw:      body,
	}, nil
}

// Get returns a workout by ID.
func (s *WorkoutService) Get(ctx context.Context, workoutID int64) (*Workout, error) {
	path := fmt.Sprintf("/workout-service/workout/%d", workoutID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workout: status %d: %s", resp.StatusCode, string(body))
	}

	var workout Workout
	if err := json.Unmarshal(body, &workout); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workout: %w", err)
	}
	workout.raw = body

	return &workout, nil
}

// Create creates a new workout and returns the created workout.
func (s *WorkoutService) Create(ctx context.Context, workout *Workout) (*Workout, error) {
	path := "/workout-service/workout"

	payload, err := json.Marshal(workout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workout: %w", err)
	}

	resp, err := s.client.doAPIWithBody(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to create workout: status %d: %s", resp.StatusCode, string(body))
	}

	var created Workout
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created workout: %w", err)
	}
	created.raw = body

	return &created, nil
}

// Update updates an existing workout.
func (s *WorkoutService) Update(ctx context.Context, workoutID int64, workout *Workout) (*Workout, error) {
	path := fmt.Sprintf("/workout-service/workout/%d", workoutID)

	payload, err := json.Marshal(workout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workout: %w", err)
	}

	resp, err := s.client.doAPIWithBody(ctx, http.MethodPut, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to update workout: status %d: %s", resp.StatusCode, string(body))
	}

	var updated Workout
	if err := json.Unmarshal(body, &updated); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated workout: %w", err)
	}
	updated.raw = body

	return &updated, nil
}

// Delete deletes a workout by ID.
func (s *WorkoutService) Delete(ctx context.Context, workoutID int64) error {
	path := fmt.Sprintf("/workout-service/workout/%d", workoutID)

	resp, err := s.client.doAPI(ctx, http.MethodDelete, path, http.NoBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete workout: status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DownloadFIT downloads a workout as a FIT file.
func (s *WorkoutService) DownloadFIT(ctx context.Context, workoutID int64) ([]byte, error) {
	path := fmt.Sprintf("/workout-service/workout/FIT/%d", workoutID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download workout FIT: status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Schedule schedules a workout for a specific date.
func (s *WorkoutService) Schedule(ctx context.Context, workoutID int64, date time.Time) (*ScheduledWorkout, error) {
	path := fmt.Sprintf("/workout-service/schedule/%d", workoutID)

	payload, err := json.Marshal(map[string]string{
		"date": date.Format("2006-01-02"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schedule request: %w", err)
	}

	resp, err := s.client.doAPIWithBody(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to schedule workout: status %d: %s", resp.StatusCode, string(body))
	}

	var scheduled ScheduledWorkout
	if err := json.Unmarshal(body, &scheduled); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scheduled workout: %w", err)
	}
	scheduled.raw = body

	return &scheduled, nil
}

// GetScheduled returns a scheduled workout by ID.
func (s *WorkoutService) GetScheduled(ctx context.Context, scheduleID int64) (*ScheduledWorkout, error) {
	path := fmt.Sprintf("/workout-service/schedule/%d", scheduleID)

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get scheduled workout: status %d: %s", resp.StatusCode, string(body))
	}

	var scheduled ScheduledWorkout
	if err := json.Unmarshal(body, &scheduled); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scheduled workout: %w", err)
	}
	scheduled.raw = body

	return &scheduled, nil
}
