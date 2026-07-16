package activity_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) CreatePR(ctx context.Context, p *model.PR) error {
	q, args, err := psql.Insert("prs").
		Columns("id", "user_id", "distance", "time", "date_label", "pending").
		Values(p.ID, p.UserID, p.Distance, p.Time, p.DateLabel, p.Pending).
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

func (r *Repo) FindPRs(ctx context.Context, userID uuid.UUID) ([]*model.PR, error) {
	q, args, err := psql.Select("id", "user_id", "distance", "time", "date_label", "pending").
		From("prs").
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
	var out []*model.PR
	for rows.Next() {
		var p model.PR
		if err := rows.Scan(&p.ID, &p.UserID, &p.Distance, &p.Time, &p.DateLabel, &p.Pending); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}
