package model

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutKind string

const (
	WorkoutPlan    WorkoutKind = "plan"
	WorkoutOwn     WorkoutKind = "own"
	WorkoutBuilder WorkoutKind = "builder"
)

type Workout struct {
	ID                  uuid.UUID
	ClubID              *uuid.UUID
	UserID              uuid.UUID
	Kind                WorkoutKind
	WorkoutType         WorkoutType
	DayLabel            string
	Tag                 string
	Title               string
	Description         string
	DistKm              float64
	HR                  string
	WeekIndex           int
	ScheduledDate       *time.Time
	Status              WorkoutStatus
	CompletedActivityID *uuid.UUID
	AssignedBy          *uuid.UUID
	IsClubTemplate      bool
	AnnounceID          *uuid.UUID
	CreatedAt           time.Time
	Segments            []Segment
}

func NewWorkout(
	userID uuid.UUID,
	kind WorkoutKind,
	dayLabel, tag, title string,
	distKm float64,
	hr string,
	weekIndex int,
) *Workout {
	return &Workout{
		ID:          uuid.New(),
		UserID:      userID,
		Kind:        kind,
		WorkoutType: WorkoutTypeEasy,
		DayLabel:    dayLabel,
		Tag:         tag,
		Title:       title,
		DistKm:      distKm,
		HR:          hr,
		WeekIndex:   weekIndex,
		Status:      WorkoutStatusPlanned,
		CreatedAt:   time.Now().UTC(),
	}
}

func (w *Workout) AddSegment(kind, title string, distKm float64, pace string, sortOrder int) {
	w.Segments = append(w.Segments, NewSegment(kind, title, distKm, pace, sortOrder))
	w.DistKm += distKm
}

type Segment struct {
	ID        uuid.UUID
	WorkoutID uuid.UUID
	Kind      string
	Title     string
	DistKm    float64
	Pace      string
	SortOrder int
}

func NewSegment(kind, title string, distKm float64, pace string, sortOrder int) Segment {
	return Segment{
		ID:        uuid.New(),
		Kind:      kind,
		Title:     title,
		DistKm:    distKm,
		Pace:      pace,
		SortOrder: sortOrder,
	}
}
