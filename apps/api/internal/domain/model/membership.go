package model

import (
	"time"

	"github.com/google/uuid"
)

type MembershipStatus string

const (
	MembershipActive  MembershipStatus = "active"
	MembershipLeft    MembershipStatus = "left"
	MembershipRemoved MembershipStatus = "removed"
)

type Membership struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ClubID    uuid.UUID
	Status    MembershipStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMembership(userID, clubID uuid.UUID) *Membership {
	now := time.Now().UTC()
	return &Membership{
		ID:        uuid.New(),
		UserID:    userID,
		ClubID:    clubID,
		Status:    MembershipActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
