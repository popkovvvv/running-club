package activity_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Update(ctx context.Context, a *model.Activity) error {
	q, args, err := psql.Update("activities").
		Set("title", a.Title).
		Set("when_label", a.WhenLabel).
		Set("dist_km", a.DistKm).
		Set("distance_meters", a.DistanceMeters).
		Set("duration", a.Duration).
		Set("pace", a.Pace).
		Set("moving_seconds", a.MovingSeconds).
		Set("elapsed_seconds", a.ElapsedSeconds).
		Set("hr", a.HR).
		Set("average_heartrate", a.AverageHeartrate).
		Set("elevation_gain", a.ElevationGain).
		Set("started_at", a.StartedAt).
		Set("updated_at", a.UpdatedAt).
		Where(sq.Eq{"id": a.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
