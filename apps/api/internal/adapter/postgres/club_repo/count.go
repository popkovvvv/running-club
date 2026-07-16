package club_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) CountActiveStudents(ctx context.Context, clubID uuid.UUID) (int, error) {
	q, args, err := psql.Select("COUNT(*)").
		From("memberships").
		Where(sq.Eq{"club_id": clubID, "status": model.MembershipActive}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ToSql: %w", err)
	}
	var n int
	err = r.pool.QueryRow(ctx, q, args...).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return n, nil
}
