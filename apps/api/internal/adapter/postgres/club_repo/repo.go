package club_repo

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var clubColumns = []string{"id", "name", "invite_code", "accent_hex", "coach_id", "created_at"}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanClub(row scannable) (*model.Club, error) {
	var c model.Club
	var created time.Time
	if err := row.Scan(&c.ID, &c.Name, &c.InviteCode, &c.AccentHex, &c.CoachID, &created); err != nil {
		return nil, err
	}
	c.CreatedAt = created
	return &c, nil
}
