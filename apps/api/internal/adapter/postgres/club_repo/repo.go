package club_repo

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

func (r *Repo) Create(ctx context.Context, c *model.Club) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO clubs (id, name, invite_code, accent_hex, coach_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		c.ID, c.Name, c.InviteCode, c.AccentHex, c.CoachID, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Club, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, invite_code, accent_hex, coach_id, created_at FROM clubs WHERE id=$1`, id)
	return scanClub(row)
}

func (r *Repo) GetByInviteCode(ctx context.Context, code string) (*model.Club, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, invite_code, accent_hex, coach_id, created_at FROM clubs WHERE invite_code=$1`, code)
	c, err := scanClub(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return c, nil
}

func (r *Repo) GetByCoachID(ctx context.Context, coachID uuid.UUID) (*model.Club, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, invite_code, accent_hex, coach_id, created_at FROM clubs WHERE coach_id=$1`, coachID)
	c, err := scanClub(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return c, nil
}

func (r *Repo) UpdateAccent(ctx context.Context, id uuid.UUID, accent string) error {
	ct, err := r.pool.Exec(ctx, `UPDATE clubs SET accent_hex=$2 WHERE id=$1`, id, accent)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *Repo) CountActiveStudents(ctx context.Context, clubID uuid.UUID) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM memberships WHERE club_id=$1 AND status='active'`, clubID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return n, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanClub(row scannable) (*model.Club, error) {
	var c model.Club
	var created time.Time
	if err := row.Scan(&c.ID, &c.Name, &c.InviteCode, &c.AccentHex, &c.CoachID, &created); err != nil {
		return nil, err
	}
	c.CreatedAt = created
	return &c, nil
}
