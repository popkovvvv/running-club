//go:build integration

package announce_repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/announce_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestAnnounceFlow(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	repo := announce_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	athlete := model.NewUser("Ath", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, coach))
	require.NoError(t, users.Create(ctx, athlete))
	club := model.NewClub("Pulse", "CODE", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))

	starts := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour)
	a := model.NewAnnounce(club.ID, "Park", "Сб", "09:00", "A", "note", &starts)
	require.NoError(t, repo.Create(ctx, a))

	got, err := repo.GetByID(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, "Park", got.Place)

	list, err := repo.FindByClub(ctx, club.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)

	require.NoError(t, repo.IncGoing(ctx, a.ID, 1))
	got, err = repo.GetByID(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, 1, got.GoingCount)

	signup := model.NewAnnounceSignup(a.ID, athlete.ID)
	require.NoError(t, repo.CreateSignup(ctx, signup))
	has, err := repo.HasSignup(ctx, a.ID, athlete.ID)
	require.NoError(t, err)
	require.True(t, has)

	athletes, err := repo.FindGoingAthletes(ctx, a.ID)
	require.NoError(t, err)
	require.Len(t, athletes, 1)
	require.Equal(t, athlete.ID, athletes[0].ID)

	label, err := repo.NextLabelForAthlete(ctx, club.ID, athlete.ID)
	require.NoError(t, err)
	require.Equal(t, "Сб · Park", label)

	require.NoError(t, repo.DeleteSignup(ctx, a.ID, athlete.ID))
	has, err = repo.HasSignup(ctx, a.ID, athlete.ID)
	require.NoError(t, err)
	require.False(t, has)

	_, err = repo.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, model.ErrNotFound)
	require.ErrorIs(t, repo.DeleteSignup(ctx, a.ID, athlete.ID), model.ErrNotFound)
	_, err = repo.NextLabelForAthlete(ctx, club.ID, athlete.ID)
	require.ErrorIs(t, err, model.ErrNotFound)
}
