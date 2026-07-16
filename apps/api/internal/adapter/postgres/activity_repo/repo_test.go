//go:build integration

package activity_repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/activity_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/membership_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestActivityCRUD(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	repo := activity_repo.NewRepo(pool)
	ctx := context.Background()

	u := model.NewUser("Nik", "nik@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, u))

	a := model.NewActivity(u.ID, "Morning", "today", 10, "50:00", "5:00", 150, 1, 0, "", 0, 0, 1, 1)
	a.Source = "strava"
	a.ExternalID = "ext-1"
	require.NoError(t, repo.Create(ctx, a))

	got, err := repo.GetByID(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, 10.0, got.DistKm)

	list, err := repo.FindByUser(ctx, u.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)

	got.Title = "Updated"
	got.DistKm = 12
	got.UpdatedAt = time.Now().UTC()
	require.NoError(t, repo.Update(ctx, got))
	updated, err := repo.GetByID(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Title)
	require.Equal(t, 12.0, updated.DistKm)

	byExt, err := repo.GetByUserSourceExternalID(ctx, u.ID, "strava", "ext-1")
	require.NoError(t, err)
	require.Equal(t, a.ID, byExt.ID)

	sum, err := repo.SumDistByUser(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, 12.0, sum)

	since := time.Now().UTC().Add(-time.Hour)
	sumSince, err := repo.SumDistByUserSince(ctx, u.ID, since)
	require.NoError(t, err)
	require.Equal(t, 12.0, sumSince)

	a.Title = "Upserted"
	a.DistKm = 15
	a.UpdatedAt = time.Now().UTC()
	require.NoError(t, repo.Upsert(ctx, a))
	byExt, err = repo.GetByUserSourceExternalID(ctx, u.ID, "strava", "ext-1")
	require.NoError(t, err)
	require.Equal(t, "Upserted", byExt.Title)

	require.NoError(t, repo.DeleteByUserSourceExternalID(ctx, u.ID, "strava", "ext-1"))
	_, err = repo.GetByID(ctx, a.ID)
	require.ErrorIs(t, err, model.ErrNotFound)

	_, err = repo.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, model.ErrNotFound)
	require.ErrorIs(t, repo.Update(ctx, &model.Activity{ID: uuid.New(), UpdatedAt: time.Now().UTC()}), model.ErrNotFound)
}

func TestActivityRelatedAndClubSum(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	memberships := membership_repo.NewRepo(pool)
	repo := activity_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	athlete := model.NewUser("Ath", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, coach))
	require.NoError(t, users.Create(ctx, athlete))
	club := model.NewClub("Pulse", "CODE", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))
	require.NoError(t, memberships.Create(ctx, model.NewMembership(athlete.ID, club.ID)))

	a := model.NewActivity(athlete.ID, "Run", "today", 8, "40:00", "5:00", 140, 0, 0, "", 0, 0, 1, 1)
	require.NoError(t, repo.Create(ctx, a))

	clubSum, err := repo.SumDistByClubAthletes(ctx, club.ID)
	require.NoError(t, err)
	require.Equal(t, 8.0, clubSum)

	pr := model.NewPR(athlete.ID, "5K", "20:00", "Jul")
	require.NoError(t, repo.CreatePR(ctx, pr))
	prs, err := repo.FindPRs(ctx, athlete.ID)
	require.NoError(t, err)
	require.Len(t, prs, 1)

	race := model.NewRace(athlete.ID, "City", "Aug", "10K", "45:00", 10)
	require.NoError(t, repo.CreateRace(ctx, race))
	races, err := repo.FindRaces(ctx, athlete.ID)
	require.NoError(t, err)
	require.Len(t, races, 1)

	require.NoError(t, repo.CreateMonthStat(ctx, athlete.ID, model.NewMonthStat("Jul", 100, 12, "5:30", "+10")))
	stats, err := repo.FindMonthStats(ctx, athlete.ID)
	require.NoError(t, err)
	require.Len(t, stats, 1)
	require.Equal(t, "Jul", stats[0].Month)
}
