package model

import "github.com/google/uuid"

func ClubFixture(id uuid.UUID, name, inviteCode, accentHex string, coachID uuid.UUID) *Club {
	c := NewClub(name, inviteCode, accentHex, coachID)
	c.ID = id
	return c
}

func MembershipFixture(id, userID, clubID uuid.UUID) *Membership {
	m := NewMembership(userID, clubID)
	m.ID = id
	return m
}

func AnnounceFixture(id, clubID uuid.UUID, place, dayLabel, timeLabel, groupName string) *Announce {
	a := NewAnnounce(clubID, place, dayLabel, timeLabel, groupName, "")
	a.ID = id
	return a
}

func UserFixture(id uuid.UUID, name, email, passwordHash string, role Role) *User {
	u := NewUser(name, email, passwordHash, role)
	u.ID = id
	return u
}
