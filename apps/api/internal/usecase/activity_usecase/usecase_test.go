//go:build unit

package activity_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase"
	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase/mocks"
)

type usecaseMocks struct {
	activityRepo       *mocks.ActivityRepo
	activityStreamRepo *mocks.ActivityStreamRepo
	workoutRepo        *mocks.WorkoutRepo
	clubRepo           *mocks.ClubRepo
	membershipRepo     *mocks.MembershipRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		activityRepo:       mocks.NewActivityRepo(t),
		activityStreamRepo: mocks.NewActivityStreamRepo(t),
		workoutRepo:        mocks.NewWorkoutRepo(t),
		clubRepo:           mocks.NewClubRepo(t),
		membershipRepo:     mocks.NewMembershipRepo(t),
	}
}

func newUC(m usecaseMocks) *activity_usecase.UseCase {
	return activity_usecase.NewUseCase(m.activityRepo, m.activityStreamRepo, m.workoutRepo, m.clubRepo, m.membershipRepo)
}
