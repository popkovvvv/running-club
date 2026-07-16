//go:build integration

package club_repo_test

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

func TestClubCRUD(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	memberships := membership_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	require.NoError(t, users.Create(ctx, coach))

	club := model.NewClub("Pulse", "INVITE-1", "#111", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))

	got, err := clubs.GetByID(ctx, club.ID)
	require.NoError(t, err)
	require.Equal(t, "Pulse", got.Name)

	byCode, err := clubs.GetByInviteCode(ctx, "INVITE-1")
	require.NoError(t, err)
	require.Equal(t, club.ID, byCode.ID)

	byCoach, err := clubs.GetByCoachID(ctx, coach.ID)
	require.NoError(t, err)
	require.Equal(t, club.ID, byCoach.ID)

	require.NoError(t, clubs.UpdateAccent(ctx, club.ID, "#c8ff34"))
	updated, err := clubs.GetByID(ctx, club.ID)
	require.NoError(t, err)
	require.Equal(t, "#c8ff34", updated.AccentHex)

	n, err := clubs.CountActiveStudents(ctx, club.ID)
	require.NoError(t, err)
	require.Equal(t, 0, n)

	athlete := model.NewUser("Ath", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, athlete))
	require.NoError(t, memberships.Create(ctx, model.NewMembership(athlete.ID, club.ID)))

	n, err = clubs.CountActiveStudents(ctx, club.ID)
	require.NoError(t, err)
	require.Equal(t, 1, n)

	_, err = clubs.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, model.ErrNotFound)
	require.ErrorIs(t, clubs.UpdateAccent(ctx, uuid.New(), "#000"), model.ErrNotFound)
}
