package user_repo

import (
	"context"
	"fmt"

	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (r *Repo) Create(ctx context.Context, u *model.User) error {
	q, args, err := psql.Insert("users").
		Columns(userColumns...).
		Values(u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt).
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
