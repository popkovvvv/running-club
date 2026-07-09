package activity_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, a *model.Activity) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO activities (id, user_id, title, when_label, dist_km, duration, pace, hr, kudos, comments, route_svg, start_x, start_y, end_x, end_y, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		a.ID, a.UserID, a.Title, a.WhenLabel, a.DistKm, a.Duration, a.Pace, a.HR, a.Kudos, a.Comments, a.RouteSVG, a.StartX, a.StartY, a.EndX, a.EndY, a.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Activity, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, when_label, dist_km, duration, pace, hr, kudos, comments, route_svg, start_x, start_y, end_x, end_y, created_at
		FROM activities WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Activity
	for rows.Next() {
		var a model.Activity
		var created time.Time
		if err := rows.Scan(&a.ID, &a.UserID, &a.Title, &a.WhenLabel, &a.DistKm, &a.Duration, &a.Pace, &a.HR, &a.Kudos, &a.Comments, &a.RouteSVG, &a.StartX, &a.StartY, &a.EndX, &a.EndY, &created); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		a.CreatedAt = created
		out = append(out, &a)
	}
	return out, rows.Err()
}

func (r *Repo) CreatePR(ctx context.Context, p *model.PR) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO prs (id, user_id, distance, time, date_label, pending) VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.UserID, p.Distance, p.Time, p.DateLabel, p.Pending)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindPRs(ctx context.Context, userID uuid.UUID) ([]*model.PR, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, distance, time, date_label, pending FROM prs WHERE user_id=$1`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.PR
	for rows.Next() {
		var p model.PR
		if err := rows.Scan(&p.ID, &p.UserID, &p.Distance, &p.Time, &p.DateLabel, &p.Pending); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *Repo) CreateRace(ctx context.Context, race *model.Race) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO races (id, club_id, user_id, name, date_label, dist, goal, days_left, finished, result)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		race.ID, race.ClubID, race.UserID, race.Name, race.DateLabel, race.Dist, race.Goal, race.DaysLeft, race.Finished, race.Result)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindRaces(ctx context.Context, userID uuid.UUID) ([]*model.Race, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, club_id, user_id, name, date_label, dist, goal, days_left, finished, result
		FROM races WHERE user_id=$1 OR user_id IS NULL ORDER BY days_left`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Race
	for rows.Next() {
		var race model.Race
		if err := rows.Scan(&race.ID, &race.ClubID, &race.UserID, &race.Name, &race.DateLabel, &race.Dist, &race.Goal, &race.DaysLeft, &race.Finished, &race.Result); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, &race)
	}
	return out, rows.Err()
}

func (r *Repo) CreateMonthStat(ctx context.Context, userID uuid.UUID, m model.MonthStat) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO month_stats (id, user_id, month, km, tr, pace, diff) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		uuid.New(), userID, m.Month, m.Km, m.Tr, m.Pace, m.Diff)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindMonthStats(ctx context.Context, userID uuid.UUID) ([]model.MonthStat, error) {
	rows, err := r.pool.Query(ctx, `SELECT month, km, tr, pace, diff FROM month_stats WHERE user_id=$1`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []model.MonthStat
	for rows.Next() {
		var m model.MonthStat
		if err := rows.Scan(&m.Month, &m.Km, &m.Tr, &m.Pace, &m.Diff); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
