package club_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, c *model.Club) error {
	q, args, err := psql.Insert("clubs").
		Columns(clubColumns...).
		Values(c.ID, c.Name, c.InviteCode, c.AccentHex, c.CoachID, c.CreatedAt).
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
