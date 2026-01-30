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

// Sport type constants.
const (
	SportTypeRunning          = 1
	SportTypeCycling          = 2
	SportTypeOther            = 3
	SportTypeSwimming         = 4
	SportTypeStrengthTraining = 5
	SportTypeCardioTraining   = 6
	SportTypeYoga             = 7
	SportTypePilates          = 8
	SportTypeHIIT             = 9
	SportTypeMultiSport       = 10
	SportTypeMobility         = 11

	SportTypeKeyRunning          = "running"
	SportTypeKeyCycling          = "cycling"
	SportTypeKeyOther            = "other"
	SportTypeKeySwimming         = "swimming"
	SportTypeKeyStrengthTraining = "strength_training"
	SportTypeKeyCardioTraining   = "cardio_training"
	SportTypeKeyYoga             = "yoga"
	SportTypeKeyPilates          = "pilates"
	SportTypeKeyHIIT             = "hiit"
	SportTypeKeyMultiSport       = "multi_sport"
	SportTypeKeyMobility         = "mobility"
)

// Step type constants.
const (
	StepTypeWarmup   = 1
	StepTypeCooldown = 2
	StepTypeInterval = 3
	StepTypeRecovery = 4
	StepTypeRest     = 5
	StepTypeRepeat   = 6
	StepTypeOther    = 7
	StepTypeMain     = 8

	StepTypeKeyWarmup   = "warmup"
	StepTypeKeyCooldown = "cooldown"
	StepTypeKeyInterval = "interval"
	StepTypeKeyRecovery = "recovery"
	StepTypeKeyRest     = "rest"
	StepTypeKeyRepeat   = "repeat"
	StepTypeKeyOther    = "other"
	StepTypeKeyMain     = "main"
)

// End condition type constants.
const (
	ConditionTypeLapButton          = 1
	ConditionTypeTime               = 2
	ConditionTypeDistance           = 3
	ConditionTypeCalories           = 4
	ConditionTypePower              = 5
	ConditionTypeHeartRate          = 6
	ConditionTypeIterations         = 7 // For repeat groups
	ConditionTypeFixedRest          = 8
	ConditionTypeFixedRepetition    = 9
	ConditionTypeReps               = 10 // For strength training
	ConditionTypeTrainingPeaksTSS   = 11
	ConditionTypeRepetitionTime     = 12
	ConditionTypeTimeAtValidCDA     = 13
	ConditionTypePowerLastLap       = 14
	ConditionTypeMaxPowerLastLap    = 15
	ConditionTypeRepetitionSwimCSS  = 16
	ConditionTypeVelocityLoss       = 17
	ConditionTypeCustomVelocity     = 18
	ConditionTypeVBTVelocityZone    = 19
	ConditionTypeVelocityMin        = 20
	ConditionTypePeakVelocityMin    = 21
	ConditionTypePeakVelocityLoss   = 22
	ConditionTypeCustomPeakVelocity = 23
	ConditionTypePowerLoss          = 24

	ConditionTypeKeyLapButton          = "lap.button"
	ConditionTypeKeyTime               = "time"
	ConditionTypeKeyDistance           = "distance"
	ConditionTypeKeyCalories           = "calories"
	ConditionTypeKeyPower              = "power"
	ConditionTypeKeyHeartRate          = "heart.rate"
	ConditionTypeKeyIterations         = "iterations"
	ConditionTypeKeyFixedRest          = "fixed.rest"
	ConditionTypeKeyFixedRepetition    = "fixed.repetition"
	ConditionTypeKeyReps               = "reps"
	ConditionTypeKeyTrainingPeaksTSS   = "training.peaks.tss"
	ConditionTypeKeyRepetitionTime     = "repetition.time"
	ConditionTypeKeyTimeAtValidCDA     = "time.at.valid.cda"
	ConditionTypeKeyPowerLastLap       = "power.last.lap"
	ConditionTypeKeyMaxPowerLastLap    = "max.power.last.lap"
	ConditionTypeKeyRepetitionSwimCSS  = "repetition.swim.css.offset"
	ConditionTypeKeyVelocityLoss       = "velocity.loss"
	ConditionTypeKeyCustomVelocity     = "custom.velocity"
	ConditionTypeKeyVBTVelocityZone    = "vbt.velocity.zone"
	ConditionTypeKeyVelocityMin        = "velocity.min"
	ConditionTypeKeyPeakVelocityMin    = "peak.velocity.min"
	ConditionTypeKeyPeakVelocityLoss   = "peak.velocity.loss"
	ConditionTypeKeyCustomPeakVelocity = "custom.peak.velocity"
	ConditionTypeKeyPowerLoss          = "power.loss"
)

// Target type constants.
const (
	TargetTypeNoTarget           = 1
	TargetTypePowerZone          = 2
	TargetTypeCadence            = 3
	TargetTypeHeartRateZone      = 4
	TargetTypeSpeedZone          = 5
	TargetTypePaceZone           = 6
	TargetTypeGrade              = 7
	TargetTypeHeartRateLap       = 8
	TargetTypePowerLap           = 9
	TargetTypePower3s            = 10
	TargetTypePower10s           = 11
	TargetTypePower30s           = 12
	TargetTypeSpeedLap           = 13
	TargetTypeSwimStroke         = 14
	TargetTypeResistance         = 15
	TargetTypePowerCurve         = 16
	TargetTypeSwimCSSOffset      = 17
	TargetTypeSwimInstruction    = 18
	TargetTypeInstruction        = 19
	TargetTypeVelocityLoss       = 20
	TargetTypeCustomVelocity     = 21
	TargetTypeVBTVelocityZone    = 22
	TargetTypePowerLoss          = 23
	TargetTypeVelocityMin        = 24
	TargetTypePeakVelocityMin    = 25
	TargetTypePeakVelocityLoss   = 26
	TargetTypeCustomPeakVelocity = 27

	TargetTypeKeyNoTarget           = "no.target"
	TargetTypeKeyPowerZone          = "power.zone"
	TargetTypeKeyCadence            = "cadence"
	TargetTypeKeyHeartRateZone      = "heart.rate.zone"
	TargetTypeKeySpeedZone          = "speed.zone"
	TargetTypeKeyPaceZone           = "pace.zone"
	TargetTypeKeyGrade              = "grade"
	TargetTypeKeyHeartRateLap       = "heart.rate.lap"
	TargetTypeKeyPowerLap           = "power.lap"
	TargetTypeKeyPower3s            = "power.3s"
	TargetTypeKeyPower10s           = "power.10s"
	TargetTypeKeyPower30s           = "power.30s"
	TargetTypeKeySpeedLap           = "speed.lap"
	TargetTypeKeySwimStroke         = "swim.stroke"
	TargetTypeKeyResistance         = "resistance"
	TargetTypeKeyPowerCurve         = "power.curve"
	TargetTypeKeySwimCSSOffset      = "swim.css.offset"
	TargetTypeKeySwimInstruction    = "swim.instruction"
	TargetTypeKeyInstruction        = "instruction"
	TargetTypeKeyVelocityLoss       = "velocity.loss"
	TargetTypeKeyCustomVelocity     = "custom.velocity"
	TargetTypeKeyVBTVelocityZone    = "vbt.velocity.zone"
	TargetTypeKeyPowerLoss          = "power.loss"
	TargetTypeKeyVelocityMin        = "velocity.min"
	TargetTypeKeyPeakVelocityMin    = "peak.velocity.min"
	TargetTypeKeyPeakVelocityLoss   = "peak.velocity.loss"
	TargetTypeKeyCustomPeakVelocity = "custom.peak.velocity"
)

// Intensity type constants.
const (
	IntensityTypeActive   = 1
	IntensityTypeRest     = 2
	IntensityTypeWarmup   = 3
	IntensityTypeCooldown = 4

	IntensityTypeKeyActive   = "active"
	IntensityTypeKeyRest     = "rest"
	IntensityTypeKeyWarmup   = "warmup"
	IntensityTypeKeyCooldown = "cooldown"
)

// Swimming stroke type constants.
const (
	StrokeTypeAny              = 1
	StrokeTypeBackstroke       = 2
	StrokeTypeBreaststroke     = 3
	StrokeTypeDrill            = 4
	StrokeTypeFly              = 5
	StrokeTypeFree             = 6
	StrokeTypeIM               = 7
	StrokeTypeMixed            = 8
	StrokeTypeIMByRound        = 9
	StrokeTypeReverseIMByRound = 10

	StrokeTypeKeyAny              = "any_stroke"
	StrokeTypeKeyBackstroke       = "backstroke"
	StrokeTypeKeyBreaststroke     = "breaststroke"
	StrokeTypeKeyDrill            = "drill"
	StrokeTypeKeyFly              = "fly"
	StrokeTypeKeyFree             = "free"
	StrokeTypeKeyIM               = "individual_medley"
	StrokeTypeKeyMixed            = "mixed"
	StrokeTypeKeyIMByRound        = "individual_medley_by_round"
	StrokeTypeKeyReverseIMByRound = "reverse_individual_medley_by_round"
)

// Swimming equipment type constants.
const (
	SwimEquipmentFins      = 1
	SwimEquipmentKickboard = 2
	SwimEquipmentPaddles   = 3
	SwimEquipmentPullBuoy  = 4
	SwimEquipmentSnorkel   = 5

	SwimEquipmentKeyFins      = "fins"
	SwimEquipmentKeyKickboard = "kickboard"
	SwimEquipmentKeyPaddles   = "paddles"
	SwimEquipmentKeyPullBuoy  = "pull_buoy"
	SwimEquipmentKeySnorkel   = "snorkel"
)

// Swimming instruction type constants.
const (
	SwimInstructionRecovery = 1
	SwimInstructionVeryEasy = 2
	SwimInstructionEasy     = 3
	SwimInstructionModerate = 4
	SwimInstructionHard     = 5
	SwimInstructionVeryHard = 6
	SwimInstructionAllOut   = 7
	SwimInstructionFast     = 8
	SwimInstructionAscend   = 9
	SwimInstructionDescend  = 10

	SwimInstructionKeyRecovery = "recovery"
	SwimInstructionKeyVeryEasy = "very_easy"
	SwimInstructionKeyEasy     = "easy"
	SwimInstructionKeyModerate = "moderate"
	SwimInstructionKeyHard     = "hard"
	SwimInstructionKeyVeryHard = "very_hard"
	SwimInstructionKeyAllOut   = "all_out"
	SwimInstructionKeyFast     = "fast"
	SwimInstructionKeyAscend   = "ascend"
	SwimInstructionKeyDescend  = "descend"
)

// Swimming drill type constants.
const (
	SwimDrillKick  = 1
	SwimDrillPull  = 2
	SwimDrillDrill = 3

	SwimDrillKeyKick  = "kick"
	SwimDrillKeyPull  = "pull"
	SwimDrillKeyDrill = "drill"
)

// Workout step DTO type constants.
const (
	StepDTOExecutable = "ExecutableStepDTO"
	StepDTORepeat     = "RepeatGroupDTO"
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

// StrokeType represents a swimming stroke type.
type StrokeType struct {
	StrokeTypeID  int    `json:"strokeTypeId"`
	StrokeTypeKey string `json:"strokeTypeKey,omitempty"`
	DisplayOrder  int    `json:"displayOrder,omitempty"`
}

// EquipmentType represents an equipment type for workout steps.
type EquipmentType struct {
	EquipmentTypeID  int    `json:"equipmentTypeId"`
	EquipmentTypeKey string `json:"equipmentTypeKey,omitempty"`
	DisplayOrder     int    `json:"displayOrder,omitempty"`
}

// UnitInfo represents unit information for distance/length measurements.
type UnitInfo struct {
	UnitID  *int64   `json:"unitId,omitempty"`
	UnitKey *string  `json:"unitKey,omitempty"`
	Factor  *float64 `json:"factor,omitempty"`
}

// WorkoutStep represents a single step in a workout.
type WorkoutStep struct {
	Type        string        `json:"type"` // ExecutableStepDTO or RepeatGroupDTO
	StepID      int64         `json:"stepId,omitempty"`
	StepOrder   int           `json:"stepOrder"`
	StepType    *StepTypeInfo `json:"stepType,omitempty"`
	ChildStepID *int64        `json:"childStepId,omitempty"`
	Description *string       `json:"description,omitempty"`

	// End condition
	EndCondition         *EndCondition `json:"endCondition,omitempty"`
	EndConditionValue    *float64      `json:"endConditionValue,omitempty"`
	PreferredEndCondUnit *UnitInfo     `json:"preferredEndConditionUnit,omitempty"`
	EndConditionCompare  *float64      `json:"endConditionCompare,omitempty"`
	EndConditionZone     *int          `json:"endConditionZone,omitempty"`

	// Primary target
	TargetType      *TargetType `json:"targetType,omitempty"`
	TargetValueOne  *float64    `json:"targetValueOne,omitempty"`
	TargetValueTwo  *float64    `json:"targetValueTwo,omitempty"`
	TargetValueUnit *UnitInfo   `json:"targetValueUnit,omitempty"`
	ZoneNumber      *int        `json:"zoneNumber,omitempty"`

	// Secondary target
	SecondaryTargetType      *TargetType `json:"secondaryTargetType,omitempty"`
	SecondaryTargetValueOne  *float64    `json:"secondaryTargetValueOne,omitempty"`
	SecondaryTargetValueTwo  *float64    `json:"secondaryTargetValueTwo,omitempty"`
	SecondaryTargetValueUnit *UnitInfo   `json:"secondaryTargetValueUnit,omitempty"`
	SecondaryZoneNumber      *int        `json:"secondaryZoneNumber,omitempty"`

	// Sport-specific
	StrokeType    *StrokeType    `json:"strokeType,omitempty"`
	EquipmentType *EquipmentType `json:"equipmentType,omitempty"`

	// Exercise info
	Category                 *string   `json:"category,omitempty"`
	ExerciseName             *string   `json:"exerciseName,omitempty"`
	WorkoutProvider          *string   `json:"workoutProvider,omitempty"`
	ProviderExerciseSourceID *int64    `json:"providerExerciseSourceId,omitempty"`
	WeightValue              *float64  `json:"weightValue,omitempty"`
	WeightUnit               *UnitInfo `json:"weightUnit,omitempty"`

	// For repeat groups (RepeatGroupDTO)
	NumberOfIterations *int          `json:"numberOfIterations,omitempty"`
	WorkoutSteps       []WorkoutStep `json:"workoutSteps,omitempty"`
	SmartRepeat        bool          `json:"smartRepeat,omitempty"`
	SkipLastRestStep   bool          `json:"skipLastRestStep,omitempty"`
}

// WorkoutSegment represents a segment of a workout.
type WorkoutSegment struct {
	SegmentOrder              int           `json:"segmentOrder"`
	SportType                 SportType     `json:"sportType"`
	WorkoutSteps              []WorkoutStep `json:"workoutSteps"`
	PoolLengthUnit            *UnitInfo     `json:"poolLengthUnit,omitempty"`
	PoolLength                *float64      `json:"poolLength,omitempty"`
	AvgTrainingSpeed          *float64      `json:"avgTrainingSpeed,omitempty"`
	EstimatedDurationInSecs   *int          `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMeters *float64      `json:"estimatedDistanceInMeters,omitempty"`
	EstimatedDistanceUnit     *UnitInfo     `json:"estimatedDistanceUnit,omitempty"`
	EstimateType              *string       `json:"estimateType,omitempty"`
	Description               *string       `json:"description,omitempty"`
}

// WorkoutAuthor represents the author of a workout.
type WorkoutAuthor struct {
	UserProfilePK       *int64  `json:"userProfilePk,omitempty"`
	DisplayName         *string `json:"displayName,omitempty"`
	FullName            *string `json:"fullName,omitempty"`
	ProfileImgNameLarge *string `json:"profileImgNameLarge,omitempty"`
	ProfileImgNameMed   *string `json:"profileImgNameMedium,omitempty"`
	ProfileImgNameSmall *string `json:"profileImgNameSmall,omitempty"`
	UserPro             bool    `json:"userPro,omitempty"`
	VivokidUser         bool    `json:"vivokidUser,omitempty"`
}

// Workout represents a Garmin workout.
type Workout struct {
	WorkoutID                int64            `json:"workoutId,omitempty"`
	OwnerID                  int64            `json:"ownerId,omitempty"`
	WorkoutName              string           `json:"workoutName"`
	Description              string           `json:"description,omitempty"`
	UpdatedDate              string           `json:"updatedDate,omitempty"`
	CreatedDate              string           `json:"createdDate,omitempty"`
	SportType                SportType        `json:"sportType"`
	SubSportType             *SportType       `json:"subSportType,omitempty"`
	TrainingPlanID           *int64           `json:"trainingPlanId,omitempty"`
	Author                   *WorkoutAuthor   `json:"author,omitempty"`
	SharedWithUsers          []WorkoutAuthor  `json:"sharedWithUsers,omitempty"`
	EstimatedDurationInSecs  int              `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMtrs  float64          `json:"estimatedDistanceInMeters,omitempty"`
	EstimateType             string           `json:"estimateType,omitempty"`
	EstimatedDistanceUnit    *UnitInfo        `json:"estimatedDistanceUnit,omitempty"`
	AvgTrainingSpeed         float64          `json:"avgTrainingSpeed,omitempty"`
	WorkoutSegments          []WorkoutSegment `json:"workoutSegments"`
	Locale                   string           `json:"locale,omitempty"`
	PoolLength               *float64         `json:"poolLength,omitempty"`
	PoolLengthUnit           *UnitInfo        `json:"poolLengthUnit,omitempty"`
	WorkoutProvider          *string          `json:"workoutProvider,omitempty"`
	WorkoutSourceID          *string          `json:"workoutSourceId,omitempty"`
	UploadTimestamp          *string          `json:"uploadTimestamp,omitempty"`
	AtpPlanID                *int64           `json:"atpPlanId,omitempty"`
	Consumer                 *string          `json:"consumer,omitempty"`
	ConsumerName             *string          `json:"consumerName,omitempty"`
	ConsumerImageURL         *string          `json:"consumerImageURL,omitempty"`
	ConsumerWebsiteURL       *string          `json:"consumerWebsiteURL,omitempty"`
	WorkoutNameI18nKey       *string          `json:"workoutNameI18nKey,omitempty"`
	DescriptionI18nKey       *string          `json:"descriptionI18nKey,omitempty"`
	WorkoutThumbnailURL      *string          `json:"workoutThumbnailUrl,omitempty"`
	SessionTransitionEnabled *bool            `json:"isSessionTransitionEnabled,omitempty"`
	Shared                   bool             `json:"shared,omitempty"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (w *Workout) RawJSON() json.RawMessage {
	return w.raw
}

// SetRaw sets the raw JSON data.
func (w *Workout) SetRaw(data json.RawMessage) {
	w.raw = data
}

// WorkoutSummary represents a workout in the list response.
type WorkoutSummary struct {
	WorkoutID               int64          `json:"workoutId"`
	OwnerID                 int64          `json:"ownerId"`
	WorkoutName             string         `json:"workoutName"`
	Description             string         `json:"description,omitempty"`
	UpdateDate              string         `json:"updateDate,omitempty"`
	CreatedDate             string         `json:"createdDate,omitempty"`
	SportType               SportType      `json:"sportType"`
	TrainingPlanID          *int64         `json:"trainingPlanId,omitempty"`
	Author                  *WorkoutAuthor `json:"author,omitempty"`
	EstimatedDurationInSecs int            `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMtrs *float64       `json:"estimatedDistanceInMeters,omitempty"`
	EstimateType            *string        `json:"estimateType,omitempty"`
	EstimatedDistanceUnit   *UnitInfo      `json:"estimatedDistanceUnit,omitempty"`
	PoolLength              *float64       `json:"poolLength,omitempty"`
	PoolLengthUnit          *UnitInfo      `json:"poolLengthUnit,omitempty"`
	WorkoutProvider         *string        `json:"workoutProvider,omitempty"`
	WorkoutSourceID         *string        `json:"workoutSourceId,omitempty"`
	Consumer                *string        `json:"consumer,omitempty"`
	AtpPlanID               *int64         `json:"atpPlanId,omitempty"`
	WorkoutNameI18nKey      *string        `json:"workoutNameI18nKey,omitempty"`
	DescriptionI18nKey      *string        `json:"descriptionI18nKey,omitempty"`
	WorkoutThumbnailURL     *string        `json:"workoutThumbnailUrl,omitempty"`
	Shared                  bool           `json:"shared,omitempty"`
	Estimated               bool           `json:"estimated,omitempty"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (w *WorkoutSummary) RawJSON() json.RawMessage {
	return w.raw
}

// SetRaw sets the raw JSON data.
func (w *WorkoutSummary) SetRaw(data json.RawMessage) {
	w.raw = data
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

// SetRaw sets the raw JSON data.
func (w *WorkoutList) SetRaw(data json.RawMessage) {
	w.raw = data
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

// SetRaw sets the raw JSON data.
func (s *ScheduledWorkout) SetRaw(data json.RawMessage) {
	s.raw = data
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
	return fetch[Workout](ctx, s.client, path)
}

// Create creates a new workout and returns the created workout.
func (s *WorkoutService) Create(ctx context.Context, workout *Workout) (*Workout, error) {
	path := "/workout-service/workout"
	return send[Workout](ctx, s.client, http.MethodPost, path, workout)
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

	// Garmin API may return empty body on successful update
	if len(body) == 0 {
		// Fetch the updated workout to return it
		return s.Get(ctx, workoutID)
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
	return sendEmpty(ctx, s.client, http.MethodDelete, path)
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
	return fetch[ScheduledWorkout](ctx, s.client, path)
}

// Unschedule removes a scheduled workout by schedule ID.
func (s *WorkoutService) Unschedule(ctx context.Context, scheduleID int64) error {
	path := fmt.Sprintf("/workout-service/schedule/%d", scheduleID)
	return sendEmpty(ctx, s.client, http.MethodDelete, path)
}
