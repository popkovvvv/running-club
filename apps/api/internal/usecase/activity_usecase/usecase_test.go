//go:build unit

package activity_usecase_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase/mocks"
)

type usecaseMocks struct {
	activityRepo *mocks.ActivityRepo
}

func newMocks(t *testing.T) usecaseMocks {
	t.Helper()
	return usecaseMocks{
		activityRepo: mocks.NewActivityRepo(t),
	}
}
