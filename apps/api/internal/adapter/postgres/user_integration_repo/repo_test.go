//go:build integration

package user_integration_repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_integration_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestUserIntegrationUpsertAndGet(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	repo := user_integration_repo.NewRepo(pool)
	ctx := context.Background()

	u := model.NewUser("Nik", "nik@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, u))

	integration := model.NewUserIntegration(u.ID, model.IntegrationProviderStrava)
	integration.Status = model.IntegrationStatusActive
	integration.ExternalAthleteID = "ext-1"
	integration.AccessToken = "access"
	integration.RefreshToken = "refresh"
	integration.ExpiresAt = time.Now().UTC().Add(time.Hour)
	integration.Scopes = []string{"read", "activity:read"}
	require.NoError(t, repo.Upsert(ctx, integration))

	got, err := repo.GetByUserProvider(ctx, u.ID, model.IntegrationProviderStrava)
	require.NoError(t, err)
	require.Equal(t, "ext-1", got.ExternalAthleteID)
	require.Equal(t, []string{"read", "activity:read"}, got.Scopes)

	byExt, err := repo.GetByProviderExternalAthleteID(ctx, model.IntegrationProviderStrava, "ext-1")
	require.NoError(t, err)
	require.Equal(t, integration.ID, byExt.ID)

	integration.AccessToken = "access-2"
	integration.UpdatedAt = time.Now().UTC()
	require.NoError(t, repo.Upsert(ctx, integration))
	got, err = repo.GetByUserProvider(ctx, u.ID, model.IntegrationProviderStrava)
	require.NoError(t, err)
	require.Equal(t, "access-2", got.AccessToken)

	_, err = repo.GetByUserProvider(ctx, uuid.New(), model.IntegrationProviderStrava)
	require.ErrorIs(t, err, model.ErrNotFound)
}
