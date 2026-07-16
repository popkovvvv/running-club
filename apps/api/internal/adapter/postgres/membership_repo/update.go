package membership_repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.MembershipStatus) error {
	q, args, err := psql.Update("memberships").
		Set("status", status).
		Set("updated_at", time.Now().UTC()).
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
