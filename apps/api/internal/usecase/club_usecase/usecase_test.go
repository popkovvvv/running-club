//go:build unit

package club_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/club_usecase/mocks"
)

type usecaseMocks struct {
	clubRepo       *mocks.ClubRepo
	membershipRepo *mocks.MembershipRepo
	userRepo       *mocks.UserRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		clubRepo:       mocks.NewClubRepo(t),
		membershipRepo: mocks.NewMembershipRepo(t),
		userRepo:       mocks.NewUserRepo(t),
	}
}
