package announce_repo

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

func (r *Repo) Create(ctx context.Context, a *model.Announce) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO announces (id, club_id, place, day_label, time, group_name, note, going_count, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		a.ID, a.ClubID, a.Place, a.DayLabel, a.Time, a.GroupName, a.Note, a.GoingCount, a.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.Announce, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, club_id, place, day_label, time, group_name, note, going_count, created_at
		FROM announces WHERE club_id=$1 ORDER BY created_at DESC`, clubID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Announce
	for rows.Next() {
		a, err := scanAnnounce(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Announce, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, club_id, place, day_label, time, group_name, note, going_count, created_at
		FROM announces WHERE id=$1`, id)
	a, err := scanAnnounce(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return a, nil
}

func (r *Repo) IncGoing(ctx context.Context, id uuid.UUID, delta int) error {
	_, err := r.pool.Exec(ctx, `UPDATE announces SET going_count = going_count + $2 WHERE id=$1`, id, delta)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) CreateSignup(ctx context.Context, s *model.AnnounceSignup) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO announce_signups (id, announce_id, athlete_id, created_at)
		VALUES ($1,$2,$3,$4)`, s.ID, s.AnnounceID, s.AthleteID, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) DeleteSignup(ctx context.Context, announceID, athleteID uuid.UUID) error {
	ct, err := r.pool.Exec(ctx, `
		DELETE FROM announce_signups WHERE announce_id=$1 AND athlete_id=$2`, announceID, athleteID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *Repo) HasSignup(ctx context.Context, announceID, athleteID uuid.UUID) (bool, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM announce_signups WHERE announce_id=$1 AND athlete_id=$2`, announceID, athleteID).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("QueryRow: %w", err)
	}
	return n > 0, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanAnnounce(row scannable) (*model.Announce, error) {
	var a model.Announce
	var created time.Time
	if err := row.Scan(&a.ID, &a.ClubID, &a.Place, &a.DayLabel, &a.Time, &a.GroupName, &a.Note, &a.GoingCount, &created); err != nil {
		return nil, err
	}
	a.CreatedAt = created
	return &a, nil
}
