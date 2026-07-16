package workout_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Update(ctx context.Context, w *model.Workout) error {
	q, args, err := psql.Update("workouts").
		Set("status", w.Status).
		Set("completed_activity_id", w.CompletedActivityID).
		Set("rpe", w.RPE).
		Set("athlete_report", w.AthleteReport).
		Set("coach_comment", w.CoachComment).
		Where(sq.Eq{"id": w.ID}).
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
