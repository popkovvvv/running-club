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

func TestListActivities(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	activity := model.NewActivity(uid, "Утренний кросс", "Сегодня", 10.5, "55 мин", "5:15", 150, 3, 1, "M0 0", 0, 0, 1, 1)

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantLen int
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindByUser(mock.Anything, uid).Return([]*model.Activity{activity}, nil).Once()
			},
			wantLen: 1,
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindByUser(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
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
			items, err := uc.ListActivities(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, items, tt.wantLen)
			require.Equal(t, "Утренний кросс", items[0].Title)
			require.Equal(t, "10.5", items[0].Dist)
		})
	}
}
