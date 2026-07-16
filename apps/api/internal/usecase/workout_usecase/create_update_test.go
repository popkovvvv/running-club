//go:build unit

package workout_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateCompleteCreatesActivity(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	wid := uuid.New()
	status := "completed"
	w := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 6, "", 0)
	w.ID = wid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()
	m.activityRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Activity")).Run(func(ctx context.Context, a *model.Activity) {
		a.ID = uuid.New()
	}).Return(nil).Once()
	m.workoutRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Workout")).Return(nil).Once()
	m.activityRepo.EXPECT().GetByID(mock.Anything, mock.Anything).Return(
		model.NewActivity(uid, "Кросс", "Ср", 6, "~46 мин", "7:40", 0, 0, 0, "", 0, 0, 0, 0),
		nil,
	).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	view, err := uc.Update(context.Background(), uid, wid, dto.UpdateWorkoutRequest{Status: &status})
	require.NoError(t, err)
	require.Equal(t, "completed", view.Status)
	require.NotNil(t, view.CompletedActivityID)
}
