package workout_repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var workoutColumns = []string{
	"id", "club_id", "user_id", "kind", "workout_type", "day_label", "tag", "title", "description",
	"dist_km", "hr", "week_index", "scheduled_date", "status",
	"completed_activity_id", "assigned_by", "is_club_template", "announce_id",
	"rpe", "athlete_report", "coach_comment", "created_at",
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
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
	var rpe *int
	var created time.Time
	if err := row.Scan(
		&w.ID, &clubID, &w.UserID, &w.Kind, &w.WorkoutType, &w.DayLabel, &w.Tag, &w.Title, &w.Description,
		&w.DistKm, &w.HR, &w.WeekIndex, &scheduledDate, &w.Status,
		&completedActivityID, &assignedBy, &w.IsClubTemplate, &announceID,
		&rpe, &w.AthleteReport, &w.CoachComment, &created,
	); err != nil {
		return nil, err
	}
	w.ClubID = clubID
	w.ScheduledDate = scheduledDate
	w.CompletedActivityID = completedActivityID
	w.AssignedBy = assignedBy
	w.AnnounceID = announceID
	w.RPE = rpe
	w.CreatedAt = created
	return &w, nil
}

func workoutValues(w *model.Workout) []any {
	return []any{
		w.ID, w.ClubID, w.UserID, w.Kind, w.WorkoutType, w.DayLabel, w.Tag, w.Title, w.Description,
		w.DistKm, w.HR, w.WeekIndex, w.ScheduledDate, w.Status,
		w.CompletedActivityID, w.AssignedBy, w.IsClubTemplate, w.AnnounceID,
		w.RPE, w.AthleteReport, w.CoachComment, w.CreatedAt,
	}
}

func (r *Repo) segments(ctx context.Context, workoutID uuid.UUID) ([]model.Segment, error) {
	q, args, err := psql.Select("id", "workout_id", "kind", "title", "dist_km", "pace", "sort_order").
		From("segments").
		Where(sq.Eq{"workout_id": workoutID}).
		OrderBy("sort_order").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
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

func (r *Repo) createSegment(ctx context.Context, s *model.Segment) error {
	q, args, err := psql.Insert("segments").
		Columns("id", "workout_id", "kind", "title", "dist_km", "pace", "sort_order").
		Values(s.ID, s.WorkoutID, s.Kind, s.Title, s.DistKm, s.Pace, s.SortOrder).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	_, err = r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
