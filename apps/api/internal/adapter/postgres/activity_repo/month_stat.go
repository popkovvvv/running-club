package activity_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) CreateMonthStat(ctx context.Context, userID uuid.UUID, m model.MonthStat) error {
	q, args, err := psql.Insert("month_stats").
		Columns("id", "user_id", "month", "km", "tr", "pace", "diff").
		Values(uuid.New(), userID, m.Month, m.Km, m.Tr, m.Pace, m.Diff).
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

func (r *Repo) FindMonthStats(ctx context.Context, userID uuid.UUID) ([]model.MonthStat, error) {
	q, args, err := psql.Select("month", "km", "tr", "pace", "diff").
		From("month_stats").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []model.MonthStat
	for rows.Next() {
		var m model.MonthStat
		if err := rows.Scan(&m.Month, &m.Km, &m.Tr, &m.Pace, &m.Diff); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
