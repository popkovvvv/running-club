package workout_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	q, args, err := psql.Delete("workouts").
		Where(sq.Eq{"id": id}).
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

func (r *Repo) DeleteByUserAndAnnounce(ctx context.Context, userID, announceID uuid.UUID) error {
	q, args, err := psql.Delete("workouts").
		Where(sq.Eq{"user_id": userID, "announce_id": announceID}).
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

func (r *Repo) DeleteClubAssignedPlans(ctx context.Context, userID uuid.UUID, weekIndex int) error {
	q, args, err := psql.Delete("workouts").
		Where(sq.And{
			sq.Eq{"user_id": userID, "week_index": weekIndex, "kind": model.WorkoutPlan, "is_club_template": false},
			sq.Expr("assigned_by IS NULL"),
			sq.Expr("announce_id IS NULL"),
		}).
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
