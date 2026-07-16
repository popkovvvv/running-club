package activity_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, a *model.Activity) error {
	q, args, err := psql.Insert("activities").
		Columns(activityColumns...).
		Values(activityValues(a)...).
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

func (r *Repo) Upsert(ctx context.Context, activity *model.Activity) error {
	q, args, err := psql.Insert("activities").
		Columns(activityColumns...).
		Values(activityValues(activity)...).
		Suffix(`ON CONFLICT (user_id, source, external_id) WHERE source <> '' AND external_id <> '' DO UPDATE SET
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
			updated_at=EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	_, err = r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	return nil
}
