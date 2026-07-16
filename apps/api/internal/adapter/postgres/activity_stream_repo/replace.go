package activity_stream_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) ReplaceByActivityID(ctx context.Context, activityID uuid.UUID, streams []*model.ActivityStream) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	delQ, delArgs, err := psql.Delete("activity_streams").
		Where(sq.Eq{"activity_id": activityID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql delete: %w", err)
	}
	if _, err := tx.Exec(ctx, delQ, delArgs...); err != nil {
		return fmt.Errorf("tx.Exec delete: %w", err)
	}
	for _, stream := range streams {
		insQ, insArgs, err := psql.Insert("activity_streams").
			Columns("id", "activity_id", "type", "data_json", "created_at", "updated_at").
			Values(stream.ID, activityID, stream.Type, sq.Expr("?::jsonb", stream.DataJSON), stream.CreatedAt, stream.UpdatedAt).
			ToSql()
		if err != nil {
			return fmt.Errorf("ToSql insert: %w", err)
		}
		if _, err := tx.Exec(ctx, insQ, insArgs...); err != nil {
			return fmt.Errorf("tx.Exec insert: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
