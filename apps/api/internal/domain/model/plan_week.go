package model

import "github.com/google/uuid"

type PlanWeek struct {
	ID         uuid.UUID
	ClubID     uuid.UUID
	WeekIndex  int
	RangeLabel string
	PlanLabel  string
}

func NewPlanWeek(clubID uuid.UUID, weekIndex int, rangeLabel, planLabel string) *PlanWeek {
	return &PlanWeek{
		ID:         uuid.New(),
		ClubID:     clubID,
		WeekIndex:  weekIndex,
		RangeLabel: rangeLabel,
		PlanLabel:  planLabel,
	}
}
