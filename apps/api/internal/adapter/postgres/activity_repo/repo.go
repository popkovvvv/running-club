package activity_repo

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

func (r *Repo) Create(ctx context.Context, a *model.Activity) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO activities (
			id, user_id, source, external_id, sport_type, title, when_label, dist_km, distance_meters,
			duration, pace, moving_seconds, elapsed_seconds, hr, average_heartrate, max_heartrate,
			elevation_gain, kudos, comments, visibility, polyline, route_svg, start_x, start_y,
			end_x, end_y, started_at, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29)`,
		a.ID, a.UserID, a.Source, a.ExternalID, a.SportType, a.Title, a.WhenLabel, a.DistKm, a.DistanceMeters,
		a.Duration, a.Pace, a.MovingSeconds, a.ElapsedSeconds, a.HR, a.AverageHeartrate, a.MaxHeartrate,
		a.ElevationGain, a.Kudos, a.Comments, a.Visibility, a.Polyline, a.RouteSVG, a.StartX, a.StartY,
		a.EndX, a.EndY, a.StartedAt, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Activity, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, source, external_id, sport_type, title, when_label, dist_km, distance_meters,
			duration, pace, moving_seconds, elapsed_seconds, hr, average_heartrate, max_heartrate,
			elevation_gain, kudos, comments, visibility, polyline, route_svg, start_x, start_y,
			end_x, end_y, started_at, created_at, updated_at
		FROM activities WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Activity
	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return nil, fmt.Errorf("scanActivity: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) Upsert(ctx context.Context, activity *model.Activity) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO activities (
			id, user_id, source, external_id, sport_type, title, when_label, dist_km, distance_meters,
			duration, pace, moving_seconds, elapsed_seconds, hr, average_heartrate, max_heartrate,
			elevation_gain, kudos, comments, visibility, polyline, route_svg, start_x, start_y,
			end_x, end_y, started_at, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29)
		ON CONFLICT (user_id, source, external_id) DO UPDATE SET
			sport_type=EXCLUDED.sport_type,
			title=EXCLUDED.title,
			when_label=EXCLUDED.when_label,
			dist_km=EXCLUDED.dist_km,
			distance_meters=EXCLUDED.distance_meters,
			duration=EXCLUDED.duration,
			pace=EXCLUDED.pace,
			moving_seconds=EXCLUDED.moving_seconds,
			elapsed_seconds=EXCLUDED.elapsed_seconds,
			hr=EXCLUDED.hr,
			average_heartrate=EXCLUDED.average_heartrate,
			max_heartrate=EXCLUDED.max_heartrate,
			elevation_gain=EXCLUDED.elevation_gain,
			kudos=EXCLUDED.kudos,
			comments=EXCLUDED.comments,
			visibility=EXCLUDED.visibility,
			polyline=EXCLUDED.polyline,
			route_svg=EXCLUDED.route_svg,
			start_x=EXCLUDED.start_x,
			start_y=EXCLUDED.start_y,
			end_x=EXCLUDED.end_x,
			end_y=EXCLUDED.end_y,
			started_at=EXCLUDED.started_at,
			updated_at=EXCLUDED.updated_at`,
		activity.ID, activity.UserID, activity.Source, activity.ExternalID, activity.SportType, activity.Title, activity.WhenLabel, activity.DistKm, activity.DistanceMeters,
		activity.Duration, activity.Pace, activity.MovingSeconds, activity.ElapsedSeconds, activity.HR, activity.AverageHeartrate, activity.MaxHeartrate,
		activity.ElevationGain, activity.Kudos, activity.Comments, activity.Visibility, activity.Polyline, activity.RouteSVG, activity.StartX, activity.StartY,
		activity.EndX, activity.EndY, activity.StartedAt, activity.CreatedAt, activity.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	return nil
}

func (r *Repo) DeleteByUserSourceExternalID(ctx context.Context, userID uuid.UUID, source, externalID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM activities WHERE user_id=$1 AND source=$2 AND external_id=$3`,
		userID, source, externalID)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
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

func (r *Repo) SumDistByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(dist_km), 0) FROM activities WHERE user_id=$1`, userID).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return sum, nil
}

func (r *Repo) SumDistByClubAthletes(ctx context.Context, clubID uuid.UUID) (float64, error) {
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(a.dist_km), 0)
		FROM activities a
		JOIN memberships m ON m.user_id = a.user_id
		WHERE m.club_id=$1 AND m.status='active'`, clubID).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return sum, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanActivity(row scannable) (*model.Activity, error) {
	var activity model.Activity
	var startedAt *time.Time
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Source,
		&activity.ExternalID,
		&activity.SportType,
		&activity.Title,
		&activity.WhenLabel,
		&activity.DistKm,
		&activity.DistanceMeters,
		&activity.Duration,
		&activity.Pace,
		&activity.MovingSeconds,
		&activity.ElapsedSeconds,
		&activity.HR,
		&activity.AverageHeartrate,
		&activity.MaxHeartrate,
		&activity.ElevationGain,
		&activity.Kudos,
		&activity.Comments,
		&activity.Visibility,
		&activity.Polyline,
		&activity.RouteSVG,
		&activity.StartX,
		&activity.StartY,
		&activity.EndX,
		&activity.EndY,
		&startedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	activity.StartedAt = startedAt
	activity.CreatedAt = createdAt
	activity.UpdatedAt = updatedAt
	return &activity, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Activity, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, source, external_id, sport_type, title, when_label, dist_km, distance_meters,
			duration, pace, moving_seconds, elapsed_seconds, hr, average_heartrate, max_heartrate,
			elevation_gain, kudos, comments, visibility, polyline, route_svg, start_x, start_y,
			end_x, end_y, started_at, created_at, updated_at
		FROM activities WHERE id=$1`, id)
	activity, err := scanActivity(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanActivity: %w", err)
	}
	return activity, nil
}

func (r *Repo) SumDistByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (float64, error) {
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(dist_km), 0) FROM activities
		WHERE user_id=$1 AND COALESCE(started_at, created_at) >= $2`,
		userID, since).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return sum, nil
}

func (r *Repo) GetByUserSourceExternalID(ctx context.Context, userID uuid.UUID, source, externalID string) (*model.Activity, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, source, external_id, sport_type, title, when_label, dist_km, distance_meters,
			duration, pace, moving_seconds, elapsed_seconds, hr, average_heartrate, max_heartrate,
			elevation_gain, kudos, comments, visibility, polyline, route_svg, start_x, start_y,
			end_x, end_y, started_at, created_at, updated_at
		FROM activities
		WHERE user_id=$1 AND source=$2 AND external_id=$3`,
		userID, source, externalID)
	activity, err := scanActivity(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanActivity: %w", err)
	}
	return activity, nil
}
