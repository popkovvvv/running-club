package plan_week_repo

import (
	"context"
	"errors"
	"fmt"

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

func (r *Repo) Create(ctx context.Context, w *model.PlanWeek) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO plan_weeks (id, club_id, week_index, range_label, plan_label)
		VALUES ($1,$2,$3,$4,$5)`,
		w.ID, w.ClubID, w.WeekIndex, w.RangeLabel, w.PlanLabel)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindByClub(ctx context.Context, clubID uuid.UUID) ([]*model.PlanWeek, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, club_id, week_index, range_label, plan_label
		FROM plan_weeks WHERE club_id=$1 ORDER BY week_index`, clubID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.PlanWeek
	for rows.Next() {
		w, err := scanPlanWeek(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (r *Repo) Upsert(ctx context.Context, w *model.PlanWeek) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO plan_weeks (id, club_id, week_index, range_label, plan_label)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (club_id, week_index) DO UPDATE SET
			range_label=EXCLUDED.range_label,
			plan_label=EXCLUDED.plan_label`,
		w.ID, w.ClubID, w.WeekIndex, w.RangeLabel, w.PlanLabel)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) GetByClubAndIndex(ctx context.Context, clubID uuid.UUID, weekIndex int) (*model.PlanWeek, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, club_id, week_index, range_label, plan_label
		FROM plan_weeks WHERE club_id=$1 AND week_index=$2`, clubID, weekIndex)
	w, err := scanPlanWeek(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return w, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanPlanWeek(row scannable) (*model.PlanWeek, error) {
	var w model.PlanWeek
	if err := row.Scan(&w.ID, &w.ClubID, &w.WeekIndex, &w.RangeLabel, &w.PlanLabel); err != nil {
		return nil, err
	}
	return &w, nil
}
