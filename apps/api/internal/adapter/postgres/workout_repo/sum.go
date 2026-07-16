package workout_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) SumPlanDistByUserWeek(ctx context.Context, userID uuid.UUID, weekIndex int) (float64, error) {
	q, args, err := psql.Select("COALESCE(SUM(dist_km), 0)").
		From("workouts").
		Where(sq.Eq{"user_id": userID, "week_index": weekIndex, "kind": model.WorkoutPlan, "is_club_template": false}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ToSql: %w", err)
	}
	var sum float64
	err = r.pool.QueryRow(ctx, q, args...).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return sum, nil
}
