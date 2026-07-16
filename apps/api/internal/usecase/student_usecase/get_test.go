//go:build unit

package student_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/student_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/student_usecase/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type usecaseMocks struct {
	clubRepo       *mocks.ClubRepo
	userRepo       *mocks.UserRepo
	membershipRepo *mocks.MembershipRepo
	activityRepo   *mocks.ActivityRepo
	workoutRepo    *mocks.WorkoutRepo
	planWeekRepo   *mocks.PlanWeekRepo
	announceRepo   *mocks.AnnounceRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		clubRepo:       mocks.NewClubRepo(t),
		userRepo:       mocks.NewUserRepo(t),
		membershipRepo: mocks.NewMembershipRepo(t),
		activityRepo:   mocks.NewActivityRepo(t),
		workoutRepo:    mocks.NewWorkoutRepo(t),
		planWeekRepo:   mocks.NewPlanWeekRepo(t),
		announceRepo:   mocks.NewAnnounceRepo(t),
	}
}

func TestGetIncludesWorkoutIDAndSummary(t *testing.T) {
	t.Parallel()
	coachID := uuid.New()
	studentID := uuid.New()
	clubID := uuid.New()
	wid := uuid.New()
	aid := uuid.New()

	club := &model.Club{ID: clubID, CoachID: coachID, Name: "Club"}
	membership := &model.Membership{UserID: studentID, ClubID: clubID, Status: model.MembershipActive}
	usr := &model.User{ID: studentID, Name: "Иван Тест", Role: model.RoleAthlete}

	started := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	act := model.NewActivity(studentID, "Run", "сегодня", 8.2, "40:00", "5:00", 0, 0, 0, "", 0, 0, 0, 0)
	act.ID = aid
	act.StartedAt = &started

	w := model.NewWorkout(studentID, model.WorkoutPlan, "Чт", "Лёгкий", "Кросс", 10, "", 0)
	w.ID = wid
	w.Status = model.WorkoutStatusCompleted
	w.CompletedActivityID = &aid
	w.ScheduledDate = &started

	m := newMocks(t)
	m.clubRepo.EXPECT().GetByCoachID(mock.Anything, coachID).Return(club, nil).Once()
	m.membershipRepo.EXPECT().GetByUserAndClub(mock.Anything, studentID, clubID).Return(membership, nil).Once()
	m.userRepo.EXPECT().GetByID(mock.Anything, studentID).Return(usr, nil).Once()
	m.activityRepo.EXPECT().SumDistByUserSince(mock.Anything, studentID, mock.Anything).Return(8.2, nil).Once()
	m.workoutRepo.EXPECT().SumPlanDistByUserWeek(mock.Anything, studentID, 0).Return(10.0, nil).Once()
	m.announceRepo.EXPECT().NextLabelForAthlete(mock.Anything, clubID, studentID).Return("Сб · Парк", nil).Once()
	m.activityRepo.EXPECT().FindByUser(mock.Anything, studentID).Return([]*model.Activity{act}, nil).Once()
	m.workoutRepo.EXPECT().FindByUser(mock.Anything, studentID).Return([]*model.Workout{w}, nil).Once()
	m.activityRepo.EXPECT().GetByID(mock.Anything, aid).Return(act, nil).Once()

	uc := student_usecase.NewUseCase(m.clubRepo, m.userRepo, m.membershipRepo, m.activityRepo, m.workoutRepo, m.planWeekRepo, m.announceRepo)
	view, err := uc.Get(context.Background(), coachID, studentID, 2026, 7)
	require.NoError(t, err)
	require.Len(t, view.PlanDays, 1)
	require.Equal(t, wid.String(), view.PlanDays[0].WorkoutID)
	require.Equal(t, "2026-07-16", view.PlanDays[0].ScheduledDate)
	require.Equal(t, 8.2, view.PlanDays[0].ActualKm)
	require.Equal(t, aid.String(), view.PlanDays[0].ActivityID)
	require.Equal(t, 2026, view.Year)
	require.Equal(t, 7, view.Month)
}
