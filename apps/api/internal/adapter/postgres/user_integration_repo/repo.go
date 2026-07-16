package user_integration_repo

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var integrationColumns = []string{
	"id", "user_id", "provider", "status", "external_athlete_id", "access_token", "refresh_token",
	"expires_at", "scopes", "last_synced_at", "last_webhook_at", "last_error", "created_at", "updated_at",
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanIntegration(row scannable) (*model.UserIntegration, error) {
	var integration model.UserIntegration
	var expiresAt time.Time
	var lastSyncedAt *time.Time
	var lastWebhookAt *time.Time
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(
		&integration.ID,
		&integration.UserID,
		&integration.Provider,
		&integration.Status,
		&integration.ExternalAthleteID,
		&integration.AccessToken,
		&integration.RefreshToken,
		&expiresAt,
		&integration.Scopes,
		&lastSyncedAt,
		&lastWebhookAt,
		&integration.LastError,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	integration.ExpiresAt = expiresAt
	integration.LastSyncedAt = lastSyncedAt
	integration.LastWebhookAt = lastWebhookAt
	integration.CreatedAt = createdAt
	integration.UpdatedAt = updatedAt
	return &integration, nil
}
