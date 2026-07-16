package user_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	q, args, err := psql.Select(userColumns...).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	u, err := scanUser(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return u, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	q, args, err := psql.Select(userColumns...).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	u, err := scanUser(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return u, nil
}
