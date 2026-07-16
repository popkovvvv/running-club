package membership_repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error) {
	q, args, err := psql.Select(membershipColumns...).
		From("memberships").
		Where(sq.Eq{"user_id": userID, "status": model.MembershipActive}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	m, err := scan(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return m, nil
}

func (r *Repo) GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error) {
	q, args, err := psql.Select(membershipColumns...).
		From("memberships").
		Where(sq.Eq{"user_id": userID, "club_id": clubID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	m, err := scan(r.pool.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return m, nil
}
