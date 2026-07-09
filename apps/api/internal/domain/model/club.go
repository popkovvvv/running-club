package model

import (
	"time"

	"github.com/google/uuid"
)

type Club struct {
	ID         uuid.UUID
	Name       string
	InviteCode string
	AccentHex  string
	CoachID    uuid.UUID
	CreatedAt  time.Time
}

func NewClub(name, inviteCode, accentHex string, coachID uuid.UUID) *Club {
	return &Club{
		ID:         uuid.New(),
		Name:       name,
		InviteCode: inviteCode,
		AccentHex:  accentHex,
		CoachID:    coachID,
		CreatedAt:  time.Now().UTC(),
	}
}
