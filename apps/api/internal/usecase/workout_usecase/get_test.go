//go:build unit

package workout_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetWithFact(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	wid := uuid.New()
	aid := uuid.New()
	w := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 10, "", 0)
	w.ID = wid
	w.Status = model.WorkoutStatusCompleted
	w.CompletedActivityID = &aid

	a := model.NewActivity(uid, "Кросс", "Ср", 9.5, "50:00", "5:15", 140, 0, 0, "", 0, 0, 0, 0)
	a.ID = aid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()
	m.activityRepo.EXPECT().GetByID(mock.Anything, aid).Return(a, nil).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	view, err := uc.Get(context.Background(), uid, wid)
	require.NoError(t, err)
	require.Equal(t, 10.0, view.PlannedKm)
	require.NotNil(t, view.ActualKm)
	require.Equal(t, 9.5, *view.ActualKm)
	require.Equal(t, "5:15", view.ActualPace)
	require.NotNil(t, view.Fact)
	require.Equal(t, aid, view.Fact.ID)
}

func TestGetWithoutFact(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	wid := uuid.New()
	w := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 10, "", 0)
	w.ID = wid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	view, err := uc.Get(context.Background(), uid, wid)
	require.NoError(t, err)
	require.Nil(t, view.Fact)
	require.Nil(t, view.ActualKm)
}
