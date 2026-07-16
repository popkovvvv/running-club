package announce_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, a *model.Announce) error {
	q, args, err := psql.Insert("announces").
		Columns(announceColumns...).
		Values(a.ID, a.ClubID, a.Place, a.DayLabel, a.Time, a.GroupName, a.Note, a.StartsOn, a.GoingCount, a.CreatedAt).
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
