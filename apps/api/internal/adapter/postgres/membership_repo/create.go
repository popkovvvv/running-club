package membership_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, m *model.Membership) error {
	q, args, err := psql.Insert("memberships").
		Columns(membershipColumns...).
		Values(m.ID, m.UserID, m.ClubID, m.Status, m.CreatedAt, m.UpdatedAt).
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
