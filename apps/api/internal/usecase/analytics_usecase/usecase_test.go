//go:build unit

package analytics_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/analytics_usecase/mocks"
)

type usecaseMocks struct {
	clubRepo *mocks.ClubRepo
	userRepo *mocks.UserRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		clubRepo: mocks.NewClubRepo(t),
		userRepo: mocks.NewUserRepo(t),
	}
}
