//go:build integration

package workout_repo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/activity_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/announce_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/club_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/testutil"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/user_repo"
	"github.com/nikpopkov/running-club/api/internal/adapter/postgres/workout_repo"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestWorkoutCRUD(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	repo := workout_repo.NewRepo(pool)
	ctx := context.Background()

	u := model.NewUser("Nik", "nik@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, u))

	w := model.NewWorkout(u.ID, model.WorkoutPlan, "Mon", "easy", "Easy run", 10, "Z2", 1)
	w.AddSegment("run", "main", 10, "5:30", 0)
	require.NoError(t, repo.Create(ctx, w))

	got, err := repo.GetByID(ctx, w.ID)
	require.NoError(t, err)
	require.Equal(t, "Easy run", got.Title)
	require.Len(t, got.Segments, 1)

	byWeek, err := repo.FindByUserWeek(ctx, u.ID, 1, model.WorkoutPlan)
	require.NoError(t, err)
	require.Len(t, byWeek, 1)

	all, err := repo.FindByUser(ctx, u.ID)
	require.NoError(t, err)
	require.Len(t, all, 1)

	own := model.NewWorkout(u.ID, model.WorkoutOwn, "Tue", "own", "Own", 5, "Z2", 0)
	require.NoError(t, repo.Create(ctx, own))
	owns, err := repo.FindOwnByUser(ctx, u.ID)
	require.NoError(t, err)
	require.Len(t, owns, 1)

	sum, err := repo.SumPlanDistByUserWeek(ctx, u.ID, 1)
	require.NoError(t, err)
	require.Equal(t, 20.0, sum)

	rpe := 7
	got.Status = model.WorkoutStatusCompleted
	got.RPE = &rpe
	got.AthleteReport = "ok"
	require.NoError(t, repo.Update(ctx, got))

	completed, err := repo.FindCompletedWithoutActivity(ctx, u.ID)
	require.NoError(t, err)
	require.Len(t, completed, 1)

	require.NoError(t, repo.Delete(ctx, own.ID))
	_, err = repo.GetByID(ctx, own.ID)
	require.ErrorIs(t, err, model.ErrNotFound)
	require.ErrorIs(t, repo.Update(ctx, &model.Workout{ID: uuid.New()}), model.ErrNotFound)
}

func TestWorkoutTemplatesAndAnnounce(t *testing.T) {
	pool := testutil.Open(t)
	users := user_repo.NewRepo(pool)
	clubs := club_repo.NewRepo(pool)
	announces := announce_repo.NewRepo(pool)
	activities := activity_repo.NewRepo(pool)
	repo := workout_repo.NewRepo(pool)
	ctx := context.Background()

	coach := model.NewUser("Coach", "coach@test.run", "hash", model.RoleCoach)
	athlete := model.NewUser("Ath", "ath@test.run", "hash", model.RoleAthlete)
	require.NoError(t, users.Create(ctx, coach))
	require.NoError(t, users.Create(ctx, athlete))
	club := model.NewClub("Pulse", "CODE", "#fff", coach.ID)
	require.NoError(t, clubs.Create(ctx, club))

	tpl := model.NewWorkout(coach.ID, model.WorkoutPlan, "Wed", "tempo", "Tempo", 12, "Z3", 2)
	tpl.ClubID = &club.ID
	tpl.IsClubTemplate = true
	tpl.AddSegment("run", "tempo", 12, "4:30", 0)
	require.NoError(t, repo.ReplaceClubTemplates(ctx, club.ID, 2, []*model.Workout{tpl}))

	templates, err := repo.FindClubTemplates(ctx, club.ID, 2)
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Len(t, templates[0].Segments, 1)

	plan := model.NewWorkout(athlete.ID, model.WorkoutPlan, "Thu", "easy", "Plan", 8, "Z2", 2)
	require.NoError(t, repo.Create(ctx, plan))
	require.NoError(t, repo.DeleteClubAssignedPlans(ctx, athlete.ID, 2))
	_, err = repo.GetByID(ctx, plan.ID)
	require.ErrorIs(t, err, model.ErrNotFound)

	ann := model.NewAnnounce(club.ID, "Park", "Sat", "09:00", "A", "", nil)
	require.NoError(t, announces.Create(ctx, ann))
	linked := model.NewWorkout(athlete.ID, model.WorkoutPlan, "Sat", "group", "Group", 10, "Z2", 3)
	linked.AnnounceID = &ann.ID
	require.NoError(t, repo.Create(ctx, linked))
	require.NoError(t, repo.DeleteByUserAndAnnounce(ctx, athlete.ID, ann.ID))
	_, err = repo.GetByID(ctx, linked.ID)
	require.ErrorIs(t, err, model.ErrNotFound)

	act := model.NewActivity(athlete.ID, "Run", "today", 10, "50:00", "5:00", 140, 0, 0, "", 0, 0, 1, 1)
	require.NoError(t, activities.Create(ctx, act))
	w := model.NewWorkout(athlete.ID, model.WorkoutPlan, "Fri", "easy", "Linked", 10, "Z2", 1)
	w.CompletedActivityID = &act.ID
	w.Status = model.WorkoutStatusCompleted
	require.NoError(t, repo.Create(ctx, w))
	found, err := repo.FindByCompletedActivity(ctx, act.ID)
	require.NoError(t, err)
	require.Equal(t, w.ID, found.ID)
}
