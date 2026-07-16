//go:build unit

package workout_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/workout_usecase/mocks"
)

type usecaseMocks struct {
	workoutRepo    *mocks.WorkoutRepo
	planWeekRepo   *mocks.PlanWeekRepo
	membershipRepo *mocks.MembershipRepo
	clubRepo       *mocks.ClubRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		workoutRepo:    mocks.NewWorkoutRepo(t),
		planWeekRepo:   mocks.NewPlanWeekRepo(t),
		membershipRepo: mocks.NewMembershipRepo(t),
		clubRepo:       mocks.NewClubRepo(t),
	}
}
