package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAthlete Role = "athlete"
	RoleCoach   Role = "coach"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

func NewUser(name, email, passwordHash string, role Role) *User {
	return &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now().UTC(),
	}
}
