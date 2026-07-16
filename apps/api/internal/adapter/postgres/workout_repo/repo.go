package workout_repo

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

const workoutColumns = `
	id, club_id, user_id, kind, workout_type, day_label, tag, title, description,
	dist_km, hr, week_index, scheduled_date, status,
	completed_activity_id, assigned_by, is_club_template, announce_id, created_at`

func (r *Repo) Create(ctx context.Context, w *model.Workout) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO workouts (
			id, club_id, user_id, kind, workout_type, day_label, tag, title, description,
			dist_km, hr, week_index, scheduled_date, status,
			completed_activity_id, assigned_by, is_club_template, announce_id, created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
		w.ID, w.ClubID, w.UserID, w.Kind, w.WorkoutType, w.DayLabel, w.Tag, w.Title, w.Description,
		w.DistKm, w.HR, w.WeekIndex, w.ScheduledDate, w.Status,
		w.CompletedActivityID, w.AssignedBy, w.IsClubTemplate, w.AnnounceID, w.CreatedAt)
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

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*model.Workout, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+workoutColumns+` FROM workouts WHERE id=$1`, id)
	w, err := scanWorkout(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanWorkout: %w", err)
	}
	segs, err := r.segments(ctx, w.ID)
	if err != nil {
		return nil, fmt.Errorf("segments: %w", err)
	}
	w.Segments = segs
	return w, nil
}

func (r *Repo) FindByUserWeek(ctx context.Context, userID uuid.UUID, week int, kind model.WorkoutKind) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT`+workoutColumns+`
		FROM workouts WHERE user_id=$1 AND week_index=$2 AND kind=$3 AND is_club_template=false
		ORDER BY scheduled_date NULLS LAST, created_at`, userID, week, kind)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT`+workoutColumns+`
		FROM workouts WHERE user_id=$1 AND is_club_template=false
		ORDER BY scheduled_date NULLS LAST, created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindOwnByUser(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT`+workoutColumns+`
		FROM workouts WHERE user_id=$1 AND kind='own' ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) FindClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT`+workoutColumns+`
		FROM workouts WHERE club_id=$1 AND week_index=$2 AND is_club_template=true
		ORDER BY scheduled_date NULLS LAST, created_at`, clubID, weekIndex)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	return r.scanWorkouts(ctx, rows)
}

func (r *Repo) ReplaceClubTemplates(ctx context.Context, clubID uuid.UUID, weekIndex int, workouts []*model.Workout) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		DELETE FROM workouts WHERE club_id=$1 AND week_index=$2 AND is_club_template=true`, clubID, weekIndex); err != nil {
		return fmt.Errorf("tx.Exec delete templates: %w", err)
	}

	for _, w := range workouts {
		_, err := tx.Exec(ctx, `
			INSERT INTO workouts (
				id, club_id, user_id, kind, workout_type, day_label, tag, title, description,
				dist_km, hr, week_index, scheduled_date, status,
				completed_activity_id, assigned_by, is_club_template, announce_id, created_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`,
			w.ID, w.ClubID, w.UserID, w.Kind, w.WorkoutType, w.DayLabel, w.Tag, w.Title, w.Description,
			w.DistKm, w.HR, w.WeekIndex, w.ScheduledDate, w.Status,
			w.CompletedActivityID, w.AssignedBy, w.IsClubTemplate, w.AnnounceID, w.CreatedAt)
		if err != nil {
			return fmt.Errorf("tx.Exec insert template: %w", err)
		}
		for i, s := range w.Segments {
			s.WorkoutID = w.ID
			s.SortOrder = i
			if s.ID == uuid.Nil {
				s.ID = uuid.New()
			}
			if _, err := tx.Exec(ctx, `
				INSERT INTO segments (id, workout_id, kind, title, dist_km, pace, sort_order)
				VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				s.ID, s.WorkoutID, s.Kind, s.Title, s.DistKm, s.Pace, s.SortOrder); err != nil {
				return fmt.Errorf("tx.Exec insert segment: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}

func (r *Repo) DeleteClubAssignedPlans(ctx context.Context, userID uuid.UUID, weekIndex int) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM workouts
		WHERE user_id=$1 AND week_index=$2 AND kind='plan' AND is_club_template=false
			AND assigned_by IS NULL AND announce_id IS NULL`,
		userID, weekIndex)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, w *model.Workout) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE workouts SET
			status=$2,
			completed_activity_id=$3
		WHERE id=$1`,
		w.ID, w.Status, w.CompletedActivityID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *Repo) FindByCompletedActivity(ctx context.Context, activityID uuid.UUID) (*model.Workout, error) {
	row := r.pool.QueryRow(ctx, `SELECT`+workoutColumns+` FROM workouts WHERE completed_activity_id=$1`, activityID)
	w, err := scanWorkout(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("scanWorkout: %w", err)
	}
	return w, nil
}

func (r *Repo) FindCompletedWithoutActivity(ctx context.Context, userID uuid.UUID) ([]*model.Workout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT`+workoutColumns+`
		FROM workouts
		WHERE user_id=$1 AND status='completed' AND completed_activity_id IS NULL AND is_club_template=false
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.Workout
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, fmt.Errorf("scanWorkout: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (r *Repo) SumPlanDistByUserWeek(ctx context.Context, userID uuid.UUID, weekIndex int) (float64, error) {
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(dist_km), 0) FROM workouts
		WHERE user_id=$1 AND week_index=$2 AND kind='plan' AND is_club_template=false`,
		userID, weekIndex).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("QueryRow: %w", err)
	}
	return sum, nil
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

func (r *Repo) DeleteByUserAndAnnounce(ctx context.Context, userID, announceID uuid.UUID) error {
	ct, err := r.pool.Exec(ctx, `
		DELETE FROM workouts WHERE user_id=$1 AND announce_id=$2`, userID, announceID)
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

func (r *Repo) scanWorkouts(ctx context.Context, rows pgx.Rows) ([]*model.Workout, error) {
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

type scannable interface {
	Scan(dest ...any) error
}

func scanWorkout(row scannable) (*model.Workout, error) {
	var w model.Workout
	var clubID *uuid.UUID
	var scheduledDate *time.Time
	var completedActivityID *uuid.UUID
	var assignedBy *uuid.UUID
	var announceID *uuid.UUID
	var created time.Time
	if err := row.Scan(
		&w.ID, &clubID, &w.UserID, &w.Kind, &w.WorkoutType, &w.DayLabel, &w.Tag, &w.Title, &w.Description,
		&w.DistKm, &w.HR, &w.WeekIndex, &scheduledDate, &w.Status,
		&completedActivityID, &assignedBy, &w.IsClubTemplate, &announceID, &created,
	); err != nil {
		return nil, err
	}
	w.ClubID = clubID
	w.ScheduledDate = scheduledDate
	w.CompletedActivityID = completedActivityID
	w.AssignedBy = assignedBy
	w.AnnounceID = announceID
	w.CreatedAt = created
	return &w, nil
}
