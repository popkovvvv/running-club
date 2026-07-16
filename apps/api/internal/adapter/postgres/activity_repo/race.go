package activity_repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) CreateRace(ctx context.Context, race *model.Race) error {
	q, args, err := psql.Insert("races").
		Columns("id", "club_id", "user_id", "name", "date_label", "dist", "goal", "days_left", "finished", "result").
		Values(race.ID, race.ClubID, race.UserID, race.Name, race.DateLabel, race.Dist, race.Goal, race.DaysLeft, race.Finished, race.Result).
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

func (r *Repo) FindRaces(ctx context.Context, userID uuid.UUID) ([]*model.Race, error) {
	q, args, err := psql.Select("id", "club_id", "user_id", "name", "date_label", "dist", "goal", "days_left", "finished", "result").
		From("races").
		Where(sq.Or{sq.Eq{"user_id": userID}, sq.Expr("user_id IS NULL")}).
		OrderBy("days_left").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := r.pool.Query(ctx, q, args...)
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
