package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateWorkoutRequest struct {
	Kind          string         `json:"kind"`
	TargetUserID  *uuid.UUID     `json:"targetUserId"`
	WorkoutType   string         `json:"workoutType"`
	DayLabel      string         `json:"dayLabel"`
	Tag           string         `json:"tag"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	DistKm        float64        `json:"distKm"`
	HR            string         `json:"hr"`
	WeekIndex     int            `json:"weekIndex"`
	ScheduledDate *string        `json:"scheduledDate"`
	Segments      []SegmentInput `json:"segments"`
}

type UpdateWorkoutRequest struct {
	Status              *string    `json:"status"`
	CompletedActivityID *uuid.UUID `json:"completedActivityId"`
	RPE                 *int       `json:"rpe"`
	AthleteReport       *string    `json:"athleteReport"`
	CoachComment        *string    `json:"coachComment"`
}

type SegmentInput struct {
	Kind   string  `json:"kind"`
	Title  string  `json:"title"`
	DistKm float64 `json:"distKm"`
	Pace   string  `json:"pace"`
}

type WorkoutView struct {
	ID                  uuid.UUID           `json:"id"`
	Kind                string              `json:"kind"`
	WorkoutType         string              `json:"workoutType"`
	DayLabel            string              `json:"dayLabel"`
	Tag                 string              `json:"tag"`
	Title               string              `json:"title"`
	Description         string              `json:"description,omitempty"`
	DistKm              float64             `json:"distKm"`
	HR                  string              `json:"hr"`
	WeekIndex           int                 `json:"weekIndex"`
	ScheduledDate       *string             `json:"scheduledDate,omitempty"`
	Status              string              `json:"status"`
	CompletedActivityID *uuid.UUID          `json:"completedActivityId,omitempty"`
	AssignedBy          *uuid.UUID          `json:"assignedBy,omitempty"`
	IsClubTemplate      bool                `json:"isClubTemplate"`
	Segments            []SegmentView       `json:"segments,omitempty"`
	PlannedKm           float64             `json:"plannedKm"`
	ActualKm            *float64            `json:"actualKm,omitempty"`
	ActualPace          string              `json:"actualPace,omitempty"`
	RPE                 *int                `json:"rpe,omitempty"`
	AthleteReport       string              `json:"athleteReport,omitempty"`
	CoachComment        string              `json:"coachComment,omitempty"`
	Fact                *ActivityDetailView `json:"fact,omitempty"`
}

func NewWorkoutView(
	id uuid.UUID,
	kind, workoutType, dayLabel, tag, title, description string,
	distKm float64,
	hr string,
	weekIndex int,
	scheduledDate *string,
	status string,
	completedActivityID, assignedBy *uuid.UUID,
	isClubTemplate bool,
	segments []SegmentView,
) WorkoutView {
	return WorkoutView{
		ID:                  id,
		Kind:                kind,
		WorkoutType:         workoutType,
		DayLabel:            dayLabel,
		Tag:                 tag,
		Title:               title,
		Description:         description,
		DistKm:              distKm,
		HR:                  hr,
		WeekIndex:           weekIndex,
		ScheduledDate:       scheduledDate,
		Status:              status,
		CompletedActivityID: completedActivityID,
		AssignedBy:          assignedBy,
		IsClubTemplate:      isClubTemplate,
		Segments:            segments,
		PlannedKm:           distKm,
	}
}

type SegmentView struct {
	ID     uuid.UUID `json:"id"`
	Kind   string    `json:"kind"`
	Title  string    `json:"title"`
	DistKm float64   `json:"distKm"`
	Pace   string    `json:"pace"`
}

func NewSegmentView(id uuid.UUID, kind, title string, distKm float64, pace string) SegmentView {
	return SegmentView{
		ID:     id,
		Kind:   kind,
		Title:  title,
		DistKm: distKm,
		Pace:   pace,
	}
}

type PlanResponse struct {
	WeekIndex int           `json:"weekIndex"`
	WeekRange string        `json:"weekRange"`
	WeekPlan  string        `json:"weekPlan"`
	WeekKm    string        `json:"weekKm"`
	Days      []WorkoutView `json:"days"`
	Mine      []WorkoutView `json:"mine"`
}

func NewPlanResponse(weekIndex int, weekRange, weekPlan, weekKm string, days, mine []WorkoutView) *PlanResponse {
	return &PlanResponse{
		WeekIndex: weekIndex,
		WeekRange: weekRange,
		WeekPlan:  weekPlan,
		WeekKm:    weekKm,
		Days:      days,
		Mine:      mine,
	}
}

type PlanWeekView struct {
	WeekIndex  int    `json:"weekIndex"`
	RangeLabel string `json:"rangeLabel"`
	PlanLabel  string `json:"planLabel"`
}

type UpsertPlanWeekRequest struct {
	RangeLabel string `json:"rangeLabel"`
	PlanLabel  string `json:"planLabel"`
}

type SaveTemplateRequest struct {
	Workouts []CreateWorkoutRequest `json:"workouts"`
}

type TemplateResponse struct {
	WeekIndex int           `json:"weekIndex"`
	Workouts  []WorkoutView `json:"workouts"`
}

func FormatDate(d *time.Time) *string {
	if d == nil {
		return nil
	}
	s := d.Format("2006-01-02")
	return &s
}
