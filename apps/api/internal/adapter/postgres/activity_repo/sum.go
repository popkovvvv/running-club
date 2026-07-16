package activity_repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *Repo) SumDistByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	q, args, err := psql.Select("COALESCE(SUM(dist_km), 0)").
		From("activities").
		Where(sq.Eq{"user_id": userID}).
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

func (r *Repo) SumDistByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (float64, error) {
	q, args, err := psql.Select("COALESCE(SUM(dist_km), 0)").
		From("activities").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Expr("COALESCE(started_at, created_at) >= ?", since),
		}).
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

func (r *Repo) SumDistByClubAthletes(ctx context.Context, clubID uuid.UUID) (float64, error) {
	q, args, err := psql.Select("COALESCE(SUM(a.dist_km), 0)").
		From("activities a").
		Join("memberships m ON m.user_id = a.user_id").
		Where(sq.Eq{"m.club_id": clubID, "m.status": "active"}).
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
