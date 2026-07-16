//go:build unit

package analytics_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/analytics_usecase/mocks"
)

type usecaseMocks struct {
	clubRepo     *mocks.ClubRepo
	userRepo     *mocks.UserRepo
	activityRepo *mocks.ActivityRepo
	planWeekRepo *mocks.PlanWeekRepo
	workoutRepo  *mocks.WorkoutRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		clubRepo:     mocks.NewClubRepo(t),
		userRepo:     mocks.NewUserRepo(t),
		activityRepo: mocks.NewActivityRepo(t),
		planWeekRepo: mocks.NewPlanWeekRepo(t),
		workoutRepo:  mocks.NewWorkoutRepo(t),
	}
}
