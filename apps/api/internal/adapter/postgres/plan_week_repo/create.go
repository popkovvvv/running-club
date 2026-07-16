package plan_week_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, w *model.PlanWeek) error {
	q, args, err := psql.Insert("plan_weeks").
		Columns(planWeekColumns...).
		Values(w.ID, w.ClubID, w.WeekIndex, w.RangeLabel, w.PlanLabel).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	_, err = r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
