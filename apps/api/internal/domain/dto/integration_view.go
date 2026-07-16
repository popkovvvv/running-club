package dto

import "time"

type IntegrationView struct {
	Provider          string     `json:"provider"`
	Status            string     `json:"status"`
	Connected         bool       `json:"connected"`
	ExternalAthleteID string     `json:"externalAthleteId,omitempty"`
	Scopes            []string   `json:"scopes,omitempty"`
	ExpiresAt         time.Time  `json:"expiresAt,omitempty"`
	LastSyncedAt      *time.Time `json:"lastSyncedAt,omitempty"`
	LastWebhookAt     *time.Time `json:"lastWebhookAt,omitempty"`
	LastError         string     `json:"lastError,omitempty"`
}
