package announce_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.Announce, error) {
	q, args, err := psql.Select(announceColumns...).
		From("announces").
		Where(sq.Eq{"club_id": clubID}).
		OrderBy(`
			CASE
				WHEN starts_on IS NULL THEN 1
				WHEN starts_on < CURRENT_DATE THEN 2
				ELSE 0
			END`,
			"starts_on ASC NULLS LAST",
			"created_at DESC",
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Announce
	for rows.Next() {
		a, err := scanAnnounce(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
