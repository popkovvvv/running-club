package plan_week_repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var planWeekColumns = []string{"id", "club_id", "week_index", "range_label", "plan_label"}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type scannable interface {
	Scan(dest ...any) error
}

func scanPlanWeek(row scannable) (*model.PlanWeek, error) {
	var w model.PlanWeek
	if err := row.Scan(&w.ID, &w.ClubID, &w.WeekIndex, &w.RangeLabel, &w.PlanLabel); err != nil {
		return nil, err
	}
	return &w, nil
}
