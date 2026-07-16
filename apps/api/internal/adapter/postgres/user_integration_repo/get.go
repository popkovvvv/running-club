package user_integration_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByUserProvider(ctx context.Context, userID uuid.UUID, provider model.IntegrationProvider) (*model.UserIntegration, error) {
	q, args, err := psql.Select(integrationColumns...).
		From("user_integrations").
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	integration, err := scanIntegration(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanIntegration: %w", err)
	}
	return integration, nil
}

func (r *Repo) GetByProviderExternalAthleteID(ctx context.Context, provider model.IntegrationProvider, externalAthleteID string) (*model.UserIntegration, error) {
	q, args, err := psql.Select(integrationColumns...).
		From("user_integrations").
		Where(sq.Eq{"provider": provider, "external_athlete_id": externalAthleteID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	integration, err := scanIntegration(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanIntegration: %w", err)
	}
	return integration, nil
}
