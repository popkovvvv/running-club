package workout_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	w, err := scanWorkout(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanWorkout: %w", err)
	}
	segs, err := r.segments(ctx, w.ID)
	if err != nil {
		return nil, fmt.Errorf("segments: %w", err)
	}
	w.Segments = segs
	return w, nil
}

func (r *Repo) FindByCompletedActivity(ctx context.Context, activityID uuid.UUID) (*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"completed_activity_id": activityID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	w, err := scanWorkout(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanWorkout: %w", err)
	}
	return w, nil
}
