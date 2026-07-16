package workout_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) ReplaceClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int, workouts []*model.Workout) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	delQ, delArgs, err := psql.Delete("workouts").
		Where(sq.Eq{"club_id": clubID, "week_index": weekIndex, "is_club_template": true}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql delete: %w", err)
	}
	if _, err := tx.Exec(ctx, delQ, delArgs...); err != nil {
		return fmt.Errorf("tx.Exec delete templates: %w", err)
	}

	for _, w := range workouts {
		insQ, insArgs, err := psql.Insert("workouts").
			Columns(workoutColumns...).
			Values(workoutValues(w)...).
			ToSql()
		if err != nil {
			return fmt.Errorf("ToSql insert: %w", err)
		}
		if _, err := tx.Exec(ctx, insQ, insArgs...); err != nil {
			return fmt.Errorf("tx.Exec insert template: %w", err)
		}
		for i, s := range w.Segments {
			s.WorkoutID = w.ID
			s.SortOrder = i
			if s.ID == uuid.Nil {
				s.ID = uuid.New()
			}
			segQ, segArgs, err := psql.Insert("segments").
				Columns("id", "workout_id", "kind", "title", "dist_km", "pace", "sort_order").
				Values(s.ID, s.WorkoutID, s.Kind, s.Title, s.DistKm, s.Pace, s.SortOrder).
				ToSql()
			if err != nil {
				return fmt.Errorf("ToSql segment: %w", err)
			}
			if _, err := tx.Exec(ctx, segQ, segArgs...); err != nil {
				return fmt.Errorf("tx.Exec insert segment: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
