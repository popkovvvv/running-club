//go:build unit

package auth_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/auth_usecase/mocks"
)

type usecaseMocks struct {
	userRepo       *mocks.UserRepo
	membershipRepo *mocks.MembershipRepo
	clubRepo       *mocks.ClubRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		userRepo:       mocks.NewUserRepo(t),
		membershipRepo: mocks.NewMembershipRepo(t),
		clubRepo:       mocks.NewClubRepo(t),
	}
}
