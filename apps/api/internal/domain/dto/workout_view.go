package dto

import "github.com/google/uuid"

type CreateWorkoutRequest struct {
	Kind      string         `json:"kind"`
	DayLabel  string         `json:"dayLabel"`
	Tag       string         `json:"tag"`
	Title     string         `json:"title"`
	DistKm    float64        `json:"distKm"`
	Duration  string         `json:"duration"`
	Pace      string         `json:"pace"`
	HR        string         `json:"hr"`
	WeekIndex int            `json:"weekIndex"`
	Segments  []SegmentInput `json:"segments"`
}

type SegmentInput struct {
	Kind   string  `json:"kind"`
	Title  string  `json:"title"`
	DistKm float64 `json:"distKm"`
	Pace   string  `json:"pace"`
}

type WorkoutView struct {
	ID       uuid.UUID     `json:"id"`
	Kind     string        `json:"kind"`
	DayLabel string        `json:"dayLabel"`
	Tag      string        `json:"tag"`
	Title    string        `json:"title"`
	DistKm   float64       `json:"distKm"`
	Duration string        `json:"duration"`
	Pace     string        `json:"pace"`
	HR       string        `json:"hr"`
	Segments []SegmentView `json:"segments,omitempty"`
}

func NewWorkoutView(
	id uuid.UUID,
	kind, dayLabel, tag, title string,
	distKm float64,
	duration, pace, hr string,
	segments []SegmentView,
) WorkoutView {
	return WorkoutView{
		ID:       id,
		Kind:     kind,
		DayLabel: dayLabel,
		Tag:      tag,
		Title:    title,
		DistKm:   distKm,
		Duration: duration,
		Pace:     pace,
		HR:       hr,
		Segments: segments,
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
	Days      []WorkoutView `json:"days"`
	Mine      []WorkoutView `json:"mine"`
}

func NewPlanResponse(weekIndex int, weekRange, weekPlan string, days, mine []WorkoutView) *PlanResponse {
	return &PlanResponse{
		WeekIndex: weekIndex,
		WeekRange: weekRange,
		WeekPlan:  weekPlan,
		Days:      days,
		Mine:      mine,
	}
}
