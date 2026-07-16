package club_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error) {
	q, args, err := psql.Select(clubColumns...).
		From("clubs").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	c, err := scanClub(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return c, nil
}

func (r *Repo) GetByInviteCode(ctx context.Context, code string) (*model.Club, error) {
	q, args, err := psql.Select(clubColumns...).
		From("clubs").
		Where(sq.Eq{"invite_code": code}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	c, err := scanClub(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return c, nil
}

func (r *Repo) GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error) {
	q, args, err := psql.Select(clubColumns...).
		From("clubs").
		Where(sq.Eq{"coach_id": coachID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	c, err := scanClub(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return c, nil
}
