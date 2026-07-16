package user_integration_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) GetByUserProvider(ctx context.Context, userID uuid.UUID, provider model.IntegrationProvider) (*model.UserIntegration, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, status, external_athlete_id, access_token, refresh_token, expires_at, scopes,
			last_synced_at, last_webhook_at, last_error, created_at, updated_at
		FROM user_integrations
		WHERE user_id=$1 AND provider=$2`,
		userID, provider)
	integration, err := scanIntegration(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanIntegration: %w", err)
	}
	return integration, nil
}

func (r *Repo) GetByProviderExternalAthleteID(ctx context.Context, provider model.IntegrationProvider, externalAthleteID string) (*model.UserIntegration, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, status, external_athlete_id, access_token, refresh_token, expires_at, scopes,
			last_synced_at, last_webhook_at, last_error, created_at, updated_at
		FROM user_integrations
		WHERE provider=$1 AND external_athlete_id=$2`,
		provider, externalAthleteID)
	integration, err := scanIntegration(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanIntegration: %w", err)
	}
	return integration, nil
}

func (r *Repo) Upsert(ctx context.Context, integration *model.UserIntegration) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_integrations (
			id, user_id, provider, status, external_athlete_id, access_token, refresh_token, expires_at, scopes,
			last_synced_at, last_webhook_at, last_error, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT (user_id, provider) DO UPDATE SET
			status=EXCLUDED.status,
			external_athlete_id=EXCLUDED.external_athlete_id,
			access_token=EXCLUDED.access_token,
			refresh_token=EXCLUDED.refresh_token,
			expires_at=EXCLUDED.expires_at,
			scopes=EXCLUDED.scopes,
			last_synced_at=EXCLUDED.last_synced_at,
			last_webhook_at=EXCLUDED.last_webhook_at,
			last_error=EXCLUDED.last_error,
			updated_at=EXCLUDED.updated_at`,
		integration.ID,
		integration.UserID,
		integration.Provider,
		integration.Status,
		integration.ExternalAthleteID,
		integration.AccessToken,
		integration.RefreshToken,
		integration.ExpiresAt,
		integration.Scopes,
		integration.LastSyncedAt,
		integration.LastWebhookAt,
		integration.LastError,
		integration.CreatedAt,
		integration.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	return nil
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
