package announce_repo

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var announceColumns = []string{
	"id", "club_id", "place", "day_label", "time", "group_name", "note", "starts_on", "going_count", "created_at",
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanAnnounce(row scannable) (*model.Announce, error) {
	var a model.Announce
	var created time.Time
	if err := row.Scan(
		&a.ID, &a.ClubID, &a.Place, &a.DayLabel, &a.Time, &a.GroupName, &a.Note,
		&a.StartsOn, &a.GoingCount, &created,
	); err != nil {
		return nil, err
	}
	a.CreatedAt = created
	return &a, nil
}
