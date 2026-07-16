package announce_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) CreateSignup(ctx context.Context, s *model.AnnounceSignup) error {
	q, args, err := psql.Insert("announce_signups").
		Columns("id", "announce_id", "athlete_id", "created_at").
		Values(s.ID, s.AnnounceID, s.AthleteID, s.CreatedAt).
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

func (r *Repo) DeleteSignup(ctx context.Context, announceID, athleteID uuid.UUID) error {
	q, args, err := psql.Delete("announce_signups").
		Where(sq.Eq{"announce_id": announceID, "athlete_id": athleteID}).
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

func (r *Repo) HasSignup(ctx context.Context, announceID, athleteID uuid.UUID) (bool, error) {
	q, args, err := psql.Select("COUNT(*)").
		From("announce_signups").
		Where(sq.Eq{"announce_id": announceID, "athlete_id": athleteID}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("ToSql: %w", err)
	}
	var n int
	err = r.pool.QueryRow(ctx, q, args...).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("QueryRow: %w", err)
	}
	return n > 0, nil
}
