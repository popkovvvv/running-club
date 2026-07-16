//go:build integration

package plan_week_repo_test

import (
	"context"
	"testing"

	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/plan_week_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestPlanWeekCRUD(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	repo := plan_week_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	require.NoError(t, users.Create(ctx, coach))
	club := model.NewClub("Pulse", "CODE", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))

	w := model.NewPlanWeek(club.ID, 1, "1–7 Jul", "40 km")
	require.NoError(t, repo.Create(ctx, w))

	got, err := repo.GetByClubAndIndex(ctx, club.ID, 1)
	require.NoError(t, err)
	require.Equal(t, "40 km", got.PlanLabel)

	list, err := repo.FindByClub(ctx, club.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)

	updated := model.NewPlanWeek(club.ID, 1, "1–7 Jul", "45 km")
	require.NoError(t, repo.Upsert(ctx, updated))
	got, err = repo.GetByClubAndIndex(ctx, club.ID, 1)
	require.NoError(t, err)
	require.Equal(t, "45 km", got.PlanLabel)

	_, err = repo.GetByClubAndIndex(ctx, club.ID, 99)
	require.ErrorIs(t, err, model.ErrNotFound)

	w2 := model.NewPlanWeek(club.ID, 2, "8–14 Jul", "30 km")
	require.NoError(t, repo.Upsert(ctx, w2))
	list, err = repo.FindByClub(ctx, club.ID)
	require.NoError(t, err)
	require.Len(t, list, 2)
}
