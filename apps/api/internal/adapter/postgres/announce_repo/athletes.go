package announce_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindGoingAthletes(ctx context.Context, announceID uuid.UUID) ([]*model.User, error) {
	q, args, err := psql.Select("u.id", "u.name", "u.email", "u.password_hash", "u.role", "u.created_at").
		From("announce_signups s").
		Join("users u ON u.id = s.athlete_id").
		Where(sq.Eq{"s.announce_id": announceID}).
		OrderBy("s.created_at").
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
		var u model.User
		var created time.Time
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &created); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		u.CreatedAt = created
		out = append(out, &u)
	}
	return out, rows.Err()
}

func (r *Repo) NextLabelForAthlete(ctx context.Context, clubID, athleteID uuid.UUID) (string, error) {
	q, args, err := psql.Select("a.day_label", "a.place").
		From("announces a").
		Join("announce_signups s ON s.announce_id = a.id").
		Where(sq.And{
			sq.Eq{"a.club_id": clubID, "s.athlete_id": athleteID},
			sq.Or{sq.Expr("a.starts_on IS NULL"), sq.Expr("a.starts_on >= CURRENT_DATE")},
		}).
		OrderBy("a.starts_on NULLS LAST", "a.created_at").
		Limit(1).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ToSql: %w", err)
	}
	var dayLabel, place string
	err = r.pool.QueryRow(ctx, q, args...).Scan(&dayLabel, &place)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", model.ErrNotFound
		}
		return "", fmt.Errorf("QueryRow: %w", err)
	}
	return dayLabel + " · " + place, nil
}
