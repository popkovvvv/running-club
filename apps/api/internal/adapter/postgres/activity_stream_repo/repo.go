package activity_stream_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) ReplaceByActivityID(ctx context.Context, activityID uuid.UUID, streams []*model.ActivityStream) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM activity_streams WHERE activity_id=$1`, activityID); err != nil {
		return fmt.Errorf("tx.Exec delete: %w", err)
	}
	for _, stream := range streams {
		if _, err := tx.Exec(ctx, `
			INSERT INTO activity_streams (id, activity_id, type, data_json, created_at, updated_at)
			VALUES ($1,$2,$3,$4::jsonb,$5,$6)`,
			stream.ID, activityID, stream.Type, stream.DataJSON, stream.CreatedAt, stream.UpdatedAt); err != nil {
			return fmt.Errorf("tx.Exec insert: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
