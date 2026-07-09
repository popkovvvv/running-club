package membership_repo

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

func (r *Repo) Create(ctx context.Context, m *model.Membership) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO memberships (id, user_id, club_id, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		m.ID, m.UserID, m.ClubID, m.Status, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.Membership, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, club_id, status, created_at, updated_at
		FROM memberships WHERE user_id=$1 AND status='active'`, userID)
	m, err := scan(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return m, nil
}

func (r *Repo) GetByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*model.Membership, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, club_id, status, created_at, updated_at
		FROM memberships WHERE user_id=$1 AND club_id=$2`, userID, clubID)
	m, err := scan(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return m, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.MembershipStatus) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE memberships SET status=$2, updated_at=$3 WHERE id=$1`, id, status, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scan(row scannable) (*model.Membership, error) {
	var m model.Membership
	if err := row.Scan(&m.ID, &m.UserID, &m.ClubID, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
