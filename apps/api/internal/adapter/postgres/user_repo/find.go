package user_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error) {
	q, args, err := psql.Select(
		"u.id", "u.name", "u.email", "u.password_hash", "u.role", "u.created_at",
	).
		From("users u").
		Join("memberships m ON m.user_id = u.id").
		Where(sq.Eq{"m.club_id": clubID, "m.status": model.MembershipActive, "u.role": model.RoleAthlete}).
		OrderBy("u.name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
