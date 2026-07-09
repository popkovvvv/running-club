package workout_repo

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

func (r *Repo) Create(ctx context.Context, w *model.Workout) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO workouts (id, club_id, user_id, kind, day_label, tag, title, dist_km, duration, pace, hr, week_index, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		w.ID, w.ClubID, w.UserID, w.Kind, w.DayLabel, w.Tag, w.Title, w.DistKm, w.Duration, w.Pace, w.HR, w.WeekIndex, w.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	for i, s := range w.Segments {
		s.WorkoutID = w.ID
		s.SortOrder = i
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}
		if err := r.createSegment(ctx, &s); err != nil {
			return fmt.Errorf("createSegment: %w", err)
		}
	}
	return nil
}

func (r *Repo) createSegment(ctx context.Context, s *model.Segment) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO segments (id, workout_id, kind, title, dist_km, pace, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.WorkoutID, s.Kind, s.Title, s.DistKm, s.Pace, s.SortOrder)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, club_id, user_id, kind, day_label, tag, title, dist_km, duration, pace, hr, week_index, created_at
		FROM workouts WHERE user_id=$1 AND week_index=$2 AND kind=$3 ORDER BY created_at`, userID, week, kind)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Workout
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		segs, err := r.segments(ctx, w.ID)
		if err != nil {
			return nil, fmt.Errorf("segments: %w", err)
		}
		w.Segments = segs
		out = append(out, w)
	}
	return out, rows.Err()
}

func (r *Repo) FindOwnByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, club_id, user_id, kind, day_label, tag, title, dist_km, duration, pace, hr, week_index, created_at
		FROM workouts WHERE user_id=$1 AND kind='own' ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Workout
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM workouts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *Repo) segments(ctx context.Context, workoutID uuid.UUID) ([]model.Segment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, workout_id, kind, title, dist_km, pace, sort_order
		FROM segments WHERE workout_id=$1 ORDER BY sort_order`, workoutID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []model.Segment
	for rows.Next() {
		var s model.Segment
		if err := rows.Scan(&s.ID, &s.WorkoutID, &s.Kind, &s.Title, &s.DistKm, &s.Pace, &s.SortOrder); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

type scannable interface {
	Scan(dest ...any) error
}

func scanWorkout(row scannable) (*model.Workout, error) {
	var w model.Workout
	var clubID *uuid.UUID
	var created time.Time
	if err := row.Scan(&w.ID, &clubID, &w.UserID, &w.Kind, &w.DayLabel, &w.Tag, &w.Title, &w.DistKm, &w.Duration, &w.Pace, &w.HR, &w.WeekIndex, &created); err != nil {
		return nil, err
	}
	w.ClubID = clubID
	w.CreatedAt = created
	return &w, nil
}
