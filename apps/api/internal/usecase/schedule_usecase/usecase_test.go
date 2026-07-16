//go:build unit

package schedule_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase/mocks"
)

type usecaseMocks struct {
	announceRepo   *mocks.AnnounceRepo
	clubRepo       *mocks.ClubResolver
	membershipRepo *mocks.MembershipRepo
	workoutRepo    *mocks.WorkoutRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		announceRepo:   mocks.NewAnnounceRepo(t),
		clubRepo:       mocks.NewClubResolver(t),
		membershipRepo: mocks.NewMembershipRepo(t),
		workoutRepo:    mocks.NewWorkoutRepo(t),
	}
}
