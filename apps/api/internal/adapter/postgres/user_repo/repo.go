package user_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, u *model.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, name, email, password_hash, role, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at FROM users WHERE email=$1`, email)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return u, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at FROM users WHERE id=$1`, id)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return u, nil
}

func (r *Repo) FindAthletesByClub(ctx context.Context, clubID uuid.UUID) ([]*model.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.name, u.email, u.password_hash, u.role, u.created_at
		FROM users u
		JOIN memberships m ON m.user_id = u.id
		WHERE m.club_id=$1 AND m.status='active' AND u.role='athlete'
		ORDER BY u.name`, clubID)
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

type scannable interface {
	Scan(dest ...any) error
}

func scanUser(row scannable) (*model.User, error) {
	var u model.User
	var created time.Time
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &created); err != nil {
		return nil, err
	}
	u.CreatedAt = created
	return &u, nil
}
