package workout_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, w *model.Workout) error {
	q, args, err := psql.Insert("workouts").
		Columns(workoutColumns...).
		Values(workoutValues(w)...).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	_, err = r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	for i, s := range w.Segments {
		s.WorkoutID = w.ID
		s.SortOrder = i
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}
		if err := r.createSegment(ctx, &s); err != nil {
			return fmt.Errorf("createSegment: %w", err)
		}
	}
	return nil
}
