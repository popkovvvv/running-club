//go:build integration

package user_repo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/membership_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGet(t *testing.T) {
	pool := testutil.Open(t)
	repo := user_repo.NewRepo(pool)
	ctx := context.Background()

	u := model.NewUser("Nik", "nik@test.run", "hash", model.RoleAthlete)
	require.NoError(t, repo.Create(ctx, u))

	got, err := repo.GetByID(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, u.Email, got.Email)
	require.Equal(t, u.Name, got.Name)

	byEmail, err := repo.GetByEmail(ctx, u.Email)
	require.NoError(t, err)
	require.Equal(t, u.ID, byEmail.ID)

	_, err = repo.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, model.ErrNotFound)

	_, err = repo.GetByEmail(ctx, "missing@test.run")
	require.ErrorIs(t, err, model.ErrNotFound)
}

func TestFindAthletesByClub(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	memberships := membership_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	athlete := model.NewUser("Athlete", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, coach))
	require.NoError(t, users.Create(ctx, athlete))

	club := model.NewClub("Pulse", "CODE-1", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))
	require.NoError(t, memberships.Create(ctx, model.NewMembership(athlete.ID, club.ID)))

	found, err := users.FindAthletesByClub(ctx, club.ID)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Equal(t, athlete.ID, found[0].ID)
}
