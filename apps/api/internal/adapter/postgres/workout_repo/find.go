package workout_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"user_id": userID, "week_index": week, "kind": kind, "is_club_template": false}).
		OrderBy("scheduled_date NULLS LAST", "created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"user_id": userID, "is_club_template": false}).
		OrderBy("scheduled_date NULLS LAST", "created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindOwnByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"user_id": userID, "kind": model.WorkoutOwn}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int) ([]*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.Eq{"club_id": clubID, "week_index": weekIndex, "is_club_template": true}).
		OrderBy("scheduled_date NULLS LAST", "created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindCompletedWithoutActivity(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	q, args, err := psql.Select(workoutColumns...).
		From("workouts").
		Where(sq.And{
			sq.Eq{"user_id": userID, "status": model.WorkoutStatusCompleted, "is_club_template": false},
			sq.Expr("completed_activity_id IS NULL"),
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Workout
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, fmt.Errorf("scanWorkout: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}
