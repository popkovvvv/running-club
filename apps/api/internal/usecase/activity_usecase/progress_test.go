//go:build unit

package activity_usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
	"github.com/nikpopkov/running-club/api/internal/usecase/activity_usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProgress(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	months := []model.MonthStat{model.NewMonthStat("Июл", 120, 12, "5:30", "+10")}

	tests := []struct {
		name    string
		before  func(m usecaseMocks)
		wantErr error
	}{
		{
			name: "ok",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindMonthStats(mock.Anything, uid).Return(months, nil).Once()
			},
		},
		{
			name: "repo_error",
			before: func(m usecaseMocks) {
				m.activityRepo.EXPECT().FindMonthStats(mock.Anything, uid).Return(nil, model.ErrNotFound).Once()
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
			uc := activity_usecase.NewUseCase(m.activityRepo)
			res, err := uc.Progress(context.Background(), uid)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 612.0, res.YearKm)
			require.Len(t, res.Months, 1)
			require.Equal(t, "Июл", res.Months[0].M)
		})
	}
}
