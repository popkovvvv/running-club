//go:build unit

package plan_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/plan_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/plan_usecase/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type usecaseMocks struct {
	planWeekRepo   *mocks.PlanWeekRepo
	workoutRepo    *mocks.WorkoutRepo
	clubRepo       *mocks.ClubRepo
	userRepo       *mocks.UserRepo
	membershipRepo *mocks.MembershipRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		planWeekRepo:   mocks.NewPlanWeekRepo(t),
		workoutRepo:    mocks.NewWorkoutRepo(t),
		clubRepo:       mocks.NewClubRepo(t),
		userRepo:       mocks.NewUserRepo(t),
		membershipRepo: mocks.NewMembershipRepo(t),
	}
}

func TestPublish(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()
	athleteID := uuid.New()
	template := &model.Workout{
		ID:          uuid.New(),
		ClubID:      &clubID,
		UserID:      coachID,
		Kind:        model.WorkoutPlan,
		WorkoutType: model.WorkoutTypeEasy,
		Title:       "Easy run",
		DistKm:      6,
		WeekIndex:   0,
		Status:      model.WorkoutStatusPlanned,
		IsClubTemplate: true,
	}

	m := newMocks(t)
	m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(&model.Club{ID: clubID}, nil).Once()
	m.workoutRepo.EXPECT().FindClubTemplates(mock.Anything, clubID, 0).Return([]*model.Workout{template}, nil).Once()
	m.userRepo.EXPECT().FindAthletesByClub(mock.Anything, clubID).Return([]*model.User{{ID: athleteID, Name: "Athlete"}}, nil).Once()
	m.workoutRepo.EXPECT().DeleteClubAssignedPlans(mock.Anything, athleteID, 0).Return(nil).Once()
	m.workoutRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Workout")).Return(nil).Once()

	uc := plan_usecase.NewUseCase(m.planWeekRepo, m.workoutRepo, m.clubRepo, m.userRepo, m.membershipRepo)
	err := uc.Publish(context.Background(), coachID, 0)
	require.NoError(t, err)
}

func TestSaveTemplate(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	clubID := uuid.New()

	m := newMocks(t)
	m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(&model.Club{ID: clubID}, nil).Once()
	m.workoutRepo.EXPECT().ReplaceClubTemplates(mock.Anything, clubID, 1, mock.Anything).Return(nil).Once()
	m.workoutRepo.EXPECT().FindClubTemplates(mock.Anything, clubID, 1).Return([]*model.Workout{}, nil).Once()

	uc := plan_usecase.NewUseCase(m.planWeekRepo, m.workoutRepo, m.clubRepo, m.userRepo, m.membershipRepo)
	res, err := uc.SaveTemplate(context.Background(), coachID, 1, dto.SaveTemplateRequest{
		Workouts: []dto.CreateWorkoutRequest{{Kind: "plan", Title: "Tempo", DistKm: 8, WeekIndex: 1}},
	})
	require.NoError(t, err)
	require.Equal(t, 1, res.WeekIndex)
}
