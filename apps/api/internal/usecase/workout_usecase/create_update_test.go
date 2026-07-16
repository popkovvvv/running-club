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
	activityID := uuid.New()
	status := "completed"
	rpe := 6
	report := "после тренировки забыл про заминку"
	w := model.NewWorkout(uid, model.WorkoutOwn, "Ср", "Лёгкий", "Кросс", 6, "", 0)
	w.ID = wid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()
	m.activityRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Activity")).Run(func(ctx context.Context, a *model.Activity) {
		a.ID = activityID
	}).Return(nil).Once()
	m.workoutRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Workout")).Run(func(ctx context.Context, updated *model.Workout) {
		require.Equal(t, model.WorkoutStatusCompleted, updated.Status)
		require.NotNil(t, updated.RPE)
		require.Equal(t, 6, *updated.RPE)
		require.Equal(t, report, updated.AthleteReport)
		require.Equal(t, activityID, *updated.CompletedActivityID)
	}).Return(nil).Once()
	fact := model.NewActivity(uid, "Кросс", "Ср", 6, "~46 мин", "7:40", 0, 0, 0, "", 0, 0, 0, 0)
	fact.ID = activityID
	m.activityRepo.EXPECT().GetByID(mock.Anything, activityID).Return(fact, nil).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	view, err := uc.Update(context.Background(), uid, model.RoleAthlete, wid, dto.UpdateWorkoutRequest{
		Status:        &status,
		RPE:           &rpe,
		AthleteReport: &report,
	})
	require.NoError(t, err)
	require.Equal(t, "completed", view.Status)
	require.NotNil(t, view.CompletedActivityID)
	require.NotNil(t, view.RPE)
	require.Equal(t, 6, *view.RPE)
	require.Equal(t, report, view.AthleteReport)
}

func TestUpdateCoachComment(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	athleteID := uuid.New()
	wid := uuid.New()
	clubID := uuid.New()
	comment := "Хорошо"
	w := model.NewWorkout(athleteID, model.WorkoutPlan, "Вт", "Кросс", "Групповая", 5, "", 29)
	w.ID = wid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()
	m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(&model.Club{ID: clubID, CoachID: coachID}, nil).Once()
	m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, athleteID, clubID).Return(
		&model.Membership{UserID: athleteID, ClubID: clubID, Status: model.MembershipActive}, nil,
	).Once()
	m.workoutRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Workout")).Run(func(ctx context.Context, updated *model.Workout) {
		require.Equal(t, comment, updated.CoachComment)
	}).Return(nil).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	view, err := uc.Update(context.Background(), coachID, model.RoleCoach, wid, dto.UpdateWorkoutRequest{
		CoachComment: &comment,
	})
	require.NoError(t, err)
	require.Equal(t, comment, view.CoachComment)
}

func TestUpdateCoachCannotSetRPE(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	athleteID := uuid.New()
	wid := uuid.New()
	clubID := uuid.New()
	rpe := 5
	comment := "ok"
	w := model.NewWorkout(athleteID, model.WorkoutPlan, "Вт", "Кросс", "Групповая", 5, "", 29)
	w.ID = wid

	m := newMocks(t)
	m.workoutRepo.EXPECT().GetByID(mock.Anything, wid).Return(w, nil).Once()
	m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(&model.Club{ID: clubID, CoachID: coachID}, nil).Once()
	m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, athleteID, clubID).Return(
		&model.Membership{UserID: athleteID, ClubID: clubID, Status: model.MembershipActive}, nil,
	).Once()

	uc := workout_usecase.NewUseCase(m.workoutRepo, m.planWeekRepo, m.membershipRepo, m.clubRepo, m.activityRepo)
	_, err := uc.Update(context.Background(), coachID, model.RoleCoach, wid, dto.UpdateWorkoutRequest{
		CoachComment: &comment,
		RPE:          &rpe,
	})
	require.ErrorIs(t, err, model.ErrForbidden)
}
