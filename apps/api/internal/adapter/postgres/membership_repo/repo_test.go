//go:build integration

package membership_repo_test

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

func TestMembershipCRUD(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	repo := membership_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	athlete := model.NewUser("Ath", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, coach))
	require.NoError(t, users.Create(ctx, athlete))

	club := model.NewClub("Pulse", "CODE", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))

	m := model.NewMembership(athlete.ID, club.ID)
	require.NoError(t, repo.Create(ctx, m))

	active, err := repo.GetActiveByUser(ctx, athlete.ID)
	require.NoError(t, err)
	require.Equal(t, m.ID, active.ID)

	byPair, err := repo.GetByUserAndClub(ctx, athlete.ID, club.ID)
	require.NoError(t, err)
	require.Equal(t, m.ID, byPair.ID)

	require.NoError(t, repo.UpdateStatus(ctx, m.ID, model.MembershipLeft))
	_, err = repo.GetActiveByUser(ctx, athlete.ID)
	require.ErrorIs(t, err, model.ErrNotFound)

	left, err := repo.GetByUserAndClub(ctx, athlete.ID, club.ID)
	require.NoError(t, err)
	require.Equal(t, model.MembershipLeft, left.Status)

	_, err = repo.GetByUserAndClub(ctx, uuid.New(), club.ID)
	require.ErrorIs(t, err, model.ErrNotFound)
	require.ErrorIs(t, repo.UpdateStatus(ctx, uuid.New(), model.MembershipActive), model.ErrNotFound)
}
