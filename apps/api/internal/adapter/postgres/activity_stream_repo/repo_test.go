//go:build integration

package activity_stream_repo_test

import (
	"context"
	"testing"

	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/activity_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/activity_stream_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestReplaceAndFind(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	activities := activity_repo.NewRepo(pool)
	repo := activity_stream_repo.NewRepo(pool)
	ctx := context.Background()

	u := model.NewUser("Nik", "nik@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, u))
	a := model.NewActivity(u.ID, "Run", "today", 5, "30:00", "6:00", 140, 0, 0, "", 0, 0, 1, 1)
	require.NoError(t, activities.Create(ctx, a))

	s1 := model.NewActivityStream(a.ID, "heartrate", `[120,130]`)
	s2 := model.NewActivityStream(a.ID, "altitude", `[10,20]`)
	require.NoError(t, repo.ReplaceByActivityID(ctx, a.ID, []*model.ActivityStream{s1, s2}))

	got, err := repo.FindByActivityID(ctx, a.ID)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, "altitude", got[0].Type)
	require.JSONEq(t, `[10,20]`, got[0].DataJSON)

	s3 := model.NewActivityStream(a.ID, "cadence", `[80]`)
	require.NoError(t, repo.ReplaceByActivityID(ctx, a.ID, []*model.ActivityStream{s3}))
	got, err = repo.FindByActivityID(ctx, a.ID)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "cadence", got[0].Type)
}
