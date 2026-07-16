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

func TestPRs(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	pr := model.NewPR(uid, "5K", "18:30", "01.06")

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantLen int
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindPRs(mock.Anything, uid).Return([]*model.PR{pr}, nil).Once()
			},
			wantLen: 1,
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindPRs(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
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
			items, err := uc.PRs(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, items, tt.wantLen)
			require.Equal(t, "5K", items[0].D)
			require.Equal(t, "18:30", items[0].T)
		})
	}
}
