package club_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) UpdateAccent(ctx context.Context, id uuid.UUID, accent string) error {
	q, args, err := psql.Update("clubs").
		Set("accent_hex", accent).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	ct, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
