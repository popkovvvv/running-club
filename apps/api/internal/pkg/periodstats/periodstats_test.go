//go:build unit

package periodstats_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/pkg/periodstats"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	weekDay := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	monthDay := time.Date(2026, 7, 2, 10, 0, 0, 0, time.UTC)
	yearDay := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)

	a1 := model.NewActivity(uid, "W", "w", 10, "1:00", "5:00", 0, 0, 0, "", 0, 0, 0, 0)
	a1.StartedAt = &weekDay
	a2 := model.NewActivity(uid, "M", "m", 5, "0:30", "5:30", 0, 0, 0, "", 0, 0, 0, 0)
	a2.StartedAt = &monthDay
	a3 := model.NewActivity(uid, "Y", "y", 20, "2:00", "6:00", 0, 0, 0, "", 0, 0, 0, 0)
	a3.StartedAt = &yearDay

	wDone := model.NewWorkout(uid, model.WorkoutPlan, "Вт", "Лёгкий", "План", 10, "", 0)
	wDone.Status = model.WorkoutStatusCompleted
	wDone.ScheduledDate = &weekDay
	wOpen := model.NewWorkout(uid, model.WorkoutPlan, "Ср", "Лёгкий", "План 2", 8, "", 0)
	wOpen.Status = model.WorkoutStatusPlanned
	wOpen.ScheduledDate = &weekDay

	summary := periodstats.Build(now, []*model.Activity{a1, a2, a3}, []*model.Workout{wDone, wOpen})

	require.Equal(t, "10.0", summary.Week.Km)
	require.Equal(t, 1, summary.Week.Workouts)
	require.Equal(t, "50%", summary.Week.PlanCompletion)
	require.Equal(t, "5:00", summary.Week.AvgPace)

	require.Equal(t, "15.0", summary.Month.Km)
	require.Equal(t, 2, summary.Month.Workouts)

	require.Equal(t, "35.0", summary.Year.Km)
	require.Equal(t, 3, summary.Year.Workouts)
}

func TestBuildEmptyPlanCompletionDash(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	summary := periodstats.Build(now, nil, nil)
	require.Equal(t, "—", summary.Week.PlanCompletion)
	require.Equal(t, "—", summary.Week.AvgPace)
	require.Equal(t, "0.0", summary.Week.Km)
	require.Equal(t, 0, summary.Week.Workouts)
}
