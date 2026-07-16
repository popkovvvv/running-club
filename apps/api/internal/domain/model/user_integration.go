package model

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationProvider string

const (
	IntegrationProviderStrava IntegrationProvider = "strava"
)

type IntegrationStatus string

const (
	IntegrationStatusActive       IntegrationStatus = "active"
	IntegrationStatusDisconnected IntegrationStatus = "disconnected"
	IntegrationStatusRevoked      IntegrationStatus = "revoked"
	IntegrationStatusError        IntegrationStatus = "error"
)

type UserIntegration struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Provider          IntegrationProvider
	Status            IntegrationStatus
	ExternalAthleteID string
	AccessToken       string
	RefreshToken      string
	ExpiresAt         time.Time
	Scopes            []string
	LastSyncedAt      *time.Time
	LastWebhookAt     *time.Time
	LastError         string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewUserIntegration(userID uuid.UUID, provider IntegrationProvider) *UserIntegration {
	now := time.Now().UTC()
	return &UserIntegration{
		ID:        uuid.New(),
		UserID:    userID,
		Provider:  provider,
		Status:    IntegrationStatusDisconnected,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
