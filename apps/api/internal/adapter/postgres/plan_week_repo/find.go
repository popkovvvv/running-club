package plan_week_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.PlanWeek, error) {
	q, args, err := psql.Select(planWeekColumns...).
		From("plan_weeks").
		Where(sq.Eq{"club_id": clubID}).
		OrderBy("week_index").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.PlanWeek
	for rows.Next() {
		w, err := scanPlanWeek(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}
