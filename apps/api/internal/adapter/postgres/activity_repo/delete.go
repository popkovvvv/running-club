package activity_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) DeleteByUserSourceExternalID(ctx context.Context, userID uuid.UUID, source, externalID string) error {
	q, args, err := psql.Delete("activities").
		Where(sq.Eq{"user_id": userID, "source": source, "external_id": externalID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
