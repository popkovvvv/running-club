//go:build unit

package activity_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRaces(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	upcoming := model.NewRace(uid, "Москва Марафон", "20.09", "42.2", "3:30", 40)
	finished := model.NewRace(uid, "Полумарафон", "01.05", "21.1", "1:40", 0)
	finished.Finished = true

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantLen int
		wantErr error
	}{
		{
			name: "ok_skips_finished",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindRaces(mock.Anything, uid).Return([]*model.Race{upcoming, finished}, nil).Once()
			},
			wantLen: 1,
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindRaces(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
			},
			wantErr: model.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := newMocks(t)
			if tt.before != nil {
				tt.before(m)
			}
			uc := newUC(m)
			items, err := uc.Races(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, items, tt.wantLen)
			require.Equal(t, "Москва Марафон", items[0].Name)
			require.Equal(t, "40", items[0].Days)
		})
	}
}
