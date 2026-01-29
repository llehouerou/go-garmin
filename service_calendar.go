package garmin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CalendarItem represents an item in the calendar (activity, workout, weight, etc.).
type CalendarItem struct {
	ID                       int64    `json:"id"`
	GroupID                  *int64   `json:"groupId,omitempty"`
	TrainingPlanID           *int64   `json:"trainingPlanId,omitempty"`
	ItemType                 string   `json:"itemType"` // "activity", "workout", "weight", etc.
	ActivityTypeID           *int     `json:"activityTypeId,omitempty"`
	WellnessActivityUUID     *string  `json:"wellnessActivityUuid,omitempty"`
	Title                    *string  `json:"title,omitempty"`
	Date                     string   `json:"date"` // YYYY-MM-DD
	Duration                 *int64   `json:"duration,omitempty"`
	Distance                 *int64   `json:"distance,omitempty"`
	Calories                 *int     `json:"calories,omitempty"`
	FloorsClimbed            *int     `json:"floorsClimbed,omitempty"`
	AvgRespirationRate       *float64 `json:"avgRespirationRate,omitempty"`
	UnitOfPoolLength         *string  `json:"unitOfPoolLength,omitempty"`
	Weight                   *float64 `json:"weight,omitempty"`
	Difference               *float64 `json:"difference,omitempty"`
	CourseID                 *int64   `json:"courseId,omitempty"`
	CourseName               *string  `json:"courseName,omitempty"`
	SportTypeKey             *string  `json:"sportTypeKey,omitempty"`
	URL                      *string  `json:"url,omitempty"`
	IsStart                  *bool    `json:"isStart,omitempty"`
	IsRace                   *bool    `json:"isRace,omitempty"`
	RecurrenceID             *int64   `json:"recurrenceId,omitempty"`
	IsParent                 *bool    `json:"isParent,omitempty"`
	ParentID                 *int64   `json:"parentId,omitempty"`
	UserBadgeID              *int64   `json:"userBadgeId,omitempty"`
	BadgeCategoryTypeID      *int     `json:"badgeCategoryTypeId,omitempty"`
	BadgeCategoryTypeDesc    *string  `json:"badgeCategoryTypeDesc,omitempty"`
	BadgeAwardedDate         *string  `json:"badgeAwardedDate,omitempty"`
	BadgeViewed              *bool    `json:"badgeViewed,omitempty"`
	HideBadge                *bool    `json:"hideBadge,omitempty"`
	StartTimestampLocal      *string  `json:"startTimestampLocal,omitempty"`
	EventTimeLocal           *string  `json:"eventTimeLocal,omitempty"`
	DiveNumber               *int     `json:"diveNumber,omitempty"`
	MaxDepth                 *float64 `json:"maxDepth,omitempty"`
	AvgDepth                 *float64 `json:"avgDepth,omitempty"`
	SurfaceInterval          *int64   `json:"surfaceInterval,omitempty"`
	ElapsedDuration          *float64 `json:"elapsedDuration,omitempty"`
	LapCount                 *int     `json:"lapCount,omitempty"`
	BottomTime               *int64   `json:"bottomTime,omitempty"`
	AtpPlanID                *int64   `json:"atpPlanId,omitempty"`
	WorkoutID                *int64   `json:"workoutId,omitempty"`
	ProtectedWorkoutSchedule bool     `json:"protectedWorkoutSchedule"`
	ActiveSets               *int     `json:"activeSets,omitempty"`
	Strokes                  *int     `json:"strokes,omitempty"`
	NoOfSplits               *int     `json:"noOfSplits,omitempty"`
	MaxGradeValue            *float64 `json:"maxGradeValue,omitempty"`
	TotalAscent              *float64 `json:"totalAscent,omitempty"`
	DifferenceStress         *int     `json:"differenceStress,omitempty"`
	ClimbDuration            *float64 `json:"climbDuration,omitempty"`
	MaxSpeed                 *float64 `json:"maxSpeed,omitempty"`
	AverageHR                *float64 `json:"averageHR,omitempty"`
	ActiveSplitSummaryDur    *float64 `json:"activeSplitSummaryDuration,omitempty"`
	ActiveSplitSummaryDist   *float64 `json:"activeSplitSummaryDistance,omitempty"`
	MaxSplitDistance         *int     `json:"maxSplitDistance,omitempty"`
	MaxSplitSpeed            *float64 `json:"maxSplitSpeed,omitempty"`
	Location                 *string  `json:"location,omitempty"`
	ShareableEventUUID       *string  `json:"shareableEventUuid,omitempty"`
	SplitSummaryMode         *string  `json:"splitSummaryMode,omitempty"`
	CompletionTarget         *any     `json:"completionTarget,omitempty"`
	WorkoutUUID              *string  `json:"workoutUuid,omitempty"`
	NapStartTimeLocal        *string  `json:"napStartTimeLocal,omitempty"`
	BeginPackWeight          *float64 `json:"beginPackWeight,omitempty"`
	HasSplits                *bool    `json:"hasSplits,omitempty"`
	PrimaryEvent             *bool    `json:"primaryEvent,omitempty"`
	ShareableEvent           bool     `json:"shareableEvent"`
	AutoCalcCalories         *bool    `json:"autoCalcCalories,omitempty"`
	Subscribed               *bool    `json:"subscribed,omitempty"`
	PhasedTrainingPlan       *bool    `json:"phasedTrainingPlan,omitempty"`
	DecoDive                 *bool    `json:"decoDive,omitempty"`
}

// Calendar represents the calendar response.
type Calendar struct {
	StartDate        string         `json:"startDate"`
	EndDate          string         `json:"endDate"`
	NumOfDaysInMonth int            `json:"numOfDaysInMonth"`
	CalendarItems    []CalendarItem `json:"calendarItems"`

	raw json.RawMessage
}

// RawJSON returns the raw JSON response.
func (c *Calendar) RawJSON() json.RawMessage {
	return c.raw
}

// SetRaw sets the raw JSON response.
func (c *Calendar) SetRaw(data json.RawMessage) {
	c.raw = data
}

// CalendarOptions specifies optional parameters for calendar retrieval.
// Parameters are hierarchical: month requires year, day requires both month and start.
type CalendarOptions struct {
	Month *int // 0-11 (January=0, December=11)
	Day   *int // Day of month (requires Start)
	Start *int // Week start day (1=Monday)
}

// Get retrieves calendar data for the given year with optional filtering.
// Options are hierarchical: month requires year, day requires both month and start.
func (s *CalendarService) Get(ctx context.Context, year int, opts *CalendarOptions) (*Calendar, error) {
	path := fmt.Sprintf("/calendar-service/year/%d", year)

	if opts != nil && opts.Month != nil {
		path += fmt.Sprintf("/month/%d", *opts.Month)

		if opts.Day != nil {
			if opts.Start == nil {
				return nil, errors.New("start is required when day is provided")
			}
			path += fmt.Sprintf("/day/%d/start/%d", *opts.Day, *opts.Start)
		}
	}

	resp, err := s.client.doAPI(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get calendar: status %d: %s", resp.StatusCode, string(raw))
	}

	var calendar Calendar
	if err := json.Unmarshal(raw, &calendar); err != nil {
		return nil, fmt.Errorf("failed to unmarshal calendar: %w", err)
	}
	calendar.raw = raw

	return &calendar, nil
}
