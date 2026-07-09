package model

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var planKmRe = regexp.MustCompile(`\d+(?:[.,]\d+)?`)

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

func (w *PlanWeek) TargetKm() (float64, bool) {
	if w == nil {
		return 0, false
	}
	raw := planKmRe.FindString(w.PlanLabel)
	if raw == "" {
		return 0, false
	}
	raw = strings.Replace(raw, ",", ".", 1)
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return 0, false
	}
	return v, true
}
