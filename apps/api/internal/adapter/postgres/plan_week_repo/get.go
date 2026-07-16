package plan_week_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error) {
	q, args, err := psql.Select(planWeekColumns...).
		From("plan_weeks").
		Where(sq.Eq{"club_id": clubID, "week_index": weekIndex}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	w, err := scanPlanWeek(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return w, nil
}
