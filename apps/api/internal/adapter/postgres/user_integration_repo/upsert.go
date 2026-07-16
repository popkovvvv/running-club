package user_integration_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Upsert(ctx context.Context, integration *model.UserIntegration) error {
	q, args, err := psql.Insert("user_integrations").
		Columns(integrationColumns...).
		Values(
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
		).
		Suffix(`ON CONFLICT (user_id, provider) DO UPDATE SET
			status=EXCLUDED.status,
			external_athlete_id=EXCLUDED.external_athlete_id,
			access_token=EXCLUDED.access_token,
			refresh_token=EXCLUDED.refresh_token,
			expires_at=EXCLUDED.expires_at,
			scopes=EXCLUDED.scopes,
			last_synced_at=EXCLUDED.last_synced_at,
			last_webhook_at=EXCLUDED.last_webhook_at,
			last_error=EXCLUDED.last_error,
			updated_at=EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	_, err = r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	return nil
}
