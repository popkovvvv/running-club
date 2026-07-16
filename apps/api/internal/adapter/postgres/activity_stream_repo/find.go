package activity_stream_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) FindByActivityID(ctx context.Context, activityID uuid.UUID) ([]*model.ActivityStream, error) {
	q, args, err := psql.Select("id", "activity_id", "type", "data_json::text", "created_at", "updated_at").
		From("activity_streams").
		Where(sq.Eq{"activity_id": activityID}).
		OrderBy("type").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	var out []*model.ActivityStream
	for rows.Next() {
		var s model.ActivityStream
		if err := rows.Scan(&s.ID, &s.ActivityID, &s.Type, &s.DataJSON, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}
