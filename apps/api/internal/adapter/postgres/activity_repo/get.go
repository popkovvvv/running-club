package activity_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Activity, error) {
	q, args, err := psql.Select(activityColumns...).
		From("activities").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	activity, err := scanActivity(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanActivity: %w", err)
	}
	return activity, nil
}

func (r *Repo) GetByUserSourceExternalID(ctx context.Context, userID uuid.UUID, source, externalID string) (*model.Activity, error) {
	q, args, err := psql.Select(activityColumns...).
		From("activities").
		Where(sq.Eq{"user_id": userID, "source": source, "external_id": externalID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	activity, err := scanActivity(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanActivity: %w", err)
	}
	return activity, nil
}
