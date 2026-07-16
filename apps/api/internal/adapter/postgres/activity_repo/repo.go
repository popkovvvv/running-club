package activity_repo

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var activityColumns = []string{
	"id", "user_id", "source", "external_id", "sport_type", "title", "when_label", "dist_km", "distance_meters",
	"duration", "pace", "moving_seconds", "elapsed_seconds", "hr", "average_heartrate", "max_heartrate",
	"elevation_gain", "kudos", "comments", "visibility", "polyline", "route_svg", "start_x", "start_y",
	"end_x", "end_y", "started_at", "created_at", "updated_at",
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

func activityValues(a *model.Activity) []any {
	return []any{
		a.ID, a.UserID, a.Source, a.ExternalID, a.SportType, a.Title, a.WhenLabel, a.DistKm, a.DistanceMeters,
		a.Duration, a.Pace, a.MovingSeconds, a.ElapsedSeconds, a.HR, a.AverageHeartrate, a.MaxHeartrate,
		a.ElevationGain, a.Kudos, a.Comments, a.Visibility, a.Polyline, a.RouteSVG, a.StartX, a.StartY,
		a.EndX, a.EndY, a.StartedAt, a.CreatedAt, a.UpdatedAt,
	}
}
