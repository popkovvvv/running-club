package user_repo

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var userColumns = []string{"id", "name", "email", "password_hash", "role", "created_at"}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanUser(row scannable) (*model.User, error) {
	var u model.User
	var created time.Time
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &created); err != nil {
		return nil, err
	}
	u.CreatedAt = created
	return &u, nil
}
