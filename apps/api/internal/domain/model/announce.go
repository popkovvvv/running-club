package model

import (
	"time"

	"github.com/google/uuid"
)

type Announce struct {
	ID         uuid.UUID
	ClubID     uuid.UUID
	Place      string
	DayLabel   string
	Time       string
	GroupName  string
	Note       string
	StartsOn   *time.Time
	GoingCount int
	CreatedAt  time.Time
}

func NewAnnounce(clubID uuid.UUID, place, dayLabel, timeLabel, groupName, note string, startsOn *time.Time) *Announce {
	return &Announce{
		ID:         uuid.New(),
		ClubID:     clubID,
		Place:      place,
		DayLabel:   dayLabel,
		Time:       timeLabel,
		GroupName:  groupName,
		Note:       note,
		StartsOn:   startsOn,
		GoingCount: 0,
		CreatedAt:  time.Now().UTC(),
	}
}

type AnnounceSignup struct {
	ID         uuid.UUID
	AnnounceID uuid.UUID
	AthleteID  uuid.UUID
	CreatedAt  time.Time
}

func NewAnnounceSignup(announceID, athleteID uuid.UUID) *AnnounceSignup {
	return &AnnounceSignup{
		ID:         uuid.New(),
		AnnounceID: announceID,
		AthleteID:  athleteID,
		CreatedAt:  time.Now().UTC(),
	}
}
