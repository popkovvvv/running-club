package membership_repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var membershipColumns = []string{"id", "user_id", "club_id", "status", "created_at", "updated_at"}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scan(row scannable) (*model.Membership, error) {
	var m model.Membership
	if err := row.Scan(&m.ID, &m.UserID, &m.ClubID, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
